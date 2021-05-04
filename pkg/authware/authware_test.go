package authware_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/models/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/stretchr/testify/require"
)

func makeTestServer() func(method, url, body, tok string) (int, string) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	repo := memorystore.MemoryStoreRepo{}
	repo.Connect()
	um := user.NewUserModel(&repo)
	pw := passwords.Passworder{KeySecret: []byte("test")}
	c := &authware.Config{
		Key:        []byte("test"),
		UserModel:  *um,
		Passworder: pw,
	}
	router.Use(authware.NewAuthware(c))
	router.POST("/token", authware.NewTokenHandler(c))
	router.POST("/register", authware.NewRegisterHandler(c))
	router.GET("/test", func(g *gin.Context) {
		u := authware.GetUserFromContext(g)
		if u == nil {
			g.Status(http.StatusUnauthorized)
			return
		}
		g.Status(http.StatusOK)
	})
	return func(method, uri, body, tok string) (int, string) {
		req := httptest.NewRequest(method, "http://inflr.app/"+strings.Trim(uri, "/"), bytes.NewBufferString(body))
		if tok != "" {
			req.Header.Add("X-Token-Verification", tok)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		return w.Code, string(bodyBytes)
	}
}

func TestUnauthorized(t *testing.T) {
	require := require.New(t)
	req := makeTestServer()

	resCode, _ := req("GET", "/test", "", "")
	require.Equal(http.StatusUnauthorized, resCode)

	resCode, _ = req("GET", "/test", "", "faketoken")
	require.Equal(http.StatusUnauthorized, resCode)
}

func TestIdiotGetToken(t *testing.T) {
	require := require.New(t)
	req := makeTestServer()

	resCode, resBody := req("POST", "/token", `{"user": "wordpress", "password": "admin"}`, "")
	require.Equal(`{"status":"missing required data"}`, resBody)
	require.Equal(http.StatusBadRequest, resCode)

	resCode, resBody = req("POST", "/token", `{"userKey": "wordpress", "passwordVerification": "admin"}`, "")
	require.Equal(`{"status":"user not found"}`, resBody)
	require.Equal(http.StatusBadRequest, resCode)
}

func TestRegister(t *testing.T) {
	require := require.New(t)
	req := makeTestServer()

	resCode, resBody := req("POST", "/register", `{"userKey": "test-register", "passwordVerification": "test"}`, "")
	require.Equal(`{"status":"user created"}`, resBody)
	require.Equal(http.StatusOK, resCode)

	resCode, resBody = req("POST", "/register", `{"userKey": "test-register", "passwordVerification": "test"}`, "")
	require.Equal(`{"status":"user exist"}`, resBody)
	require.Equal(http.StatusBadRequest, resCode)
}

func TestGetToken(t *testing.T) {
	require := require.New(t)
	req := makeTestServer()

	resCode, _ := req("GET", "/test", "", "")
	require.Equal(http.StatusUnauthorized, resCode)

	resCode, _ = req("POST", "/register", `{"userKey": "test-login", "passwordVerification": "test"}`, "")
	require.Equal(http.StatusOK, resCode)
	resCode, resBody := req("POST", "/token", `{"userKey": "test-login", "passwordVerification": "test"}`, "")
	require.Equal(http.StatusOK, resCode)
	message := &struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}{}
	err := json.Unmarshal([]byte(resBody), message)
	require.NoError(err)
	require.Equal("token created", message.Status)
	require.NotEmpty(message.Token)

	resCode, _ = req("GET", "/test", "", message.Token)
	require.Equal(http.StatusOK, resCode)
}
