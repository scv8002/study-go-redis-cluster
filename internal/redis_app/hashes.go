package redis_app

import (
	"context"
	"fmt"
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
