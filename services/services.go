package services

import (
	"errors"

	"github.com/lum8rjack/redcompass/services/namecheap"
	"github.com/lum8rjack/redcompass/services/porkbun"
	"github.com/lum8rjack/redcompass/services/types"
)

func NewService(provider string, settings string) (types.Service, error) {

	switch provider {
	case "Namecheap":
		return namecheap.NewClient(settings)
	case "Porkbun":
		return porkbun.NewClient(settings)
	case "Cloudflare":
		return nil, errors.New("service provider Cloudflare is not supported yet")
	}

	return nil, errors.New("invalid provider")
}
