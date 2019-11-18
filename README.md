[![Build Status](https://travis-ci.org/steinfletcher/aws-secrets-manager-conf.svg?branch=master)](https://travis-ci.org/steinfletcher/aws-secrets-manager-conf)

# aws-secrets-manager-conf

Provides AWS Secrets Manager support for [conf](https://github.com/steinfletcher/conf).

## Usage

Initialize the secrets manager provider by passing the AWS secrets manager instance

```go
provider := secretsmanager.NewProvider(secretsManager)
```

Then parse the configuration with [conf](https://github.com/steinfletcher/conf)

```go
err := conf.Parse(test.config, provider)
```
