package controllers_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thatnerdjosh/example-webservices/internal/controllers"
)

func baseChecks(t *testing.T, rr *httptest.ResponseRecorder) {
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"expected content-type %s, got \"%s\"",
			"application/json",
			contentType)
	}
}

func TestNotFoundHandler(t *testing.T) {
	controller := controllers.TaskController{}
	t.Run("404 - Not Found", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST", "/abc123", nil)
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.Handler(controller)

		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf(
				"expected status code %d, received %d",
				http.StatusNotFound,
				rr.Code,
			)
		}
	})
}

func TestExecuteTask(t *testing.T) {
	controller := controllers.TaskController{}

	// TODO: Extract to contract tests
	t.Run("403 - not authenticated", func(t *testing.T) {
		req, err := http.NewRequest(
			"POST", "/", nil)
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.Handler(controller)

		handler.ServeHTTP(rr, req)

		baseChecks(t, rr)
		if rr.Code != http.StatusForbidden {
			t.Errorf(
				"expected status code %d, received %d",
				http.StatusForbidden,
				rr.Code,
			)
		}
	})

	t.Run("400 - run invalid task", func(t *testing.T) {
		// TODO: Extract to fixtures
		data := `{"foo":1}`
		req, err := http.NewRequest(
			"POST", "/", bytes.NewBuffer([]byte(data)))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("Authorization", "foobar")

		rr := httptest.NewRecorder()
		handler := http.Handler(controller)

		handler.ServeHTTP(rr, req)

		baseChecks(t, rr)
		if rr.Code != http.StatusBadRequest {
			t.Errorf(
				"expected status code %d, received %d",
				http.StatusBadRequest,
				rr.Code,
			)
		}
	})

	// t.Run("200 - execute valid task", func(t *testing.T) {
	// 	// TODO: Extract to fixtures
	// 	data := `{"id": "abc123"}`
	// 	req, err := http.NewRequest(
	// 		"POST", "/", bytes.NewBuffer([]byte(data)))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	rr := httptest.NewRecorder()
	// 	handler := http.Handler(controller)

	// 	handler.ServeHTTP(rr, req)

	// 	baseChecks(t, rr)
	// 	if rr.Code != http.StatusOK {
	// 		t.Errorf(
	// 			"expected status code %d, received %d",
	// 			http.StatusOK,
	// 			rr.Code,
	// 		)
	// 	}
	// })

}
