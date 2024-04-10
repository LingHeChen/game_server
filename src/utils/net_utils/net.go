package net_utils

import (
	"game_server/src/frame"
	"game_server/src/protocol"
	"game_server/src/protocol/game_proto"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ws 网络数据结构
// Conn ws链接
// StrMD5 房间加密信息
type NetDataConn struct {
  Conn *websocket.Conn
  StrMD5 string
}

type MessageData struct {
  msgType int
  raw []byte
  data *sync.Map
}

func (this *MessageData) LoadJson() error {
  if err := json.Unmarshal(this.raw, &this.data); err != nil {
    return err
  }
  return nil
}

func (this *MessageData) ToString() string {
  return string(this.raw)
}

func (this *MessageData) GetData() *sync.Map {
  return this.data
}

func (this *MessageData) GetMsgType() int {
  return this.msgType
}

func (this *MessageData) GetRaw() []byte {
  return this.raw
}

func (this *NetDataConn) PullFromClient()  {
  for {
    messageType, content, err := this.Conn.ReadMessage()
    if err != nil {
      break
    }
    if len(content) == 0 {
      break
    }
    switch messageType {
      case websocket.TextMessage:
        msg := MessageData{
          msgType: messageType,
          raw: content,
        }
        go this.SyncMessageFunc(&msg)
      case websocket.BinaryMessage:
      case websocket.CloseMessage:
      case websocket.PingMessage:
      case websocket.PongMessage:
      default:
    }
  }
  return
}

func (this *NetDataConn) SyncMessageFunc(msg *MessageData) {
  fmt.Println(msg.msgType, msg.ToString())
  if err := msg.LoadJson(); err != nil {
    frame.Logger.Errorln("解析数据时出现问题", err)
  }
  data := msg.GetData()
  proto, ok := data.Load("Proto")
  if !ok {
    frame.Logger.Errorln(errors.New("未能解析主协议"))
    return
  }
  proto2, ok := data.Load("Proto2")
  if !ok {
    frame.Logger.Errorln(errors.New("未能解析子协议"))
    return
  }
  this.HandleCltProtocol(proto.(int), proto2.(int), data)
}

func (
  this *NetDataConn,
) HandleCltProtocol(
  proto int,
  proto2 int,
  protocolData *sync.Map,
) error {
  switch proto {
  case protocol.GAME_DATA_PROTO:
    switch proto2 {
    case game_proto.C2S_PLAYER_LOGIN_PROTO:
    case game_proto.S2C_PLAYER_LOGIN_PROTO:
    case game_proto.C2S_CHOOSE_ROOM_PROTO:
    case game_proto.S2C_CHOOSE_ROOM_PROTO:
    default:
      panic(errors.New(fmt.Sprintf("未知的子协议: %d", proto2)))
    }
  case protocol.GAME_DB_PROTO:
  default:
    panic(errors.New(fmt.Sprintf("未知的主协议: %d", proto)))
  }
  return nil
}


func ApplyMiddlewareBefore(
  next func ( w http.ResponseWriter, r *http.Request, ),
  middleware func( w http.ResponseWriter, r *http.Request, ),
) func( w http.ResponseWriter, r *http.Request, ) {
  return func (w http.ResponseWriter, r *http.Request)  {
    middleware(w, r)
    next(w, r)
  }
}

func ApplyMiddlewareAfter(
  before func ( w http.ResponseWriter, r *http.Request, ),
  middleware func( w http.ResponseWriter, r *http.Request, ),
) func ( w http.ResponseWriter, r *http.Request, ) {
  return func (w http.ResponseWriter, r *http.Request) {
    before(w, r)
    middleware(w, r)
  }
}
