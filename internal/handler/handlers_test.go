package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	mock_service "github.com/cyberdr0id/referral/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestServer_LogIn(t *testing.T) {

	testTable := []struct {
		testName             string
		inputUserLogin       string
		inputUserPassword    string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
		mock                 func(s *mock_service.MockAuth, name, password string)
	}{
		{
			testName:          "Success",
			inputUserLogin:    "Alexander",
			inputUserPassword: "password",
			inputBody: `{
				"name":"Alexander",
				"password":"password"
			}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{
				"id":"75"
			}`,
			mock: func(s *mock_service.MockAuth, name, password string) {
				s.EXPECT().LogIn(name, password).Return("75", nil)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuth(ctrl)
			tc.mock(auth, tc.inputUserLogin, tc.inputUserPassword)

			s := Server{
				Auth: auth,
			}

			r := mux.NewRouter()
			r.HandleFunc("/user/login", s.LogIn).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(tc.inputBody))

			s.Router.ServeHTTP(w, req)

			assertStatusCode(t, w.Code, tc.expectedStatusCode)
			assertBody(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func assertBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got body \n %s \n want \n %s \n", got, want)
	}
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d status code, want %d", got, want)
	}
}
