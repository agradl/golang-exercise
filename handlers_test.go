package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/agradl/golang-exercise/testutil"
)

func TestHandlerWrapperHelperChecksShutdownState(t *testing.T) {
	registrar := &mockRegistrar{}
	state := &mockState{isShutdownReturn: true}
	makeHandlerWithState(registrar, "/test", state, methods(http.MethodGet), func(otherState Server, writer http.ResponseWriter, request *http.Request) {})
	req := httptest.NewRequest("GET", "http://foo/test", nil)
	w := httptest.NewRecorder()
	registrar.handler(w, req)
	assertResponse(t, w, 503, "text/plain; charset=utf-8", "Server shutting down\n")
}
func TestHandlerWrapperHelperValidatesMethod(t *testing.T) {
	registrar := &mockRegistrar{}
	state := &mockState{}
	makeHandlerWithState(registrar, "/test", state, methods(http.MethodGet), func(otherState Server, writer http.ResponseWriter, request *http.Request) {})
	req := httptest.NewRequest("POST", "http://foo/test", nil)
	w := httptest.NewRecorder()
	registrar.handler(w, req)
	assertResponse(t, w, 405, "text/plain; charset=utf-8", "Invalid request method.\n")
}
func TestHandlerWrapperHelperPassesStateToHandler(t *testing.T) {
	registrar := &mockRegistrar{}
	state := &mockState{}
	var handlerState Server
	makeHandlerWithState(registrar, "/test", state, methods(http.MethodGet), func(otherState Server, writer http.ResponseWriter, request *http.Request) {
		handlerState = otherState
	})
	req := httptest.NewRequest("GET", "http://foo/test", nil)
	w := httptest.NewRecorder()
	registrar.handler(w, req)
	if handlerState != state {
		t.Errorf("Expected handler state to equal passed state")
	}
}
func TestHashHandlerNonIntIndex(t *testing.T) {
	state := &mockState{}
	req := httptest.NewRequest("POST", "http://foo/hash/abc", nil)
	w := httptest.NewRecorder()
	getHashHandler(state, w, req)
	assertResponse(t, w, 400, "text/plain; charset=utf-8", "Invalid hash index.\n")
}
func TestHashHandlerNonExistantIndex(t *testing.T) {
	state := &mockState{getHashReturn: "not found"}
	req := httptest.NewRequest("POST", "http://foo/hash/123", nil)
	w := httptest.NewRecorder()
	getHashHandler(state, w, req)
	assertResponse(t, w, 400, "text/plain; charset=utf-8", "Invalid hash index.\n")
}
func TestHashHandlerHappyPath(t *testing.T) {
	state := &mockState{getHashReturn: "some hash"}
	req := httptest.NewRequest("POST", "http://foo/hash/123", nil)
	w := httptest.NewRecorder()
	getHashHandler(state, w, req)
	assertResponse(t, w, 200, "text/plain; charset=utf-8", "some hash")
}
func TestComputeHashHandlerHappyPath(t *testing.T) {
	data := url.Values{}
	data.Set("password", "foo")
	req := httptest.NewRequest("POST", "http://foo/hash", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	state := &mockState{doHashReturn: 22}
	computeHashHandler(state, w, req)

	testutil.AssertStringEqual(t, "foo", state.doHashArg1, "Password to Hash")
	testutil.AssertIntsEqual(t, 5, state.doHashArg2, "Hash Delay Seconds")
	assertResponse(t, w, 200, "text/plain; charset=utf-8", "22")

}
func TestComputeHashHandlerNoPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://foo/hash", nil)
	w := httptest.NewRecorder()
	state := &mockState{}
	computeHashHandler(state, w, req)
	assertResponse(t, w, 400, "text/plain; charset=utf-8", "Invalid request, missing param 'password'\n")
}
func TestShutdownHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://foo/shutdown", nil)
	w := httptest.NewRecorder()
	state := &mockState{}
	shutdownHandler(state, w, req)
	testutil.AssertTrue(t, state.shutdownCalled, "shutdown called")
	assertResponse(t, w, 200, "text/plain; charset=utf-8", "initiating shutdown")
}

func TestStatsHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://foo/stats", nil)
	w := httptest.NewRecorder()
	state := &mockState{getStatsReturn: &statsObj{Total: 2, Average: 3.0}}
	statsHandler(state, w, req)
	testutil.AssertStringEqual(t, "/hash", state.getStatsArg, "getStatsArg")
	assertResponse(t, w, 200, "application/json", `{"total":2,"average":3}`)
}

func assertResponse(t *testing.T, w *httptest.ResponseRecorder, statusCode int, contentType string, expectedBody string) {
	resp := w.Result()
	testutil.AssertIntsEqual(t, statusCode, resp.StatusCode, "Status Code")
	testutil.AssertStringEqual(t, contentType, resp.Header.Get("Content-Type"), "Content-Type")
	body, _ := ioutil.ReadAll(resp.Body)
	testutil.AssertStringEqual(t, expectedBody, string(body), "Body")
}

type mockRegistrar struct {
	handler func(writer http.ResponseWriter, request *http.Request)
}

func (registrar *mockRegistrar) registerHandler(pattern string, handler func(writer http.ResponseWriter, request *http.Request)) {
	registrar.handler = handler
}

type mockState struct {
	getPendingHashCtReturn int
	getHashArg             int
	getHashReturn          string
	shutdownCalled         bool
	isShutdownReturn       bool
	logResponseArg1        string
	logResponseArg2        int
	getStatsArg            string
	getStatsReturn         *statsObj
	doHashArg1             string
	doHashArg2             int
	doHashReturn           int
}

func (state *mockState) getPendingHashCt() int {
	return state.getPendingHashCtReturn
}

func (state *mockState) getHash(index int) string {
	state.getHashArg = index
	return state.getHashReturn
}

func (state *mockState) shutdown() {
	state.shutdownCalled = true
}

func (state *mockState) isShutdown() bool {
	return state.isShutdownReturn
}

func (state *mockState) logResponse(pattern string, time int) {
	state.logResponseArg1 = pattern
	state.logResponseArg2 = time
}

func (state *mockState) getStats(pattern string) *statsObj {
	state.getStatsArg = pattern
	return state.getStatsReturn
}

func (state *mockState) doHash(password string, delayS int) int {
	state.doHashArg1 = password
	state.doHashArg2 = delayS
	return state.doHashReturn
}
