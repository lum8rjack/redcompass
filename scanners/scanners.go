package scanners

import (
	"errors"

	"github.com/lum8rjack/redcompass/scanners/types"
	"github.com/lum8rjack/redcompass/scanners/virustotal"
)

func NewScanner(provider string, settings string) (types.Scanner, error) {

	if provider == "VirusTotal" {
		return virustotal.NewClient(settings)
	}

	return nil, errors.New("invalid scanner")
}
