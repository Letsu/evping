package ping_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/letsu/evping/pkg/ping"
)

func TestGetPingData(t *testing.T) {
	// Create a temporary CSV file for testing
	tempFile, err := os.CreateTemp(t.TempDir(), "test_ping.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary file
	testData := []string{"2022-01-01T12:00:00Z,192.168.0.1,example.com,10ms"}
	for _, data := range testData {
		_, err := tempFile.WriteString(data + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temporary file: %v", err)
		}
	}

	// Create a PingCsv instance with the temporary folder
	pingCsv := ping.PingCsv{Folder: t.TempDir()}

	// Test case: Retrieve ping data for an existing host
	host := "example.com"
	pingData, err := pingCsv.GetPingData(host)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert the expected number of ping data entries
	expectedCount := 1
	if len(pingData) != expectedCount {
		t.Errorf("Expected %d ping data entries, but got %d", expectedCount, len(pingData))
	}

	// Assert the expected ping data values
	expectedPingData := ping.StructPingData{
		Time: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		Ip:   "192.168.0.1",
		Host: "example.com",
		Rtt:  10 * time.Millisecond,
	}
	if !reflect.DeepEqual(pingData[0], expectedPingData) {
		t.Errorf("Expected ping data %+v, but got %+v", expectedPingData, pingData[0])
	}

	// Test case: Retrieve ping data for a non-existent host
	nonExistentHost := "nonexistent.com"
	pingData, err = pingCsv.GetPingData(nonExistentHost)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the ping data slice is empty
	if len(pingData) != 0 {
		t.Errorf("Expected empty ping data slice, but got %d entries", len(pingData))
	}
}
