package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Email  string `json:"username"`
	Role   string `json:"roles"`
	UserID uint
	jwt.StandardClaims
}

func JwtTokenStart(ctx *gin.Context, userId uint, email string, role string) string {
	tokenString, err := createToken(userId, email, role)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "Fail",
			"Error":  "Failed To Create Token",
			"code":   500,
		})
	}
	return tokenString
}

func createToken(userId uint, email string, role string) (string, error) {
	claims := Claims{
		Email:  email,
		Role:   role,
		UserID: uint(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		fmt.Println("==============", err, tokenString)
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("jwtToken" + requiredRole)
		if err != nil {
			ctx.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Can't find Cookie",
				"code":    401,
			})
			ctx.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRETKEY")), nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Invalid Or Expired Token",
				"code":    401,
			})
			ctx.Abort()
			return
		}
		if claims.Role != requiredRole {
			ctx.JSON(403, gin.H{
				"status": "Forbidden",
				"Error":  "Insufficient Permission",
				"code":   403,
			})
			ctx.Abort()
			return
		}
		ctx.Set("userId", claims.UserID)
		ctx.Next()
	}
}
