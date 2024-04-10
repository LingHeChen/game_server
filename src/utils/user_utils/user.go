package user_utils

import (
	"fmt"
	"math/rand"
)


func GenerateCaptcha() string {
	min := 100000
	max := 999999
	randomNumber := min + rand.Intn(max-min+1)
  return fmt.Sprint(randomNumber)
}
