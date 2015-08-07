package TDApiTest_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jeanlebrument/TDApiTest"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
	"testing"
)

const (
	GET       = "GET"
	PathHello = "/hello"
	urlKey    = "keys"
)

var headerValueArr = []string{"Content-Type", "application/json"}
var headerValueMap = map[string]string{"Content-Type": "application/json"}
var urlValues = []string{"Happy", "Coding"}
var requestResponse = map[string]string{"Hello": "World"}

type LoggerTest struct{}

func NewLoggerTest() *LoggerTest {
	return &LoggerTest{}
}

func (l LoggerTest) Log(message string) error {
	fmt.Printf("%s\n", message)

	return nil
}

func Hello(w http.ResponseWriter, r *http.Request) {
	NewLoggerTest().Log(fmt.Sprintf("%s\n", r.FormValue(urlKey)))
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
	td := TDApiTest.NewTDApiTest(newRouter(), NewLoggerTest())

	assert.NotNil(t, td)

	td.BeforeEach()

	td.TestContainers = TDApiTest.TestContainers{
		{
			Method: GET,
			Path:   PathHello,
			TestsToRun: TDApiTest.TestsToRun{
				{
					Desc:     "Hello",
					Status:   http.StatusOK,
					Header:   headerValueMap,
					Params:   url.Values{urlKey: urlValues},
					TestFunc: checkHelloRoute,
				},
			},
		},
	}

	td.RunTests(t)
}
