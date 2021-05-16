package profile_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/controllers"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/views"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/kiselev-nikolay/inflr-be/pkg/tools/testhelper"
	"github.com/stretchr/testify/assert"
)

func getTestPlayer() *testhelper.Player {
	repo := &memorystore.MemoryStoreRepo{}
	repo.Connect()

	m := models.New(repo)
	ctrls := controllers.New(m)
	views := views.New(m)

	testplayer := testhelper.New(repo)
	testplayer.Router.POST("/ctrl/new", ctrls.New)
	testplayer.Router.POST("/ctrl/add-youtube", ctrls.AddYoutube)
	testplayer.Router.GET("/view/get", views.Get)
	return testplayer
}

func TestCtrlNew(t *testing.T) {
	assert := assert.New(t)
	testplayer := getTestPlayer()

	code, res := testplayer.TestPost("/ctrl/new", controllers.NewReq{
		Name: "Hello",
	})
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"created"}`, res)

	code, res = testplayer.TestGet("/view/get")
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"profile":{"name":"Hello","about":"","availability":0,"country":"","links":[],"telegram":null,"tiktok":null,"youtube":null},"status":"found"}`, res)
}

func TestCtrlAddYoutube(t *testing.T) {
	assert := assert.New(t)
	testplayer := getTestPlayer()

	code, res := testplayer.TestPost("/ctrl/new", controllers.NewReq{
		Name: "Hello",
	})
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"created"}`, res)

	file, err := ioutil.ReadFile("../../integrations/youtube/test/res.json")
	if err != nil {
		t.Fatal(err)
	}
	url := "https://www.googleapis.com/youtube/v3/channels?key=AIzaSyBQQ-zTp3e4o0GkJEbnnmH35hTMOSxsW_E&part=statistics&part=snippet&id=UC-lHJZR3Gqxm24_Vd_AJ5Yw"
	httpmock.Activate()
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, file))
	code, res = testplayer.TestPost("/ctrl/add-youtube", controllers.AddYoutubeReq{
		Link: "https://www.youtube.com/channel/UC-lHJZR3Gqxm24_Vd_AJ5Yw",
	})
	httpmock.DeactivateAndReset()
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"added"}`, res)

	code, res = testplayer.TestGet("/view/get")
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"profile":{"name":"Hello","about":"","availability":0,"country":"","links":[],"telegram":null,"tiktok":null,"youtube":{"UC-lHJZR3Gqxm24_Vd_AJ5Yw":{"title":"PewDiePie","description":"I make videos.","register":"2010-04-29T10:54:00Z","imageUrl":"https://yt3.ggpht.com/ytc/AAUvwnga3eXKkQgGU-3j1_jccZ0K9m6MbjepV0ksd7eBEw=s800-c-k-c0x00ffffff-no-rj","subs":110000000,"views":27225286026,"videos":4318}}},"status":"found"}`, res)
}
