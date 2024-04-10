package main

import (
	"game_server/src/frame"
	"game_server/src/handler"
	"game_server/src/middlewares"
	"net/http"
	"runtime"
)

// 游戏服务器的初始化工作
func init()  {
  
}

func main()  {
  // fmt.Println("本机最大核心数: ", runtime.NumCPU())
  runtime.GOMAXPROCS(runtime.NumCPU() * 2)
  frame.Logger.InitLogger("GameTest", frame.DEBUG)

  serverConfig := frame.ServerConfig{}
  server := serverConfig.SpawnServer()

  group1 := server.NewRouterGroup("/group1")

  server.BeforeHandle(func (w http.ResponseWriter, r *http.Request)  {
    frame.Logger.Infof("%s -> %s", r.Method, r.URL.Path)
  })

  server.BeforeInitial(func (s *frame.Server) {
    frame.Logger.Debugln(s.BaseRouter.Children["/user"])
    return
  })

  server.IncludeRouter(handler.UserRouter)
  server.Apply(middlewares.CheckSession)

  server.Get("/hello", nil, handler.HelloHandler)
  group1.Get("/hello", nil, handler.HelloHandler)

  group2 := frame.NewRouterGroup("/group2")
  group2.Get("/hello", nil, handler.HelloHandler)
  group1.IncludeRouter(group2)
  server.Run()
}
