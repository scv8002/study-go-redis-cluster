package main

import (
	"os"
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

	RedisTest()
	GoTest()
}
