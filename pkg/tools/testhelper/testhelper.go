package testhelper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

type Player struct {
	Router *gin.Engine
	repo   repository.Repo
}

func (p *Player) TestDo(method string, uri string, bodyJson interface{}) (int, string) {
	body := bytes.NewBuffer([]byte(""))
	if bodyJson != nil {
		rawBody, _ := json.Marshal(bodyJson)
		body = bytes.NewBuffer(rawBody)
	}
	req := httptest.NewRequest(method, "http://inflr.app/"+strings.Trim(uri, "/"), body)
	w := httptest.NewRecorder()
	p.Router.ServeHTTP(w, req)
	bodyBytes, _ := ioutil.ReadAll(w.Body)
	return w.Code, string(bodyBytes)
}
func (p *Player) TestGet(uri string) (int, string) {
	return p.TestDo("GET", uri, nil)
}
func (p *Player) TestPost(uri string, body interface{}) (int, string) {
	return p.TestDo("POST", uri, body)
}

func New(repo repository.Repo) *Player {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(g *gin.Context) {
		g.Set("authware:UserToken", authware.Token{
			Expire: time.Now().Add(15 * time.Minute),
			User: user.User{
				Login: "test",
			},
			TraceId: "test",
		})
		g.Next()
	})
	return &Player{
		Router: router,
		repo:   repo,
	}
}
