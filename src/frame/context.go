package frame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Context struct {
  W http.ResponseWriter
  R *http.Request
  Config *HandlerConfig
  SessionData HashResult
}


func (ctx *Context) ResponseJson(data M)  {
  jsonBytes, err := json.Marshal(data)
  if err != nil {
    http.Error(ctx.W, "json 序列化失败", http.StatusInternalServerError)
  }
  jsonString := string(jsonBytes)
  fmt.Fprintf(ctx.W, jsonString)
}


func (ctx *Context) URL() *url.URL{
  return ctx.R.URL
}


func (ctx *Context) Query() url.Values {
  return ctx.R.URL.Query()
}


func (ctx *Context) Json(data *M) error {
  bytesData, err := io.ReadAll(ctx.R.Body)
  if err != nil {
    return err
  }
  if err := json.Unmarshal([]byte(bytesData), data); err != nil {
    return err
  }
  return nil
}


func (ctx *Context) ParseBody()  {
  return 
}

func (ctx *Context) ErrorJson(data map[string]any, status int) {
  jsonBytes, err := json.Marshal(data)
  if err != nil {
    ctx.W.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(ctx.W, "json 序列化失败")
    return
  }
  jsonString := string(jsonBytes)
  ctx.W.WriteHeader(status)
  fmt.Fprintf(ctx.W, jsonString)
}
