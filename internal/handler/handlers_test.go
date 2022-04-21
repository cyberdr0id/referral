package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
	mock_service "github.com/cyberdr0id/referral/internal/service/mock"
	myjwt "github.com/cyberdr0id/referral/pkg/jwt"
	mylog "github.com/cyberdr0id/referral/pkg/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: handle 500 error login/signup

const (
	defaultName         = "testName"
	defaultPassword     = "password"
	defaultID           = "1"
	defaultUserID       = "1"
	defaultIsAdmin      = false
	longName            = "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnn"
	longPassword        = "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnn"
	nonExistentUserName = "abcedefg"
	shortName           = "n"
	shortPassword       = "p"
	emptyParameter      = ""
	token               = "token"

	defaultJWT             = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	defaultAuthHeaderValue = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	defaultAuthHeaderKey   = "Authorization"
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
		isErrorExpected       bool
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
			isErrorExpected:       false,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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
			isErrorExpected:    true,
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

			if tc.isErrorExpected {
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

func TestServer_SendCandidate(t *testing.T) {
	const (
		defaultName    = "name"
		defaultSurname = "surname"
	)
	testTable := []struct {
		testName           string
		filename           string
		requestBody        service.CandidateSubmittingRequest
		serviceData        service.CandidateSubmittingRequest
		expectedResponse   CandidateSubmittingResponse
		expectedStatusCode int
		isErrorExpected    bool
		errorResponse      ErrorResponse
		mock               func(s *mock_service.MockReferral, r service.CandidateSubmittingRequest)
	}{
		{
			testName:    "Success",
			filename:    "file.pdf",
			requestBody: service.CandidateSubmittingRequest{},
			serviceData: service.CandidateSubmittingRequest{},
			expectedResponse: CandidateSubmittingResponse{
				CandidateID: defaultID,
			},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, r service.CandidateSubmittingRequest) {
				s.EXPECT().AddCandidate(gomock.Any(), r).Return(defaultID, nil)
			},
		},
	}

	for _, tc := range testTable {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMock := func(s *mock_service.MockAuth, token string) {
			claims := &myjwt.Claims{}
			claims.Subject = defaultUserID
			claims.IsAdmin = false

			s.EXPECT().ParseToken(token).Return(claims, nil)
		}

		referral := mock_service.NewMockReferral(ctrl)

		auth := mock_service.NewMockAuth(ctrl)
		authMock(auth, defaultJWT)

		logger, err := mylog.NewLogger()
		if err != nil {
			t.Fatalf("error with logger creating: %s", err.Error())
		}

		s := NewServer(auth, referral, logger)

		w := httptest.NewRecorder()

		file, err := os.Create("./file.pdf")
		require.NoError(t, err)
		defer file.Close()
		defer os.Remove(file.Name())

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", file.Name())
		require.NoError(t, err)

		_, _ = io.Copy(part, file)

		_ = writer.WriteField("candidateName", defaultName)
		_ = writer.WriteField("candidateSurname", defaultSurname)

		err = writer.Close()
		require.NoError(t, err)

		fileStat, err := file.Stat()
		if err != nil {
			t.Fatalf("something wrong with file.Stat(): %s", err.Error())
		}

		fh := &multipart.FileHeader{
			Filename: fileStat.Name(),
			Size:     fileStat.Size(),
		}
		t.Log(fh.Filename)

		req, _ := http.NewRequest("POST", "/references", body)
		req.Header.Set(defaultAuthHeaderKey, defaultAuthHeaderValue)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		f, err := fh.Open()
		if err != nil {
			t.Fatalf("cannot open file: %s", err.Error())
		}

		jija := service.CandidateSubmittingRequest{
			File:             f,
			CandidateName:    defaultName,
			CandidateSurname: defaultSurname,
		}

		tc.serviceData = jija
		tc.requestBody = jija

		tc.mock(referral, tc.serviceData)

		s.Router.ServeHTTP(w, req)

		if tc.isErrorExpected {
			var response ErrorResponse
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, tc.errorResponse, response)
		} else {
			var response CandidateSubmittingResponse
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedResponse, response)
		}

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	}
}

func TestServer_GetRequests(t *testing.T) {
	const (
		defaultUserID     = "1"
		defaultPageNumber = "1"
		defaultPageSize   = "10"

		defaultPageNumberInt = 1
		defaultPageSizeInt   = 10

		defaultStatus = "submitted"

		pageNumberParameter = "page"
		pageSizeParameter   = "size"
		statusParameter     = "status"
		idParameter         = "user_id"

		invalidStatus           = "q"
		invalidNumericParameter = "q"
		emptyStatus             = ""
	)

	var (
		invalidStatusMessage     = ErrInvalidParameter.Error() + ": request status"
		invalidPageNumberMessage = ErrInvalidParameter.Error() + ": page number has bad format"
		invalidPageSizeMesssage  = ErrInvalidParameter.Error() + ": page size has bad format"
	)

	testTable := []struct {
		testName           string
		params             map[string]string
		serviceUserID      string
		serviceStatus      string
		servicePageSize    int
		servicePageNumber  int
		expectedResponse   []repository.UserRequests
		isErrorExpected    bool
		errorResponse      ErrorResponse
		expectedStatusCode int
		mock               func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int)
	}{
		{
			testName: "Success",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				pageSizeParameter:   defaultPageSize,
				statusParameter:     defaultStatus,
			},
			serviceUserID:      defaultUserID,
			serviceStatus:      defaultStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return([]repository.UserRequests{}, nil)
			},
		},
		{
			testName: "Success, no page number in request",
			params: map[string]string{
				pageSizeParameter: defaultPageSize,
				statusParameter:   defaultStatus,
			},
			serviceUserID:      defaultUserID,
			serviceStatus:      defaultStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return([]repository.UserRequests{}, nil)
			},
		},
		{
			testName: "Success, no page size in request",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				statusParameter:     defaultStatus,
			},
			serviceUserID:      defaultUserID,
			serviceStatus:      defaultStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return([]repository.UserRequests{}, nil)
			},
		},
		{
			testName: "Success, no status in request",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				pageSizeParameter:   defaultPageSize,
			},
			serviceUserID:      defaultUserID,
			serviceStatus:      emptyStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return([]repository.UserRequests{}, nil)
			},
		},
		{
			testName:           "Success, no request parameters",
			params:             map[string]string{},
			serviceUserID:      defaultUserID,
			serviceStatus:      emptyStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return([]repository.UserRequests{}, nil)
			},
		},
		{
			testName: "Failure: invalid request status",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				pageSizeParameter:   defaultPageSize,
				statusParameter:     invalidStatus,
			},
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidStatusMessage},
			expectedStatusCode: http.StatusBadRequest,
			mock:               func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {},
		},
		{
			testName: "Failure: invalid page number",
			params: map[string]string{
				pageNumberParameter: invalidNumericParameter,
				pageSizeParameter:   defaultPageSize,
				statusParameter:     defaultStatus,
			},
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidPageNumberMessage},
			expectedStatusCode: http.StatusBadRequest,
			mock:               func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {},
		},
		{
			testName: "Failure: invalid page size",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				pageSizeParameter:   invalidNumericParameter,
				statusParameter:     defaultStatus,
			},
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidPageSizeMesssage},
			expectedStatusCode: http.StatusBadRequest,
			mock:               func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {},
		},
		{
			testName: "Failure: internal server error",
			params: map[string]string{
				pageNumberParameter: defaultPageNumber,
				pageSizeParameter:   defaultPageSize,
				statusParameter:     defaultStatus,
			},
			serviceUserID:      defaultUserID,
			serviceStatus:      defaultStatus,
			servicePageSize:    defaultPageSizeInt,
			servicePageNumber:  defaultPageNumberInt,
			expectedResponse:   []repository.UserRequests{},
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: errInternalServerError.Error()},
			expectedStatusCode: http.StatusInternalServerError,
			mock: func(s *mock_service.MockReferral, userID, status string, pageSize, pageNumber int) {
				s.EXPECT().GetRequests(userID, status, pageNumber, pageSize).Return(nil, errInternalServerError)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := func(s *mock_service.MockAuth, token string) {
				claims := &myjwt.Claims{}
				claims.Subject = defaultUserID
				claims.IsAdmin = false

				s.EXPECT().ParseToken(token).Return(claims, nil)
			}

			referral := mock_service.NewMockReferral(ctrl)
			tc.mock(referral, tc.serviceUserID, tc.serviceStatus, tc.servicePageSize, tc.servicePageNumber)

			auth := mock_service.NewMockAuth(ctrl)
			authMock(auth, defaultJWT)

			logger, err := mylog.NewLogger()
			if err != nil {
				t.Fatalf("error with logger creating: %s", err.Error())
			}

			s := NewServer(auth, referral, logger)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/references", nil)
			req.Header.Set(defaultAuthHeaderKey, defaultAuthHeaderValue)

			query := req.URL.Query()

			for k, v := range tc.params {
				query.Add(k, v)
			}

			req.URL.RawQuery = query.Encode()

			s.Router.ServeHTTP(w, req)

			if tc.isErrorExpected {
				var response ErrorResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.errorResponse, response)
			} else {
				var response []repository.UserRequests
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedResponse, response)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_DownloadCV(t *testing.T) {
	const (
		invalidCVID = "qw"

		successDownloadMessage = "All right! Check 'Downloads' folder."
	)

	var (
		invalidCVIDMesage = ErrInvalidParameter.Error() + ": id has bad format"
	)

	testTable := []struct {
		testName           string
		CVID               string
		userID             string
		serviceCVID        string
		serviceUserID      string
		expectedResponse   DownloadResponse
		expectedStatusCode int
		isErrorExpected    bool
		errorResponse      ErrorResponse
		mock               func(s *mock_service.MockReferral, id, userID string)
	}{
		{
			testName:           "Success",
			CVID:               defaultID,
			userID:             defaultUserID,
			serviceCVID:        defaultID,
			serviceUserID:      defaultUserID,
			expectedResponse:   DownloadResponse{Message: successDownloadMessage},
			expectedStatusCode: http.StatusOK,
			isErrorExpected:    false,
			errorResponse:      ErrorResponse{},
			mock: func(s *mock_service.MockReferral, id, userID string) {
				s.EXPECT().DownloadFile(gomock.Any(), id, userID).Return(nil)
			},
		},
		{
			testName:           "Failure: invalid id",
			CVID:               invalidCVID,
			userID:             defaultUserID,
			expectedStatusCode: http.StatusBadRequest,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidCVIDMesage},
			mock:               func(s *mock_service.MockReferral, id, userID string) {},
		},
		{
			testName:           "Failure: internal server error",
			CVID:               defaultID,
			userID:             defaultUserID,
			serviceCVID:        defaultID,
			serviceUserID:      defaultUserID,
			expectedStatusCode: http.StatusInternalServerError,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: errInternalServerError.Error()},
			mock: func(s *mock_service.MockReferral, id, userID string) {
				s.EXPECT().DownloadFile(gomock.Any(), id, userID).Return(errInternalServerError)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := func(s *mock_service.MockAuth, token string) {
				claims := &myjwt.Claims{}
				claims.Subject = defaultUserID
				claims.IsAdmin = false

				s.EXPECT().ParseToken(token).Return(claims, nil)
			}

			referral := mock_service.NewMockReferral(ctrl)
			tc.mock(referral, tc.serviceCVID, tc.serviceUserID)

			auth := mock_service.NewMockAuth(ctrl)
			authMock(auth, defaultJWT)

			logger, err := mylog.NewLogger()
			if err != nil {
				t.Fatalf("error with logger creating: %s", err.Error())
			}

			s := NewServer(auth, referral, logger)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/cvs", nil)
			req.Header.Set(defaultAuthHeaderKey, defaultAuthHeaderValue)

			query := req.URL.Query()

			query.Add("id", tc.CVID)

			req.URL.RawQuery = query.Encode()

			s.Router.ServeHTTP(w, req)

			if tc.isErrorExpected {
				var response ErrorResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.errorResponse, response)
			} else {
				var response DownloadResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedResponse, response)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_UpdateRequest(t *testing.T) {
	const (
		newStatus        = "accepted"
		nonExistentID    = "100"
		invalidParameter = "q"
	)

	var (
		updateResponseMessage      = fmt.Sprintf("request status with %s ID has been updated", defaultID)
		invalidStatusMessage       = ErrInvalidParameter.Error() + ": request status"
		invalidIDMessage           = ErrInvalidParameter.Error() + ": id has bad format"
		requestDoesntExistsMessage = service.ErrNoResult.Error()
		internalServerErrorMessage = errInternalServerError.Error()
	)

	testTable := []struct {
		testName           string
		updateRequest      UpdateRequest
		serviceID          string
		serviceStatus      string
		expectedResponse   UpdateResponse
		expectedStatusCode int
		isErrorExpected    bool
		errorResponse      ErrorResponse
		mock               func(s *mock_service.MockReferral, id, status string)
	}{
		{
			testName: "Success",
			updateRequest: UpdateRequest{
				ID:        defaultID,
				NewStatus: newStatus,
			},
			serviceID:          defaultID,
			serviceStatus:      newStatus,
			expectedResponse:   UpdateResponse{Message: updateResponseMessage},
			expectedStatusCode: http.StatusOK,
			mock: func(s *mock_service.MockReferral, id, status string) {
				s.EXPECT().UpdateRequest(id, status).Return(nil)
			},
		},
		{
			testName: "Failure: invalid request status",
			updateRequest: UpdateRequest{
				ID:        defaultID,
				NewStatus: invalidParameter,
			},
			expectedStatusCode: http.StatusBadRequest,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidStatusMessage},
			mock:               func(s *mock_service.MockReferral, id, status string) {},
		},
		{
			testName: "Failure: invalid user id",
			updateRequest: UpdateRequest{
				ID:        invalidParameter,
				NewStatus: newStatus,
			},
			expectedStatusCode: http.StatusBadRequest,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: invalidIDMessage},
			mock:               func(s *mock_service.MockReferral, id, status string) {},
		},
		{
			testName: "Failure: request doesn't exists",
			updateRequest: UpdateRequest{
				ID:        nonExistentID,
				NewStatus: newStatus,
			},
			serviceID:          nonExistentID,
			serviceStatus:      newStatus,
			expectedStatusCode: http.StatusBadRequest,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: requestDoesntExistsMessage},
			mock: func(s *mock_service.MockReferral, id, status string) {
				s.EXPECT().UpdateRequest(id, status).Return(service.ErrNoResult)
			},
		},
		{
			testName: "Failure: internal server error",
			updateRequest: UpdateRequest{
				ID:        defaultID,
				NewStatus: newStatus,
			},
			serviceID:          defaultID,
			serviceStatus:      newStatus,
			expectedStatusCode: http.StatusInternalServerError,
			isErrorExpected:    true,
			errorResponse:      ErrorResponse{Message: internalServerErrorMessage},
			mock: func(s *mock_service.MockReferral, id, status string) {
				s.EXPECT().UpdateRequest(id, status).Return(errInternalServerError)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := func(s *mock_service.MockAuth, token string) {
				claims := &myjwt.Claims{}
				claims.Subject = defaultUserID
				claims.IsAdmin = true

				// it works
				s.EXPECT().ParseToken(token).Return(claims, nil)
				s.EXPECT().ParseToken(token).Return(claims, nil)
			}

			referral := mock_service.NewMockReferral(ctrl)
			tc.mock(referral, tc.serviceID, tc.serviceStatus)

			auth := mock_service.NewMockAuth(ctrl)
			authMock(auth, defaultJWT)

			logger, err := mylog.NewLogger()
			if err != nil {
				t.Fatalf("error with logger creating: %s", err.Error())
			}

			s := NewServer(auth, referral, logger)

			w := httptest.NewRecorder()

			request, _ := json.Marshal(tc.updateRequest)

			req, _ := http.NewRequest("PUT", "/admin/references", bytes.NewBuffer(request))
			req.Header.Set(defaultAuthHeaderKey, defaultAuthHeaderValue)

			s.Router.ServeHTTP(w, req)

			if tc.isErrorExpected {
				var response ErrorResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.errorResponse, response)
			} else {
				var response UpdateResponse
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tc.expectedResponse, response)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
