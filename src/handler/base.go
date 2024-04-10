package handler

import (
	"class1/src/frame"
	"class1/src/utils/net_utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestWSHandler(w http.ResponseWriter, r *http.Request)  {
  params := r.URL.Query()
  data := params.Get("data")
  fmt.Println("data: ", data)
  conn, err := upgrader.Upgrade(w, r, nil)
  defer conn.Close()
  if err != nil {
    log.Fatal(err)
  }
  conn.WriteMessage(websocket.TextMessage, []byte("Hello world"))
  netDataConnmp := &net_utils.NetDataConn{
    Conn: conn,
    StrMD5: "",
  }
  netDataConnmp.PullFromClient()
}

func HelloHandler(ctx *frame.Context)  {
  ctx.ResponseJson(map[string]any{
    "res": "Hello World",
  })
}
