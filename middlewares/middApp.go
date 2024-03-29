package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Max-Age", "1728000")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

// 用户验证
type customClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"IsAdmin"`
	jwt.RegisteredClaims
}

var jwtKey []byte = []byte("secret")

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) { return jwtKey, nil })
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": fmt.Sprintf("access token parse error: %v", err)})
			return
		}
		//验证过期
		if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
			if !claims.VerifyExpiresAt(time.Now(), false) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": "access token expired"})
				return
			}
			ctx.Set("claims", claims)
		} else {
			//token解析失败
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": fmt.Sprintf("Claims parse error: %v", err)})
			return
		}
		ctx.Next()
	}
}

//func InterceptRequests() gin.HandlerFunc {
//	return func(context *gin.Context) {
//
//		get := redisServer.RedisDb.Get()
//		defer func(get redis.Conn) {
//			err := get.Close()
//			if err != nil {
//
//			}
//		}(get)
//
//		ip := context.ClientIP()
//		if len(ip) == 0 {
//			panic("IP错误")
//			return
//		}
//
//		keyRedis := fmt.Sprintf("%s-%s", context.Request.URL, ip)
//		existsRedis := redisServer.ExistsRedis(keyRedis, get)
//		if existsRedis == false {
//			redisServer.SetRedis(keyRedis, 1, 60, get)
//		}
//		redisServer.GetRedis(keyRedis, get)
//		redisServer.ExpireRedis(keyRedis, 60, get)
//		redisServer.IncrRedis(keyRedis, get)
//
//		fmt.Println()
//	}
//}
