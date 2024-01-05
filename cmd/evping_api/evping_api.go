package main

import (
	"fmt"
	"github.com/letsu/evping/pkg/api"
	"github.com/letsu/evping/pkg/hosts"
	"log"
	"os"
)

func main() {
	wrkDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error by getting the Working Directory: %v", err)
	}
	host := hosts.HostsCsv{File: fmt.Sprintf("%s/../../data/hosts.csv", wrkDir)}

	api.Router(&host)
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.1", PingFrequency: 10})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.2", PingFrequency: 1})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.3", PingFrequency: 1})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.4", PingFrequency: 1})
	//fmt.Println(host.GetHosts())
}
