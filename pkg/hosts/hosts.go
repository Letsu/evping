package hosts

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Host defines a host used for pinging with IP/Host and needed sesttings
type Host struct {
	IpAddress     string `json:"ip_address"`
	PingFrequency int    `json:"ping_frequency"`
}

// Hosts defines the interface for hosts
type Hosts interface {
	GetHosts() ([]Host, error)
	AddHost(Host) error
	DeleteHost(string) error
}

type HostsCsv struct {
	File string
}

// GetHosts returns all hosts from the hosts csv file
// befor returning reding the file is check if it exists
func (h HostsCsv) GetHosts() ([]Host, error) {
	// Check if file exists
	absPath, err := filepath.Abs(h.File)
	if err != nil {
		return []Host{}, err
	}

	if _, err := os.Stat(h.File); os.IsNotExist(err) {
		return []Host{}, fmt.Errorf("host file not found: searched in %v", absPath)
	}

	// Read file
	f, err := os.OpenFile(absPath, os.O_RDONLY, 0755)
	if err != nil {
		return []Host{}, err
	}

	// Parse file
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []Host{}, err
	}

	// Convert to struct
	var hosts []Host
	for _, row := range data {
		frequenvy, err := strconv.Atoi(row[1])
		if err != nil {
			return []Host{}, err
		}
		s := Host{
			IpAddress:     row[0],
			PingFrequency: frequenvy,
		}
		hosts = append(hosts, s)
	}

	return hosts, nil
}

// AddHost adds a new host to the hosts csv file
// will return a error if the host already exists
func (h HostsCsv) AddHost(newHost Host) error {
	// First open the file in case the file dosent exist
	f, err := os.OpenFile(h.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if host already exists
	hosts, err := h.GetHosts()
	if err != nil {
		return err
	}

	for _, host := range hosts {
		if host.IpAddress == newHost.IpAddress {
			return fmt.Errorf("host already exists: %v", newHost.IpAddress)
		}
	}

	// Write host to file
	w := csv.NewWriter(f)
	w.WriteAll([][]string{{newHost.IpAddress, strconv.Itoa(newHost.PingFrequency)}})
	w.Flush()

	return nil
}

// DeleteHost deletes a host from the hosts csv file
// will return a error if the host dosent exists
func (h HostsCsv) DeleteHost(delHost string) error {
	// Check if file exists
	absPath, err := filepath.Abs(h.File)
	if err != nil {
		return err
	}

	if _, err := os.Stat(h.File); os.IsNotExist(err) {
		return fmt.Errorf("host file not found: searched in %v", absPath)
	}

	// Open file
	f, err := os.OpenFile(absPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read file, remove host and write file
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	found := false
	new_data := [][]string{}
	for i, row := range data {
		if row[0] == delHost {
			found = true
			new_data = append(data[:i], data[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("host not found: %v", delHost)
	}

	// Write file
	f, err = os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(new_data)
	w.Flush()

	return nil
}
