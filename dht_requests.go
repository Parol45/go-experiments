package main

import (
	"fmt"
	"log/slog"
	"magnet-parser/utils"
)

const id = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const version = "aaaaaaaa"
var transactionId = 0

const pingFmt = `{"a":{"id":"%s"},"q":"ping","t":"%s","v":"%s","y":"q"}`
const getPeersFmt = `{"a":{"id":"%s","info_hash":"%s"},"q":"get_peers","t":"%s","v":"%s","y":"q"}`
const findNodeFmt = `{"a":{"id":"%s","target":"%s"},"q":"find_node","t":"%s","v":"%s","y":"q"}`
const announcePeerFmt = `{"a":{"id":"%s","info_hash":"%s","port":%d,"token":"%s"},"q":"announce_peers","t":"%s","v":"%s","y":"q"}`

func transIdToStr() string {
	transactionId = (transactionId + 1) % 100000000
	return fmt.Sprintf("%08d", transactionId)
}

func Ping() []byte {
	resultingJson := fmt.Sprintf(pingFmt, id, transIdToStr(), version)
	result, err := utils.JSONToBencode(resultingJson)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func GetPeers(hash string) []byte {
	resultingJson := fmt.Sprintf(getPeersFmt, id, hash, transIdToStr(), version)
	result, err := utils.JSONToBencode(resultingJson)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func FindNode(foreignId string) []byte {
	resultingJson := fmt.Sprintf(findNodeFmt, id, foreignId, transIdToStr(), version)
	result, err := utils.JSONToBencode(resultingJson)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func AnnouncePeer(hash string, port int, token string) []byte {
	resultingJson := fmt.Sprintf(announcePeerFmt, id, hash, port, token, transIdToStr(), version)
	result, err := utils.JSONToBencode(resultingJson)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}