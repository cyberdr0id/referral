package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
	defaultName     = "testName"
	defaultPassword = "password"
	defaultID       = "1"

	emptyParameter = ""
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func TestServer_SignUp(t *testing.T) {
	testTable := []struct {
		testName           string
		serviceName        string
		servicePassword    string
		requestBody        SignUpRequest
		expectedStatusCode int
		expectedResponse   SignUpResponse
		isErrorExpeced     bool
		errorResponse      ErrorResponse
		mock               func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName:        "Success: status 201",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   SignUpResponse{ID: defaultID},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().CreateUser(name, password).Return(defaultID, nil)
			},
		},
		{
			testName:        "Failure: user already exists, status 409",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   SignUpResponse{},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().CreateUser(name, password).Return("", ErrUserAlreadyExists)
			},
		},
		{
			testName:        "Failure: empty password, status 400",
			serviceName:     defaultName,
			servicePassword: emptyParameter,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: emptyParameter,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: empty name, status 400",
			request: SignUpRequest{
				Name:     defaultName,
				Password: emptyParameter,
			},
			requestBody:        `{"name":"","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": name\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: wrong name length(too long), status 400",
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
			testName: "Failure: wrong password length(too long), status 400",
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
			testName: "Failure: wrong name length(too small), status 400",
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
			testName: "Failure: wrong password length(too small), status 400",
			request: SignUpRequest{
				Name:     name,
				Password: "p",
			},
			requestBody:        `{"name":"name","password":"p"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   service.ErrInvalidParameter.Error() + ": wrong length\n",
			mock:               func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName: "Failure: internal server error, status 500",
			request: SignUpRequest{
				Name:     name,
				Password: password,
			},
			requestBody:        `{"name":"testName","password":"password"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   errors.New("internal server error").Error() + "\n",
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().CreateUser(name, password).Return("1", errors.New("internal server error"))
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.serviceName, tc.servicePassword)

			s := NewServer(auth, nil)

			w := httptest.NewRecorder()

			request, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(request))

			s.Router.ServeHTTP(w, req)

			var response SignUpResponse
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}

// func TestServer_LogIn(t *testing.T) {
// 	testTable := []struct {
// 		testName           string
// 		request            LogInRequest
// 		requestBody        string
// 		expectedStatusCode int
// 		expectedResponse   interface{}
// 		mock               func(s *mock_service.MockAuth, name, password string)
// 	}{
// 		{
// 			testName: "Success: status 200",
// 			request: LogInRequest{
// 				Name:     name,
// 				Password: password,
// 			},
// 			requestBody:        `{"name":"testName","password":"password"}`,
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   LogInResponse{AccessToken: "token", RefreshToken: "token"}.String(),
// 			mock: func(s *mock_service.MockAuth, name, password string) {
// 				s.EXPECT().LogIn(name, password).Return("token", "token", nil)
// 			},
// 		},
// 		{
// 			testName: "Failure: empty password, status 401",
// 			request: LogInRequest{
// 				Name:     name,
// 				Password: emptyParameter,
// 			},
// 			requestBody:        `{"name":"testName","password":""}`,
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedResponse:   service.ErrInvalidParameter.Error() + ": password\n",
// 			mock:               func(s *mock_service.MockAuth, name, password string) {},
// 		},
// 		{
// 			testName: "Failure: empty name, status 401",
// 			request: LogInRequest{
// 				Name:     "",
// 				Password: password,
// 			},
// 			requestBody:        `{"name":"","password":"password"}`,
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedResponse:   service.ErrInvalidParameter.Error() + ": name\n",
// 			mock:               func(s *mock_service.MockAuth, name, password string) {},
// 		},
// 		{
// 			testName: "Failure: user doesn't exists, status 401",
// 			request: LogInRequest{
// 				Name:     name + "aaaaaaaaa",
// 				Password: password,
// 			},
// 			requestBody:        `{"name":"testNameaaaaaaaaa","password":"password"}`,
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedResponse:   service.ErrNoUser.Error() + "\n",
// 			mock: func(s *mock_service.MockAuth, name, password string) {
// 				s.EXPECT().LogIn(name, password).Return("", "", service.ErrNoUser)
// 			},
// 		},
// 		{
// 			testName: "Failure: internal server error, status 500",
// 			request: LogInRequest{
// 				Name:     name,
// 				Password: password,
// 			},
// 			requestBody:        `{"name":"testName","password":"password"}`,
// 			expectedStatusCode: http.StatusInternalServerError,
// 			expectedResponse:   errors.New("internal server error").Error() + "\n",
// 			mock: func(s *mock_service.MockAuth, name, password string) {
// 				s.EXPECT().LogIn(name, password).Return("token", "token", errors.New("internal server error"))
// 			},
// 		},
// 	}

// 	for _, tc := range testTable {
// 		t.Run(tc.testName, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			auth := mock_service.NewMockAuth(ctrl)
// 			tc.mock(auth, tc.request.Name, tc.request.Password)

// 			s := NewServer(auth, nil)

// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(tc.requestBody))

// 			s.Router.ServeHTTP(w, req)

// 			assert.Equal(t, tc.expectedStatusCode, w.Code)
// 			assert.Equal(t, tc.expectedResponse, w.Body.String())
// 		})
// 	}
// }
