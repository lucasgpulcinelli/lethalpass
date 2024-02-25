package main

import (
	"net"
)

func main() {
	localServerV6Addr, _ := net.ResolveUDPAddr("udp6", "[::]:7777")
	realServerV4Addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:7777")

	localServerV6, _ := net.ListenUDP("udp6", localServerV6Addr)

	fwdConnV4, _ := net.DialUDP("udp4", nil, realServerV4Addr)

	b := make([]byte, 4096)
	didConnect := false
	for {
		n, remoteAddrV6, _ := localServerV6.ReadFromUDP(b)
		fwdConnV4.Write(b[:n])

		if didConnect {
			continue
		}
		go func() {
			b2 := make([]byte, 4096)
			for {
				m, _ := fwdConnV4.Read(b2)
				localServerV6.WriteToUDP(b2[:m], remoteAddrV6)
			}
		}()
	}
}
