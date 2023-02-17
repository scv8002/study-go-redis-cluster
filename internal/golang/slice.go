package golang

import "github.com/rs/zerolog/log"

func SliceAppend() {
	var fields []string
	ifMatch := 1

	fields = append(fields, "a")
	fields = append(fields, "b")
	fields = append(fields, "c")

	argv := []interface{}{ifMatch}

	for _, v := range fields {
		argv = append(argv, v)
	}

	for i, v := range argv {
		log.Debug().Int("index", i).Interface("value", v).Msg("SliceAppend")
	}
}
