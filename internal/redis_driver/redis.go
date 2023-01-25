package redis_driver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var (
	_universalClient redis.UniversalClient
	_nodes           []RedisNode
)

func Start(csvAddrs string) {
	addrs := strings.Split(csvAddrs, ",")
	_universalClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addrs,
		DB:    0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := _universalClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("redia.Ping()")
	}

	_nodes = ClusterNodes()
	ShowClusterNodes(_nodes)
}

func Client() redis.UniversalClient {
	return _universalClient
}

func KeyNode(key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	retval, err := Client().ClusterKeySlot(ctx, key).Result()
	if err != nil {
		log.Warn().Err(err).Str("key", key).Msg("redis.ClusterKeySlot")
		return 0, err
	}

	slotId := int(retval)
	for _, node := range _nodes {
		if (node.SlotRangeStart <= slotId) && (node.SlotRangeEnd >= slotId) {
			return node.Id, nil
		}
	}

	return 0, errors.New("not-found node")
}
