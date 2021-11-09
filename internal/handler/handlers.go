package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/cyberdr0id/referral/internal/repository"
)

var (
	// errInvalidLength presents an error when user enters data with invalid length.
	errInvalidLength = errors.New("input parameters has invalid length")

	// errParameterRequired presents an error when user didn't fill some information.
	errParameterRequired = errors.New("input parameter is required")

	// errInvalidFile presents an error when user load file with invalid name or wrong extension.
	errInvalidFile = errors.New("invalid format or name of input file")

	// errInvalidName presents an error when user send candidate with invalid name/surname.
	errInvalidName = errors.New("input name didn't match to the desired format")

	// errInvalidName presents an error when user try to login with wrong password.
	errWrongPassword = errors.New("wrong password for inputed user")

	currentUserID = ""
)

// AuthRequest presents request for login.
type LogInRequest struct {
	Name     string
	Password string
}

// AuthRequest presents request for signup.
type SignUpRequest struct {
	Name     string
	Password string
}

// CandidateRequest presents request for sending candidate.
type CandidateSendingRequest struct {
	FileName         string
	CandidateName    string
	CandidateSurname string
}

// UserRequests type presents structure which contains all user requests.
type UserRequestsResponse struct {
	Requests []repository.Request `json:"requests"`
}

// CandidateResponse type presents candidate sending response.
type CandidateSendingResponse struct {
	CandidateID string `json:"candidateid"`
}

// DownloadResponse type presents a type which contains downloaded candidate cv.
type DownloadCVResponse struct {
	FileLink string `json:"filelink"`
}

// LogInResponse type presents structure of the log in response.
type LogInResponse struct {
	Token string `json:"token"`
}

// SignUpResponse type presents structure of the sign up response.
type SignUpResponse struct {
	ID string `json:"id"`
}

// UpdateResponse type presents message about success of updating.
type UpdateCandidateResponse struct {
	Message string `json:"message"`
}

// SignUp registers user.
func (s *Server) SignUp(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var request SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.ValidateSignUpRequest(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.Service.SignUp(request.Name, request.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(SignUpResponse{ID: id}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LogIn logs in user
func (s *Server) LogIn(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var request LogInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if request.Name == "" || request.Password == "" {
		http.Error(rw, errParameterRequired.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.Service.LogIn(request.Name, request.Password)
	if errors.Is(err, repository.ErrNoUser) {
		http.Error(rw, repository.ErrNoUser.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID = id

	if err := json.NewEncoder(rw).Encode(LogInResponse{Token: "something"}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SendCandidate sends candidate info and his cv.
func (s *Server) SendCandidate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var request CandidateSendingRequest

	file, header, err := r.FormFile("fileName")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	request.FileName = header.Filename
	request.CandidateName = r.FormValue("candidateName")
	request.CandidateSurname = r.FormValue("candidateSurname")

	if err := request.ValidateCandidateSendingRequest(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: adding file to object storage
	fileID := "1"

	id, err := s.Service.SendCandidate(request.CandidateName, request.CandidateSurname, fileID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(CandidateSendingResponse{CandidateID: id}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetRequests inputs all user requests.
func (s *Server) GetRequests(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	t := r.URL.Query().Get("type")

	userRequests, err := s.Service.GetRequests(currentUserID, t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(UserRequestsResponse{Requests: userRequests}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DownloadCV downloads CV of a particular candidate.
func (s *Server) DownloadCV(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	link, err := s.Service.DownloadCV(id)
	if errors.Is(err, repository.ErrNoFile) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(rw).Encode(DownloadCVResponse{FileLink: link}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateRequest updated status of request by id.
func (s *Server) UpdateRequest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	status := r.URL.Query().Get("status")
	requestId := r.URL.Query().Get("id")

	err := s.Service.UpdateRequest(requestId, status)
	if errors.Is(err, repository.ErrNoResult) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(rw).Encode(UpdateCandidateResponse{Message: "request update was successful"}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ValidateCandidateSendingRequest validates data after sending a candidate.
func (r *CandidateSendingRequest) ValidateCandidateSendingRequest() error {
	if len(r.CandidateName) == 0 || len(r.CandidateSurname) == 0 || len(r.FileName) == 0 {
		return errParameterRequired
	}
	fileExp := "([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$"
	isRightFile, _ := regexp.MatchString(fileExp, r.FileName)
	if !isRightFile {
		return errInvalidFile
	}

	nameExp := "(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?"
	isValidName, _ := regexp.MatchString(nameExp, r.CandidateName)
	if !isValidName {
		return errInvalidName
	}

	surnameExp := "(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?"
	isValidSurname, _ := regexp.MatchString(surnameExp, r.CandidateSurname)
	if !isValidSurname {
		return errInvalidName
	}

	return nil
}

// ValidateSignUpRequest validates data after signup.
func (r *SignUpRequest) ValidateSignUpRequest() error {
	if r.Name == "" || r.Password == "" {
		return errParameterRequired
	}

	if len(r.Name) < 6 || len(r.Name) > 18 {
		return errInvalidLength
	}

	if len(r.Password) < 6 || len(r.Password) > 18 {
		return errInvalidLength
	}

	return nil
}
