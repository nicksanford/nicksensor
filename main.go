package main

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"go.uber.org/atomic"
	goutils "go.viam.com/utils"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

var Model = resource.NewModel("ncs", "sensor", "nicksensor")
var (
	Reset = "\033[0m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
)

type fake struct {
	mu       sync.Mutex
	bootTime time.Time
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	counter atomic.Uint64
	logger  logging.Logger
}
type Config struct {
	StartupTime string `json:"startup_time,omitempty"`
}

func (c *Config) Validate(path string) ([]string, error) {
	if c.StartupTime == "" {
		return nil, nil
	}

	if _, err := time.ParseDuration(c.StartupTime); err != nil {
		return nil, err
	}

	return nil, nil
}

func maybeSleep(c *Config, logger logging.Logger) {
	if c.StartupTime == "" {
		return
	}
	d, err := time.ParseDuration(c.StartupTime)
	if err != nil {
		return
	}
	logger.Infof(Cyan+"startup sleeping for %s"+Reset, d)
	time.Sleep(d)
}

func newSensor(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	c, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, err
	}

	maybeSleep(c, logger)
	return &fake{
		Named:    conf.ResourceName().AsNamed(),
		bootTime: time.Now(),
		logger:   logger,
	}, nil
}

func (f *fake) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(extra) > 0 {
		f.logger.Debugf("extra: %#v", extra)
	}
	count := f.counter.Add(1)
	now := time.Now()
	return map[string]interface{}{
		"now_unix":       now.Unix(),
		"boot_time_unix": f.bootTime.Unix(),
		"now_unix_micro": now.UnixMicro(),
		"call_count":     count,
	}, nil
}

func (f *fake) DoCommand(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, boom := extra["boom"]
	duration, hang := extra["hang"]
	switch {
	case boom:
		f.logger.Info(Cyan + "Boom" + Reset)
		os.Exit(1)
	case hang:
		dStr, ok := duration.(string)
		if !ok {
			return nil, errors.New("expected hang duration to be a string")
		}

		dur, err := time.ParseDuration(dStr)
		if err != nil {
			return nil, err
		}
		f.logger.Infof(Cyan+"DoCommand hanging for %s"+Reset, dur)
		time.Sleep(dur)
	}
	return nil, nil
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *Config]{Constructor: newSensor})

	module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	if err := module.AddModelFromRegistry(ctx, sensor.API, Model); err != nil {
		return err
	}

	err = module.Start(ctx)
	defer module.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func main() {
	goutils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs(Model.String()))
}
