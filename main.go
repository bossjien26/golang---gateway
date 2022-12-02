package main

import (
	"apigateway/authentication"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Kong/go-pdk"
	"github.com/joho/godotenv"
)

// var redisConfig = &authentication.RedisConfig{
// }

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// func main() {
// 	user := authentication.User{
// 		ID:   1,
// 		Name: "test",
// 		Key:  "test",
// 	}
// 	token := authentication.GenerateToken(user)
// 	fmt.Println(token)
// 	// // get value
// 	authentication.SetVal(user.Key, token)
// 	authentication.GetVal(user.Key)
// 	id, name, key := authentication.ValidateToken(token)
// 	fmt.Println(id, name, key)
// }

type Config struct {
	Apikey string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	// serviceKey, err := kong.Request.GetQueryArg("serviceKey")
	serviceKey, serviceKeyErr := kong.Request.GetHeader("serviceKey")
	publicKey, publicKeyErr := kong.Request.GetHeader("publicKey")
	if serviceKeyErr != nil || publicKeyErr != nil {
		kong.Log.Err(serviceKeyErr.Error())
		kong.Log.Err(publicKeyErr.Error())
	}

	// println(kong.Request.GetHeader("key"))
	id, name, redisKey := authentication.ValidateToken(publicKey)
	fmt.Println(id, name, redisKey)

	apiKey := conf.Apikey

	response := make([]string, 1)
	response[0] = "You have no correct key"

	j, _ := json.Marshal(response)
	x := make(map[string][]string)

	x["Content-Type"] = append(x["Content-Type"], "application/json")
	if apiKey != serviceKey || id == 0 || authentication.GetVal(redisKey) == "" {
		kong.Response.Exit(403, string(j), x)
	}
}
