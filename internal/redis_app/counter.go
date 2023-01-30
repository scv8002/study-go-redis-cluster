package redis_app

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"study-redis-cluster/internal/redis_driver"
)

var revisionFieldName string = "$rev"

type CounterDoc struct {
	documentFullPath string
	parentId         string
	createdAt        string
	fieldName        string
	value            float64
}

func TestCounter() {
	d := CounterDoc{
		documentFullPath: "test/counter/a",
		parentId:         "123",
		createdAt:        "456",
		fieldName:        "good",
		value:            1,
	}

	log.Debug().Msg("-------------------------------------------------")
	log.Debug().Msg("CREATE")
	Counter(d)
	log.Debug().Msg("-------------------------------------------------")
	log.Debug().Msg("UPDATE")
	d.parentId = "333"
	d.createdAt = "444"
	Counter(d)
}

func Counter(doc CounterDoc) (map[string]interface{}, error) {
	ltag := "Counter"

	keys := []string{doc.documentFullPath} // 1개 초과 등록되면 "CROSSSLOT Keys in request don't hash to the same slot" 발생함
	vals := []interface{}{}
	vals = append(vals, doc.parentId, doc.createdAt, doc.fieldName, doc.value)

	retval, err := _counterScript.Run(context.TODO(), redis_driver.Client(), keys, vals...).Result()
	if err != nil {
		log.Warn().Err(err).Msg(ltag)
		return nil, err
	}

	switch val := retval.(type) {
	case int64:
		log.Warn().Int64("result", val).Msg(ltag)
	case []interface{}:
		path, _ := val[0].(string)
		data, _ := val[1].([]interface{})

		var ret = map[string]interface{}{}
		for i := 0; i < len(data); i += 2 {
			k, _ := data[i].(string)
			v, _ := data[i+1].(string)
			ret[k] = v
		}
		log.Debug().Str("doc_path", path).Interface("kv", ret).Msg(ltag)
		return ret, nil
	default:
		log.Warn().Str("error", "unsupported data type").Msg(ltag)
	}

	return nil, errors.New("fail")
}
