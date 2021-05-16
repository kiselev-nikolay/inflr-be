package views

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
)

//go:embed templates/get.html
var templateGetText string

var (
	templateGet *template.Template
)

func init() {
	var err error
	templateGet, err = template.New("landing/get").Parse(templateGetText)
	if err != nil {
		panic(err)
	}
}

type View struct {
	Model *models.Model
}

func New(model *models.Model) *View {
	return &View{Model: model}
}

func (c *View) Get(g *gin.Context) {
	u := &user.User{
		Login: g.Param("user"),
	}

	landing, findErr := c.Model.Find(u, g.Param("landing"))
	if findErr != nil {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}

	htmlLanding := bytes.NewBuffer([]byte{})
	executeErr := templateGet.Execute(htmlLanding, landing)
	if executeErr != nil {
		// TODO logging
		fmt.Println(executeErr)
		g.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	g.Data(http.StatusOK, "text/html; charset=utf-8", htmlLanding.Bytes())
}
