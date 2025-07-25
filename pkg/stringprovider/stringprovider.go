package stringprovider

import (
	"fmt"

	"github.com/nateschererforks/vals/pkg/api"
	"github.com/nateschererforks/vals/pkg/log"
	"github.com/nateschererforks/vals/pkg/providers/awskms"
	"github.com/nateschererforks/vals/pkg/providers/awssecrets"
	"github.com/nateschererforks/vals/pkg/providers/azurekeyvault"
	"github.com/nateschererforks/vals/pkg/providers/conjur"
	"github.com/nateschererforks/vals/pkg/providers/doppler"
	"github.com/nateschererforks/vals/pkg/providers/gcpsecrets"
	"github.com/nateschererforks/vals/pkg/providers/gcs"
	"github.com/nateschererforks/vals/pkg/providers/gitlab"
	"github.com/nateschererforks/vals/pkg/providers/gkms"
	"github.com/nateschererforks/vals/pkg/providers/hcpvaultsecrets"
	"github.com/nateschererforks/vals/pkg/providers/httpjson"
	"github.com/nateschererforks/vals/pkg/providers/k8s"
	"github.com/nateschererforks/vals/pkg/providers/onepassword"
	"github.com/nateschererforks/vals/pkg/providers/onepasswordconnect"
	"github.com/nateschererforks/vals/pkg/providers/pulumi"
	"github.com/nateschererforks/vals/pkg/providers/s3"
	"github.com/nateschererforks/vals/pkg/providers/sops"
	"github.com/nateschererforks/vals/pkg/providers/ssm"
	"github.com/nateschererforks/vals/pkg/providers/tfstate"
	"github.com/nateschererforks/vals/pkg/providers/vault"
)

func New(l *log.Logger, provider api.StaticConfig) (api.LazyLoadedStringProvider, error) {
	tpe := provider.String("name")

	switch tpe {
	case "s3":
		return s3.New(l, provider), nil
	case "gcs":
		return gcs.New(provider), nil
	case "ssm":
		return ssm.New(l, provider), nil
	case "vault":
		return vault.New(l, provider), nil
	case "awskms":
		return awskms.New(provider), nil
	case "awssecrets":
		return awssecrets.New(l, provider), nil
	case "sops":
		return sops.New(l, provider), nil
	case "gcpsecrets":
		return gcpsecrets.New(provider), nil
	case "tfstate":
		return tfstate.New(provider, ""), nil
	case "tfstategs":
		return tfstate.New(provider, "gs"), nil
	case "tfstates3":
		return tfstate.New(provider, "s3"), nil
	case "tfstateazurerm":
		return tfstate.New(provider, "azurerm"), nil
	case "tfstateremote":
		return tfstate.New(provider, "remote"), nil
	case "azurekeyvault":
		return azurekeyvault.New(provider), nil
	case "gitlab":
		return gitlab.New(provider), nil
	case "onepassword":
		return onepassword.New(provider), nil
	case "onepasswordconnect":
		return onepasswordconnect.New(provider), nil
	case "doppler":
		return doppler.New(l, provider), nil
	case "pulumistateapi":
		return pulumi.New(l, provider, "pulumistateapi"), nil
	case "gkms":
		return gkms.New(l, provider), nil
	case "k8s":
		return k8s.New(l, provider)
	case "conjur":
		return conjur.New(l, provider), nil
	case "hcpvaultsecrets":
		return hcpvaultsecrets.New(l, provider), nil
	case "httpjson":
		return httpjson.New(l, provider), nil
	}

	return nil, fmt.Errorf("failed initializing string provider from config: %v", provider)
}
