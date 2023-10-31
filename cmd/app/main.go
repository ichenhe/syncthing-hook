package main

import (
	"errors"
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/exevent"
	"github.com/ichenhe/syncthing-hook/hook"
	"github.com/ichenhe/syncthing-hook/stclient"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	var err error
	var appProfile *domain.AppProfile

	// read configuration
	configLoader := newConfigurationLoader(&defaultArgumentFetcher{})
	appProfile, err = configLoader.loadConfiguration()
	if err != nil {
		if errors.Is(err, errNoProfileSpecified) {
			log.Printf("please specify the profile with cmd or environment variable '%s'", _ProfileEnv)
			configLoader.printUsage()
			os.Exit(1)
			return
		}
		log.Fatalf("failed to load confirguration: %s", err)
		return
	}

	// create logger
	var logger *zap.SugaredLogger
	if logger_, err := createLogger(&appProfile.Log); err != nil {
		log.Fatalf("failed to create logger: %s", err)
		return
	} else {
		defer func(logger *zap.Logger) {
			_ = logger.Sync()
		}(logger_)
		logger = logger_.Sugar()
	}

	// create syncthing client
	logger.Infow("trying to connect Syncthing...", zap.String("url", appProfile.Syncthing.Url), zap.Bool("withToken", len(appProfile.Syncthing.ApiKey) > 0))
	var stClient domain.SyncthingClient
	stClient, err = stclient.NewSyncthing(appProfile.Syncthing.Url, appProfile.Syncthing.ApiKey, logger)
	if err != nil {
		logger.Fatal("failed to create syncthing client: ", err.Error())
		return
	}
	if status, err := stClient.GetSystemStatus(); err != nil {
		var e domain.StApiError
		if errors.As(err, &e) && e.HttpStatusCode == 403 {
			logger.Fatal("invalid syncthing api token, exit...")
			return
		}
		logger.Fatalf("failed to connect to Syncthing: %s, exit...", err)
		return
	} else {
		logger.Infow("connected to Syncthing", zap.String("id", status.MyID))
	}

	eventSource := exevent.NewEventSource(stClient, logger)
	actionExecutorCreator := exevent.NewActionExecutorCreator(logger)
	hookManager := hook.NewManager(eventSource, actionExecutorCreator, logger)

	logger.Info("loading hooks...")
	for i, h := range appProfile.Hooks {
		hookDef := domain.NewHookDefinition(h.Name, i)
		if err := hookManager.RegisterHook(&h, hookDef); err != nil {
			hookDef.AddToLogger(logger).Errorf("failed to regiter hook: %s", err)
		}
	}
	logger.Infof("%d hook loaded", len(appProfile.Hooks))

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	logger.Infow("received signal, shutting down...", zap.String("sig", sig.String()))
	hookManager.UnregisterAll()
	logger.Info("bye")
}

func createLogger(config *domain.LogConfig) (*zap.Logger, error) {
	parseLevel := func(l string) zapcore.LevelEnabler {
		switch strings.ToLower(l) {
		case "debug":
			return zapcore.DebugLevel
		case "info":
			return zapcore.InfoLevel
		case "warn":
			return zapcore.WarnLevel
		case "error":
			return zapcore.ErrorLevel
		default:
			return nil
		}
	}

	var cores []zapcore.Core
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if config.Stdout.Enabled {
		level := parseLevel(config.Stdout.Level)
		if level == nil {
			return nil, errors.New("invalid log level for stdout")
		}
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	}
	if config.File.Enabled {
		level := parseLevel(config.File.Level)
		if level == nil {
			return nil, errors.New("invalid log level for file")
		}
		if config.File.MaxSize <= 1 {
			return nil, errors.New("log.file.max-size must >= 1 (1MB)")
		} else if config.File.MaxBackups < 0 {
			return nil, errors.New("log.file.max-backups must >= 0")
		}
		ws := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(config.File.Dir, "sthook.log"),
			MaxSize:    config.File.MaxSize,
			MaxBackups: config.File.MaxBackups,
			LocalTime:  true,
		})
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, level))
	}
	return zap.New(zapcore.NewTee(cores...)), nil
}
