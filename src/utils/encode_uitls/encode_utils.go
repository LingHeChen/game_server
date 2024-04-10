package encode_uitls

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/xdg-go/pbkdf2"
)


func EncodePwd(password string, salt []byte) ([]byte, error) {
  _, err := rand.Read(salt)
  if err != nil {
      return nil, err
  }

  hashedPassword := pbkdf2.Key([]byte(password), salt, 10000, 128, sha256.New)

  return hashedPassword, nil
}
