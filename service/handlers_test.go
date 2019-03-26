package service_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	// Change this for your service
	service "github.com/orov.io/GoBackbone/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDataWrongPath(t *testing.T) {
	Convey("Given a HTTP request for /invalid", t, func() {
		req, resp := getTestDataWrongPathPetitionParameters()

		Convey("When the request is handler by the Router", func() {
			serveTest(req, resp)

			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func getTestDataWrongPathPetitionParameters() (*http.Request, *httptest.ResponseRecorder) {
	return getRequestAndResponseRecorder("GET", "/api/data/invalid", nil)
}

func TestPong(t *testing.T) {
	Convey("Given a ping HTTP request (/ping)", t, func() {
		req, resp := getTestPongPetitionParameters()

		Convey("When the request is handler by the Router", func() {
			serveTest(req, resp)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func getTestPongPetitionParameters() (*http.Request, *httptest.ResponseRecorder) {
	return getRequestAndResponseRecorder("GET", "/ping", nil)
}

func TestPrivatePong(t *testing.T) {
	Convey("Given a private ping HTTP request (/ping)", t, func() {
		req, resp := getTestPrivatePongPetitionParameters()

		Convey("When the request is handler by the Router", func() {
			serveTest(req, resp)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func getTestPrivatePongPetitionParameters() (*http.Request, *httptest.ResponseRecorder) {
	return getRequestAndResponseRecorder("GET", "/v1/ping", nil)
}

func serveTest(req *http.Request, resp *httptest.ResponseRecorder) *service.App {
	app := &service.App{}
	app.Initialize()
	app.ServeTestHTTP(resp, req)
	return app
}

func getRequestAndResponseRecorder(method, target string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, target, body)
	resp := httptest.NewRecorder()
	return req, resp
}
