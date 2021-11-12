package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_GetRequests(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `SELECT id, userid, candidateid, status, created, updated 
			  FROM requests`

	testTable := []struct {
		testName        string
		inputID         string
		inputFilterType string
		expectedResult  []Request
		expectedError   error
		isErrorExpected bool
		mock            func(string, string)
	}{
		{
			testName:        "Success",
			inputID:         "1",
			inputFilterType: "ID",
			expectedResult: []Request{
				{
					ID:          "1",
					UserID:      "1",
					CandidateID: "1",
					Status:      "Accepted",
					Created:     time.Time{},
					Updated:     time.Time{},
				},
				{
					ID:          "2",
					UserID:      "2",
					CandidateID: "2",
					Status:      "Accepted",
					Created:     time.Time{},
					Updated:     time.Time{},
				},
			},
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1, s2 string) {
				rows := sqlmock.NewRows([]string{"id", "userid", "candidateid", "status", "created", "updated"}).
					AddRow("1", "1", "1", "Accepted", time.Time{}, time.Time{}).
					AddRow("2", "2", "2", "Accepted", time.Time{}, time.Time{})

				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(rows)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputID, tc.inputFilterType)

			requests, err := r.GetRequests(tc.inputID, tc.inputFilterType)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertRequests(t, requests, tc.expectedResult)
		})
	}
}

func TestRepository_AddCandidate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := "INSERT INTO candidates"

	testTable := []struct {
		testName        string
		inputName       string
		inputSurname    string
		inputFileID     string
		expectedResult  string
		expectedError   error
		isErrorExpected bool
		mock            func(string, string, string)
	}{
		{
			testName:        "Success",
			inputName:       "Grove",
			inputSurname:    "Street",
			inputFileID:     "1",
			expectedResult:  "1",
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1, s2, s3 string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).WithArgs(s1, s2, s3).WillReturnRows(rows)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputName, tc.inputSurname, tc.inputFileID)

			result, err := r.AddCandidate(tc.inputName, tc.inputSurname, tc.inputFileID)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertEquals(t, result, tc.expectedResult)
		})
	}
}

func TestRepository_UpdateRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `UPDATE requests SET status`

	testTable := []struct {
		testName        string
		inputID         string
		inputNewState   string
		expectedError   error
		isErrorExpected bool
		mock            func(string, string)
	}{
		{
			testName:        "Success",
			inputID:         "1",
			inputNewState:   "Accepted",
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1, s2 string) {
				mock.ExpectExec(query).WithArgs(s1, s2).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			testName:        "Failure: updating with no results",
			inputID:         "1",
			inputNewState:   "Accepted",
			expectedError:   ErrNoResult,
			isErrorExpected: true,
			mock: func(s1, s2 string) {
				mock.ExpectExec(query).WithArgs(s1, s2).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputNewState, tc.inputID)

			err := r.UpdateRequest(tc.inputID, tc.inputNewState)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}
		})
	}
}

func TestRepository_GetCVID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `SELECT cvosfileid FROM candidates`

	testTable := []struct {
		testName        string
		inputID         string
		expectedResult  string
		expectedError   error
		isErrorExpected bool
		mock            func(string)
	}{
		{
			testName:        "Success",
			inputID:         "1",
			expectedResult:  "1",
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).WithArgs(s).WillReturnRows(rows)
			},
		},
		{
			testName:        "Failure: invalid id",
			inputID:         "10000000",
			expectedError:   ErrNoFile,
			isErrorExpected: true,
			mock: func(s string) {
				mock.ExpectQuery(query).WithArgs(s).WillReturnError(ErrNoFile)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputID)

			id, err := r.GetCVID(tc.inputID)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertEquals(t, id, tc.expectedResult)
		})
	}
}

func TestRepository_CreateRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO requests`

	testTable := []struct {
		testName         string
		inputUserID      string
		inputCandidateID string
		expectedResult   string
		expectedError    error
		isErrorExpected  bool
		mock             func(string, string)
	}{
		{
			testName:         "Success",
			inputUserID:      "1",
			inputCandidateID: "1",
			expectedResult:   "1",
			expectedError:    nil,
			isErrorExpected:  false,
			mock: func(s1, s2 string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(rows)
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputUserID, tc.inputCandidateID)

			id, err := r.CreateRequest(tc.inputUserID, tc.inputCandidateID)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertEquals(t, id, tc.expectedResult)
		})
	}
}

func TestRepository_IsUserRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `SELECT id
			  FROM requests`

	testTable := []struct {
		testName        string
		inputUserID     string
		inputRequestID  string
		expectedResult  string
		expectedError   error
		isErrorExpected bool
		mock            func(string, string)
	}{
		{
			testName:        "Success",
			inputUserID:     "1",
			inputRequestID:  "1",
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1, s2 string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(rows)
			},
		},
		{
			testName:        "Failure: user not an admin",
			inputUserID:     "1",
			inputRequestID:  "1",
			expectedError:   ErrNoAccess,
			isErrorExpected: true,
			mock: func(s1, s2 string) {
				mock.ExpectQuery(query).WithArgs(s1, s2).WillReturnRows(&sqlmock.Rows{})
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputUserID, tc.inputRequestID)

			err := r.IsUserRequest(tc.inputUserID, tc.inputRequestID)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}
		})
	}
}

func TestRepository_IsUserAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	query := `SELECT isadmin
			  FROM users`

	testTable := []struct {
		testName        string
		inputUserID     string
		expectedResult  bool
		expectedError   error
		isErrorExpected bool
		mock            func(string)
	}{
		{
			testName:        "Success",
			inputUserID:     "1",
			expectedResult:  true,
			expectedError:   nil,
			isErrorExpected: false,
			mock: func(s1 string) {
				rows := sqlmock.NewRows([]string{"isadmin"}).AddRow("true")
				mock.ExpectQuery(query).WithArgs(s1).WillReturnRows(rows)
			},
		},
		{
			testName:        "Failre: user not an admin",
			inputUserID:     "1",
			expectedError:   ErrNoAccess,
			isErrorExpected: true,
			mock: func(s1 string) {
				mock.ExpectQuery(query).WithArgs(s1).WillReturnRows(&sqlmock.Rows{})
			},
		},
	}

	r := NewRepository(db)

	for _, tc := range testTable {
		t.Run(tc.testName, func(t *testing.T) {
			tc.mock(tc.inputUserID)

			result, err := r.IsUserAdmin(tc.inputUserID)
			if tc.isErrorExpected {
				assertError(t, err, tc.expectedError)
			} else {
				assertNoError(t, err)
			}

			assertBool(t, result, tc.expectedResult)
		})
	}
}

func assertBool(t *testing.T, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("Error: expected %v,got %v", want, got)
	}
}

func assertRequests(t *testing.T, got, want []Request) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error: fake requests")
	}
}
