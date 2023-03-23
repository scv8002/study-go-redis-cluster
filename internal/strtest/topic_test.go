package strtest

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSuiteRunDataStoreDelete(t *testing.T) {
	suite.Run(t, new(TestSuiteStr))
}

type TestSuiteStr struct {
	suite.Suite
}

func (suite *TestSuiteStr) SetupSuite() {
}

func (suite *TestSuiteStr) TearDownSuite() {
}

func (suite *TestSuiteStr) TestTopicGet() {
	tctag := "TestTopicGet"

	type args struct {
		serviceId string
		path      string
	}
	type expected struct {
		topic string
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name:     tctag + "-01",
			args:     args{serviceId: "SVC001", path: "/a/b"},
			expected: expected{topic: "de_datastore_SVC001/a/b"},
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			topic := GetDataStoreTopicId(tc.args.serviceId, tc.args.path)
			if tc.expected.topic != topic {
				t.Errorf("%s, topic actual %v, expected %v", tc.name, topic, tc.expected.topic)
				return
			}
		})
	}
}

func (suite *TestSuiteStr) TestTopicParse() {
	tctag := "TestTopicParse"

	type args struct {
		topic string
	}
	type expected struct {
		hasError  bool
		serviceId string
		path      string
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name:     tctag + "-01",
			args:     args{topic: "asd_sdf_d/24/g"},
			expected: expected{hasError: true},
		},
		{
			name:     tctag + "-02",
			args:     args{topic: "de_datastore_SVC002/1/b"},
			expected: expected{serviceId: "SVC002", path: "/1/b"},
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			serviceId, path, err := ParseDataStoreTopicId(tc.args.topic)
			if tc.expected.hasError && err == nil {
				t.Errorf("%s, body actual %v, expected %v", tc.name, err, tc.expected.hasError)
				return
			}
			if !tc.expected.hasError && err != nil {
				t.Errorf("%s, body actual %v, expected %v", tc.name, err, tc.expected.hasError)
				return
			}
			if tc.expected.hasError {
				return
			}

			if tc.expected.serviceId != serviceId {
				t.Errorf("%s, serviceId actual %v, expected %v", tc.name, serviceId, tc.expected.serviceId)
				return
			}
			if tc.expected.path != path {
				t.Errorf("%s, path actual %v, expected %v", tc.name, path, tc.expected.path)
				return
			}
		})
	}
}
