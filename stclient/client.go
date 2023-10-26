package stclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
	"strings"
)

// SyncthingClient represents a SyncthingClient instance. It is used to interact with SyncthingClient.
type SyncthingClient struct {
	Url               string
	client            *resty.Client
	eventSubscription *eventSubscription
	logger            *zap.SugaredLogger

	// allEventTypes used as the parameter of getEventsWithStringTypes to get all types of event, including XChangedDetected.
	allEventTypes string
}

var _ domain.SyncthingClient = (*SyncthingClient)(nil)

func NewSyncthing(url string, apikey string, logger *zap.SugaredLogger) (*SyncthingClient, error) {
	client := resty.New()
	client.SetBaseURL(url)
	client.SetHeader("X-API-Key", apikey)

	allEvents := make([]string, 0)
	for _, e := range []events.EventType{
		events.Starting, events.StartupComplete, events.DeviceDiscovered, events.DeviceConnected,
		events.DeviceDisconnected, events.PendingDevicesChanged,
		events.DevicePaused, events.DeviceResumed, events.ClusterConfigReceived, events.LocalChangeDetected,
		events.RemoteChangeDetected, events.LocalIndexUpdated, events.RemoteIndexUpdated, events.ItemStarted,
		events.ItemFinished, events.StateChanged, events.PendingFoldersChanged, events.ConfigSaved,
		events.DownloadProgress, events.RemoteDownloadProgress, events.FolderSummary, events.FolderCompletion,
		events.FolderErrors, events.FolderScanProgress, events.FolderPaused, events.FolderResumed,
		events.FolderWatchStateChanged, events.ListenAddressesChanged, events.LoginAttempt, events.Failure,
	} {
		allEvents = append(allEvents, e.String())
	}

	server := &SyncthingClient{
		Url:               url,
		client:            client,
		eventSubscription: newEventSubscription(),
		logger:            logger,
		allEventTypes:     strings.Join(allEvents, ","),
	}
	return server, nil
}

func (s *SyncthingClient) String() string {
	return fmt.Sprintf("SyncthingClient[%s]", s.Url)
}
