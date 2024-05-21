package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mubashir/e-commerce/models"
)

var Userdetails models.User

type Claims struct {
	ID    uint
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func JwtToken(c *gin.Context, id uint, email string, role string) (string, error) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to sign token"})
		return "", nil
	}

	return signedToken, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring, err := c.Cookie("Authorization" + requiredRole)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRETKEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Please Login !!!"})
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
