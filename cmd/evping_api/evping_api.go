package main

import (
	"github.com/letsu/evping/pkg/api"
	"github.com/letsu/evping/pkg/hosts"
	"log"
	"os"
	"path/filepath"
)

func main() {
	wrkDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error by getting the Working Directory: %v", err)
	}
	pathToHostFile := filepath.Join(wrkDir, "..", "..", "data", "hosts.csv")
	host := hosts.HostsCsv{File: pathToHostFile}
	log.Println(host.File)

	api.Router(&host)
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.1", PingFrequency: 10})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.2", PingFrequency: 1})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.3", PingFrequency: 1})
	//host.AddHost(hosts.Host{IpAddress: "1.1.1.4", PingFrequency: 1})
	//fmt.Println(host.GetHosts())
}
