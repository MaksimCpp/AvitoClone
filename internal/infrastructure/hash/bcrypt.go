package hash

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (hasher *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

func (hasher *BcryptHasher) Compare(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
