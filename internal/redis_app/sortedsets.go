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

func TestZRange() {
	// 10.182.58.125:7201> zadd myzip 2 aa 3 bb 5 cc 4 dd 6 ff 1 gg
	// (integer) 6
	// 10.182.58.125:7201> zrange myzip 0 -1
	// 1) "gg"
	// 2) "aa"
	// 3) "bb"
	// 4) "dd"
	// 5) "cc"
	// 6) "ff"
	// 10.182.58.125:7200> zrange myzip 0 -1 withscores
	// 1) "gg"
	// 2) "1"
	// 3) "aa"
	// 4) "2"
	// 5) "bb"
	// 6) "3"
	// 7) "dd"
	// 8) "4"
	// 9) "cc"
	// 10) "5"
	// 11) "ff"
	// 12) "6"
	// 10.182.58.125:7201> zrange myzip 1 2
	// 1) "aa"
	// 2) "bb"
	// 10.182.58.125:7201> zrange myzip 1 2 withscores
	// 1) "aa"
	// 2) "2"
	// 3) "bb"
	// 4) "3"
	// 10.182.58.125:7201> zrange myzip 1 2 withscores rev
	// 1) "cc"
	// 2) "5"
	// 3) "dd"
	// 4) "4"
	// 10.182.58.125:7201>

	// zadd "/test/zset" 2 aa 3 bb 5 cc 4 dd 6 ff 1 gg
	path := "/test/zset"
	zrangeStart := 1
	zrangeStop := 3
	descending := false

	args := []interface{}{zrangeStart, zrangeStop}
	if descending {
		args = append(args, "rev")
	}

	vals, err := _zrange.Run(context.TODO(), redis_driver.Client(), []string{path}, args...).Result()
	if err != nil {
		log.Warn().Err(err).Msg("TestZRange")
		return
	}

	switch val := vals.(type) {
	case int64:
		switch val {
		case -1:
			log.Warn().Int64("val", val).Msg("TestZRange")
		default:
			log.Warn().Int64("val", val).Msg("TestZRange")
		}
	case []interface{}:
		for i, v := range val {
			log.Debug().Int("index", i).Interface("val", v).Msg("TestZRange")
		}
	default:
		err = fmt.Errorf("MUST NOT reach here [TestZRange:2]")
	}
}

var _zrange = redis.NewScript(`
	local path = KEYS[1]
	local rev = ARGV[3]

	local args = {}
	table.insert(args, ARGV[1])
	table.insert(args, ARGV[2])
	table.insert(args, "withscores")
	if rev == "rev" then
		table.insert(args, "rev")
	end

	local vals = redis.call("ZRANGE", path, unpack(args))
	if #vals == 0 then
		return -1
	end

	return vals
`)
