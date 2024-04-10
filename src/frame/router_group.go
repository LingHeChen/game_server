package frame

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)


const (
  IN_QUERY = "query"
  IN_HEADER = "header"
  IN_PATH = "path"
  IN_COOKIES = "cookies"

  TYPE_STRING = "string"
  TYPE_NUMBER = "number"
  TYPE_INTEGER = "integer"
  TYPE_BOOLEAN = "boolean"
  TYPE_ARRAY = "array"
  TYPE_OBJECT = "null"
  TYPE_ANY = "any"
)


type RouterGroup struct {
  level int
  fullPath string
  basePath string
  middlewares []*func(*Context) bool
  Group *http.ServeMux
  Children map[string]*RouterGroup
  handlers map[string]*Handler
}


func NewRouterGroup(basePath string) *RouterGroup {
  level := ROUTER_DEFAULT_LEVEL
  return &RouterGroup {
    level: level,
    fullPath: basePath,
    middlewares: make([]*func(*Context) bool, 0),
    basePath: basePath,
    handlers: make(map[string]*Handler),
    Group: http.NewServeMux(),
    Children: make(map[string]*RouterGroup),
  }
}


func (rg *RouterGroup) NewRouterGroup(basePath string) *RouterGroup {
  group := NewRouterGroup(basePath)

  var fullPath string
  if rg.fullPath == ROOT_PATH {
    fullPath = basePath
  } else {
    fullPath = rg.fullPath + basePath
  }
  group.fullPath = fullPath
  group.level = rg.level + 1
  if rg.Children == nil {
    rg.Children = make(map[string]*RouterGroup)
  }

  rg.Children[basePath] = group

  return group
}


func ComposePath(paths ...string) string {
  if paths[0] == ROOT_PATH {
    paths = paths[1:]
  }

  return strings.Join(paths, "")
}


func (rg *RouterGroup) IncludeRouter(group *RouterGroup)  {
  group.level = rg.level + 1
  group.fullPath = ComposePath(rg.fullPath, group.basePath)
  group.middlewares = rg.middlewares
  if rg.Children == nil {
    rg.Children = make(map[string]*RouterGroup)
  }
  rg.Children[group.basePath] = group
}


func (rg *RouterGroup) Apply(middleware func (*Context) bool) {
  for _, child := range rg.Children {
    child.Apply(middleware)
  }
  rg.middlewares = append(rg.middlewares, &middleware)
}


func (rg *RouterGroup) GetLevelPath(url string, bias int) string {
  level := rg.level + bias
  if level < 1 {
    return ROOT_PATH
  }
  parts := strings.Split(url, PATH_SEP)
  return PATH_SEP + parts[level]
}


func (rg *RouterGroup) Handle(w http.ResponseWriter, r *http.Request) error {
  levelPath := rg.GetLevelPath(r.URL.Path, 1)
  group := rg.Children[levelPath]
  if group != nil {
    group.Handle(w, r)
    return nil
  }
  Logger.Debugln(r.URL.Path, levelPath, group)
  if rg.Group != nil {
    rg.Group.ServeHTTP(w, r)
    return nil
  }
  return errors.New(fmt.Sprintf("Unknown Address: %s", r.URL.Path))
}


func (rg *RouterGroup) Serve()  {
  if len(rg.Children) > 0{
    for _, group := range rg.Children {
      group.Serve()
    }
  }
  if len(rg.handlers) > 0 {
    for path, handler := range rg.handlers {
      path = ComposePath(rg.fullPath, path)
      rg.Group.HandleFunc(path, func (w http.ResponseWriter, r *http.Request) {
        con := &Context{
          W: w,
          R: r,
        }
        Logger.Infof("[%s] (%s)%s -> (%s)%s\n", r.RemoteAddr, r.Method, r.RequestURI, handler.method, path)
        if r.Method != handler.method {
          con.ErrorJson(M{
            "msg": "Method Not Allowed",
          }, http.StatusMethodNotAllowed)
          return
        }
        for _, callable := range rg.middlewares {
          if callable != nil {
            if ok := (*callable)(con); !ok {
              return
            }
          }
        }
        (*handler.callable)(con)
      })
    }
  }
}


func (rg *RouterGroup) Response(method string, path string, config *HandlerConfig, handler func (*Context))  {
  rg.handlers[path] = &Handler {
    method: method,
    callable: &handler,
    config: config,
  }
}

func (rg *RouterGroup) Get(path string, config *HandlerConfig, handler func (*Context)) {
  rg.Response(GET_METHOD, path, config, handler)
}


func (rg *RouterGroup) Post(path string, config *HandlerConfig, handler func (*Context))  {
  rg.Response(POST_METHOD, path, config, handler)
}


func (rg *RouterGroup) Put(path string, config *HandlerConfig, handler func (*Context))  {
  rg.Response(PUT_METHOD, path, config, handler)
}


func (rg *RouterGroup) Delete(path string, config *HandlerConfig, handler func (*Context))  {
  rg.Response(DELETE_METHOD, path, config, handler)
}
