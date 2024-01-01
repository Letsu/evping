package main

import (
	"fmt"

	"github.com/letsu/evping/pkg/hosts"
)

func main() {
	// api.Router()
	host := hosts.HostsCsv{File: "data\\hosts.csv"}
	host.AddHost(hosts.Host{IpAddress: "1.1.1.1", PingFrequency: 10})
	host.AddHost(hosts.Host{IpAddress: "1.1.1.2", PingFrequency: 1})
	host.AddHost(hosts.Host{IpAddress: "1.1.1.3", PingFrequency: 1})
	host.AddHost(hosts.Host{IpAddress: "1.1.1.4", PingFrequency: 1})
	fmt.Println(host.GetHosts())
}
