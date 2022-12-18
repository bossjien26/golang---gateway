package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt"
)

var _redisConfig = &RedisConfig{}

func ConnectRedisClient() {
	ctx := context.Background()
	client := NewClient(ctx, "redis-pod-service:6379")
	_redisConfig = &RedisConfig{
		ctx:    ctx,
		client: client,
	}
	_redisConfig.ctx = ctx
	_redisConfig.client = client
}

func GetVal(key string) string {
	// ConnectRedisClient()
	if _redisConfig == nil {
		return ""
	}

	val, err := _redisConfig.client.Get(_redisConfig.ctx, key).Result()
	if err != nil {
		// fmt.Println("client.Get failed", err)
		return ""
	}
	return val
}

type RedisConfig struct {
	ctx    context.Context
	client *redis.Client
}

type User struct {
	ID   uint
	Name string
	Key  string
}

type AuthClaims struct {
	jwt.StandardClaims
	UserID uint `json:"userId"`
	Key    string
	Name   string
}

type Config struct {
	WaitTime int
}

type RequestHeader struct {
	serviceKey string

	publicKey string
}

func New() interface{} {
	return &Config{}
}

var ctx = context.Background()

// var key = []byte("test")

func NewClient(ctx context.Context, addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	client.Ping(ctx).Result()
	// if err != nil {
	// 	// panic(err)
	// }

	// fmt.Println(pong)
	return client
}

//	func SetVal(key string, val string) bool {
//		ConnectRedisClient()
//		if err := redisConfig.client.Set(redisConfig.ctx, key, val, 0).Err(); err != nil {
//			fmt.Println("client.Set failed", err)
//			return false
//		}
//		return true
//	}
var _requests = make(map[string]time.Time)
var _key = []byte("test")

func ValidateToken(tokenString string) (uint, string, string) {
	claims := AuthClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return _key, nil
	})

	if !token.Valid {
		return 0, "", ""
	}
	id := claims.UserID
	name := claims.Name
	key := claims.Key
	return id, name, key
}

func (config Config) Access(kong *pdk.PDK) {
	requestHeader := &RequestHeader{}
	// serviceKey, err := kong.Request.GetQueryArg("serviceKey")
	// _ = kong.Response.SetHeader("x-wait-time", fmt.Sprintf("%d seconds", config.WaitTime))

	// host, _ := kong.Request.GetHost()
	// lastRequest, exists := requests[host]

	// if exists && time.Now().Sub(lastRequest) < time.Duration(config.WaitTime)*time.Second {
	// 	kong.Response.Exit(400, "Maximum Requests Reached", make(map[string][]string))
	// } else {
	// 	requests[host] = time.Now()
	// }

	serviceKey, _ := kong.Request.GetHeader("serviceKey")
	publicKey, _ := kong.Request.GetHeader("publicKey")
	requestHeader.serviceKey = serviceKey
	requestHeader.publicKey = publicKey
	// if serviceKeyErr != nil || publicKeyErr != nil {
	// 	fmt.Println("test")
	// 	// kong.Log.Err(serviceKeyErr.Error())
	// 	// kong.Log.Err(publicKeyErr.Error())
	// }

	// println(kong.Request.GetHeader("key"))
	id, name, redisKey := ValidateToken(requestHeader.publicKey)
	// fmt.Println(id, name, redisKey)

	// apiKey := config.WaitTime
	host, _ := kong.Request.GetHost()
	// keys := GetVal("key")
	ConnectRedisClient()
	var token = GetVal("key")
	// if serviceKey == "" || token == "" {
	if serviceKey == "" || token == "" || id == 0 || name == "" || redisKey == "" {
		// if apiKey != serviceKey || id == 0 {

		response := make([]string, 1)
		response[0] = "You have no correct key"
		j, _ := json.Marshal(response)
		x := make(map[string][]string)
		x["Content-Type"] = append(x["Content-Type"], "application/json")
		kong.Response.Exit(403, string(j), x)
	} else {
		_requests[host] = time.Now()
	}
}

// func GenerateToken(user User) string {
// 	expiresAt := time.Now().Add(24 * time.Hour).Unix()
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &authClaims{
// 		StandardClaims: jwt.StandardClaims{
// 			// Name:      user.Name,
// 			ExpiresAt: expiresAt,
// 		},
// 		UserID: user.ID,
// 		Key:    user.Key,
// 	})
// 	tokenString, _ := token.SignedString(key)
// 	return tokenString
// }

func main() {
	Version := "1.1"
	Priority := 1000
	_ = server.StartServer(New, Version, Priority)
}
