package cryptionservice

import "golang.org/x/crypto/bcrypt"

type CryptionService interface {
	BcryptEncode(in string) (string, error)
	BcryptCheck(hash string, in string) bool
}

type CryptionServiceImpl struct {
}

func NewCryptionService() CryptionService {
	return &CryptionServiceImpl{}
}

func (s *CryptionServiceImpl) BcryptEncode(
	in string,
) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

// check hash and decode input
func (s *CryptionServiceImpl) BcryptCheck(
	hash string,
	in string,
) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(in))
	return err == nil
}
