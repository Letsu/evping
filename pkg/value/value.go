package value

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/letsu/evping/pkg/hosts"

	"github.com/gin-gonic/gin"
)

type structGetHosts struct {
	IpAddress     string `json:"ip_address"`
	PingFrequency int    `json:"ping_frequency"`
}
type structDataOfHost struct {
	IpAddress string               `json:"ip_address"`
	Rows      []structRowOfHostCsv `json:"rows"`
}
type structRowOfHostCsv struct {
	Timestamp time.Time `json:"timestamp"`
	IpAddress string    `json:"ip_address"`
	DnsName   string    `json:"dns_name"`
	RTT       string    `json:"rtt"`
}
type structInquiryHost struct {
	IpAddresses []string  `json:"ip_addresses"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}
type structDeleteHost struct {
	IpAddress string `json:"ip_address"`
}

type Router struct {
	Hosts *hosts.Hosts
}

func (r Router) GetHosts(c *gin.Context) {
	listOfHost, err := r.Hosts.GetHosts()
	if err != nil {
		log.Printf("Error by getting Hosts from CSV-File: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	c.JSON(http.StatusOK, listOfHost)
}
func DataOfHost(c *gin.Context) {
	var (
		inquiryHosts structInquiryHost
		t            time.Time
		dataOfHosts  []structDataOfHost
	)
	err := c.BindJSON(&inquiryHosts)
	log.Println(inquiryHosts)
	if err != nil {
		log.Printf("Failed to bind JSON to variable: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	wrkDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error get Working Directory: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	for _, row := range inquiryHosts.IpAddresses {
		var dataOfHost structDataOfHost
		filename := fmt.Sprintf("host_%s.csv", strings.ReplaceAll(row, ".", "-"))
		pathToFile := filepath.Join(wrkDir, "..", "..", "data", filename)
		file, _ := os.OpenFile(pathToFile, os.O_RDONLY, 0755)
		data, _ := csv.NewReader(file).ReadAll()
		for _, column := range data {
			err := t.UnmarshalText([]byte(column[0]))
			if err != nil {
				log.Printf("Failed to Unmarshal Text to time: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				c.Next()
			}
			if t.After(inquiryHosts.StartTime) && t.Before(inquiryHosts.EndTime) || t.Equal(inquiryHosts.StartTime) || t.Equal(inquiryHosts.EndTime) {
				fmt.Println(column)
				d := structRowOfHostCsv{
					Timestamp: t,
					IpAddress: column[1],
					DnsName:   column[2],
					RTT:       column[3],
				}
				dataOfHost.Rows = append(dataOfHost.Rows, d)
			}
		}
		dataOfHost.IpAddress = row
		dataOfHosts = append(dataOfHosts, dataOfHost)
	}
	c.JSON(http.StatusOK, dataOfHosts)
}
func AddHost(c *gin.Context) {
	host, ok := c.MustGet("hostKey").(*hosts.HostsCsv)
	if !ok {
		log.Printf("Error by getting *hosts.HostsCsv")
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	var getHost structGetHosts
	err := c.BindJSON(&getHost)
	if err != nil {
		log.Printf("Failed to bind JSON to variable: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	h := hosts.Host{
		Host:          getHost.IpAddress,
		PingFrequency: getHost.PingFrequency,
	}
	err = host.AddHost(h)
	if err != nil {
		log.Printf("Error by adding host via host.AddHost: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	c.Status(http.StatusCreated)
}

func DeleteHost(c *gin.Context) {
	host, ok := c.MustGet("hostKey").(*hosts.HostsCsv)
	if !ok {
		log.Printf("Error by getting *hosts.HostsCsv")
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	var delHost structDeleteHost
	err := c.BindJSON(&delHost)
	if err != nil {
		log.Printf("Failed to bind JSON to variable: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	err = host.DeleteHost(delHost.IpAddress)
	if err != nil {
		log.Printf("Failed to delete Host from CSV: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	c.Status(http.StatusAccepted)
}
