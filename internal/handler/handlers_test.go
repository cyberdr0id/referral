package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyberdr0id/referral/internal/service"
	mock_service "github.com/cyberdr0id/referral/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TODO: handle 500 error login/signup

const (
	name           = "testName"
	password       = "password"
	emptyParameter = ""
)

func TestServer_SignUp(t *testing.T) {
	testTable := []struct {
		testName           string
		request            SignUpRequest
		requestBody        string
		expectedStatusCode int
		expectedResponse   string
		mock               func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName: "Success",
			request: SignUpRequest{
				Name:     name,
				Password: password,
			},
			requestBody:        `{"name":"testName","password":"password"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   SignUpResponse{ID: "1"}.String(),
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().CreateUser(name, password).Return("1", nil)
			},
		},
		{
			testName: "Failure: user already exists",
			request: SignUpRequest{
				Name:     name,
				Password: password,
			},
			requestBody:        `{"name":"testName","password":"password"}`,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   service.ErrUserAlreadyExists.Error() + "\n",
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().CreateUser(name, password).Return("", service.ErrUserAlreadyExists)
			},
		},
		{
			testName: "Failure: empty password",
			request: SignUpRequest{
				Name:     name,
				Password: emptyParameter,
			},
			requestBody:        `{"name":"testName","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": password\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: empty name",
			request: SignUpRequest{
				Name:     name,
				Password: emptyParameter,
			},
			requestBody:        `{"name":"","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": name\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: wrong name length",
			request: SignUpRequest{
				Name:     name + "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
				Password: password,
			},
			requestBody:        `{"name":"nameqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": wrong length\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: wrong password length",
			request: SignUpRequest{
				Name:     name,
				Password: password + "qqqqqqqqqqqqqqqqqqqqqqqqqqq",
			},
			requestBody:        `{"name":"name","password":"passwordqqqqqqqqqqqqqqqqqqqqqqqqqqq"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": wrong length\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: wrong name length",
			request: SignUpRequest{
				Name:     "n",
				Password: password,
			},
			requestBody:        `{"name":"n","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": wrong length\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: wrong password length",
			request: SignUpRequest{
				Name:     name,
				Password: "p",
			},
			requestBody:        `{"name":"name","password":"p"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": wrong length\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.request.Name, tc.request.Password)

			s := NewServer(auth, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(tc.requestBody))

			s.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestServer_LogIn(t *testing.T) {
	testTable := []struct {
		testName           string
		request            LogInRequest
		requestBody        string
		expectedStatusCode int
		expectedResponse   interface{}
		mock               func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName: "Success",
			request: LogInRequest{
				Name:     name,
				Password: password,
			},
			requestBody:        `{"name":"testName","password":"password"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   LogInResponse{AccessToken: "token", RefreshToken: "token"}.String(),
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return("token", "token", nil)
			},
		},
		{
			testName: "Failure: empty password",
			request: LogInRequest{
				Name:     name,
				Password: emptyParameter,
			},
			requestBody:        `{"name":"testName","password":""}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": password\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: empty name",
			request: LogInRequest{
				Name:     "",
				Password: password,
			},
			requestBody:        `{"name":"","password":"password"}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": name\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: user doesn't exists",
			request: LogInRequest{
				Name:     name + "aaaaaaaaa",
				Password: password,
			},
			requestBody:        `{"name":"testNameaaaaaaaaa","password":"password"}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   service.ErrNoUser.Error() + "\n",
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return("", "", service.ErrNoUser)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.request.Name, tc.request.Password)

			s := NewServer(auth, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(tc.requestBody))

			w.WriteHeader(tc.expectedStatusCode)

			s.Router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponse, w.Body.String())
		})
	}
}
