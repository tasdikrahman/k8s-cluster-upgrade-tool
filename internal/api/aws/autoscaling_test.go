package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockAutoScalingGroupApi struct {
	mock.Mock
}

func (m *mockAutoScalingGroupApi) UpdateAutoScalingGroupCount(ctx context.Context, cfg aws.Config) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	args := m.Called(ctx, cfg)
	return args.Get(0).(*autoscaling.UpdateAutoScalingGroupOutput), args.Error(1)
}

func TestAutoscalingGroup_SetMaxInstanceSizeToCurrentDesiredSize(t *testing.T) {
	t.Run("when the autoscaling group update call is successful", func(t *testing.T) {
		m := new(mockAutoScalingGroupApi)

		m.On("UpdateAutoScalingGroupCount",
			mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("aws.Config")).
			Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil).
			Once()

		s := AutoscalingGroupUpdater{m}

		result, err := s.Update(context.TODO(), aws.Config{})

		assert.Nil(t, err)
		assert.IsType(t, result, &autoscaling.UpdateAutoScalingGroupOutput{})
	})

	t.Run("when the autoscaling group update call is not successful", func(t *testing.T) {
		m := new(mockAutoScalingGroupApi)

		m.On("UpdateAutoScalingGroupCount",
			mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("aws.Config")).
			Return(&autoscaling.UpdateAutoScalingGroupOutput{}, errors.New("some error")).
			Once()

		s := AutoscalingGroupUpdater{m}

		result, err := s.Update(context.TODO(), aws.Config{})

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
