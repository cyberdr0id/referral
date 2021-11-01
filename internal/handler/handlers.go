package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
)

var (
	// errInvalidLength presents an error when user enters data with invalid length.
	errInvalidLength = errors.New("input parameter has invalid length")

	// errParameterRequired presents an error when user didn't fill some information.
	errParameterRequired = errors.New("input parameter is required")

	// errInvalidFile presents an error when user load file with invalid name or wrong extension.
	errInvalidFile = errors.New("invalid format or name of input file")

	// errInvalidName presents an error when user send candidate with invalid name/surname.
	errInvalidName = errors.New("input name didn't match to the desired format")
)

// AuthRequest presents request for authorization.
type AuthRequest struct {
	name     string `json:"name"`
	password string `json:"password"`
}

// CandidateSendingRequest presents request for sending candidate.
type CandidateRequest struct {
	fileName         string
	candidateName    string
	candidateSurname string
}

// Request presents model of request.
type Request struct {
	id          int       `json:"id"`
	userID      int       `json:"userid"`
	candidateID int       `json:"candidateid"`
	status      string    `json:"status"`
	created     time.Time `json:"created"`
	updated     time.Time `json:"updated"`
}

// UserRequests type presents structure which contains all user requests.
type UserRequests struct {
	requests []Request `json:"requests"`
}

// CandidateResponse type presents candidate sending response.
type CandidateResponse struct {
	candidateID int `json:"candidateid"`
}

// DownloadResponse type presents a type which contains downloaded candidate cv.
type DownloadResponse struct {
	file []byte
}

// LogInResponse type presents structure of the log in response.
type LogInResponse struct {
	token string `json:"token"`
}

// SignUpResponse type presents strcuture of the sign up response.
type SignUpResponse struct {
	id int `json:"id"`
}

// SignUp registers user.
func SignUp(rw http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.CheckSignUpRequest(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: database operations

	if err := json.NewEncoder(rw).Encode(SignUpResponse{id: 1}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LogIn logs in user
func LogIn(rw http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// TODO: validation, database operations etc.

	err = json.NewEncoder(rw).Encode(LogInResponse{token: "something"})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SendCandidate sends candidate info and his cv.
func SendCandidate(rw http.ResponseWriter, r *http.Request) {
	var request CandidateRequest

	file, header, err := r.FormFile("fileName")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	request.fileName = header.Filename
	request.candidateName = r.FormValue("candidateName")
	request.candidateSurname = r.FormValue("candidateSurname")

	if err := request.CheckCandidateSendingRequest(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: database interaction
	// TODO: file sending

	if err := json.NewEncoder(rw).Encode(CandidateResponse{candidateID: 1}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetRequests inputs all user requests.
func GetRequests(rw http.ResponseWriter, r *http.Request) {
	// if parameters count == 0
	// getting requests by user id - database
	// else if there is parameter type
	// getting filtered requests by user id - database
	err := json.NewEncoder(rw).Encode(UserRequests{requests: []Request{}})
	_ = err
}

// LoadCV downloads CV of a particular candidate.
func DownloadCV(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_ = id
	// check if entered id is valid - database
	// getting image by id - database
	// download id from storage - storage

}

// CheckCandidateSendingRequest validates data after sending a candidate
func (r *CandidateRequest) CheckCandidateSendingRequest() error {
	if len(r.candidateName) == 0 || len(r.candidateSurname) == 0 || len(r.fileName) == 0 {
		return errParameterRequired
	}

	isRightFile, _ := regexp.MatchString("([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$", r.fileName)
	if !isRightFile {
		return errInvalidFile
	}

	isValidName, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", r.candidateName)
	if !isValidName {
		return errInvalidName
	}

	isValidSurname, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", r.candidateSurname)
	if !isValidSurname {
		return errInvalidName
	}

	return nil
}

// CheckSignUpRequest validates data after login/signup
func (r *AuthRequest) CheckSignUpRequest() error {
	if r.name == "" || r.password == "" {
		return errParameterRequired
	}

	if len(r.name) < 6 || len(r.name) > 18 {
		return errInvalidLength
	}

	if len(r.password) < 6 || len(r.password) > 18 {
		return errInvalidLength
	}

	return nil
}
