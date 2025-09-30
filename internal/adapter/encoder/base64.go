package encoder

import (
	"encoding/base64"
	"errors"

	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Base64Encoder struct{}

func NewBase64Encoder() out.Encoder {
	return &Base64Encoder{}
}

func (e *Base64Encoder) Encode(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return []byte(encoded), nil
}

type Base64URLEncoder struct{}

func NewBase64URLEncoder() out.Encoder {
	return &Base64URLEncoder{}
}

func (e *Base64URLEncoder) Encode(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	encoded := base64.URLEncoding.EncodeToString(data)
	return []byte(encoded), nil
}
