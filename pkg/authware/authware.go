package authware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
)

type Config struct {
	Key        []byte
	UserModel  user.UserModel
	Passworder passwords.Passworder
	traceIds   map[string]time.Time
}

func (c *Config) validate(t *Token) bool {
	if c.traceIds == nil {
		return false
	}
	v, ok := c.traceIds[t.TraceId]
	if !ok {
		return false
	}
	if v.Equal(t.Expire) {
		return true
	}
	return false
}

func (c *Config) saveTraceId(t *Token) {
	if c.traceIds == nil {
		c.traceIds = make(map[string]time.Time)
	}
	c.traceIds[t.TraceId] = t.Expire
}

type Token struct {
	Expire  time.Time `json:"e"`
	User    user.User `json:"u"`
	TraceId string    `json:"t"`
}

func (t *Token) Valid() error {
	return nil
}

func NewAuthware(c *Config) gin.HandlerFunc {
	return func(g *gin.Context) {
		tokenVerification := g.GetHeader("X-Token-Verification")
		if tokenVerification == "" {
			g.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenVerification, &Token{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return c.Key, nil
		})

		if err != nil {
			g.Next()
			return
		}
		tokenValue, ok := token.Claims.(*Token)
		if !ok || !token.Valid {
			g.Next()
			return
		}

		if !c.validate(tokenValue) {
			g.Next()
			return
		}
		g.Set("authware:UserToken", *tokenValue)
		g.Next()
	}
}

func RequiredLoginPassed(g *gin.Context) bool {
	_, exists := g.Get("authware:UserToken")
	if !exists {
		g.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func GetUserFromContext(g *gin.Context) *user.User {
	value, exists := g.Get("authware:UserToken")
	if !exists {
		return nil
	}
	token := value.(Token)
	return &token.User
}

type UserAuthBody struct {
	UserKey      string `json:"userKey" binding:"required"`
	UserPassword string `json:"passwordVerification" binding:"required"`
}

func NewTokenHandler(c *Config) gin.HandlerFunc {
	return func(g *gin.Context) {
		if GetUserFromContext(g) == nil {
			body := &UserAuthBody{}
			if err := g.BindJSON(body); err != nil {
				g.JSON(http.StatusBadRequest, gin.H{
					"status": "missing required data",
				})
				return
			}
			user, err := c.UserModel.Find(body.UserKey)
			if err != nil {
				g.JSON(http.StatusBadRequest, gin.H{
					"status": "user not found",
				})
				return
			}
			passwordCorrect, err := c.Passworder.IsCorrect(user.SecretPassword, body.UserPassword)
			if err != nil {
				g.JSON(http.StatusBadRequest, gin.H{
					"status": "verification failed",
				})
				return
			}
			if !passwordCorrect {
				g.JSON(http.StatusBadRequest, gin.H{
					"status": "verification wrong",
				})
				return
			}
			traceId, err := c.Passworder.Hash(user.Login)
			if err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{
					"status": "cannot do verification",
				})
				return
			}
			expire := time.Now().Add(1 * time.Hour)
			tokenValue := &Token{
				Expire:  expire,
				User:    *user,
				TraceId: string(traceId),
			}
			c.saveTraceId(tokenValue)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenValue)

			tokenString, err := token.SignedString(c.Key)
			if err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{
					"status": "cannot make token",
				})
				return
			}

			g.JSON(http.StatusOK, gin.H{
				"status": "token created",
				"token":  tokenString,
				"expire": expire.UTC().Unix(),
			})
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"status": "already valid token",
		})
	}
}

func NewRegisterHandler(c *Config) gin.HandlerFunc {
	return func(g *gin.Context) {
		body := &UserAuthBody{}
		if err := g.BindJSON(body); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"status": "missing required data",
			})
			return
		}
		_, err := c.UserModel.Find(body.UserKey)
		if err == nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"status": "user exist",
			})
			return
		}
		user := user.User{
			Login: body.UserKey,
		}
		hash, err := c.Passworder.Hash(body.UserPassword)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"status": "verification failed",
			})
			return
		}
		user.SecretPassword = hash
		err = c.UserModel.Send(user.Login, &user)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"status": "user unsaved",
			})
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"status": "user created",
		})
	}
}
