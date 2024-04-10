package handler

import (
	"game_server/src/frame"
	"game_server/src/service/user_service"
	"game_server/src/utils/email_utils"
	"game_server/src/utils/session_utils"
	"game_server/src/utils/user_utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


var UserRouter *frame.RouterGroup


func LoginHandler(ctx *frame.Context)  {
  var data frame.M
  if err := ctx.Json(&data);err != nil{
    ctx.ResponseJson(frame.M{
      "code": 20000,
      "success": false,
      "data": frame.M{
        "err": err,
      },
    })
    return
  }
  username := data["username"].(string)
  password := data["password"].(string)
  userInfo, err := user_service.Login(username, password)
  if err != nil {
    ctx.ResponseJson(frame.M{
      "code": 20000,
      "success": false,
      "data": frame.M{
        "err": err,
      },
    })
    return
  }
  ctx.ResponseJson(frame.M{
    "code": 20000,
    "success": true,
    "data": userInfo,
  })
}


func RegisterHandler(ctx *frame.Context)  {
  var data frame.M
  if err := ctx.Json(&data); err != nil {
    ctx.ErrorJson(frame.M{
      "code": 200,
      "success": false,
      "data": frame.M{
        "msg": err,
      },
    }, http.StatusInternalServerError)
  }
  captchaData := strings.Split(ctx.SessionData["captcha"], ",")
  captchaExp, err := strconv.ParseInt(captchaData[0], 10, 64)
  if err != nil {
    frame.Logger.Errorln(err)
      ctx.ErrorJson(frame.M{
        "code": 20000,
        "success": false,
        "data": frame.M{
          "msg": "网络问题，请重试",
        },
      }, http.StatusInternalServerError)

  }
  if time.Now().Unix() > captchaExp {
    ctx.ErrorJson(frame.M{
      "code": 20000,
      "success": false,
      "data": frame.M{
        "msg": "验证码已过期",
      },
    }, http.StatusNetworkAuthenticationRequired)
  }

  trueCaptcha := captchaData[1]

  username := data["username"].(string)
  password := data["password"].(string)
  captcha := data["captcha"].(string)
  rawGender := data["gender"].(float64)
  email := data["email"].(string)

  gender := int(rawGender)

  frame.Logger.Debugln(captcha, " ", trueCaptcha)

  if captcha == trueCaptcha {
    if err := user_service.Register(
      username,
      password,
      gender,
      email,
    ); err != nil {
      log.Println(err)
      ctx.ErrorJson(frame.M{
        "code": 20000,
        "success": false,
        "data": frame.M{
          "msg": "网络问题，请重试",
        },
      }, http.StatusInternalServerError)
    }
    ctx.ResponseJson(frame.M{
      "code": 20000,
      "success": true,
      "data": frame.M{
        "msg": "注册成功！",
      },
    })
  } else {
    ctx.ErrorJson(frame.M{
      "code": 20000,
      "success": false,
      "data": frame.M{
        "msg": "验证码错误",
      },
    }, http.StatusNonAuthoritativeInfo)
  }
}


func CaptchaHandler(ctx *frame.Context)  {
  captcha := user_utils.GenerateCaptcha()
  exp := time.Now().Add(5 * time.Minute)
  sessionId := ctx.SessionData["sessionId"]
  data := fmt.Sprintf("%s,%d", captcha, exp.Unix())
  if err := session_utils.SetSessionData(sessionId, "captcha", data, 10 * time.Second, 0); err != nil {
    frame.Logger.Errorln(err)
    ctx.ErrorJson(frame.M{
      "code": 20000,
      "success": false,
      "msg": "生成验证码失败, 请稍后重试",
    }, http.StatusInternalServerError)
    return
  }
  reciver := "1369324180@qq.com"
  subject := "这是你的验证码"
  msg := "这是你的验证码：" + captcha
  if err := email_utils.EmailSender.Send(reciver, subject, msg); err != nil {
    frame.Logger.Errorln("Error Sending Captcha: ", err)
  }

  ctx.ResponseJson(frame.M{
    "code": 20000,
    "success": true,
  })
}


func init()  {
  UserRouter = *frame.NewRouterGroup("/user")
  UserRouter.Post("/login", nil, LoginHandler)
  UserRouter.Get("/captcha", nil, CaptchaHandler)
  UserRouter.Post("/register", nil, RegisterHandler)
}
