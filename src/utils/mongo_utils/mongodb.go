package mongo_utils

import (
	"game_server/src/frame"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var MongoClient *mongo.Client

type Query[TModel any] struct {
  DbName string
  CollectionName string
}


func (q *Query[TModel]) FindOne(target TModel) (TModel, error) {
  var result TModel
  db := MongoClient.Database(q.DbName)
  collection := db.Collection(q.CollectionName)
  if err := collection.FindOne(context.Background(), target).Decode(result); err != nil {
    return *new(TModel), err
  }
  return result, nil
}


func (q *Query[TModel]) InsertOne(newModel TModel) error {
  db := MongoClient.Database(q.DbName)
  collection := db.Collection(q.CollectionName)
  if _, err := collection.InsertOne(context.Background(), newModel); err != nil {
    return err
  }
  return nil
}


func init()  {
  paramsStr := strings.Join(frame.Config.Mongo.Params, "&")
  mongoUri := fmt.Sprintf(
    "mongodb://%s:%s@%s:%d/?%s",
    frame.Config.Mongo.UserName,
    frame.Config.Mongo.Password,
    frame.Config.Mongo.Host,
    frame.Config.Mongo.Port,
    paramsStr,
  )
  // log.Println(mongoUri)
  clientOptions := options.Client().ApplyURI(mongoUri)
  var err error
  MongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		frame.Logger.Errorln(err)
	}
}
