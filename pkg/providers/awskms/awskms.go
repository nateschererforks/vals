package awskms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/service/kms"
	"gopkg.in/yaml.v3"

	"github.com/nateschererforks/vals/pkg/api"
	"github.com/nateschererforks/vals/pkg/awsclicompat"
)

type provider struct {
	// Keeping track of KMS services since we need a service per region
	client *kms.KMS

	// AWS KMS configuration
	Region, Profile, RoleARN                      string
	KeyId, EncryptionAlgorithm, EncryptionContext string
}

func New(cfg api.StaticConfig) *provider {
	p := &provider{}
	p.Region = cfg.String("region")
	p.Profile = cfg.String("profile")
	p.RoleARN = cfg.String("role_arn")
	p.KeyId = cfg.String("key")
	p.EncryptionAlgorithm = cfg.String("alg")
	p.EncryptionContext = cfg.String("context")
	return p
}

func (p *provider) GetString(key string) (string, error) {
	cli := p.getClient()

	blob, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	in := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	if p.KeyId != "" {
		in = in.SetKeyId(p.KeyId)
	}

	if p.EncryptionAlgorithm != "" {
		in = in.SetEncryptionAlgorithm(p.EncryptionAlgorithm)
	}

	if p.EncryptionContext != "" {
		m := map[string]*string{}

		if err := yaml.Unmarshal([]byte(p.EncryptionContext), &m); err != nil {
			return "", err
		}

		in = in.SetEncryptionContext(m)
	}

	result, err := cli.Decrypt(in)
	if err != nil {
		return "", err
	}

	return string(result.Plaintext), nil
}

func (p *provider) GetStringMap(key string) (map[string]interface{}, error) {
	yamlData, err := p.GetString(key)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(yamlData), &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (p *provider) getClient() *kms.KMS {
	if p.client != nil {
		return p.client
	}

	sess := awsclicompat.NewSession(p.Region, p.Profile, p.RoleARN)

	p.client = kms.New(sess)
	return p.client
}
