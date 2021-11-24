package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_service "github.com/cyberdr0id/referral/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type Request struct {
	name     string
	password string
}

func TestServer_SignUp(t *testing.T) {
	testTable := []struct {
		testName             string
		request              Request
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
		mock                 func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName: "Success",
			request: Request{
				name:     "Nameee",
				password: "Password",
			},
			requestBody:          `{"name":"Nameee","password":"password"}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":"1"}`,
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().SignUp(name, password).Return("1", nil)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.request.name, tc.request.password)

			s := Server{
				Auth: auth,
			}

			r := mux.NewRouter()
			r.HandleFunc("/auth/signup", s.SignUp).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(tc.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestServer_LogIn(t *testing.T) {
	testTable := []struct {
		testName             string
		request              Request
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
		mock                 func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName: "Success",
			request: Request{
				name:     "Alexander",
				password: "password",
			},
			requestBody:          `{"name":"Alexander","password":"password"}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"accessToken":"token","refreshToken":"token"}`,
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return("token", "token", nil)
			},
		},
		{
			testName: "Failure: empty password",
			request: Request{
				name:     "Alexander",
				password: "",
			},
			requestBody:          `{"name":"Alexander","password":""}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid parameter: password\n",
			mock:                 func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: empty name",
			request: Request{
				name:     "",
				password: "password",
			},
			requestBody:          `{"name":"","password":""}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid parameter: name\n",
			mock:                 func(s *mock_service.MockAuth, name, password string) {},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.request.name, tc.request.password)

			s := Server{
				Auth: auth,
			}

			r := mux.NewRouter()
			r.HandleFunc("/auth/login", s.LogIn).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(tc.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
