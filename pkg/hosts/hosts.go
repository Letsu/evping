package hosts

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	File  string
	mutex sync.Mutex
}

// GetHosts retrieves the list of hosts from the hosts CSV file.
// It reads the file line by line and parses each row to create a list of Host objects.
// If the file doesn't exist, it returns an empty slice and no error.
// If there is an error while reading or parsing the file, it returns nil and the error.
func (h *HostsCsv) GetHosts() ([]Host, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Open the file for reading
	f, err := os.OpenFile(h.File, os.O_RDONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty slice if the file doesn't exist
			return []Host{}, nil
		}
		return nil, err
	}
	defer f.Close()

	// Read the file line by line
	reader := csv.NewReader(f)
	var hosts []Host
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse the row data
		frequency, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}
		host := Host{
			IpAddress:     row[0],
			PingFrequency: frequency,
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

// AddHost adds a new host to the hosts csv file
// will return a error if the host already exists
func (h *HostsCsv) AddHost(newHost Host) error {
	// Check if data is valid
	if newHost.IpAddress == "" {
		return fmt.Errorf("ip address is empty")
	}
	if newHost.PingFrequency <= 0 {
		return fmt.Errorf("ping frequency is invalid: %v", newHost.PingFrequency)
	}

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

	h.mutex.Lock()
	defer h.mutex.Unlock()
	// First open the file in case the file dosent exist
	f, err := os.OpenFile(h.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write host to file
	w := csv.NewWriter(f)
	defer w.Flush()
	err = w.WriteAll([][]string{{newHost.IpAddress, strconv.Itoa(newHost.PingFrequency)}})
	if err != nil {
		return err
	}

	return nil
}

// DeleteHost deletes a host from the hosts csv file
// will return a error if the host dosent exists
func (h *HostsCsv) DeleteHost(delHost string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
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
