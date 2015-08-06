package TDApiTest_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jeanlebrument/TDApiTest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type LoggerTest struct{}

func NewLoggerTest() *LoggerTest {
	return &LoggerTest{}
}

func (l LoggerTest) Log(message string) error {
	fmt.Printf("%s\n", message)

	return nil
}

func SetContentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func SetHeaderOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

type HelloStruct struct {
	Str1 string `json:"str1"`
	Str2 string `json:"str2"`
}

func Hello(res http.ResponseWriter, req *http.Request) {
	SetContentTypeJson(res)
	SetHeaderOk(res)

	if err := json.NewEncoder(res).Encode(HelloStruct{
		Str1: "Hello",
		Str2: "World",
	}); err != nil {
		panic(err)
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").
		Path("/hello").
		Name("Hello").
		Handler(http.HandlerFunc(Hello))

	return router
}

func checkHelloRoute(t *testing.T, result string) {
	expected, err := json.Marshal(HelloStruct{
		Str1: "Hello",
		Str2: "World",
	})

	assert.Nil(t, err)
	assert.Equal(t, result, string(expected))
}

func TestGetRoutes(t *testing.T) {
	td := TDApiTest.NewTDApiTest(NewRouter(), NewLoggerTest())

	assert.NotNil(t, td)

	td.BeforeEach()

	td.TestContainers = TDApiTest.TestContainers{
		{
			Method: "GET",
			Path:   "/hello",
			TestsToRun: TDApiTest.TestsToRun{
				{
					Desc:     "Hello",
					Status:   http.StatusOK,
					TestFunc: checkHelloRoute,
				},
			},
		},
	}

	td.RunTests(t)
}
