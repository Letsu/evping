package hosts

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

// HostsCsv represents a CSV file containing a list of hosts.
// It implemates the Hosts interface.
type HostsCsv struct {
	File  string // The path to the CSV file where all hosts are stored
	mutex sync.Mutex
}

// GetHosts retrieves the list of hosts from the hosts CSV file.
// It reads the file line by line and parses each row to create a list of Host objects.
// If the file doesn't exist, it returns an empty slice and no error.
// If there is an error while reading or parsing the file, it returns nil and the error.
// The function uses a mutex to ensure concurrent-safe access to the shared resources.
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
			Host:          row[0],
			PingFrequency: frequency,
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

// AddHost adds a new host to the hosts CSV file.
// It checks if the provided host data is valid, including a non-empty host name and a positive ping frequency.
// If the host already exists in the CSV file, it returns an error.
// It opens the file, creates it if it doesn't exist, and writes the new host data to the file.
// The function uses a mutex to ensure concurrent-safe access to the shared resources.
func (h *HostsCsv) AddHost(newHost Host) error {
	// Check if data is valid
	if newHost.Host == "" {
		return fmt.Errorf("host is empty")
	}
	if newHost.PingFrequency <= 0 {
		return fmt.Errorf("ping frequency is invalid: %v", newHost.PingFrequency)
	}

	// Check if host already exists
	hosts, err := h.GetHosts()
	if err != nil {
		return fmt.Errorf("failed to get existing hosts: %v", err)
	}

	for _, host := range hosts {
		if host.Host == newHost.Host {
			return fmt.Errorf("host already exists: %v", newHost.Host)
		}
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()
	// Open file and create if it dosent exists
	f, err := os.OpenFile(h.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Write host to file
	// Write host to file
	w := csv.NewWriter(f)
	defer w.Flush()
	err = w.WriteAll([][]string{{newHost.Host, strconv.Itoa(newHost.PingFrequency)}})
	if err != nil {
		return fmt.Errorf("failed to write host to file: %v", err)
	}

	return nil
}

// DeleteHost deletes a host from the hosts CSV file.
// It retrieves the existing hosts from the file, searches for the specified host to delete,
// and removes it from the list of hosts. If the host is not found, it returns an error.
// The function then opens the file, truncates its contents, and writes the updated list of hosts to the file.
// The function uses a mutex to ensure concurrent-safe access to the shared resources.
func (h *HostsCsv) DeleteHost(delHost string) error {
	hosts, err := h.GetHosts()
	if err != nil {
		return fmt.Errorf("failed to get existing hosts: %v", err)
	}

	found := false
	updatedHosts := []Host{}
	for i, host := range hosts {
		if host.Host == delHost {
			found = true
			updatedHosts = append(hosts[:i], hosts[i+1:]...)
			break
		}
	}
	fmt.Println(updatedHosts)

	if !found {
		return fmt.Errorf("host not found: %v", delHost)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Write new host file without the deleted host
	f, err := os.OpenFile(h.File, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	// Write updatedHosts
	for _, host := range updatedHosts {
		err = w.Write([]string{host.Host, strconv.Itoa(host.PingFrequency)})
		if err != nil {
			return fmt.Errorf("failed to write host to file: %v", err)
		}
	}

	return nil
}
