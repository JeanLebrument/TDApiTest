package TDApiTest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type TestToRun struct {
	Desc     string
	Header   map[string]string
	Params   url.Values
	Status   int
	TestFunc func(*testing.T, string)
}

type TestsToRun []TestToRun

type TestContainer struct {
	Method     string
	Path       string
	TestsToRun TestsToRun
}

type TestContainers []TestContainer

type AbstractRouter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Logger interface {
	Log(message string) error
}

type TDApiTest struct {
	router         AbstractRouter
	logger         Logger
	RespRec        *httptest.ResponseRecorder
	TestContainers TestContainers
}

func NewTDApiTest(router AbstractRouter, logger Logger) *TDApiTest {
	return &TDApiTest{router: router, logger: logger}
}

func (td *TDApiTest) beforeEach() {
	td.RespRec = httptest.NewRecorder()
}

func (td *TDApiTest) RunTests(t *testing.T) {
	td.beforeEach()

	for _, route := range td.TestContainers {
		for _, testToRun := range route.TestsToRun {
			req, err := http.NewRequest(route.Method, route.Path,
				strings.NewReader(testToRun.Params.Encode()))

			assert.Nil(t, err)

			for k, v := range testToRun.Header {
				req.Header.Set(k, v)
			}

			td.router.ServeHTTP(td.RespRec, req)
			assert.Equal(t, td.RespRec.Code, testToRun.Status)
			content, err := ioutil.ReadAll(td.RespRec.Body)
			assert.Nil(t, err)
			td.logger.Log(fmt.Sprintf("Executing test: %s, function called: %s", testToRun.Desc,
				runtime.FuncForPC(reflect.ValueOf(testToRun.TestFunc).Pointer()).Name()))
			testToRun.TestFunc(t, strings.TrimSpace(string(content)))
		}
	}
}
