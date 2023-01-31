package main

import (
	"os"
	"study-redis-cluster/internal/redis_app"
	"study-redis-cluster/internal/redis_driver"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorFieldName = "reason"

	if true {
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
}

func main() {
	log.Print("hello world")

	redis_driver.Start("10.182.58.125:7200, 10.182.58.125:7201")

	if false {
		redis_app.KeyTest()
	}

	if false {
		redis_app.TestWriteSingle()
		redis_app.TestWriteAll()
		redis_app.TestWriteSlotNode()
		redis_app.TestWriteSlotNode2() // 동일한 slot이 아니라면 일괄작업 불가능
		//2023-01-25T22:34:00+09:00 DBG TestWriteSingle result=0
		//2023-01-25T22:34:00+09:00 WRN TestWriteAll reason="CROSSSLOT Keys in request don't hash to the same slot"
		//2023-01-25T22:34:00+09:00 WRN TestWriteSlotNode reason="CROSSSLOT Keys in request don't hash to the same slot"
		//2023-01-25T22:34:00+09:00 DBG TestWriteSlotNode result=0
	}

	if false {
		redis_app.TestWriteSlotNode2()
	}

	if false {
		redis_app.TestCounter()
	}

	if false {
		// redis_app.WriteAllKeys() // 실패 : Lua script attempted to access a non local key in a cluster node
		redis_app.WriteNodeKeys()
	}

	if true {
		redis_app.ReadNodeKeys()
	}
}
