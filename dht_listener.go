package main

import (
	"fmt"
	"log/slog"
	"magnet-parser/utils"
	"net"
	"strings"
)

const serverPort = 14888
var udpServer, _ = net.ListenPacket("udp", fmt.Sprintf(":%d", serverPort))

func listenUDP() {
	defer udpServer.Close()
	for {
		buf := make([]byte, 65507)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go processInputRequest(udpServer, addr, buf[:n])
	}
}

// todo ответы на разные пакеты (сохранение токенов и сидов)
func processInputRequest(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(buf), "\r", "")))
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
		slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(bytes), "\r", "")))
	}
	json, _ = utils.BencodeToJSON(bytes)
	slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
	bytes, err = utils.JSONToBencode(json)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(bytes), "\r", "")))
	}
}
func main() {
	// TODO utils.SetupLogger()
	addr := "95.211.60.23:58672"
	udpServer.WriteTo(Ping(), utils.StringToUDPAddr(addr))
	udpServer.WriteTo(GetPeers("0f52fb0190dcd61381cdd0893173b1309ba9868f"), utils.StringToUDPAddr(addr))
	udpServer.WriteTo(FindNode("0f43110b4bdf604508cb75dfd326dcd68ac936c7"), utils.StringToUDPAddr(addr))
	udpServer.WriteTo(AnnouncePeer("0f52fb0190dcd61381cdd0893173b1309ba9868f", serverPort, "1dcbde3b"), utils.StringToUDPAddr(addr))
	listenUDP()
}
