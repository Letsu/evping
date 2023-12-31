package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingData struct {
	Ip     string   `json:"ip"`
	Host   string   `json:"host"`
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

func formatCsvData(data [][]string) PingData {
	var newData PingData
	newData.Ip = data[0][1]
	newData.Host = data[0][2]
	for _, row := range data {
		newData.Labels = append(newData.Labels, row[0])
		if row[3] == "-1" {
			newData.Data = append(newData.Data, -1)
		} else {
			rtt, err := time.ParseDuration(row[3])
			if err != nil {
				fmt.Println(err)
			}
			seconds := int(rtt.Milliseconds())
			newData.Data = append(newData.Data, seconds)
		}
	}
	return newData
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	f, _ := os.OpenFile("hosts.csv", os.O_RDONLY, 0755)
	cr := csv.NewReader(f)
	data, _ := cr.ReadAll()

	// Get the current time
	now := time.Now()

	// Filter the data to only include the last 10 minutes
	var filteredData [][]string
	var t time.Time
	for _, row := range data {
		// Parse the time from the first column
		err := t.UnmarshalText([]byte(row[0]))
		if err != nil {
			log.Fatalf("Failed to parse time: %v", err)
		}

		// If the time is within the last 10 minutes, add the row to filteredData
		if now.Sub(t).Minutes() <= 5 {
			filteredData = append(filteredData, row)
		}
	}

	newData := []PingData{formatCsvData(filteredData)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newData)
}

func ping() {
	//tc qdisc replace dev eth0 root netem loss 25%
	pinger, err := probing.NewPinger("1.1.1.1")
	//When on Windows must set the follow Line
	pinger.SetPrivileged(true)
	//pinger.Debug = true
	if err != nil {
		log.Fatal(err)
	}
	last := 0

	pinger.OnSend = func(pkt *probing.Packet) {
		if pkt.Seq-last > 1 {
			for i := last + 1; i < pkt.Seq; i++ {
				t := time.Now()
				curTime, _ := t.MarshalText()
				f, err := os.OpenFile("hosts.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				w := csv.NewWriter(f)
				w.Write([]string{string(curTime), pkt.IPAddr.String(), pkt.Addr, "-1"})
				w.Flush()
			}
			last++
		}
	}
	pinger.OnRecv = func(pkt *probing.Packet) {
		last = pkt.Seq
		t := time.Now()
		curTime, _ := t.MarshalText()
		f, err := os.OpenFile("hosts.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{string(curTime), pkt.IPAddr.String(), pkt.Addr, pkt.Rtt.String()})
		w.Flush()

		fmt.Println("normal: ", pkt)
		//fmt.Printf("%v: %d bytes from %s: icmp_seq=%d time=%v\n",
		//	time.Now(), pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnSendError = func(pkt *probing.Packet, err error) {
		t := time.Now()
		curTime, _ := t.MarshalText()
		f, err := os.OpenFile("hosts.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{string(curTime), pkt.IPAddr.String(), pkt.Addr, "-1"})
		w.Flush()
		fmt.Println("error: ", pkt)
	}

	//pinger.OnRecvError = func(err error) { fmt.Println("error: ", err) }

	err = pinger.Run()
	fmt.Println(err)
}

func main() {
	fmt.Println("Starting server...")
	go ping()
	http.HandleFunc("/api/", handler)
	fs := http.FileServer(http.Dir("website/dist"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
