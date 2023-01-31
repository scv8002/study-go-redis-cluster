package redis_app

import (
	"context"
	"github.com/rs/zerolog/log"

	"study-redis-cluster/internal/redis_driver"
)

func KeyTest() {
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	for _, key := range keys {
		id, err := redis_driver.KeyNode(key)
		if err != nil {
			log.Fatal().Err(err).Msg("KeyTest")
		}
		log.Debug().Str("key", key).Int("node-id", id).Msg("KeyTest")
	}
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=a node-id=4
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=b node-id=1
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=c node-id=2
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=d node-id=3
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=e node-id=4
	//2023-01-25T21:47:04+09:00 DBG internal/redis_app/nodes.go:17 > KeyTest key=f node-id=1
}

func TestWriteSingle() {
	ltag := "TestWriteSingle"

	keys := []string{"a"}
	vals := []interface{}{}
	vals = append(vals, "v:a", "v:a")

	result, err := _writeScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)
}

func TestWriteAll() {
	ltag := "TestWriteAll"

	keys := []string{"a", "b", "c", "d", "e", "f"}
	vals := []interface{}{}
	vals = append(vals, "v:a", "v:a")
	vals = append(vals, "v:b", "v:b")
	vals = append(vals, "v:c", "v:c")
	vals = append(vals, "v:d", "v:d")
	vals = append(vals, "v:e", "v:e")
	vals = append(vals, "v:f", "v:f")

	result, err := _writeScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)
}

func TestWriteSlotNode() {
	ltag := "TestWriteSlotNode"

	keys := []string{"a", "e"}
	vals := []interface{}{}
	vals = append(vals, "v:a", "v:a")
	vals = append(vals, "v:e", "v:e")

	result, err := _writeScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)
}

func TestWriteSlotNode2() {
	ltag := "TestWriteSlotNode"

	keys := []string{"a", "a"}
	vals := []interface{}{}
	vals = append(vals, "v:a", "v:a")
	vals = append(vals, "v:a-2", "v:a-2")

	result, err := _writeScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)
}
