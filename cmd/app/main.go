package main

import (
	"errors"
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/exevent"
	"github.com/ichenhe/syncthing-hook/hook"
	"github.com/ichenhe/syncthing-hook/stclient"
	"github.com/syncthing/syncthing/lib/sync"
	"go.uber.org/zap"
	"log"
	"os"
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
	if logger_, err := zap.NewDevelopment(); err != nil {
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
			logger.Errorf("failed to regiter hook: %s", err)
		}
	}
	logger.Infof("%d hook loaded", len(appProfile.Hooks))

	group := sync.NewWaitGroup()
	group.Add(1)
	group.Wait()
}
