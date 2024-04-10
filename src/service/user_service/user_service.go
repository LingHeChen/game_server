package user_service

import (
	"bytes"
	"game_server/src/model"
	"game_server/src/utils/encode_uitls"
	"class1/src/utils/mongo_utils"
	"errors"
)



var UserQuery mongo_utils.Query[model.UserSt]


func Login(username string, password string) (*model.UserInfoSt, error) {
  target := model.UserSt{
    Username: username,
  }

  may, err := UserQuery.FindOne(target)
  if err != nil {
    return nil, err
  }

  try, err := encode_uitls.EncodePwd(password, may.Salt)
  if err != nil {
    return nil, err
  }

  if !bytes.Equal(try, may.HashedPassword) {
    return nil, errors.New("Wrong Password")
  }

  userInfo := &model.UserInfoSt{
    Username: may.Username,
    Gender: may.Gender,
    Email: may.Email,
  }

  return userInfo, nil
}


func Register(
  username string,
  password string,
  gender int,
  email string,
) error {
  salt := make([]byte, 16)
  hashedPassword, err := encode_uitls.EncodePwd(password, salt)
  if err != nil {
    return err
  }

  newUser := model.UserSt{
    Username: username,
    HashedPassword: hashedPassword,
    Salt: salt,
    Gender: gender,
    Email: email,
  }

  if err := UserQuery.InsertOne(newUser); err != nil {
    return err
  }

  return nil
}


func init()  {
  UserQuery = mongo_utils.Query[model.UserSt]{
    DbName: "game",
    CollectionName: "users",
  }
}
