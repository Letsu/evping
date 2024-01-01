package hosts

// Host defines a host used for pinging with IP/Host and needed sesttings
type Host struct {
	Host          string `json:"host"`
	PingFrequency int    `json:"ping_frequency"`
}

// Hosts defines the interface for hosts
type Hosts interface {
	GetHosts() ([]Host, error)
	AddHost(Host) error
	DeleteHost(string) error
}
