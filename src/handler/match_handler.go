package handler

import "game_server/src/frame"


var MatchRouter *frame.RouterGroup


func MatchHandler(ctx *frame.Context)  {
  
}


func init()  {
  MatchRouter = frame.NewRouterGroup("/match")
  MatchRouter.Get("/normal", nil,MatchHandler)
}
