package profile_test

import (
	"net/http"
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/kiselev-nikolay/inflr-be/pkg/tools/testhelper"
	"github.com/stretchr/testify/assert"
)

func getTestPlayer() *testhelper.Player {
	repo := &memorystore.MemoryStoreRepo{}
	repo.Connect()

	model := profile.NewModel(repo)
	ctrl := profile.NewController(model)
	view := profile.NewView(model)

	testplayer := testhelper.New(repo)
	testplayer.Router.POST("/ctrl/new", ctrl.New)
	testplayer.Router.GET("/view/get", view.Get)
	return testplayer
}

func TestCtrlNew(t *testing.T) {
	assert := assert.New(t)
	testplayer := getTestPlayer()

	code, res := testplayer.TestPost("/ctrl/new", profile.NewReq{
		Name: "Hello",
	})
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"created"}`, res)

	code, res = testplayer.TestGet("/view/get")
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"profile":{"name":"Hello","about":"","availability":0,"country":"","links":[],"telegram":null,"tiktok":null,"youtube":null},"status":"found"}`, res)
}
