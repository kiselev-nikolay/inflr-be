package landing_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/controllers"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/kiselev-nikolay/inflr-be/pkg/tools/testhelper"
	"github.com/stretchr/testify/assert"
)

const testTitle = "Hello world I am in test suit right now"

func getTestPlayer() *testhelper.Player {
	repo := &memorystore.MemoryStoreRepo{}
	repo.Connect()

	testplayer := testhelper.New(repo)
	landing.Connect(testplayer.Router, "", repo)
	return testplayer
}

func TestLanding(t *testing.T) {
	assert := assert.New(t)
	testplayer := getTestPlayer()

	code, res := testplayer.TestPost("/new", controllers.CreateReq{
		Key:   "test",
		Title: testTitle,
	})
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"created"}`, res)

	code, res = testplayer.TestPost("/add-card", controllers.AddCardReq{
		LandingKey: "test",
		Title:      "My website",
		Text:       "Full-stack web developer. High-load SRE, 1 year expirence. Go developer, 1+ year expirence. Previously Python Senior Developer, Team Lead, 4 year expirence. Typescript developer, 4+ year expirence. React and React Native, 1+ year expirence. Previously VueJS developer, 3 year expirence. Well know Kubernetes, Linux, Docker, PostgreSQL, MySQL, InfluxDB, Scylla, Redis, MongoDB, RabbitMQ, Traefik, Ansible, C++, Arduino, Tensorflow, PyTorch. ",
		Link:       "https://nikolay.works",
	})
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"status":"updated"}`, res)

	code, res = testplayer.TestGet("/l/test/test")
	assert.Equal(http.StatusOK, code)
	assert.True(strings.Contains(res, testTitle))

	// ioutil.WriteFile("example.html", []byte(res), 0644)
}
