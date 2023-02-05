package redis_app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

func TestZScore() {
	path := "/asdf/q23r/111"
	val, err := redis_driver.Client().ZScore(context.TODO(), path, "$parentId").Result()
	if err != nil {
		// redis에 ZSets가 존재하지 않는다면 error를 리턴한다.
		log.Warn().Err(err).Msg("TestZScore")
		// log : TestZScore reason="redis: nil"
		// ----
		// 10.182.58.125:7200> zscore /asdf/q23r/111 "$parentId"
		// -> Redirected to slot [14981] located at 10.182.58.125:7203
		// (nil)
		// 10.182.58.125:7203> zscore /asdf/q23r/111 "$parentId"
		// (nil)
	} else {
		log.Debug().Float64("score", val).Msg("TestZScore")
	}

	// 데이터 추가
	// 10.182.58.125:7203> zadd /asdf/q23r/111 9876 "$parentId"
	// (integer) 1
	// 10.182.58.125:7203> zscore /asdf/q23r/111 "$parentId"
	// "9876"
	// member 삭제
	// 10.182.58.125:7203> zrem /asdf/q23r/111 "$parentId"
	// (integer) 1
}

var _zsets_not_found_member = redis.NewScript(`
	local key = KEYS[1]
	local member = ARGV[1]
	local score = redis.call("ZSCORE", key, member)

	if not score then
		return "<<NOT-FOUND>>"
	end

	return score
`)

func TestZScoreNotFound() {
	val, err := _zsets_not_found_member.Run(context.TODO(), redis_driver.Client(), []string{"/asdf1234"}, "qwer").Result()
	if err != nil {
		log.Warn().Err(err).Msg("TestZScoreNotFound")
		return
	}
	switch v := val.(type) {
	case int64:
		log.Debug().Str("type", "int64").Int64("val", v).Msg("TestZScoreNotFound")
	case string:
		log.Debug().Str("type", "string").Str("val", v).Msg("TestZScoreNotFound")
	default:
		log.Debug().Str("unknown", fmt.Sprintf("%v", val)).Msg("TestZScoreNotFound")
	}
}
