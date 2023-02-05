package redis_app

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

func ReadAllKeys() {
	ltag := "ReadAllKeys"

	keys := []string{"a"}
	vals := []interface{}{}
	vals = append(vals, "a")
	vals = append(vals, "b")
	vals = append(vals, "c")
	vals = append(vals, "d")
	vals = append(vals, "e")
	vals = append(vals, "f")

	_, err := _readScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Result()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	// 에러발생
	// Lua script attempted to access a non local key in a cluster node
}

func ReadNodeKeys() {
	ltag := "ReadNodeKeys"

	keys := []string{"a"}
	vals := []interface{}{}
	vals = append(vals, "a")
	vals = append(vals, "e")
	vals = append(vals, "i")

	retval, err := _readScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Result()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}

	switch val := retval.(type) {
	case int64:
		log.Warn().Int64("result", val).Msg(ltag)
	case []interface{}:
		for i := 0; i < len(val); i += 2 {

			path, _ := val[i].(string)
			data, _ := val[i+1].([]interface{})

			var ret = map[string]interface{}{}
			for j := 0; j < len(data); j += 2 {
				k, _ := data[j].(string)
				v, _ := data[j+1].(string)
				ret[k] = v
				log.Info().Str("doc-path", path).Str("field", k).Str("value", v).Msg(ltag)
			}
		}
	default:
		log.Warn().Str("error", "unsupported data type").Msg(ltag)
	}
}

func ReadMeta() {
	ltag := "ReadNodeKeys"
	vals, err := _testHashesIpair.Run(context.TODO(), redis_driver.Client(), []string{"/129vn12/234"}).Result()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}

	log.Debug().Interface("test", vals).Msg("ReadMeta")
	switch val := vals.(type) {
	case int64:
		log.Warn().Int64("result", val).Msg(ltag)
	case []interface{}:
		for i := 0; i < len(val); i++ {
			log.Debug().Int("index", i).Str("value", fmt.Sprintf("%v", val[i])).Msg("ReadMeta")
		}
	default:
		log.Warn().Str("error", "unsupported data type").Msg(ltag)
	}
}
