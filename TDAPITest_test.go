package TDApiTest_test

import (
	"encoding/json"
	"github.com/jeanlebrument/TDApiTest"
	"github.com/jeanlebrument/TDApiTest/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/jeanlebrument/TDApiTest/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/jeanlebrument/TDApiTest/Godeps/_workspace/src/github.com/unrolled/render"
	"net/http"
	"net/url"
	"testing"
)

const (
	GET       = "GET"
	PathHello = "/hello"
)

var headerValueArr = []string{"Content-Type", "application/x-www-form-urlencoded; param=value"}
var headerValueMap = map[string]string{"Content-Type": "application/x-www-form-urlencoded; param=value"}
var urlParams = url.Values{"key": {"Happy", "Coding"}}

type UT struct {
	t *testing.T
}

func newUT(t *testing.T) *UT {
	return &UT{t: t}
}

func (ut *UT) Hello(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	assert.Nil(ut.t, err)

	render.New().JSON(w, http.StatusOK, r.Form)
}

func newRouter(t *testing.T) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods(GET).
		Path(PathHello).
		Name("Hello").
		Headers(headerValueArr...).
		Handler(http.HandlerFunc(newUT(t).Hello))

	return router
}

func checkHelloRoute(t *testing.T, result string) {
	expected, err := json.Marshal(urlParams)

	assert.Nil(t, err)
	assert.Equal(t, result, string(expected))
}

func TestGetRoutes(t *testing.T) {
	td := TDApiTest.NewTDApiTest(newRouter(t))

	assert.NotNil(t, td)

	td.TestContainers = TDApiTest.TestContainers{
		{
			Method: GET,
			Path:   PathHello,
			TestsToRun: TDApiTest.TestsToRun{
				{
					Desc:     "Hello route",
					Status:   http.StatusOK,
					Header:   headerValueMap,
					Params:   urlParams,
					TestFunc: checkHelloRoute,
				},
				{
					Desc:   "Hello route without test func",
					Status: http.StatusOK,
					Params: urlParams,
					Header: headerValueMap,
				},
				{
					Desc:   "Hello route without header",
					Status: http.StatusNotFound,
				},
			},
		},
	}

	td.RunTests(t)
}
