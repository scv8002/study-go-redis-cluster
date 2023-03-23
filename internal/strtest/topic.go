package strtest

import (
	"fmt"
	"strings"
)

const (
	dataStoreTopicPrefix = "de_datastore_"
)

func GetDataStoreTopicId(serviceId, collectionPath string) string {
	return fmt.Sprintf("%s%s%s", dataStoreTopicPrefix, serviceId, collectionPath)
}

func ParseDataStoreTopicId(topic string) (serviceId, collectionPath string, err error) {
	if !strings.HasPrefix(topic, dataStoreTopicPrefix) {
		err = fmt.Errorf("invalid DataStoreTopic path:%s", topic)
		return
	}

	topic = strings.TrimPrefix(topic, dataStoreTopicPrefix)

	pos := strings.Index(topic, "/")
	if pos == -1 {
		err = fmt.Errorf("invalid DataStoreTopic path:%s", topic)
		return
	}

	serviceId = topic[0:pos]
	collectionPath = topic[pos:]

	return
}
