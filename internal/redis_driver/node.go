package redis_driver

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type RedisNode struct {
	Id             int
	Master         bool
	Host           string
	Port           int
	Connected      bool
	SlotRangeStart int
	SlotRangeEnd   int
}

func ClusterNodes() []RedisNode {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	retval, err := Client().ClusterNodes(ctx).Result()
	if err != nil {
		log.Warn().Err(err).Msg("redis.ClusterNodes")
		return nil
	}

	//log.Info().Str("nodes", retval).Send()
	// "0f27fd406afbd372d169492f99a10887cb16ee13 10.182.58.125:7201@17201 myself,master - 0 1674631613000 2 connected 4096-8191\nf944e680d2d6afa8536d8deb91132887f5a074ca 10.182.58.125:7203@17203 master - 0 1674631615030 4 connected 12288-16383\nbc3007f6838b912fdf826ec7f18e48304d9f6865 10.182.58.125:7202@17202 master - 0 1674631615000 3 connected 8192-12287\nb77678a6f31234ba9b5f974c528ee2b7efdb3730 10.182.58.125:7200@17200 master - 0 1674631614027 1 connected 0-4095\n"

	var nodes []RedisNode
	lines := strings.Split(retval, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "# ") {
			continue
		}

		//log.Debug().Str(strconv.Itoa(index), line).Msg("redis-node")
		node := parseNode(line)
		if node.Master && node.Connected {
			nodes = append(nodes, node)
			//log.Debug().Int("id", node.Id).Str("host", node.Host).Int("port", node.Port).
			//	Int("slot-start", node.SlotRangeStart).Int("slot-end", node.SlotRangeEnd).
			//	Msg("RedisNodes")
		}
	}

	return nodes
}

func ShowClusterNodes(nodes []RedisNode) {
	for _, node := range nodes {
		log.Debug().Int("id", node.Id).Str("host", node.Host).Int("port", node.Port).
			Int("slot-start", node.SlotRangeStart).Int("slot-end", node.SlotRangeEnd).
			Msg("RedisNodes")
	}
}

func parseNode(line string) RedisNode {
	xs := strings.Split(line, " ")
	addr, flags := xs[1], xs[2]

	id, err := strconv.Atoi(xs[len(xs)-3])
	if err != nil {
		return RedisNode{}
	}

	host, port := parseNodeAddr(addr)
	if (host == "") || (port == 0) {
		return RedisNode{}
	}

	var master bool
	if strings.Contains(flags, "master") {
		master = true
	}

	var connected bool
	if xs[len(xs)-2] == "connected" {
		connected = true
	}

	slotStart, slotEnd := parseSlotRange(xs[len(xs)-1])
	//log.Debug().Str("addr", addr).Bool("master", master).
	//	Str("host", host).Int("port", port).Str("status", status).
	//	Int("slot-start", slotStart).Int("slot-end", slotEnd).
	//	Msg("parseNode")

	return RedisNode{
		Id:             id,
		Master:         master,
		Host:           host,
		Port:           port,
		SlotRangeStart: slotStart,
		SlotRangeEnd:   slotEnd,
		Connected:      connected,
	}
}

func parseNodeAddr(addr string) (string, int) {
	xs := strings.Split(addr, ":")
	if len(xs) != 2 {
		return "", 0
	}
	host := xs[0]

	ports := strings.Split(xs[1], "@")
	port, err := strconv.Atoi(ports[0])
	if err != nil {
		return "", 0
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return host, port
}

func parseSlotRange(ranges string) (int, int) {
	xs := strings.Split(ranges, "-")
	if len(xs) != 2 {
		return 0, 0
	}

	start, err := strconv.Atoi(xs[0])
	if err != nil {
		return 0, 0
	}

	end, err := strconv.Atoi(xs[1])
	if err != nil {
		return 0, 0
	}

	return start, end
}
