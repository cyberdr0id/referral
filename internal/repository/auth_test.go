package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := "INSERT INTO users"

	testTable := []struct {
		testName        string
		inputName       string
		inputPassword   string
		expectedResult  string
		expectedError   error
		isErrorExpected bool
		mock            func(string, string)
	}{
		{
			testName:        "Success",
			inputName:       "nameee",
			inputPassword:   "password",
			expectedResult:  "1",
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1, s2 string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(rows)
			},
		},
		// TODO: handle error when user already exists
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputName, tc.inputPassword)

			result, err := r.CreateUser(tc.inputName, tc.inputPassword)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertEquals(t, result, tc.expectedResult)
		})
	}
}

func TestRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `SELECT id, name, password, isadmin, created, updated 
			  FROM users`

	testTable := []struct {
		testName        string
		inputName       string
		expectedResult  User
		expectedError   error
		isErrorExpected bool
		mock            func(string)
	}{
		{
			testName:  "Success",
			inputName: "CJ",
			expectedResult: User{
				ID:       "1",
				Name:     "CJ",
				Password: "grove-street",
				IsAdmin:  false,
				Created:  time.Time{},
				Updated:  time.Time{},
			},
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s string) {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "isadmin", "created", "updated"}).
					AddRow("1", "CJ", "grove-street", false, time.Time{}, time.Time{})

				mock.ExpectQuery(query).WithArgs(s).WillReturnRows(rows)
			},
		},
		{
			testName:        "Failure: fake username",
			inputName:       "DJ",
			expectedError:   ErrNoUser,
			isErrorExpected: true,
			mock: func(s string) {
				mock.ExpectQuery(query).WithArgs(s).WillReturnError(ErrNoUser)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputName)

			user, err := r.GetUser(tc.inputName)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertEquals(t, user.ID, tc.expectedResult.ID)
		})
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("Error: expected %v, got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Error: nothing expected, but got %v", err)
	}
}

func assertEquals(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("Error: got %v, want %v", got, want)
	}
}
