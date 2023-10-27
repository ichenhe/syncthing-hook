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
	var logger *zap.Logger
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger: %s", err)
		return
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	stClient, _ := stclient.NewSyncthing(appProfile.Syncthing.Url, appProfile.Syncthing.ApiKey, logger.Sugar())
	eventSource := exevent.NewEventSource(stClient, logger.Sugar())
	actionExecutorCreator := exevent.NewActionExecutorCreator(logger.Sugar())
	hookManager := hook.NewManager(eventSource, actionExecutorCreator, logger.Sugar())

	for i, h := range appProfile.Hooks {
		hookDef := domain.NewHookDefinition(h.Name, i)
		if err := hookManager.RegisterHook(&h, hookDef); err != nil {
			logger.Sugar().Errorf("failed to regiter hook: %s", err)
		}
	}

	group := sync.NewWaitGroup()
	group.Add(1)
	group.Wait()
}
