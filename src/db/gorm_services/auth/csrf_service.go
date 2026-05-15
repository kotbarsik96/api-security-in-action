package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type CsrfService struct {
	Secret []byte
}

func NewCsrfService(secret []byte) *CsrfService {
	return &CsrfService{
		Secret: secret,
	}
}

func (s *CsrfService) GenerateToken(id string) string {
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write([]byte(id))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *CsrfService) CompareToken(id string, token string) bool {
	idDecoded, err := base64.URLEncoding.DecodeString(s.GenerateToken(id))
	if err != nil {
		return false
	}

	tokenDecoded, err := base64.URLEncoding.DecodeString(token)

	return hmac.Equal(idDecoded, tokenDecoded)
}
