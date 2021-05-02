package authware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

type Config struct {
	Key        string
	UserModel  repository.UserModel
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
	if v == t.Expire {
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
	jwt.RegisteredClaims
	Expire  time.Time       `json:"e"`
	User    repository.User `json:"u"`
	TraceId string          `json:"t"`
}

func getJwt(key string) (*jwt.Builder, jwt.Verifier) {
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	builder := jwt.NewBuilder(signer)
	if err != nil {
		log.Fatal(err)
	}
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	return builder, verifier
}

func NewAuthware(c *Config) gin.HandlerFunc {
	_, verifier := getJwt(c.Key)
	return func(g *gin.Context) {
		tokenVerification := g.GetHeader("X-Token-Verification")
		if tokenVerification == "" {
			g.Next()
			return
		}

		token, err := jwt.ParseAndVerifyString(tokenVerification, verifier)
		if err != nil {
			log.Println(err)
			g.Next()
			return
		}

		tokenValue := &Token{}
		err = json.Unmarshal(token.RawClaims(), tokenValue)
		if err != nil {
			g.Next()
			return
		}

		if !c.validate(tokenValue) {
			g.Next()
			return
		}
		g.Set("authware:UserToken", tokenValue)
		g.Next()
	}
}

func GetUserFromContext(g *gin.Context) *repository.User {
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
	builder, _ := getJwt(c.Key)
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
				log.Printf("user login is unhashable: %v\n", err)
			}
			expire := time.Now().Add(1 * time.Hour)
			tokenValue := &Token{
				Expire:  expire,
				User:    *user,
				TraceId: string(traceId),
			}
			c.saveTraceId(tokenValue)
			newToken, err := builder.Build(tokenValue)
			if err != nil {
				g.JSON(http.StatusBadRequest, gin.H{
					"status": "cannot create token",
				})
				return
			}
			tokenString := newToken.String()
			g.JSON(http.StatusOK, gin.H{
				"status": "token created",
				"token":  tokenString,
				"expire": expire.UTC().Unix(),
			})
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"status": "you got already valid token",
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
		user := repository.User{
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
