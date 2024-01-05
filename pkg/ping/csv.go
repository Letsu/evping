package ping

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type PingCsv struct {
	Folder string
	mutex  sync.Mutex
}

// GetPingData retrieves the ping data for a specific host from the CSV file.
// It opens the file for reading, constructs the file path based on the host name,
// and reads the file line by line using a CSV reader.
// The function returns a slice of StructPingData containing the parsed ping data.
// If the file doesn't exist, it returns an empty slice and no error.
// If there is an error while opening or reading the file, it returns an error with a descriptive message.
// The function uses a mutex to ensure concurrent-safe access to the shared resources.
func (p *PingCsv) GetPingData(host string) ([]StructPingData, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	// Acquire the mutex lock for concurrent-safe access
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Open the file for reading
	filePath := filepath.Join(p.Folder, host+".csv")
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty slice if the file doesn't exist
			return []StructPingData{}, nil
		}
		return nil, fmt.Errorf("error while opening file: %w", err)
	}
	defer f.Close()

	// Read the file line by line using a CSV reader
	reader := csv.NewReader(f)
	var pingData []StructPingData
	curRow := 0
	for {
		curRow++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error while reading file in row %d: %w", curRow, err)
		}

		// Validate the row length
		if len(row) != 4 {
			return nil, fmt.Errorf("malformed row in row %d: expected 4 fields, got %d", curRow, len(row))
		}

		// Parse the row data and append it to the pingData slice
		t := time.Time{}
		err = t.UnmarshalText([]byte(row[0]))
		if err != nil {
			return nil, fmt.Errorf("error while parsing time in row %d: %w", curRow, err)
		}
		rtt, err := time.ParseDuration(row[3])
		if err != nil {
			return nil, fmt.Errorf("error while parsing RTT in row %d: %w", curRow, err)
		}
		data := StructPingData{
			Time: t,
			Ip:   row[1],
			Host: row[2],
			Rtt:  rtt,
		}
		pingData = append(pingData, data)
	}

	return pingData, nil
}

func (p *PingCsv) AddPingData(StructPingData) error {
	return nil
}
