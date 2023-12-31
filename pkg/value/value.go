package value

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
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
	var inquiryHosts structInquiryHost
	err := c.BindJSON(&inquiryHosts)
	if err != nil {
		log.Fatalf("Failed to bind JSON to variable: %v", err)
	}
	// for _, row := range inquiryHosts.IpAddresses {
	// 	//TODO Datei auslesen
	// }
}
