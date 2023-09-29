package main

import (
	"fmt"
	"log/slog"
	"magnet-parser/utils"
	"net"
)

const serverPort = "14888"
var udpServer, _ = net.ListenPacket("udp", ":" + serverPort)

func listenUDP() {
	defer udpServer.Close()
	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go processInputRequest(udpServer, addr, buf[:n])
	}
}

func processInputRequest(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	slog.Info(fmt.Sprintf("Encoded json is: %s\n", string(buf)))
	json, err := utils.BencodeToJSON(buf)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
	}
	var bytes []byte
	bytes, err = utils.JSONToBencode(json)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s\n", string(bytes)))
	}
	json, _ = utils.BencodeToJSON(bytes)
	slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
}

func main() {
	// TODO utils.SetupLogger()
	udpServer.WriteTo([]byte("d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t4:34d81:y1:qe"), utils.StringToUDPAddr("91.204.96.231:13343"))
	listenUDP()
}
