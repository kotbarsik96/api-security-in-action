package auth

import (
	"crypto/hmac"
	"crypto/sha256"
)

type CsrfService struct {
	Secret []byte
}

func NewCsrfService(secret []byte) *CsrfService {
	return &CsrfService{
		Secret: secret,
	}
}

func (s *CsrfService) GenerateToken(sessID []byte) []byte {
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write(sessID)
	return mac.Sum(nil)
}

func (s *CsrfService) CompareToken(sessionID []byte, token []byte) bool {
	expected := s.GenerateToken(sessionID)
	return hmac.Equal(expected, token)
}

func (s *CsrfService) GetCsrfProtectedMethods() []string {
	return []string{"POST", "PUT", "DELETE"}
}
