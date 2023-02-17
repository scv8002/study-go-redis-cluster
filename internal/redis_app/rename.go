package redis_app

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

var _renameLua = redis.NewScript(`
	local src = KEYS[1]
	local dst = ARGV[1]

	local ret = redis.call("rename", KEYS[1], ARGV[1])

	return ret
`)

func KeyRename() {
	ltag := "KeyRename"

	keys := []string{"a"}
	vals := []interface{}{"{a}deleted"}

	// key가 존재하지 않아도 err를 리턴한다.
	// 성공일 경우만 result를 리턴한다.
	result, err := _renameLua.Run(context.TODO(), redis_driver.Client(), keys, vals...).Result()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Interface("result", result).Msg(ltag)

}
