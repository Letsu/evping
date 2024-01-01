package hosts_test

import (
	"os"
	"testing"

	"github.com/letsu/evping/pkg/hosts"
)

func TestGetHosts(t *testing.T) {
	// Create a temporary hostsData CSV file for testing
	tempFile, err := os.CreateTemp("", "test_hosts.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary file
	testData := []string{"192.168.0.1,10", "192.168.0.2,5", "192.168.0.3,15"}
	for _, data := range testData {
		_, err := tempFile.WriteString(data + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temporary file: %v", err)
		}
	}

	// Create a HostsCsv instance with the temporary file
	hostsData := hosts.HostsCsv{File: tempFile.Name()}

	// Call the GetHosts function
	result, err := hostsData.GetHosts()

	// Assert that there is no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert the expected number of hostsData
	expectedCount := len(testData)
	if len(result) != expectedCount {
		t.Errorf("Expected %d hostsData, but got %d", expectedCount, len(result))
	}

	// Assert the expected host data
	expectedHosts := []hosts.Host{
		{IpAddress: "192.168.0.1", PingFrequency: 10},
		{IpAddress: "192.168.0.2", PingFrequency: 5},
		{IpAddress: "192.168.0.3", PingFrequency: 15},
	}
	for i, expected := range expectedHosts {
		if result[i] != expected {
			t.Errorf("Expected host %v, but got %v", expected, result[i])
		}
	}
}

func TestGetHosts_EmptyFile(t *testing.T) {
	// Create a temporary empty hosts CSV file for testing
	tempFile, err := os.CreateTemp("", "test_hosts.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a HostsCsv instance with the temporary file
	hostsData := hosts.HostsCsv{File: tempFile.Name()}

	// Call the GetHosts function
	result, err := hostsData.GetHosts()

	// Assert that there is no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the result is an empty slice
	if len(result) != 0 {
		t.Errorf("Expected an empty slice, but got %v", result)
	}
}
