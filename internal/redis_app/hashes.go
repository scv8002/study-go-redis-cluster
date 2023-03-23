package redis_app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

func TestHMGet() {
	path := "/asdf/q23r"
	sli, err := redis_driver.Client().HMGet(context.TODO(), path, "$parentId", "$createdAt").Result()
	if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("TestHMGet")
		return
	}
	log.Debug().Int("slice len", len(sli)).Str("path", path).Msg("TestHMGet")
	// redis에 hashes가 존재하지 않아도 요청 field 쌍에 대응하는 nil 값을 갖는 value 2개 리턴한다.

	for i := 0; i < len(sli); i++ {
		if sli[i] != nil {
			log.Debug().Int("hmget index", i).Str("slice", fmt.Sprintf("%v", sli[i])).Str("path", path).Msg("TestHMGet")
		} else {
			log.Debug().Int("hmget index", i).Str("val", "<<nil>>").Str("path", path).Msg("TestHMGet")
		}
	}

	//	val := sli[0].(string) // panic: interface conversion: interface {} is nil, not string
	//	if len(val) == 0 {
	//		log.Debug().Msg("value[0] string len is 0")
	//	}
	//
	//	conval, err := strconv.ParseInt(val, 10, 64)
	//	if err != nil {
	//		log.Warn().Err(err).Msg("strconv.ParseInt")
	//	} else {
	//		log.Debug().Int64("conv", conval).Msg("strconv.ParseInt")
	//	}
}

var _hashes_scan_fields = redis.NewScript(`
	local cursor = 0
	local fields = {}
	local results = {}
	repeat
		local result = redis.call("HSCAN", KEYS[1], cursor, 'MATCH', ARGV[1])
		cursor = tonumber(result[1])
		fields = result[2]
		for i, v in ipairs(fields) do
			table.insert(results, v)
		end
	until cursor == 0
	return results
`)

func TestHSCAN() {
	// 62개 fields를 등록해도 cursor가 나뉘지 않고 한번에 scan됨
	// hmset user-1 email charlie@redisgate.com language English gender m
	// hmset user-1 "$col:a" "12345" "$col:b" "23456"
	// hscan user-1 0 match $col:*
	// 1) "0"
	// 2) 1) "$col:a"
	//    2) "12345"
	//    3) "$col:b"
	//    4) "23456"
	// hmset user-1 "$col:c" "c1" "$col:d" "d1" "$col:e" "e1" "$col:f" "f1" "$col:g" "g1"
	// hmset user-1 "$col:h" "h1" "$col:i" "i1" "$col:j" "j1" "$col:k" "k1" "$col:l" "l1"
	// hmset user-1 "$col:m" "m1" "$col:n" "n1" "$col:o" "o1" "$col:p" "p1" "$col:q" "q1"
	// hmset user-1 "$col:r" "r1" "$col:s" "s1" "$col:t" "t1" "$col:u" "u1" "$col:v" "v1"
	// hmset user-1 "$col:1c" "1c1" "$col:1d" "1d1" "$col:1e" "1e1" "$col:1f" "1f1" "$col:1g" "1g1"
	// hmset user-1 "$col:1h" "1h1" "$col:1i" "1i1" "$col:1j" "1j1" "$col:1k" "1k1" "$col:1l" "1l1"
	// hmset user-1 "$col:1m" "1m1" "$col:1n" "1n1" "$col:1o" "1o1" "$col:1p" "1p1" "$col:1q" "1q1"
	// hmset user-1 "$col:1r" "1r1" "$col:1s" "1s1" "$col:1t" "1t1" "$col:1u" "1u1" "$col:1v" "1v1"
	// hmset user-1 "$col:2c" "1c1" "$col:2d" "1d1" "$col:2e" "1e1" "$col:2f" "1f1" "$col:2g" "1g1"
	// hmset user-1 "$col:2h" "1h1" "$col:2i" "1i1" "$col:2j" "1j1" "$col:2k" "1k1" "$col:2l" "1l1"
	// hmset user-1 "$col:2m" "1m1" "$col:2n" "1n1" "$col:2o" "1o1" "$col:2p" "1p1" "$col:2q" "1q1"
	// hmset user-1 "$col:2r" "1r1" "$col:2s" "1s1" "$col:2t" "1t1" "$col:2u" "1u1" "$col:2v" "1v1"
	vals, err := _hashes_scan_fields.Run(context.TODO(), redis_driver.Client(), []string{"user-1"}, "$col:*").Result()
	if err != nil {
		log.Warn().Err(err).Msg("TestHSCAN")
		return
	}
	switch val := vals.(type) {
	case int64:
		switch val {
		case -1:
			log.Warn().Int64("val", val).Msg("TestHSCAN")
		default:
			log.Warn().Int64("val", val).Msg("TestHSCAN")
		}
	case []interface{}:
		for i := 0; (i + 2) <= len(val); i += 2 {
			field, _ := val[i].(string)
			value, _ := val[i+1].(string)
			log.Debug().Str("field", field).Str("val", value).Msg("TestHSCAN")
		}
	default:
		err = fmt.Errorf("MUST NOT reach here [TestHSCAN:2]")
	}
}

func TestHSCANCast() {
	tcname := "TestHSCANCast"
	// 62개 fields를 등록해도 cursor가 나뉘지 않고 한번에 scan됨
	// hmset user-1 email charlie@redisgate.com language English gender m
	// hmset user-1 "$col:a" "12345" "$col:b" "23456"
	// hscan user-1 0 match $col:*
	// 1) "0"
	// 2) 1) "$col:a"
	//    2) "12345"
	//    3) "$col:b"
	//    4) "23456"
	// hmset user-1 "$col:c" "c1" "$col:d" "d1" "$col:e" "e1" "$col:f" "f1" "$col:g" "g1"
	// hmset user-1 "$col:h" "h1" "$col:i" "i1" "$col:j" "j1" "$col:k" "k1" "$col:l" "l1"
	// hmset user-1 "$col:m" "m1" "$col:n" "n1" "$col:o" "o1" "$col:p" "p1" "$col:q" "q1"
	// hmset user-1 "$col:r" "r1" "$col:s" "s1" "$col:t" "t1" "$col:u" "u1" "$col:v" "v1"
	// hmset user-1 "$col:1c" "1c1" "$col:1d" "1d1" "$col:1e" "1e1" "$col:1f" "1f1" "$col:1g" "1g1"
	// hmset user-1 "$col:1h" "1h1" "$col:1i" "1i1" "$col:1j" "1j1" "$col:1k" "1k1" "$col:1l" "1l1"
	// hmset user-1 "$col:1m" "1m1" "$col:1n" "1n1" "$col:1o" "1o1" "$col:1p" "1p1" "$col:1q" "1q1"
	// hmset user-1 "$col:1r" "1r1" "$col:1s" "1s1" "$col:1t" "1t1" "$col:1u" "1u1" "$col:1v" "1v1"
	// hmset user-1 "$col:2c" "1c1" "$col:2d" "1d1" "$col:2e" "1e1" "$col:2f" "1f1" "$col:2g" "1g1"
	// hmset user-1 "$col:2h" "1h1" "$col:2i" "1i1" "$col:2j" "1j1" "$col:2k" "1k1" "$col:2l" "1l1"
	// hmset user-1 "$col:2m" "1m1" "$col:2n" "1n1" "$col:2o" "1o1" "$col:2p" "1p1" "$col:2q" "1q1"
	// hmset user-1 "$col:2r" "1r1" "$col:2s" "1s1" "$col:2t" "1t1" "$col:2u" "1u1" "$col:2v" "1v1"
	vals, err := _hashes_scan_fields.Run(context.TODO(), redis_driver.Client(), []string{"user-1"}, "$col:*").Result()
	if err != nil {
		log.Warn().Err(err).Msg(tcname)
		return
	}
	switch val := vals.(type) {
	case int64:
		switch val {
		case -1:
			log.Warn().Int64("val", val).Msg(tcname)
		default:
			log.Warn().Int64("val", val).Msg(tcname)
		}
	case []string:
		for i := 0; (i + 2) <= len(val); i += 2 {
			field := val[i]
			value := val[i+1]
			log.Debug().Str("field", field).Str("val", value).Msg(tcname)
		}
	default:
		log.Warn().Err(fmt.Errorf("MUST NOT reach here [" + tcname + ":2]")).Msg(tcname)
	}
}

var _hashes_not_found_field = redis.NewScript(`
	local result = redis.call("HGET", KEYS[1], ARGV[1])
	if result then
		return -2
	end
	if not result then
		return -1
	end

	return result
`)

func HGetNotFound() {
	vals, err := _hashes_not_found_field.Run(context.TODO(), redis_driver.Client(), []string{"user-1"}, "$expires").Result()
	if err != nil {
		log.Warn().Err(err).Msg("HGetNotFound")
		return
	}
	switch val := vals.(type) {
	case int64:
		switch val {
		case -1: // field 가 없는 경우
			log.Warn().Str("type", "int64").Int64("val", val).Msg("HGetNotFound")
		case -2: // field 가 존재하는 경우
			log.Warn().Str("type", "int64").Int64("val", val).Msg("HGetNotFound")
		default:
			log.Warn().Msg("MUST NOT reach here")
		}
	case string: // field가 존재하는 경우
		log.Warn().Str("type", "string").Str("val", val).Msg("HGetNotFound")
	default:
		err = fmt.Errorf("MUST NOT reach here [TestHSCAN:2]")
	}
}

var _hash_mget_test = redis.NewScript(`
	local docFullPath = KEYS[1]
	local docIndex = tonumber(ARGV[1])
	local ifMatch = tonumber(ARGV[2])
	local pattern = ARGV[3]

	return docIndex
`)

func HMGetTest() {
	keys := []string{"asdf"}
	args := []interface{}{"0", 0, "$col:*"}
	vals, err := _hash_mget_test.Run(context.TODO(), redis_driver.Client(), keys, args).Result()

	if err != nil {
		log.Warn().Err(err).Msg("HMGetTest")
		return
	}
	switch val := vals.(type) {
	case int64:
		switch val {
		case 0:
			log.Debug().Str("type", "int64").Int64("val", val).Msg("HMGetTest")
		case -1: // field 가 없는 경우
			log.Warn().Str("type", "int64").Int64("val", val).Msg("HMGetTest")
		case -2: // field 가 존재하는 경우
			log.Warn().Str("type", "int64").Int64("val", val).Msg("HMGetTest")
		default:
			log.Warn().Msg("MUST NOT reach here")
		}
	case string: // field가 존재하는 경우
		log.Warn().Str("type", "string").Str("val", val).Msg("HMGetTest")
	default:
		err = fmt.Errorf("MUST NOT reach here [HMGetTest:2]")
	}
}
