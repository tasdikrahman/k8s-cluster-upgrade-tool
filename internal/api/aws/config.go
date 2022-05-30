package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct{}

type ConfigClientInterface interface {
	LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)
}

func (a *Config) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, optFns[0], optFns[1])
	return cfg, err
}

type ConfigGetter struct {
	ConfigClientInterface
}

func (a *ConfigGetter) GetConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	cfg, err := a.LoadDefaultConfig(ctx, optFns[0], optFns[1])

	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}
