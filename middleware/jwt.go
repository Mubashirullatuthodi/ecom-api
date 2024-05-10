package middleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/models"
)

var SecretKey = []byte("qwertuiouplkhgfdsazxcvbnm")
var Userdetails models.User
var BlacklistedToken = make(map[string]bool)

type Claims struct {
	ID    uint
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func JwtToken(c *gin.Context, id uint, email string, role string) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Token": signedToken})
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.GetHeader("Authorization")
		if tokenstring == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Token not provided"})
			c.Abort()
			return
		}
		if BlacklistedToken[tokenstring] {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "please login!!"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No permission"})
			c.Abort()
			return
		}

		c.Set("userid", claims.ID)
		c.Next()
	}
}
