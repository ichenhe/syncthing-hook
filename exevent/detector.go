package exevent

import (
	"SyncthingHook/domain"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
)

func detectLocalFolderContentChanged(stClient domain.SyncthingClient, params *domain.HookParameters, logger *zap.SugaredLogger) (domain.EventCh, eventUnsubscriber, *domain.IllegalEventParamError) {
	folderId := params.GetString("st-folder", "")
	if len(folderId) == 0 {
		return nil, nil, domain.NewIllegalEventParamError("'st-folder' can not be empty")
	}
	path := params.GetString("path", "/")
	if len(path) == 0 {
		return nil, nil, domain.NewIllegalEventParamError("'path' can not be empty")
	}

	sourceEventCh := stClient.SubscribeEvent([]events.EventType{events.LocalChangeDetected}, 0)

	unsubscriber := eventUnsubscriberFunc(func() {
		stClient.UnsubscribeEvent(sourceEventCh)
	})

	return applyHandlers(domain.ConvertNativeEventChannel(sourceEventCh),
		NewFolderIdFilter(folderId, logger),
		NewDirPathFilter(path, logger),
	), unsubscriber, nil
}
