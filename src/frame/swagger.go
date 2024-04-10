package frame

import (
	"reflect"
)


var swaggerConfig SwaggerConfig


type SwaggerConfig struct {
  Openapi string `json:"openapi"`
  Info string `json:"info"`
  Paths map[string]PathConfig `json:"paths"`
}


type ArrayItems struct {
  Type string `json:"type"`
}


type JsonSchema struct {
  Type string `json:"type"`
  Required []string `json:"required"`
  Properties map[string]JsonSchema `json:"properties"`
  Items ArrayItems `json:"items"`
  AdditionalItems ArrayItems `json:"addtionalItems"`
  Example map[string]any `json:"example"`
}


type ParamConfig struct {
  Name string `json:"name"`
  In string `json:"in"`
  Description string `json:"description"`
  Required bool `json:"required"`
  Schema map[string]JsonSchema `json:"schema"`
  Deprecated bool `json:"deprecated"`
  AllowEmptyValue bool `json:"allowEmptyValue"`
}


type ExampleConfig struct {
  Summary string `json:"summary"`
  Description string `json:"description"`
  Value any `json:"value"`
  ExternalValue string `json:"externalValue"`
}


type HeaderConfig struct {
  Description string `json:"description"`
  Schema struct { Type string } `json:"schema"`
}


type EncodingConfig struct {
  ContentType string `json:"contentType"`
  Headers map[string]HeaderConfig `json:"headers"`
  Style string `json:"style"`
  Explode bool `json:"explode"`
  AllowReserved bool `json:"allowReserved"`
}


type MediaType struct {
  Schema JsonSchema `json:"schema"`
  Example any `json:"example"`
  Examples []ExampleConfig `json:"examples"`
  Encoding EncodingConfig `json:"encoding"`
}


type BodyConfig struct {
  Decription string `json:"description"`
  Content map[string]MediaType `json:"content"`
  Required bool `json:"required"`
}


type ResponseConfig struct {
  Description string `json:"description"`
  Headers map[string]HeaderConfig `json:"headers"`
  Content map[string]MediaType `json:"content"`
}


type HandlerConfig struct {
  path string
  method string
  Summary string `json:"summary"`
  Data any
  Deprecated bool `json:"deprecated"`
}

type PathConfig struct {
  Summary string `json:"summary"`
  Params []ParamConfig `json:"parameters"`
  Body BodyConfig `json:"requestBody"`
  Reponse ResponseConfig `json:"response"`
  Deprecated bool `json:"deprecated"`
}

func (h *HandlerConfig) GenerateConfig()  {
  dataTpye := reflect.TypeOf(h.Data)
  dataKind := dataTpye.Kind()
  if dataKind == reflect.Struct {

  } else if dataKind == reflect.Map {

  }
  swaggerConfig.Paths[h.path] = PathConfig {

  }
}

func init()  {
  swaggerConfig = SwaggerConfig{}
}
