package game_proto

import (
	"game_server/src/model"
)

const (
  INIT_PROTO = iota

  C2S_PLAYER_LOGIN_PROTO
  S2C_PLAYER_LOGIN_PROTO

  C2S_CHOOSE_ROOM_PROTO
  S2C_CHOOSE_ROOM_PROTO
)


// 功能结构

type HeaderProto struct {
  Proto int
  Proto2 int
}

type C2S_PlayerLogin struct {
  HeaderProto
  Code string // 微信授权code
}

type S2C_PlayerLogin struct {
  HeaderProto
  PlayerData *model.PlayerSt
}
