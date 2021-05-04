package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/apierrors"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
)

type Controller struct {
	Model *Model
}

func NewController(model *Model) *Controller {
	return &Controller{Model: model}
}

type NewReq struct {
	Name string `json:"name" bind:"required"`
}

func (ctrl *Controller) New(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)
	body := &NewReq{}
	if err := g.BindJSON(body); err != nil {
		apierrors.MissingRequiredData.Send(g)
		return
	}

	_, err := ctrl.Model.Find(u.Login)
	if err == nil {
		apierrors.AlreadyHave.Send(g)
		return
	}

	p := &Profile{}
	p.Bio.Name = body.Name
	err = ctrl.Model.Send(u.Login, p)
	if err != nil {
		apierrors.CannotCreate.Send(g)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "created",
	})
}
