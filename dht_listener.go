package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"magnet-parser/utils"
	"net"
)

func sendUDPPacket(addr string, contents string) []byte {
	p := make([]byte, 2048)
	var n int
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return nil
	}
	fmt.Fprintf(conn, contents)
	n, err = bufio.NewReader(conn).Read(p)
	p = p[:n]
	defer conn.Close()
	if err == nil {
		return p
	} else {
		fmt.Printf("Some error %v\n", err)
		return nil
	}
}

func main() {
	//utils.SetupLogger()
	response := sendUDPPacket("46.147.238.182:14184", "d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t1:01:y1:qe")
	json, err := utils.BencodeToJSON(response)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
	}
}
