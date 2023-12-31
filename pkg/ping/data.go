package ping

import "time"

type StructPingData struct {
	Time time.Time     `json:"time"`
	Ip   string        `json:"ip"`
	Host string        `json:"host"`
	Rtt  time.Duration `json:"rtt"`
}

type PingData interface {
	GetPingData(string, time.Time, time.Time) ([]StructPingData, error)
	AddPingData(StructPingData) error
}
