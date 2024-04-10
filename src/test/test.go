package main

import (
	"game_server/src/utils/encode_uitls"
	"fmt"
)

func main() {
  encoded_password, salt := encode_uitls.EncodePwd("123456")
  fmt.Println(encoded_password, salt)
}
