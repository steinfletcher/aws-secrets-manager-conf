package secretsmanager_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/steinfletcher/conf-aws-secrets-manager/secretsmanager"

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

type configBool struct {
	SecretBool bool `secret:"/my-group/my-secret"`
}

type configWithRequired struct {
	SecretPlaintext string `secret:"/my-group/my-secret,required"`
}

type configWithDefault struct {
	SecretPlaintext string `secret:"/my-group/my-secret" secretDefault:"12345"`
}

type configWithJSON struct {
	SecretJSON jsonSecret `secret:"/my-group/my-secret"`
}

type jsonSecret struct {
	ApiKey    string `json:"api_key"`
	IsEnabled bool   `json:"is_enabled"`
}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		config      interface{}
		awsErr      error
		secret      string
		expected    interface{}
		expectedErr error
	}{
		"success": {
			config: &config{},
			secret: "myPlaintextSecretValue",
			expected: &config{
				SecretPlaintext: "myPlaintextSecretValue",
			},
		},
		"success with bool": {
			config: &configBool{},
			secret: "true",
			expected: &configBool{
				SecretBool: true,
			},
		},
		"success with json payload": {
			config: &configWithJSON{},
			secret: `{"api_key": "12345", "is_enabled": false}`,
			expected: &configWithJSON{
				SecretJSON: jsonSecret{
					ApiKey:    "12345",
					IsEnabled: false,
				},
			},
		},
		"default value": {
			config: &configWithDefault{},
			secret: "",
			awsErr: awserr.New(awsSecretsManager.ErrCodeResourceNotFoundException, "", nil),
			expected: &configWithDefault{
				SecretPlaintext: "12345",
			},
		},
		"not required": {
			config: &config{},
			secret: "",
			expected: &config{
				SecretPlaintext: "",
			},
		},
		"error if required and not present": {
			config:      &configWithRequired{},
			awsErr:      awserr.New(awsSecretsManager.ErrCodeResourceNotFoundException, "", nil),
			expected:    &configWithRequired{},
			expectedErr: errors.New(`conf: required variable "/my-group/my-secret" is not set`),
		},
		"aws error": {
			config:      &configWithDefault{},
			secret:      "",
			awsErr:      awserr.New(awsSecretsManager.ErrCodeInternalServiceError, "", nil),
			expected:    &configWithDefault{},
			expectedErr: awserr.New(awsSecretsManager.ErrCodeInternalServiceError, "", nil),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			secretsManager := mocks.NewMockSecretsManager(ctrl)
			secretsManager.EXPECT().
				GetSecretValue(&awsSecretsManager.GetSecretValueInput{SecretId: aws.String("/my-group/my-secret")}).
				Return(&awsSecretsManager.GetSecretValueOutput{SecretString: &test.secret}, test.awsErr)
			provider := secretsmanager.NewSecretsManagerProvider(secretsManager)

			err := conf.Parse(test.config, provider)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, test.config)
		})
	}
}
