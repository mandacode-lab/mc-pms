// Package domain
package shared

import (
	"errors"
)

type Provider string

const (
	ProviderUnknown   Provider = ""
	ProviderGoogle    Provider = "google"
	ProviderApple     Provider = "apple"
	ProviderNaver     Provider = "naver"
	ProviderKakao     Provider = "kakao"
	ProviderMandacode Provider = "mandacode"
)

func NewProvider(provider string) (Provider, error) {
	switch provider {
	case string(ProviderUnknown):
		return ProviderUnknown, nil
	case string(ProviderGoogle):
		return ProviderGoogle, nil
	case string(ProviderApple):
		return ProviderApple, nil
	case string(ProviderNaver):
		return ProviderNaver, nil
	case string(ProviderKakao):
		return ProviderKakao, nil
	case string(ProviderMandacode):
		return ProviderMandacode, nil
	default:
		return "", errors.New("invalid provider: " + provider)
	}
}

func (p Provider) IsValid() bool {
	switch p {
	case ProviderGoogle, ProviderApple, ProviderNaver, ProviderKakao, ProviderMandacode:
		return true
	default:
		return false
	}
}

func (Provider) Values() (providers []string) {
	for _, v := range []Provider{ProviderGoogle, ProviderApple, ProviderNaver, ProviderKakao, ProviderMandacode} {
		providers = append(providers, string(v))
	}
	return
}
