[![Build Status](https://travis-ci.org/steinfletcher/aws-secrets-manager-conf.svg?branch=master)](https://travis-ci.org/steinfletcher/aws-secrets-manager-conf)

# aws-secrets-manager-conf

Provides AWS Secrets Manager support for [conf](https://github.com/steinfletcher/conf).

## Usage

Initialize the secrets manager provider by passing the AWS secrets manager instance

```go
import "github.com/steinfletcher/aws-secrets-manager-conf/secretsmanager"
import awsSecretsManager "github.com/aws/aws-sdk-go/service/secretsmanager"

...

provider := secretsmanager.NewProvider(awsSecretsManager.New(session.New()))
```

Then parse the configuration with [conf](https://github.com/steinfletcher/conf)

```go
err := conf.Parse(test.config, provider)
```

Use the `secret` struct tag to resolve secrets

```go
type Config struct {
	// github.com/caarlos0/env properties
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`

	// secrets manager properties
	MySecret     string        `secret:"/my_secret,required"`
}
```
