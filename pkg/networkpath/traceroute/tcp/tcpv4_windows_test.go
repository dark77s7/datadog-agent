//go:build test

package tcp

import (
	"fmt"
	"net"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func Test_TracerouteSequential(t *testing.T) {
	currentTTL := 1
	target := net.ParseIP("8.8.8.8")
	srcIP := net.ParseIP("192.168.1.1")

	sendto = func(s windows.Handle, buf []byte, flags int, to windows.Sockaddr) (err error) {
		// TODO: Check the sent packets
		// check for increasing TTL up to where we want
		// check the IP header fields
		return nil
	}
	recvFrom = func(_ windows.Handle, buf []byte, _ int) (n int, from windows.Sockaddr, err error) {
		for currentTTL <= 10 {
			// Mock ICMP packet
			fakeHopIP := net.Parse(fmt.Sprintf("1.1.1.%d", currentTTL))
			packetBytes := createMockICMPPacket(createMockIPv4Layer(fakeHopIP, srcIP, layers.IPProtocolICMPv4), createMockICMPLayer(layers.ICMPv4CodeTTLExceeded), createMockIPv4Layer(srcIP, target, layers.IPProtocolTCP), createMockTCPLayer(12345, 443, 28394, 12737, true, true, true), false)
			copy(buf, packetBytes)
			return len(packetBytes), nil, nil
		}
		currentTTL++
		return 0, nil, nil
	}

	tracer := &TCPv4{
		Target:   net.ParseIP("8.8.8.8"),
		srcIP:    net.ParseIP("192.168.1.1"),
		srcPort:  12345,
		DestPort: 443,
		MinTTL:   1,
		MaxTTL:   15,
		Delay:    time.Millisecond * 100,
		Timeout:  time.Second * 1,
	}

	_, err := tracer.TracerouteSequential()
	assert.NoError(t, err)
}

func decodeICMPPacket(data []byte) (*layers.ICMPv4, error) {
	packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
	if err := packet.ErrorLayer(); err != nil {
		return nil, err.Error()
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return nil, fmt.Errorf("No IPv4 layer found in packet")
	}

	ipPacket, _ := ipLayer.(*layers.IPv4)
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer == nil {
		return nil, fmt.Errorf("No ICMP layer found in packet")
	}

	icmpPacket, _ := icmpLayer.(*layers.ICMPv4)
	return icmpPacket, nil
}
