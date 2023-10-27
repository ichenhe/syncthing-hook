package stclient

import (
	"github.com/go-resty/resty/v2"
	"github.com/ichenhe/syncthing-hook/domain"
)

func (s *SyncthingClient) GetSystemStatus() (*domain.SystemStatus, error) {
	result := &domain.SystemStatus{}
	if resp, err := s.newRequest(result).Get("/rest/system/status"); err != nil {
		return nil, domain.NewStApiReqError(err)
	} else if resp.IsError() {
		return nil, domain.NewStApiHttpError(resp)
	} else {
		return result, nil
	}
}

func (s *SyncthingClient) newRequest(result interface{}) *resty.Request {
	return s.client.NewRequest().SetResult(result)
}
