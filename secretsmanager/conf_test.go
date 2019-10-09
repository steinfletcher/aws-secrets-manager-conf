package secretsmanager_test

import (
	"github.com/steinfletcher/conf-aws-secrets-manager/secretsmanager"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awsSecretsManager "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/conf"
	"github.com/steinfletcher/conf-aws-secrets-manager/mocks"
	"github.com/stretchr/testify/assert"
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
		GetSecretValue(&awsSecretsManager.GetSecretValueInput{SecretId: aws.String("/my-group/my-secret")}).
		Return(&awsSecretsManager.GetSecretValueOutput{SecretString: &plaintextSecretPayload}, nil)
	provider := secretsmanager.NewSecretsManagerProvider(secretsManager)
	var cfg config

	err := conf.Parse(&cfg, provider)

	assert.NoError(t, err)
	assert.Equal(t, plaintextSecretPayload, cfg.SecretPlaintext)
}
