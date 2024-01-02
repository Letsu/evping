package ping

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Ping interface {
	Ping(context.Context) error
}

type PingerImpl struct {
	host       string
	Privileged bool
}

// NewPinger
func NewPinger(host string) Ping {
	return PingerImpl{
		Privileged: true,
	}
}

func (p PingerImpl) Ping(ctx context.Context) error {
	pinger, err := probing.NewPinger(p.host)
	if err != nil {
		return err
	}
	pinger.SetPrivileged(p.Privileged)

	// On send a check is performed if the last ping was reviced successfully
	// All missing pings are saved with -1
	last := 0
	pinger.OnSend = func(pkt *probing.Packet) {
		if pkt.Seq-last > 1 {
			for i := last + 1; i < pkt.Seq; i++ {
				t := time.Now()
				curTime, _ := t.MarshalText()
				// TODO: move file functions to a own interface
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

	// On revice the RTT is saved
	pinger.OnRecv = func(pkt *probing.Packet) {
		last = pkt.Seq
		t := time.Now()
		curTime, _ := t.MarshalText()
		// TODO: move file functions to a own interface
		f, err := os.OpenFile("hosts.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{string(curTime), pkt.IPAddr.String(), pkt.Addr, pkt.Rtt.String()})
		w.Flush()

		fmt.Println("normal: ", pkt)
	}

	// On error while sending the RTT is saved as -1
	pinger.OnSendError = func(pkt *probing.Packet, send_err error) {
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

	return nil
}
