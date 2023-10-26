package main

import (
	"SyncthingHook/domain"
	"SyncthingHook/exevent"
	"SyncthingHook/hook"
	"SyncthingHook/stclient"
	"github.com/syncthing/syncthing/lib/sync"
	"go.uber.org/zap"
)

func main() {
	appProfile, _ := domain.LoadAppProfile("/Users/chenhe/Developer/SyncthingHook/config/config.sthook.yaml")
	logger, _ := zap.NewDevelopment()
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
