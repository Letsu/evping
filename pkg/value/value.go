package value

import (
	"encoding/csv"
	"fmt"
	"github.com/letsu/evping/pkg/hosts"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type structGetHosts struct {
	IpAddress     string `json:"ip_address"`
	PingFrequency int    `json:"ping_frequency"`
}
type structDataOfHost struct {
}
type structInquiryHost struct {
	IpAddresses []string  `json:"ip_addresses"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

func GetHosts(c *gin.Context) {
	//IP Adresse	Wie oft gepingt wird
	f, _ := os.OpenFile(".\\data\\hosts.csv", os.O_RDONLY, 0755)
	data, _ := csv.NewReader(f).ReadAll()
	var js []structGetHosts
	for _, row := range data {
		intPingFreq, err := strconv.Atoi(row[1])
		if err != nil {
			log.Fatalf("Failed to convert string to int: %v", err)
		}
		s := structGetHosts{
			IpAddress:     row[0],
			PingFrequency: intPingFreq,
		}
		js = append(js, s)
	}
	c.JSON(http.StatusOK, js)
}
func DataOfHost(c *gin.Context) {
	var (
		inquiryHosts structInquiryHost
		t            time.Time
	)
	err := c.BindJSON(&inquiryHosts)
	if err != nil {
		log.Fatalf("Failed to bind JSON to variable: %v", err)
	}
	log.Println(inquiryHosts)
	for _, row := range inquiryHosts.IpAddresses {
		wrkDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error get Working Directory: %v", err)
		}
		file, _ := os.OpenFile(fmt.Sprintf("%s/../../data/host_%s.csv", wrkDir, strings.ReplaceAll(row, ".", "-")), os.O_RDONLY, 0755)
		data, _ := csv.NewReader(file).ReadAll()
		for _, column := range data {
			err := t.UnmarshalText([]byte(column[0]))
			if err != nil {
				log.Fatalf("Failed to Unmarshal Text to time: %v", err)
			}
			if t.After(inquiryHosts.StartTime) && t.Before(inquiryHosts.EndTime) || t.Equal(inquiryHosts.StartTime) || t.Equal(inquiryHosts.EndTime) {
				fmt.Println(column)
				//TODO Return Value to User
			}
		}
	}
}
func AddHost(c *gin.Context) {
	host, ok := c.MustGet("hostKey").(*hosts.HostsCsv)
	if !ok {
		log.Fatalf("Error by getting *hosts.HostsCsv")
	}
	var getHost structGetHosts
	err := c.BindJSON(&getHost)
	if err != nil {
		log.Fatalf("Failed to bind JSON to variable: %v", err)
	}
	h := hosts.Host{
		Host:          getHost.IpAddress,
		PingFrequency: getHost.PingFrequency,
	}
	err = host.AddHost(h)
	if err != nil {
		log.Fatalf("Error by adding host via host.AddHost: %v", err)
	}
	//TODO Return HTTP Code to User
}
