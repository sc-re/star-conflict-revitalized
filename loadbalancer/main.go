package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"log"
	"math"
	"net"
	"os"
	"starconflict/lib/protocol"
)

const (
	cvarsType     = 0x8108
	shardAddrType = 0x8109
)

type Config struct {
	ShardIP   string `json:"shard_ip"`
	ShardPort int    `json:"shard_port"`
	ChatIP    string `json:"chat_ip"`
	ChatPort  int    `json:"chat_port"`
}

type CvarsMap map[string]float32

// encodeCvarName applies the carry-bit right-shift encoding from the LB protocol:
// each char is shifted right by 1; the dropped LSB becomes the MSB of the next byte.
func encodeCvarName(name string) []byte {
	carry := byte(0)
	out := make([]byte, 0, len(name)+1)
	for _, c := range []byte(name) {
		out = append(out, (c>>1)|(carry<<7))
		carry = c & 1
	}
	out = append(out, carry<<7)
	return out
}

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
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(len(*cvarsMap)))
	for k, v := range *cvarsMap {
		out = append(out, encodeCvarName(k)...)
		ev := encodeCvarFloat(v)
		out = append(out, 0x02, byte(ev>>8), byte(ev&0xFF), 0x00, 0x00)
	}
	out = append(out, 0x00)
	return out
}

type bitWriter struct {
	buf []byte
	bit int
}

func (bw *bitWriter) writeBit(b bool) {
	if bw.bit == 0 {
		bw.buf = append(bw.buf, 0)
	}
	if b {
		bw.buf[len(bw.buf)-1] |= 1 << (7 - bw.bit)
	}
	bw.bit = (bw.bit + 1) % 8
}

func (bw *bitWriter) writeBool(v bool) { bw.writeBit(v) }

func (bw *bitWriter) writeUint8(v uint8) {
	for i := 7; i >= 0; i-- {
		bw.writeBit((v>>i)&1 == 1)
	}
}

func (bw *bitWriter) writeUint16(v uint16) {
	for i := 15; i >= 0; i-- {
		bw.writeBit((v>>i)&1 == 1)
	}
}

func (bw *bitWriter) writeString(s string) {
	for _, c := range []byte(s) {
		bw.writeUint8(c)
	}
	bw.writeUint8(0)
}

// XXX: Why is this padded by two single bits...
func buildShardBody(shardIP string, shardPort int, chatIP string, chatPort int) []byte {
	bw := &bitWriter{}
	bw.writeBool(true)
	bw.writeUint8(1)
	bw.writeString(shardIP)
	bw.writeUint16(uint16(shardPort))
	bw.writeBool(true)
	bw.writeString(chatIP)
	bw.writeUint16(uint16(chatPort))
	return bw.buf
}

func handle(conn net.Conn, cfg *Config, cvars *CvarsMap) {
	defer conn.Close()
	log.Printf("%s: new connection", conn.RemoteAddr())

	if _, err := conn.Write(protocol.MakePacket(cvarsType, 2, buildCvarsBody(cvars))); err != nil {
		log.Printf("error sending cvars: %v", err)
		return
	}
	log.Printf("%s: sent cvars", conn.RemoteAddr())

	body := buildShardBody(cfg.ShardIP, cfg.ShardPort, cfg.ChatIP, cfg.ChatPort)
	if _, err := conn.Write(protocol.MakePacket(shardAddrType, 0, body)); err != nil {
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
