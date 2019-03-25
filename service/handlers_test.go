package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// Change this for your service
	service "github.com/orov.io/GoBackbone/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDataWrongPath(t *testing.T) {
	Convey("Given a HTTP request for /invalid", t, func() {
		req := httptest.NewRequest("GET", "/api/data/invalid", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handler by the Router", func() {
			app := service.App{}
			app.Initialize()
			app.ServeTestHTTP(resp, req)

			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestPong(t *testing.T) {
	Convey("Given a ping HTTP request (/ping)", t, func() {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handler by the Router", func() {
			app := service.App{}
			app.Initialize()
			app.ServeTestHTTP(resp, req)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestPrivatePong(t *testing.T) {
	Convey("Given a private ping HTTP request (/ping)", t, func() {
		req := httptest.NewRequest("GET", "/v1/ping", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handler by the Router", func() {
			app := service.App{}
			app.Initialize()
			app.ServeTestHTTP(resp, req)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
