package model


type UserSt struct {
  Id string `bson:"id"`
  Username string `bson:"username"`
  HashedPassword []byte `bson:"hashPassword"`
  Salt []byte `bson:"salt"`
  Gender int `bson:"gender"`
  Email string `bson:"email"`
}


type UserLogin struct {
  Username string
  Password string
}


type UserRegister struct {
  Username string
  Password string
  Captcha string
}


type UserInfoSt struct {
  Username string `json:"username"`
  Gender int `json:"gender"`
  Email string `json:"email"`
}




func init() {

}
