package performance

import (
	"fmt"
	"testing"
)

func TestGetNetValues(t *testing.T) {
	c := make(chan NetUsage)
	go GetNetValues(c, 10)
	msg := <-c
	for _, i := range msg.Ifaces {
		fmt.Println(i.String())
	}
}
