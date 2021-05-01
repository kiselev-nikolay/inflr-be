package authware

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

const authErrorText = "Auth error"

type Config struct {
	Key        []byte
	UserModel  repository.UserModel
	Passworder passwords.Passworder
	Verbose    bool
}

type UserAuthBody struct {
	UserKey      string `json:"userKey"`
	UserPassword string `json:"passwordVerification"`
}

func New(configuration *Config) gin.HandlerFunc {
	errorResp := func(g *gin.Context, message string) {
		g.AbortWithStatus(401)
	}
	if configuration.Verbose {
		errorResp = func(g *gin.Context, message string) {
			g.AbortWithStatusJSON(401, gin.H{
				"status":  401,
				"message": message,
			})
		}
	}
	return func(g *gin.Context) {
		if g.GetHeader("X-Token") == "" {
			body := &UserAuthBody{}
			g.Bind(body)
			user, err := configuration.UserModel.Find(body.UserKey)
			if err != nil {
				errorResp(g, "User not found")
				return
			}
			passwordCorrect, err := configuration.Passworder.IsCorrect(user.SecretPassword, body.UserPassword)
			if err != nil {
				errorResp(g, "Verification failed")
				return
			}
			if !passwordCorrect {
				errorResp(g, "Verification wrong")
				return
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"e": time.Now().Add(1 * time.Hour),
				"u": body.UserKey,
			})
			tokenString, err := token.SignedString(configuration.Key)
			if err != nil {
				errorResp(g, "Cannot create token")
				return
			}
			g.Header("X-Token", tokenString)
		}
		g.Next()
	}
}
