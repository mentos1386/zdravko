package k6

import (
	"crypto/x509"
	"fmt"
	"log/slog"
	"time"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

// loadedTest contains all of data, details and dependencies of a loaded
// k6 test, but without any config consolidation.
type loadedTest struct {
	source         *loader.SourceData
	fs             fsext.Fs
	fileSystems    map[string]fsext.Fs
	preInitState   *lib.TestPreInitState
	initRunner     lib.Runner // TODO: rename to something more appropriate
	moduleResolver *modules.ModuleResolver

	logger       *slog.Logger
	loggerCompat *logrus.Logger
}

func (lt *loadedTest) initializeFirstRunner() error {
	testPath := lt.source.URL.String()
	logger := lt.logger.With("test_path", testPath)

	logger.Debug("Trying to load as a JS test...")
	runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems)
	// TODO: should we use common.UnwrapGojaInterruptedError() here?
	if err != nil {
		return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
	}
	lt.initRunner = runner
	lt.moduleResolver = runner.Bundle.ModuleResolver
	return nil
}

func (lt *loadedTest) consolidateDeriveAndValidateConfig() (*loadedAndConfiguredTest, error) {
	lt.logger.Debug("Consolidating config layers...")

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

	lt.logger.Debug("Parsing thresholds and validating config...")
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

	config, err := executor.DeriveScenariosFromShortcuts(config, lt.loggerCompat)
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
