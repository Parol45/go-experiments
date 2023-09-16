package main

import (
	"bufio"
	"fmt"
	"magnet-parser/lib"
	"net"
)

func sendUDPPacket(addr string, contents string) string {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return ""
	}
	fmt.Fprintf(conn, contents)
	_, err = bufio.NewReader(conn).Read(p)
	defer conn.Close()
	if err == nil {
		fmt.Printf("%s\n", p)
		return string(p)
	} else {
		fmt.Printf("Some error %v\n", err)
		return ""
	}
}

func main() {
	println(lib.BencodeToJSON("d1:rd2:id20:abcdefghij01234567895:nodes9:def456...5:token8:aoeusnthe1:ti0e1:y1:re"))
	sendUDPPacket("188.243.244.136:20484", "d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t1:01:y1:qe")
}
