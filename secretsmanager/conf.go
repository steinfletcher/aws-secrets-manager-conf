package secretsmanager

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/steinfletcher/conf"
)

type asmConf struct {
	secretsManager secretsmanageriface.SecretsManagerAPI
}

func NewProvider(secretsManager secretsmanageriface.SecretsManagerAPI) conf.Provider {
	return asmConf{secretsManager: secretsManager}
}

func (o asmConf) Provide(field reflect.StructField) (string, error) {
	key, opts := parseTag(field, "secret")
	if key == "" {
		return "", nil
	}

	defaultValue, _ := parseTag(field, "secretDefault")
	isRequired := hasOption(opts, "required")

	return getValue(o.secretsManager, key, defaultValue, isRequired)
}

func hasOption(opts []string, name string) bool {
	for _, opt := range opts {
		if opt == name {
			return true
		}
	}

	return false
}

// split the tag into the key and options, e.g. ("/my-secret-key", []{"required"})
func parseTag(field reflect.StructField, name string) (string, []string) {
	tagValue := field.Tag.Get(name)
	opts := strings.Split(tagValue, ",")

	return opts[0], opts[1:]
}

func getValue(secretsManager secretsmanageriface.SecretsManagerAPI, key, defaultValue string, isRequired bool) (string, error) {
	value, err := fetchSecret(secretsManager, key)
	if err != nil {
		return "", err
	}

	if value == "" && defaultValue == "" && isRequired {
		return "", fmt.Errorf(`conf: required variable %q is not set`, key)
	}

	if value == "" && defaultValue != "" {
		value = defaultValue
	}

	return value, nil
}

func fetchSecret(secretsManager secretsmanageriface.SecretsManagerAPI, key string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	output, err := secretsManager.GetSecretValue(input)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
				return "", nil
			}
		}

		return "", err
	}

	return *output.SecretString, nil
}
