package redis_app

import (
	"context"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

func WriteAllKeys() {
	ltag := "MultiKeysWrite"

	keys := []string{"a"}
	vals := []interface{}{}
	vals = append(vals, "a", "v:a", "v:a")
	vals = append(vals, "b", "v:b", "v:b")
	vals = append(vals, "c", "v:c", "v:c")
	vals = append(vals, "d", "v:d", "v:d")
	vals = append(vals, "e", "v:e", "v:e")
	vals = append(vals, "f", "v:f", "v:f")

	result, err := _writeScript2.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)

	// 에러 발생
	// "ERR Error running script (call to f_b663be011db0173179af274d2dad5a11fb8094f4): @user_script:10: @user_script: 10: Lua script attempted to access a non local key in a cluster node"
}

func WriteNodeKeys() {
	ltag := "WriteNodeKeys"

	keys := []string{"a"}
	vals := []interface{}{}
	vals = append(vals, "a", "v:a", "v:a")
	vals = append(vals, "e", "v:e", "v:e")

	result, err := _writeScript2.Run(context.TODO(), redis_driver.Client(), keys, vals...).Int()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return
	}
	log.Debug().Int("result", result).Msg(ltag)
}
