package main

import (
	"log"
	"net"
	"net/netip"
)

var (
	localServerV6    *net.UDPConn
	realServerV4Addr *net.UDPAddr
	connMapV6ToV4    = map[netip.AddrPort]*net.UDPConn{}
)

func dieif(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func handleConn(remoteAddrV6 *net.UDPAddr, fwdConnV4 *net.UDPConn) {
	b := make([]byte, 4096)
	for {
		n, err := fwdConnV4.Read(b)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("read", n, "bytes from", realServerV4Addr.String())

		m, err := localServerV6.WriteToUDP(b[:n], remoteAddrV6)
		if err != nil {
			log.Println(err)
		}
		if m != n {
			log.Println("wrote less v6")
		}

		log.Println("forwarded", n, "bytes to", remoteAddrV6.String())
	}
}

func handlePacket(p []byte, remoteAddrV6 *net.UDPAddr) {
	var err error

	fwdConnV4, exists := connMapV6ToV4[remoteAddrV6.AddrPort()]
	if !exists {
		fwdConnV4, err = net.DialUDP("udp4", nil, realServerV4Addr)
		if err != nil {
			log.Println(err)
			return
		}
		connMapV6ToV4[remoteAddrV6.AddrPort()] = fwdConnV4

		log.Println("created new connection for", remoteAddrV6.String(), ":", fwdConnV4)

		go handleConn(remoteAddrV6, fwdConnV4)
	}

	n, err := fwdConnV4.Write(p)
	if err != nil {
		log.Println(err)
		return
	}
	if n < len(p) {
		log.Println("wrote less v4")
		return
	}

	log.Println("forwarded", n, "bytes to", realServerV4Addr.String())
}

func main() {
	var err error

	localServerV6Addr, err := net.ResolveUDPAddr("udp6", "[::]:7777")
	dieif(err)

	realServerV4Addr, err = net.ResolveUDPAddr("udp4", "127.0.0.1:7777")
	dieif(err)

	localServerV6, err = net.ListenUDP("udp6", localServerV6Addr)
	dieif(err)

	b := make([]byte, 4096)
	for {
		n, remoteAddrV6, err := localServerV6.ReadFromUDP(b)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("read", n, "bytes from", remoteAddrV6.String())

		handlePacket(b[:n], remoteAddrV6)
	}
}
