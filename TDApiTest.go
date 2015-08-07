package TDApiTest

import (
	"github.com/jeanlebrument/TDApiTest/Godeps/_workspace/src/github.com/stretchr/testify/assert"
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

type TDApiTest struct {
	router         AbstractRouter
	RespRec        *httptest.ResponseRecorder
	TestContainers TestContainers
}

func init() {
	var _ assert.Assertions // Tricks for Godep to import packages needed by tests.
}

func NewTDApiTest(router AbstractRouter) *TDApiTest {
	return &TDApiTest{router: router}
}

func (td *TDApiTest) beforeEach() {
	td.RespRec = httptest.NewRecorder()
}

func (td *TDApiTest) RunTests(t *testing.T) {
	for _, route := range td.TestContainers {
		for _, testToRun := range route.TestsToRun {
			req, err := http.NewRequest(route.Method, route.Path,
				strings.NewReader(testToRun.Params.Encode()))

			assert.Nil(t, err)

			for k, v := range testToRun.Header {
				req.Header.Set(k, v)
			}

			td.beforeEach()
			td.router.ServeHTTP(td.RespRec, req)
			t.Logf("Executing test: %s, function called: %s", testToRun.Desc,
				runtime.FuncForPC(reflect.ValueOf(testToRun.TestFunc).Pointer()).Name())
			assert.Equal(t, testToRun.Status, td.RespRec.Code)
			content, err := ioutil.ReadAll(td.RespRec.Body)
			assert.Nil(t, err)
			if testToRun.TestFunc != nil {
				testToRun.TestFunc(t, strings.TrimSpace(string(content)))
			}
		}
	}
}
