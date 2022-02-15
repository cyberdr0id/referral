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
	mylog "github.com/cyberdr0id/referral/pkg/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TODO: handle 500 error login/signup

const (
	defaultName         = "testName"
	defaultPassword     = "password"
	defaultID           = "1"
	longName            = "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnn"
	longPassword        = "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnn"
	nonExistentUserName = "abcedefg"
	shortName           = "n"
	shortPassword       = "p"
	emptyParameter      = ""
	token               = "token"
)

var (
	userAlreadyExistsMessage     = service.ErrUserAlreadyExists.Error()
	invalidNameMessage           = ErrInvalidParameter.Error() + ": name"
	invalidPasswordMessage       = ErrInvalidParameter.Error() + ": password"
	invalidNameLengthMessage     = ErrInvalidParameter.Error() + ": name must be between 6 and 18 symbols"
	invalidPasswordLengthMessage = ErrInvalidParameter.Error() + ": password must be between 6 and 18 symbols"

	errInternalServerError = errors.New("internal server error")
)

func TestServer_SignUp(t *testing.T) {
	testTable := []struct {
		testName              string
		serviceName           string
		servicePassword       string
		requestBody           SignUpRequest
		expectedStatusCode    int
		expectedResponse      SignUpResponse
		isErrorExpeced        bool
		expectedErrorResponse ErrorResponse
		mock                  func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName:        "Success: status 201",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode:    http.StatusCreated,
			expectedResponse:      SignUpResponse{ID: defaultID},
			isErrorExpeced:        false,
			expectedErrorResponse: ErrorResponse{},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().SignUp(name, password).Return(defaultID, nil)
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
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: userAlreadyExistsMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().SignUp(name, password).Return("", service.ErrUserAlreadyExists)
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
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidPasswordMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: empty name, status 400",
			serviceName:     emptyParameter,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     emptyParameter,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidNameMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: wrong name length(too long), status 400",
			serviceName:     longName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     longName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidNameLengthMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: wrong password length(too long), status 400",
			serviceName:     defaultName,
			servicePassword: longPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: longPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidPasswordLengthMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: wrong name length(too small), status 400",
			serviceName:     shortName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     shortName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidNameLengthMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: wrong password length(too small), status 400",
			serviceName:     defaultName,
			servicePassword: shortPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: shortPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidPasswordLengthMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: internal server error, status 500",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: SignUpRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   SignUpResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: errInternalServerError.Error(),
			},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().SignUp(name, password).Return("", errInternalServerError)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.serviceName, tc.servicePassword)

			logger, err := mylog.NewLogger()
			if err != nil {
				t.Fatalf("error with logger creating: %s", err.Error())
			}

			s := NewServer(auth, nil, logger)

			w := httptest.NewRecorder()

			request, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(request))

			s.Router.ServeHTTP(w, req)

			if tc.isErrorExpeced {
				var response ErrorResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedErrorResponse, response)
			} else {
				var response SignUpResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedResponse, response)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_LogIn(t *testing.T) {
	testTable := []struct {
		testName              string
		serviceName           string
		servicePassword       string
		requestBody           LogInRequest
		expectedStatusCode    int
		expectedResponse      LogInResponse
		isErrorExpeced        bool
		expectedErrorResponse ErrorResponse
		mock                  func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName:        "Success: status 200",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: LogInRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: LogInResponse{
				Token: token,
			},
			isErrorExpeced:        false,
			expectedErrorResponse: ErrorResponse{},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return(token, nil)
			},
		},
		{
			testName:        "Failure: empty password, status 401",
			serviceName:     defaultName,
			servicePassword: emptyParameter,
			requestBody: LogInRequest{
				Name:     defaultName,
				Password: emptyParameter,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   LogInResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidPasswordMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: empty name, status 401",
			serviceName:     emptyParameter,
			servicePassword: defaultPassword,
			requestBody: LogInRequest{
				Name:     emptyParameter,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   LogInResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: invalidNameMessage,
			},
			mock: func(s *mock_service.MockAuth, name, password string) {},
		},
		{
			testName:        "Failure: user doesn't exists, status 401",
			serviceName:     nonExistentUserName,
			servicePassword: defaultPassword,
			requestBody: LogInRequest{
				Name:     nonExistentUserName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   LogInResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: service.ErrNoUser.Error(),
			},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return(emptyParameter, service.ErrNoUser)
			},
		},
		{
			testName:        "Failure: wrong password for existent user, status 401",
			serviceName:     defaultName,
			servicePassword: shortPassword,
			requestBody: LogInRequest{
				Name:     defaultName,
				Password: shortPassword,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   LogInResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: service.ErrNoUser.Error(),
			},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return(emptyParameter, service.ErrNoUser)
			},
		},
		{
			testName:        "Failure: internal server error, status 500",
			serviceName:     defaultName,
			servicePassword: defaultPassword,
			requestBody: LogInRequest{
				Name:     defaultName,
				Password: defaultPassword,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   LogInResponse{},
			isErrorExpeced:     true,
			expectedErrorResponse: ErrorResponse{
				Message: errInternalServerError.Error(),
			},
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return(token, errInternalServerError)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.serviceName, tc.servicePassword)

			logger, err := mylog.NewLogger()
			if err != nil {
				t.Fatalf("error with logger creating: %s", err.Error())
			}

			s := NewServer(auth, nil, logger)

			w := httptest.NewRecorder()

			request, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(request))

			s.Router.ServeHTTP(w, req)

			if tc.isErrorExpeced {
				var response ErrorResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedErrorResponse, response)
			} else {
				var response LogInResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedResponse, response)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
