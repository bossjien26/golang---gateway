package authentication

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt"
)

var redisConfig = &RedisConfig{}

func init() {
	ctx := context.Background()
	client := NewClient(ctx, os.Getenv("redisConnect"))
	// redisConfig := &RedisConfig{
	// 	ctx:    ctx,
	// 	client: client,
	// }
	redisConfig.ctx = ctx
	redisConfig.client = client
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

type authClaims struct {
	jwt.StandardClaims
	UserID uint `json:"userId"`
	Key    string
}

// var ctx = context.Background()
var key = []byte(os.Getenv("token_pwd"))

func NewClient(ctx context.Context, addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(pong)
	return client
}

func SetVal(key string, val string) bool {
	if err := redisConfig.client.Set(redisConfig.ctx, key, val, 0).Err(); err != nil {
		log.Println("client.Set failed", err)
		return false
	}
	return true
}

func GetVal(key string) string {
	val, err := redisConfig.client.Get(redisConfig.ctx, key).Result()
	if err != nil {
		log.Println("client.Get failed", err)
		return ""
	}
	return val
}

func GenerateToken(user User) string {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &authClaims{
		StandardClaims: jwt.StandardClaims{
			Name:      user.Name,
			ExpiresAt: expiresAt,
		},
		UserID: user.ID,
		Key:    user.Key,
	})
	tokenString, _ := token.SignedString(key)
	return tokenString
}

func ValidateToken(tokenString string) (uint, string, string) {
	var claims authClaims
	token, _ := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	// if err != nil {
	// 	return 0, "", err
	// }
	if !token.Valid {
		return 0, "", ""
	}
	id := claims.UserID
	name := claims.Name
	key := claims.Key
	return id, name, key
}
