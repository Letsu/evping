package main

import (
	"fmt"

	"github.com/letsu/evping/pkg/hosts"
)

func main() {
	// api.Router()
	host := hosts.HostsCsv{File: "data\\hosts.csv"}
	host.AddHost(hosts.Host{IpAddress: "1.1.1.1", PingFrequency: 10})
	fmt.Println(host.GetHosts())
}
