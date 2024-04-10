package middlewares

import (
	"game_server/src/frame"
	"game_server/src/utils/session_utils"
	"log"
	"net/http"
)


func CheckSession(ctx *frame.Context) bool {
  sessionId, err := ctx.R.Cookie("session-id")
  if err != nil || sessionId == nil || sessionId.Value == "" {
    ctx.ErrorJson(frame.M{
      "code": 302,
      "success": false,
      "msg": "会话已过期！",
    }, http.StatusTemporaryRedirect)
    return false
  }
  sessionData, err := session_utils.GetSessionData(sessionId.Value)
  if err != nil {
    log.Println(err)
    ctx.ErrorJson(frame.M{
      "code": 20000,
      "success": false,
      "msg": err.Error(),
    }, http.StatusInternalServerError)
    return false
  }

  sessionData["sessionId"] = sessionId.Value

  ctx.SessionData = sessionData
  return true
}
