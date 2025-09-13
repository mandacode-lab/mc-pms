package encoder

import (
	"encoding/base64"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type Base64Encoder struct{}

func NewBase64Encoder() out.Encoder {
	return &Base64Encoder{}
}

func (e *Base64Encoder) Encode(data []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(data)
	return []byte(encoded)
}

type Base64URLEncoder struct{}

func NewBase64URLEncoder() out.Encoder {
	return &Base64URLEncoder{}
}

func (e *Base64URLEncoder) Encode(data []byte) []byte {
	encoded := base64.URLEncoding.EncodeToString(data)
	return []byte(encoded)
}