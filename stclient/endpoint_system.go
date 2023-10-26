package stclient

import (
	"SyncthingHook/domain"
	"github.com/go-resty/resty/v2"
)

func (s *SyncthingClient) GetSystemStatus() (*domain.SystemStatus, error) {
	result := &domain.SystemStatus{}
	if resp, err := s.newRequest(result).Get("/rest/system/status"); err != nil {
		return nil, newApiError(err)
	} else if resp.IsError() {
		return nil, newHttpApiError(resp)
	} else {
		return result, nil
	}
}

func (s *SyncthingClient) newRequest(result interface{}) *resty.Request {
	return s.client.NewRequest().SetResult(result)
}
