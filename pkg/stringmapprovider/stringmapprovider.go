package stringmapprovider

import (
	"fmt"

	"github.com/nateschererforks/vals/pkg/api"
	"github.com/nateschererforks/vals/pkg/log"
	"github.com/nateschererforks/vals/pkg/providers/awskms"
	"github.com/nateschererforks/vals/pkg/providers/awssecrets"
	"github.com/nateschererforks/vals/pkg/providers/azurekeyvault"
	"github.com/nateschererforks/vals/pkg/providers/doppler"
	"github.com/nateschererforks/vals/pkg/providers/gcpsecrets"
	"github.com/nateschererforks/vals/pkg/providers/gkms"
	"github.com/nateschererforks/vals/pkg/providers/httpjson"
	"github.com/nateschererforks/vals/pkg/providers/k8s"
	"github.com/nateschererforks/vals/pkg/providers/onepasswordconnect"
	"github.com/nateschererforks/vals/pkg/providers/sops"
	"github.com/nateschererforks/vals/pkg/providers/ssm"
	"github.com/nateschererforks/vals/pkg/providers/vault"
)

func New(l *log.Logger, provider api.StaticConfig) (api.LazyLoadedStringMapProvider, error) {
	tpe := provider.String("name")

	switch tpe {
	case "s3":
		return ssm.New(l, provider), nil
	case "ssm":
		return ssm.New(l, provider), nil
	case "vault":
		return vault.New(l, provider), nil
	case "awssecrets":
		return awssecrets.New(l, provider), nil
	case "sops":
		return sops.New(l, provider), nil
	case "gcpsecrets":
		return gcpsecrets.New(provider), nil
	case "azurekeyvault":
		return azurekeyvault.New(provider), nil
	case "awskms":
		return awskms.New(provider), nil
	case "onepasswordconnect":
		return onepasswordconnect.New(provider), nil
	case "doppler":
		return doppler.New(l, provider), nil
	case "gkms":
		return gkms.New(l, provider), nil
	case "k8s":
		return k8s.New(l, provider)
	case "httpjson":
		return httpjson.New(l, provider), nil
	}

	return nil, fmt.Errorf("failed initializing string-map provider from config: %v", provider)
}
