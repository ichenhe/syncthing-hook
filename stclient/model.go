package stclient

import "time"

type ConnectionServiceStatus struct {
	Error        *string  `json:"error"`
	LanAddresses []string `json:"lanAddresses"`
	WanAddresses []string `json:"wanAddresses"`
}

type DiscoveryStatus struct {
	Error *string `json:"error"`
}

type LastDialStatus struct {
	When  time.Time `json:"when"`
	Error *string   `json:"error"`
}

type SystemStatus struct {
	Alloc                   int                                `json:"alloc"`
	ConnectionServiceStatus map[string]ConnectionServiceStatus `json:"connectionServiceStatus"`
	DiscoveryEnabled        bool                               `json:"discoveryEnabled"`
	DiscoveryErrors         map[string]string                  `json:"discoveryErrors"`
	DiscoveryStatus         map[string]DiscoveryStatus         `json:"discoveryStatus"`
	DiscoveryMethods        int                                `json:"discoveryMethods"`
	Goroutines              int                                `json:"goroutines"`
	LastDialStatus          map[string]LastDialStatus          `json:"lastDialStatus"`
	MyID                    string                             `json:"myID"`
	PathSeparator           string                             `json:"pathSeparator"`
	StartTime               time.Time                          `json:"startTime"`
	Sys                     int                                `json:"sys"`
	Themes                  []string                           `json:"themes"`
	Tilde                   string                             `json:"tilde"`
	Uptime                  int                                `json:"uptime"`
}
