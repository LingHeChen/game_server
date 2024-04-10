package frame

import (
	"fmt"
	"net/http"
)


const (
  GET_METHOD = "GET"
  POST_METHOD = "POST"
  PUT_METHOD = "PUT"
  DELETE_METHOD = "DELETE"
  PATH_SEP = "/"
  ROOT_PATH = "/"
  ROUTER_DEFAULT_LEVEL = 0
)


type Server struct {
  Addr string
  Server *http.Server
  BaseRouter *RouterGroup
  beforeHandle *func(http.ResponseWriter, *http.Request)
  beforeInitail *func(*Server)
}


type ServerConfig struct {
  Host string
  Port int16
}


type Handler struct {
  method string
  callable *func(*Context)
  config *HandlerConfig
}


func (c *ServerConfig) SpawnServer() Server {
  server := Server{}
  host := "127.0.0.1"
  if len(c.Host) != 0 {
    host = c.Host
  }
  var port int16 = 8888
  if c.Port != 0 {
    port = c.Port
  }
  addr := fmt.Sprintf("%s:%d", host, port)

  server.Addr = addr
  server.Server = &http.Server{
    Addr: server.Addr,
  }
  server.BaseRouter = NewRouterGroup(ROOT_PATH)
  return server
}

func (s *Server) Run() error {
  if s.beforeInitail != nil {
    (*s.beforeInitail)(s)
  }

  s.BaseRouter.Serve()

  s.Server.Handler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
    if s.beforeHandle != nil {
      (*s.beforeHandle)(w, r)
    }
    if err := s.BaseRouter.Handle(w, r); err != nil {
      http.Error(w, err.Error(), http.StatusNotFound)
      return
    }
  })
  if err := s.Server.ListenAndServe(); err != nil {
    return err
  }

  return nil
}


func (s *Server) IncludeRouter(group *RouterGroup) {
  s.BaseRouter.IncludeRouter(group)
}


func (s *Server) BeforeHandle(callable func (http.ResponseWriter, *http.Request))  {
  s.beforeHandle = &callable
}


func (s *Server) BeforeInitial(callable func (*Server))  {
  s.beforeInitail = &callable
}


func (s *Server) NewRouterGroup(basePath string) *RouterGroup {
  return s.BaseRouter.NewRouterGroup(basePath)
}

func (s *Server) Apply(middleware func (*Context) bool) {
  s.BaseRouter.Apply(middleware)
}


func (s *Server) Response(method string, path string, config *HandlerConfig, handler func (*Context))  {
  s.BaseRouter.Response(method, path, config, handler)
}

func (s *Server) Get(path string, config *HandlerConfig, handler func (*Context)) {
  s.BaseRouter.Response(GET_METHOD, path, config, handler)
}


func (s *Server) Post(path string, config *HandlerConfig, handler func (*Context))  {
  s.BaseRouter.Response(POST_METHOD, path, config, handler)
}


func (s *Server) Put(path string, config *HandlerConfig, handler func (*Context))  {
  s.BaseRouter.Response(PUT_METHOD, path, config, handler)
}


func (s *Server) Delete(path string, config *HandlerConfig, handler func (*Context))  {
  s.BaseRouter.Response(DELETE_METHOD, path, config, handler)
}

