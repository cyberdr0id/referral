package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	idColumn = "id"

	defaultName     = "Denzel"
	defaultPassword = "password"
	defaultID       = "1"
	emptyResult     = ""
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
		isErrorExpected bool
		expectedError   error
		mock            func(string, string)
	}{
		{
			testName:        "Success",
			inputName:       defaultName,
			inputPassword:   defaultPassword,
			expectedResult:  defaultID,
			isErrorExpected: false,
			expectedError:   nil,
			mock: func(s1, s2 string) {
				rows := sqlmock.NewRows([]string{idColumn}).AddRow(defaultID)
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(rows)
			},
		},
		{
			testName:        "Failure: user already exists",
			inputName:       defaultName,
			inputPassword:   defaultPassword,
			expectedResult:  emptyResult,
			isErrorExpected: true,
			expectedError:   ErrUserAlreadyExists,
			mock: func(s1, s2 string) {
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnError(ErrUserAlreadyExists)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputName, tc.inputPassword)

			result, err := r.CreateUser(tc.inputName, tc.inputPassword)
			if tc.isErrorExpected {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, result, tc.expectedResult)
		})
	}
}
