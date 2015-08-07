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
	urlKey    = "keys"
)

var headerValueArr = []string{"Content-Type", "application/test.1.0+json"}
var headerValueMap = map[string]string{"Content-Type": "application/test.1.0+json"}
var urlValues = []string{"Happy", "Coding"}
var requestResponse = map[string]string{"Hello": "World"}

func Hello(w http.ResponseWriter, r *http.Request) {
	render.New().JSON(w, http.StatusOK, requestResponse)
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods(GET).
		Path(PathHello).
		Name("Hello").
		Headers(headerValueArr...).
		Handler(http.HandlerFunc(Hello))

	return router
}

func checkHelloRoute(t *testing.T, result string) {
	expected, err := json.Marshal(requestResponse)

	assert.Nil(t, err)
	assert.Equal(t, result, string(expected))
}

func TestGetRoutes(t *testing.T) {
	td := TDApiTest.NewTDApiTest(newRouter())

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
					Params:   url.Values{urlKey: urlValues},
					TestFunc: checkHelloRoute,
				},
				{
					Desc:   "Hello route without test func",
					Status: http.StatusOK,
					Params: url.Values{urlKey: urlValues},
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
