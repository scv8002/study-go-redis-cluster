package redis_app

import (
	"github.com/rs/zerolog/log"

	"study-redis-cluster/internal/redis_driver"
)

func KeyTest() {
	keys := []string{"a", "b", "c", "d", "e", "f"}

	for _, key := range keys {
		id, err := redis_driver.KeyNode(key)
		if err != nil {
			log.Fatal().Err(err).Msg("KeyTest")
		}
		log.Debug().Str("key", key).Int("node-id", id).Msg("KeyTest")
	}
}
