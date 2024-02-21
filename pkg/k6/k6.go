package k6

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/event"
	"go.k6.io/k6/execution"
	"go.k6.io/k6/execution/local"
	"go.k6.io/k6/js"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/trace"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/metrics/engine"
	"go.k6.io/k6/output"
)

const (
	testTypeJS = "js"
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

// loadedTest contains all of data, details and dependencies of a loaded
// k6 test, but without any config consolidation.
type loadedTest struct {
	sourceRootPath string // contains the raw string the user supplied
	pwd            string
	source         *loader.SourceData
	fs             fsext.Fs
	fileSystems    map[string]fsext.Fs
	preInitState   *lib.TestPreInitState
	initRunner     lib.Runner // TODO: rename to something more appropriate
	keyLogger      io.Closer
	moduleResolver *modules.ModuleResolver
}

func loadLocalTest(gs *state.GlobalState, script string) (*loadedTest, error) {
	src, fileSystems, pwd, err := readSource(gs, script)
	if err != nil {
		return nil, err
	}
	resolvedPath := src.URL.String()
	gs.Logger.Debugf(
		"successfully loaded %d bytes!",
		len(src.Data),
	)

	gs.Logger.Debugf("Gathering k6 runtime options...")
	runtimeOptions := lib.RuntimeOptions{}

	registry := metrics.NewRegistry()
	state := &lib.TestPreInitState{
		Logger:         gs.Logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Events:         gs.Events,
		LookupEnv: func(key string) (string, bool) {
			val, ok := gs.Env[key]
			return val, ok
		},
	}

	test := &loadedTest{
		pwd:            pwd,
		sourceRootPath: "-",
		source:         src,
		fs:             gs.FS,
		fileSystems:    fileSystems,
		preInitState:   state,
	}

	gs.Logger.Debugf("Initializing k6 runner for (%s)...", resolvedPath)
	if err := test.initializeFirstRunner(gs); err != nil {
		return nil, fmt.Errorf("could not initialize: %w", err)
	}
	gs.Logger.Debug("Runner successfully initialized!")
	return test, nil
}

func (lt *loadedTest) initializeFirstRunner(gs *state.GlobalState) error {
	testPath := lt.source.URL.String()
	logger := gs.Logger.WithField("test_path", testPath)

	testType := lt.preInitState.RuntimeOptions.TestType.String
	if testType == "" {
		testType = testTypeJS
	}

	// TODO: k6-cli also has TAR support which might be nice to have.
	switch testType {
	case testTypeJS:
		logger.Debug("Trying to load as a JS test...")
		runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems)
		// TODO: should we use common.UnwrapGojaInterruptedError() here?
		if err != nil {
			return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
		}
		lt.initRunner = runner
		lt.moduleResolver = runner.Bundle.ModuleResolver
		return nil
	default:
		return fmt.Errorf("unknown or unspecified test type '%s' for '%s'", testType, testPath)
	}
}

// readSource is a small wrapper around loader.ReadSource returning
// result of the load and filesystems map
func readSource(gs *state.GlobalState, script string) (*loader.SourceData, map[string]fsext.Fs, string, error) {
	data := []byte(script)

	filesystems := loader.CreateFilesystems(gs.FS)

	err := fsext.WriteFile(filesystems["file"].(fsext.CacheLayerGetter).GetCachingFs(), "/-", data, 0o644)
	if err != nil {
		return nil, nil, "", fmt.Errorf("caching data read from -: %w", err)
	}

	return &loader.SourceData{URL: &url.URL{Path: "/-", Scheme: "file"}, Data: data}, filesystems, "/", err
}

func (lt *loadedTest) consolidateDeriveAndValidateConfig(
	gs *state.GlobalState,
) (*loadedAndConfiguredTest, error) {
	gs.Logger.Debug("Consolidating config layers...")

	config := lib.Options{}

	config.Apply(lt.initRunner.GetOptions())
	if config.SystemTags == nil {
		config.SystemTags = &metrics.DefaultSystemTagSet
	}
	if config.SummaryTrendStats == nil {
		config.SummaryTrendStats = lib.DefaultSummaryTrendStats
	}
	defDNS := types.DefaultDNSConfig()
	if !config.DNS.TTL.Valid {
		config.DNS.TTL = defDNS.TTL
	}
	if !config.DNS.Select.Valid {
		config.DNS.Select = defDNS.Select
	}
	if !config.DNS.Policy.Valid {
		config.DNS.Policy = defDNS.Policy
	}
	if !config.SetupTimeout.Valid {
		config.SetupTimeout.Duration = types.Duration(60 * time.Second)
	}
	if !config.TeardownTimeout.Valid {
		config.TeardownTimeout.Duration = types.Duration(60 * time.Second)
	}

	gs.Logger.Debug("Parsing thresholds and validating config...")
	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.preInitState.RuntimeOptions.NoThresholds.Bool {
		for metricName, thresholdsDefinition := range config.Thresholds {
			err := thresholdsDefinition.Parse()
			if err != nil {
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			}

			err = thresholdsDefinition.Validate(metricName, lt.preInitState.Registry)
			if err != nil {
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			}
		}
	}

	config, err := executor.DeriveScenariosFromShortcuts(config, gs.Logger)
	if err == nil {
		errors := config.Validate()
		// FIXME: should show them all.
		if len(errors) > 0 {
			err = errors[0]
		}
	}
	if err != nil {
		return nil, err
	}

	return &loadedAndConfiguredTest{
		loadedTest: lt,
		config:     config,
	}, nil
}

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type loadedAndConfiguredTest struct {
	*loadedTest
	config lib.Options
}

func loadAndConfigureLocalTest(
	gs *state.GlobalState,
	script string,
) (*loadedAndConfiguredTest, error) {
	test, err := loadLocalTest(gs, script)
	if err != nil {
		return nil, err
	}

	return test.consolidateDeriveAndValidateConfig(gs)
}

// loadSystemCertPool attempts to load system certificates.
func loadSystemCertPool(logger logrus.FieldLogger) {
	if _, err := x509.SystemCertPool(); err != nil {
		logger.WithError(err).Warning("Unable to load system cert pool")
	}
}

func (lct *loadedAndConfiguredTest) buildTestRunState(
	configToReinject lib.Options,
) (*lib.TestRunState, error) {
	// This might be the full derived or just the consodlidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil {
		return nil, err
	}

	// it pre-loads system certificates to avoid doing it on the first TLS request.
	// This is done async to avoid blocking the rest of the loading process as it will not stop if it fails.
	go loadSystemCertPool(lct.preInitState.Logger)

	return &lib.TestRunState{
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.config,
		RunTags:          lct.preInitState.Registry.RootTagSet().WithTagsFromMap(configToReinject.RunTags),
	}, nil
}

type Execution struct {
	gs *state.GlobalState

	// TODO: figure out something more elegant?
	loadConfiguredTest func() (*loadedAndConfiguredTest, execution.Controller, error)
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

func NewExecution(gs *state.GlobalState, script string) *Execution {
	return &Execution{
		gs: gs,
		loadConfiguredTest: func() (*loadedAndConfiguredTest, execution.Controller, error) {
			test, err := loadAndConfigureLocalTest(gs, script)
			return test, local.NewController(), err
		},
	}
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
		waitDone := e.gs.Events.Emit(evt)
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
		e.gs.Events.UnsubscribeAll()
	}()

	test, controller, err := e.loadConfiguredTest()
	if err != nil {
		return err
	}
	if test.keyLogger != nil {
		defer func() {
			if klErr := test.keyLogger.Close(); klErr != nil {
				logger.WithError(klErr).Warn("Error while closing the SSLKEYLOGFILE")
			}
		}()
	}

	if err = e.setupTracerProvider(globalCtx, test); err != nil {
		return err
	}
	waitTracesFlushed := func() {
		ctx, cancel := context.WithTimeout(globalCtx, waitForTracerProviderStopTimeout)
		defer cancel()
		if tpErr := test.preInitState.TracerProvider.Shutdown(ctx); tpErr != nil {
			logger.Errorf("The tracer provider didn't stop gracefully: %v", tpErr)
		}
	}

	// Write the full consolidated *and derived* options back to the Runner.
	conf := test.config
	testRunState, err := test.buildTestRunState(conf)
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
			summaryResult, hsErr := test.initRunner.HandleSummary(globalCtx, &lib.Summary{
				Metrics:         metricsEngine.ObservedMetrics,
				RootGroup:       testRunState.Runner.GetDefaultGroup(),
				TestRunDuration: executionState.GetCurrentTestRunDuration(),
				NoColor:         e.gs.Flags.NoColor,
				UIState: lib.UIState{
					IsStdOutTTY: e.gs.Stdout.IsTTY,
					IsStdErrTTY: e.gs.Stderr.IsTTY,
				},
			})
			if hsErr == nil {
				for _, o := range summaryResult {
					_, err := io.Copy(e.gs.Stdout, o)
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
	samples := make(chan metrics.SampleContainer, test.config.MetricSamplesBufferSize.Int64)
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
