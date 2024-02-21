package k6

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/event"
	"go.k6.io/k6/execution"
	"go.k6.io/k6/execution/local"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/trace"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/metrics/engine"
	"go.k6.io/k6/output"
)

const (
	// We use an excessively high timeout to wait for event processing to complete,
	// since prematurely proceeding before it is done could create bigger problems.
	// In practice, this effectively acts as no timeout, and the user will have to
	// kill k6 if a hang happens, which is the behavior without events anyway.
	waitEventDoneTimeout = 30 * time.Minute

	// This timeout should be long enough to flush all remaining traces, but still
	// provides a safeguard to not block indefinitely.
	waitForTracerProviderStopTimeout = 3 * time.Minute
)

type Execution struct {
	FS     fsext.Fs
	Env    map[string]string
	Events *event.System

	LoggerCompat *logrus.Logger

	Script string

	Logger *slog.Logger
}

func NewExecution(logger *slog.Logger, script string) *Execution {
	loggerCompat := logrus.StandardLogger()

	return &Execution{
		FS:           fsext.NewOsFs(),
		Env:          BuildEnvMap(os.Environ()),
		Events:       event.NewEventSystem(100, loggerCompat),
		LoggerCompat: loggerCompat,
		Logger:       logger,
		Script:       script,
	}
}

func (e *Execution) loadLocalTest() (*loadedAndConfiguredTest, *local.Controller, error) {
	data := []byte(e.Script)

	fileSystems := loader.CreateFilesystems(e.FS)

	err := fsext.WriteFile(fileSystems["file"].(fsext.CacheLayerGetter).GetCachingFs(), "/-", data, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("caching data read from -: %w", err)
	}

	src := &loader.SourceData{URL: &url.URL{Path: "/-", Scheme: "file"}, Data: data}
	e.Logger.Debug(
		"successfully loaded bytes!",
		"bytes", len(src.Data),
	)

	e.Logger.Debug("Gathering k6 runtime options...")
	runtimeOptions := lib.RuntimeOptions{}

	registry := metrics.NewRegistry()
	state := &lib.TestPreInitState{
		Logger:         e.LoggerCompat,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Events:         e.Events,
		LookupEnv: func(key string) (string, bool) {
			val, ok := e.Env[key]
			return val, ok
		},
	}

	test := &loadedTest{
		source:       src,
		fs:           e.FS,
		fileSystems:  fileSystems,
		preInitState: state,
		logger:       e.Logger,
		loggerCompat: e.LoggerCompat,
	}

	e.Logger.Debug("Initializing k6 runner...")
	if err := test.initializeFirstRunner(); err != nil {
		return nil, nil, fmt.Errorf("could not initialize: %w", err)
	}
	e.Logger.Debug("Runner successfully initialized!")

	configuredTest, err := test.consolidateDeriveAndValidateConfig()
	if err != nil {
		return nil, nil, err
	}

	controller := local.NewController()
	return configuredTest, controller, nil
}

func (e *Execution) setupTracerProvider(ctx context.Context, test *loadedAndConfiguredTest) error {
	ro := test.preInitState.RuntimeOptions
	if ro.TracesOutput.String == "" {
		test.preInitState.TracerProvider = trace.NewNoopTracerProvider()
		return nil
	}

	tp, err := trace.TracerProviderFromConfigLine(ctx, ro.TracesOutput.String)
	if err != nil {
		return err
	}
	test.preInitState.TracerProvider = tp

	return nil
}

func (e *Execution) Start(ctx context.Context) error {
	var err error
	var logger logrus.FieldLogger = logrus.StandardLogger()

	globalCtx, globalCancel := context.WithCancel(ctx)
	defer globalCancel()

	// lingerCtx is cancelled by Ctrl+C, and is used to wait for that event when
	// k6 was started with the --linger option.
	lingerCtx, lingerCancel := context.WithCancel(globalCtx)
	defer lingerCancel()

	// runCtx is used for the test run execution and is created with the special
	// execution.NewTestRunContext() function so that it can be aborted even
	// from sub-contexts while also attaching a reason for the abort.
	runCtx, runAbort := execution.NewTestRunContext(lingerCtx, logger)

	emitEvent := func(evt *event.Event) func() {
		waitDone := e.Events.Emit(evt)
		return func() {
			waitCtx, waitCancel := context.WithTimeout(globalCtx, waitEventDoneTimeout)
			defer waitCancel()
			if werr := waitDone(waitCtx); werr != nil {
				logger.WithError(werr).Warn()
			}
		}
	}

	defer func() {
		waitExitDone := emitEvent(&event.Event{
			Type: event.Exit,
			Data: &event.ExitData{Error: err},
		})
		waitExitDone()
		e.Events.UnsubscribeAll()
	}()

	configuredTest, controller, err := e.loadLocalTest()
	if err != nil {
		return err
	}

	if err = e.setupTracerProvider(globalCtx, configuredTest); err != nil {
		return err
	}
	waitTracesFlushed := func() {
		ctx, cancel := context.WithTimeout(globalCtx, waitForTracerProviderStopTimeout)
		defer cancel()
		if tpErr := configuredTest.preInitState.TracerProvider.Shutdown(ctx); tpErr != nil {
			logger.Errorf("The tracer provider didn't stop gracefully: %v", tpErr)
		}
	}

	// Write the full consolidated *and derived* options back to the Runner.
	conf := configuredTest.config
	testRunState, err := configuredTest.buildTestRunState(conf)
	if err != nil {
		return err
	}

	// Create a local execution scheduler wrapping the runner.
	logger.Debug("Initializing the execution scheduler...")
	execScheduler, err := execution.NewScheduler(testRunState, controller)
	if err != nil {
		return err
	}

	backgroundProcesses := &sync.WaitGroup{}
	defer backgroundProcesses.Wait()

	// Create all outputs.
	// executionPlan := execScheduler.GetExecutionPlan()
	outputs := []output.Output{}

	metricsEngine, err := engine.NewMetricsEngine(testRunState.Registry, logger)
	if err != nil {
		return err
	}

	// We'll need to pipe metrics to the MetricsEngine and process them if any
	// of these are enabled: thresholds, end-of-test summary
	shouldProcessMetrics := (!testRunState.RuntimeOptions.NoSummary.Bool ||
		!testRunState.RuntimeOptions.NoThresholds.Bool)
	var metricsIngester *engine.OutputIngester
	if shouldProcessMetrics {
		err = metricsEngine.InitSubMetricsAndThresholds(conf, testRunState.RuntimeOptions.NoThresholds.Bool)
		if err != nil {
			return err
		}
		// We'll need to pipe metrics to the MetricsEngine if either the
		// thresholds or the end-of-test summary are enabled.
		metricsIngester = metricsEngine.CreateIngester()
		outputs = append(outputs, metricsIngester)
	}

	executionState := execScheduler.GetState()
	if !testRunState.RuntimeOptions.NoSummary.Bool {
		defer func() {
			logger.Debug("Generating the end-of-test summary...")
			summaryResult, hsErr := configuredTest.initRunner.HandleSummary(globalCtx, &lib.Summary{
				Metrics:         metricsEngine.ObservedMetrics,
				RootGroup:       testRunState.Runner.GetDefaultGroup(),
				TestRunDuration: executionState.GetCurrentTestRunDuration(),
				NoColor:         true,
				UIState: lib.UIState{
					IsStdOutTTY: false,
					IsStdErrTTY: false,
				},
			})
			if hsErr == nil {
				for _, o := range summaryResult {
					_, err := io.Copy(os.Stdout, o)
					if err != nil {
						logger.WithError(err).Error("failed to write summary output")
					}
				}
			}
			if hsErr != nil {
				logger.WithError(hsErr).Error("failed to handle the end-of-test summary")
			}
		}()
	}

	waitInitDone := emitEvent(&event.Event{Type: event.Init})

	outputManager := output.NewManager(outputs, logger, func(err error) {
		if err != nil {
			logger.WithError(err).Error("Received error to stop from output")
		}
		// TODO: attach run status and exit code?
		runAbort(err)
	})
	samples := make(chan metrics.SampleContainer, configuredTest.config.MetricSamplesBufferSize.Int64)
	waitOutputsFlushed, stopOutputs, err := outputManager.Start(samples)
	if err != nil {
		return err
	}
	defer func() {
		logger.Debug("Stopping outputs...")
		// We call waitOutputsFlushed() below because the threshold calculations
		// need all of the metrics to be sent to the MetricsEngine before we can
		// calculate them one last time. We need the threshold calculated here,
		// since they may change the run status for the outputs.
		stopOutputs(err)
	}()

	if !testRunState.RuntimeOptions.NoThresholds.Bool {
		finalizeThresholds := metricsEngine.StartThresholdCalculations(
			metricsIngester, runAbort, executionState.GetCurrentTestRunDuration,
		)
		handleFinalThresholdCalculation := func() {
			// This gets called after the Samples channel has been closed and
			// the OutputManager has flushed all of the cached samples to
			// outputs (including MetricsEngine's ingester). So we are sure
			// there won't be any more metrics being sent.
			logger.Debug("Finalizing thresholds...")
			breachedThresholds := finalizeThresholds()
			if len(breachedThresholds) == 0 {
				return
			}
			tErr := errext.WithAbortReasonIfNone(
				errext.WithExitCodeIfNone(
					fmt.Errorf("thresholds on metrics '%s' have been crossed", strings.Join(breachedThresholds, ", ")),
					exitcodes.ThresholdsHaveFailed,
				), errext.AbortedByThresholdsAfterTestEnd)

			if err == nil {
				err = tErr
			} else {
				logger.WithError(tErr).Debug("Crossed thresholds, but test already exited with another error")
			}
		}
		if finalizeThresholds != nil {
			defer handleFinalThresholdCalculation()
		}
	}

	defer func() {
		logger.Debug("Waiting for metrics and traces processing to finish...")
		close(samples)

		ww := [...]func(){
			waitOutputsFlushed,
			waitTracesFlushed,
		}
		var wg sync.WaitGroup
		wg.Add(len(ww))
		for _, w := range ww {
			w := w
			go func() {
				w()
				wg.Done()
			}()
		}
		wg.Wait()

		logger.Debug("Metrics and traces processing finished!")
	}()

	// Initialize the VUs and executors
	stopVUEmission, err := execScheduler.Init(runCtx, samples)
	if err != nil {
		return err
	}
	defer stopVUEmission()

	waitInitDone()

	waitTestStartDone := emitEvent(&event.Event{Type: event.TestStart})
	waitTestStartDone()

	// Start the test! However, we won't immediately return if there was an
	// error, we still have things to do.
	err = execScheduler.Run(globalCtx, runCtx, samples)

	waitTestEndDone := emitEvent(&event.Event{Type: event.TestEnd})
	defer waitTestEndDone()

	// Check what the execScheduler.Run() error is.
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		logger.WithError(err).Debug("Test finished with an error")
		return err
	}

	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 {
		logger.Warn("No script iterations fully finished, consider making the test duration longer")
	}

	logger.Debug("Test finished cleanly")

	return nil
}
