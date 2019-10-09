package asm_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/conf"
	"github.com/steinfletcher/conf-aws-secrets-manager/asm"
	"github.com/steinfletcher/conf-aws-secrets-manager/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

type config struct {
	SecretPlaintext string `secret:"/my-group/my-secret"`
}

func TestParse_PlaintextSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	secretsManager := mocks.NewMockSecretsManager(ctrl)
	plaintextSecretPayload := "myPlaintextSecretValue"
	secretsManager.EXPECT().
		GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String("/my-group/my-secret")}).
		Return(&secretsmanager.GetSecretValueOutput{SecretString: &plaintextSecretPayload}, nil)
	provider := asm.NewSecretsManagerProvider(secretsManager)
	var cfg config

	err := conf.Parse(&cfg, provider)

	assert.NoError(t, err)
	assert.Equal(t, plaintextSecretPayload, cfg.SecretPlaintext)
}
