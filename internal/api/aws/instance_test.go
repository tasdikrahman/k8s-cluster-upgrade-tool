package aws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAwsInstances_Count(t *testing.T) {
	tests := []struct {
		name  string
		input AwsInstances
		want  int
	}{
		{
			"when there are 2 Aws Instances, it should return 2",
			AwsInstances{
				{"instanceID1", "privdns.1", "asgname1"},
				{"instanceID2", "privdns.1", "asgname1"},
			}, 2,
		},
		{
			"when the Aws Instances are empty it should return 0",
			AwsInstances{}, 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.input.Count(), tt.want)
		})
	}
}
