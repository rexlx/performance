package performance

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	// _ = iota
	rxBytes = iota
	rxPackets
	rxErrs
	rxDrop
	rxFifo
	rxFrame
	rxCompressed
	rxMulticast
	txBytes
	txPackets
	txErrs
	txDrop
	txFifo
	txColls
	txCarrier
	txCompressed
)

type RxData struct {
	Bytes      int
	Packets    int
	Errs       int
	Drop       int
	Fifo       int
	Frame      int
	Compressed int
	Multicast  int
}

type TxData struct {
	Bytes      int
	Packets    int
	Errs       int
	Drop       int
	Fifo       int
	Colls      int
	Carrier    int
	Compressed int
}

type Iface struct {
	Name string
	Rx   *RxData
	Tx   *TxData
}

func (i *Iface) BytesTransmittedToHumanReadable() string {
	return ByteConverter(i.Tx.Bytes)
}

func (i *Iface) BytesReceivedToHumanReadable() string {
	return ByteConverter(i.Rx.Bytes)
}

func (i *Iface) TotalBytesToHumanReadable() string {
	return ByteConverter(i.Tx.Bytes + i.Rx.Bytes)
}

func (i *Iface) String() string {
	return fmt.Sprintf("Name: %s\tRx: %+v Tx: %+v (total) %v", i.Name, i.BytesReceivedToHumanReadable(), i.BytesTransmittedToHumanReadable(), i.TotalBytesToHumanReadable())
}

// NetUsage type represents the utilization of a given network interface
type NetUsage struct {
	Name   string    `json:"name"`
	Time   time.Time `json:"time"`
	Ifaces []Iface   `json:"ifaces"`
}

// GetNetValues polls the network statistics for a given interval in seconds
func GetNetValues(c chan NetUsage, refresh int) {
	var netVals []Iface
	if refresh < 1 {
		refresh = 5
	}

	initialPoll, err := PollNetworkStatistics()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get and sort keys
	var keys []string
	for k := range initialPoll {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	time.Sleep(time.Duration(refresh) * time.Second)

	poll, err := PollNetworkStatistics()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// calculate the difference between the initial poll and the current poll
	for _, key := range keys {
		newIface := Iface{
			Rx: &RxData{},
			Tx: &TxData{},
		}
		initial := initialPoll[key]
		current := poll[key]
		newIface.Name = key
		newIface.Rx.Bytes = current.Rx.Bytes - initial.Rx.Bytes
		newIface.Rx.Packets = current.Rx.Packets - initial.Rx.Packets
		newIface.Rx.Errs = current.Rx.Errs - initial.Rx.Errs
		newIface.Rx.Drop = current.Rx.Drop - initial.Rx.Drop
		newIface.Rx.Fifo = current.Rx.Fifo - initial.Rx.Fifo
		newIface.Rx.Frame = current.Rx.Frame - initial.Rx.Frame
		newIface.Rx.Compressed = current.Rx.Compressed - initial.Rx.Compressed
		newIface.Rx.Multicast = current.Rx.Multicast - initial.Rx.Multicast
		newIface.Tx.Bytes = current.Tx.Bytes - initial.Tx.Bytes
		newIface.Tx.Packets = current.Tx.Packets - initial.Tx.Packets
		newIface.Tx.Errs = current.Tx.Errs - initial.Tx.Errs
		newIface.Tx.Drop = current.Tx.Drop - initial.Tx.Drop
		newIface.Tx.Fifo = current.Tx.Fifo - initial.Tx.Fifo
		newIface.Tx.Colls = current.Tx.Colls - initial.Tx.Colls
		newIface.Tx.Carrier = current.Tx.Carrier - initial.Tx.Carrier
		newIface.Tx.Compressed = current.Tx.Compressed - initial.Tx.Compressed
		netVals = append(netVals, newIface)
	}

	c <- NetUsage{Name: "net", Time: time.Now(), Ifaces: netVals}
}

func PollNetworkStatistics() (ifaces map[string]*Iface, err error) {
	out := make(map[string]*Iface)
	contents, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return out, err
	}
	lines := strings.Split(string(contents), "\n")
	// lines[2:] is the actual data, ignore the header
	for _, line := range lines[2:] {
		// split the line into the interface name and the stats
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		stats := strings.Fields(parts[1])
		if len(stats) < 16 {
			continue
		}
		iface := &Iface{
			Name: name,
			Rx: &RxData{
				Bytes:      ValueToInteger(stats[rxBytes]),
				Packets:    ValueToInteger(stats[rxPackets]),
				Errs:       ValueToInteger(stats[rxErrs]),
				Drop:       ValueToInteger(stats[rxDrop]),
				Fifo:       ValueToInteger(stats[rxFifo]),
				Frame:      ValueToInteger(stats[rxFrame]),
				Compressed: ValueToInteger(stats[rxCompressed]),
				Multicast:  ValueToInteger(stats[rxMulticast]),
			},
			Tx: &TxData{
				Bytes:      ValueToInteger(stats[txBytes]),
				Packets:    ValueToInteger(stats[txPackets]),
				Errs:       ValueToInteger(stats[txErrs]),
				Drop:       ValueToInteger(stats[txDrop]),
				Fifo:       ValueToInteger(stats[txFifo]),
				Colls:      ValueToInteger(stats[txColls]),
				Carrier:    ValueToInteger(stats[txCarrier]),
				Compressed: ValueToInteger(stats[txCompressed]),
			},
		}
		out[name] = iface
	}
	return out, nil
}
