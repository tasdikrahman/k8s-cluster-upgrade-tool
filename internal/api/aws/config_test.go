package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockAwsConfig struct {
	mock.Mock
}

func (m *mockAwsConfig) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	args := m.Called(ctx, optFns[0], optFns[1])
	return args.Get(0).(aws.Config), args.Error(1)
}

func TestAwsConfig_LoadDefaultConfig(t *testing.T) {
	t.Run("when the LoadDefaultConfig is passed with the right configuration and it returns the aws config without any error", func(t *testing.T) {
		m := new(mockAwsConfig)

		m.On("LoadDefaultConfig",
			mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("func(*config.LoadOptions) error"), mock.AnythingOfType("func(*config.LoadOptions) error")).
			Return(aws.Config{Region: "correct-region"}, nil).
			Once()

		s := ConfigGetter{m}

		cfg, err := s.GetConfig(context.TODO(), config.WithRegion("correct-region"), config.WithSharedConfigProfile("correct-aws-profile"))

		assert.Nil(t, err)
		assert.IsType(t, cfg, aws.Config{})
	})

	t.Run("when the LoadDefaultConfig is passed with incorrect configuration and it returns an empty aws config with any error", func(t *testing.T) {
		m := new(mockAwsConfig)

		m.On("LoadDefaultConfig",
			mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("func(*config.LoadOptions) error"), mock.AnythingOfType("func(*config.LoadOptions) error")).
			Return(aws.Config{}, errors.New("some aws config error")).
			Once()

		s := ConfigGetter{m}

		cfg, err := s.GetConfig(context.TODO(), config.WithRegion("incorrect-region"), config.WithSharedConfigProfile("incorrect-aws-profile"))

		assert.NotNil(t, err)
		assert.IsType(t, cfg, aws.Config{})
	})
}
