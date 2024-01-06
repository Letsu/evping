package value

import (
	"github.com/letsu/evping/pkg/ping"
	"log"
	"net/http"
	"time"

	"github.com/letsu/evping/pkg/hosts"

	"github.com/gin-gonic/gin"
)

type structGetHosts struct {
	IpAddress     string `json:"ip_address"`
	PingFrequency int    `json:"ping_frequency"`
}
type structDataOfHost struct {
	IpAddress string                `json:"ip_address"`
	Rows      []ping.StructPingData `json:"rows"`
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
	Hosts hosts.Hosts
	Ping  ping.PingData
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
func (r Router) DataOfHost(c *gin.Context) {
	var (
		inquiryHosts structInquiryHost
		dataOfHosts  []structDataOfHost
	)
	err := c.BindJSON(&inquiryHosts)
	if err != nil {
		log.Printf("Failed to bind JSON to variable: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	for _, ipAddress := range inquiryHosts.IpAddresses {
		data, err := r.Ping.GetPingData(ipAddress, inquiryHosts.StartTime, inquiryHosts.EndTime)
		if err != nil {
			log.Printf("Failed to get Data of Host (%s): %v", ipAddress, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Next()
		}
		dataOfHost := structDataOfHost{
			IpAddress: ipAddress,
			Rows:      data,
		}
		dataOfHosts = append(dataOfHosts, dataOfHost)
	}
	c.JSON(http.StatusOK, dataOfHosts)
}
func (r Router) AddHost(c *gin.Context) {
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
	err = r.Hosts.AddHost(h)
	if err != nil {
		log.Printf("Error by adding host via host.AddHost: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	c.Status(http.StatusCreated)
}
func (r Router) DeleteHost(c *gin.Context) {
	var delHost structDeleteHost
	err := c.BindJSON(&delHost)
	if err != nil {
		log.Printf("Failed to bind JSON to variable: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	err = r.Hosts.DeleteHost(delHost.IpAddress)
	if err != nil {
		log.Printf("Failed to delete Host from CSV: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Next()
	}
	c.Status(http.StatusAccepted)
}
