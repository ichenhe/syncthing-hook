package main

import (
	"SyncthingHook/stclient"
	"SyncthingHook/sthook"
	"github.com/syncthing/syncthing/lib/sync"
	"go.uber.org/zap"
)

func main() {
	appProfile, _ := sthook.LoadAppProfile("/Users/chenhe/Developer/SyncthingHook/config/config.sthook.yaml")
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	st, _ := stclient.NewSyncthing(appProfile.Syncthing.Url, appProfile.Syncthing.ApiKey, logger.Sugar())

	hookManager := sthook.NewHookManager(st, logger)
	if err := hookManager.RegisterHook(&appProfile.Hooks[0]); err != nil {
		logger.Sugar().Errorln(err)
	}

	group := sync.NewWaitGroup()
	group.Add(1)
	group.Wait()

}
