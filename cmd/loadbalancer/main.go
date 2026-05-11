package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"net"
	"os"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type Config struct {
	ShardIP   string `json:"shard_ip"`
	ShardPort uint16 `json:"shard_port"`
	ChatIP    string `json:"chat_ip"`
	ChatPort  uint16 `json:"chat_port"`
}

type CvarsMap map[string]float32

// float32ToFloat16 converts a float32 to IEEE 754 half-precision (float16).
func float32ToFloat16(f float32) uint16 {
	if f == 0 {
		return 0
	}
	bits := math.Float32bits(f)
	sign := uint16((bits >> 31) & 1)
	exp := int((bits>>23)&0xFF) - 127 + 15
	mant := uint16((bits >> 13) & 0x3FF)
	if exp <= 0 {
		return sign << 15
	}
	if exp >= 31 {
		return sign<<15 | 0x7C00
	}
	return sign<<15 | uint16(exp)<<10 | mant
}

// encodeCvarFloat converts float32 to the custom 16-bit LB encoding:
// encoded = ((f16 & 0xC000) >> 1) | ((f16 & 0x3FFF) >> 4)
func encodeCvarFloat(f float32) uint16 {
	h := float32ToFloat16(f)
	return ((h & 0xC000) >> 1) | ((h & 0x3FFF) >> 4)
}

func buildCvarsBody(cvarsMap *CvarsMap) []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 90))
	bw.WriteBeUint32(uint32(len(*cvarsMap)))
	for k, v := range *cvarsMap {
		bw.WriteBool(false)
		bw.WriteString(k)
		// XXX: bw.Align()
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteBool(false)
		bw.WriteByte(0x2)
		ev := encodeCvarFloat(v)
		bw.WriteBeUint16(ev)
		bw.WriteBeUint16(0x00)
	}
	bw.WriteByte(0x00)
	return bw.ReturnSlice()
}

// XXX: Why is this padded by two single bits...
func buildShardBody(shardIP string, shardPort uint16, chatIP string, chatPort uint16) []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 40))
	bw.WriteBool(true)
	bw.WriteByte(1)
	bw.WriteCString(shardIP)
	bw.WriteBeUint16(shardPort)
	bw.WriteBool(true)
	bw.WriteCString(chatIP)
	bw.WriteBeUint16(chatPort)
	return bw.ReturnSlice()
}

func handle(conn net.Conn, cfg *Config, cvars *CvarsMap) {
	defer conn.Close()
	log.Printf("%s: new connection", conn.RemoteAddr())

	if _, err := conn.Write(protocol.MakeMessage(types.SCMD_LB_CVARS, 0, 0, buildCvarsBody(cvars))); err != nil {
		log.Printf("error sending cvars: %v", err)
		return
	}
	log.Printf("%s: sent cvars", conn.RemoteAddr())

	body := buildShardBody(cfg.ShardIP, cfg.ShardPort, cfg.ChatIP, cfg.ChatPort)
	if _, err := conn.Write(protocol.MakeMessage(types.SCMD_ASSIGNED_SHARD, 2, 0, body)); err != nil {
		log.Printf("error sending shard address: %v", err)
		return
	}
	log.Printf("%s: sent shard -> %s:%d chat -> %s:%d", conn.RemoteAddr(),
		cfg.ShardIP, cfg.ShardPort, cfg.ChatIP, cfg.ChatPort)
}

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	cvarsPath := flag.String("cvars", "cvars.json", "path to cvars.json")
	listenAddress := flag.String("listen", ":3801", "Address to listen on")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	var cvars CvarsMap
	data, err = os.ReadFile(*cvarsPath)
	if err != nil {
		log.Fatalf("Failed to read cvars config: %v", err)
	}
	if err := json.Unmarshal(data, &cvars); err != nil {
		log.Fatalf("invalid cvars config: %v", err)
	}

	ln, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("load balancer listening on :%s", *listenAddress)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn, &cfg, &cvars)
	}
}
