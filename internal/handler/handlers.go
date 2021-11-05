package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/pkg/hash"
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

	// errInvalidName presents an error when user try to login with wrong password.
	errWrongPassword = errors.New("wrong password for inputed user")

	currentUserID = "1"
)

// AuthRequest presents request for authorization.
type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CandidateRequest presents request for sending candidate.
type CandidateRequest struct {
	FileName         string
	CandidateName    string
	CandidateSurname string
}

// UserRequests type presents structure which contains all user requests.
type UserRequests struct {
	Requests []repository.Request `json:"requests"`
}

// CandidateResponse type presents candidate sending response.
type CandidateResponse struct {
	CandidateID string `json:"candidateid"`
}

// DownloadResponse type presents a type which contains downloaded candidate cv.
type DownloadResponse struct {
	File []byte
}

// LogInResponse type presents structure of the log in response.
type LogInResponse struct {
	Token string `json:"token"`
}

// SignUpResponse type presents structure of the sign up response.
type SignUpResponse struct {
	ID string `json:"id"`
}

// SignUp registers user.
func (s *Server) SignUp(rw http.ResponseWriter, r *http.Request) {
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
	pass, _ := hash.HashPassword(request.Password)
	id, err := s.Repo.CreateUser(request.Name, pass)
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
	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(request)

	user, err := s.Repo.GetUser(request.Name)
	if errors.Is(err, repository.ErrNoUser) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := request.CheckLogInRequest(user); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID = user.ID

	if err := json.NewEncoder(rw).Encode(LogInResponse{Token: "something"}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SendCandidate sends candidate info and his cv.
func (s *Server) SendCandidate(rw http.ResponseWriter, r *http.Request) {
	var request CandidateRequest

	file, header, err := r.FormFile("fileName")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	request.FileName = header.Filename
	request.CandidateName = r.FormValue("candidateName")
	request.CandidateSurname = r.FormValue("candidateSurname")

	if err := request.CheckCandidateSendingRequest(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: adding file to object storage
	var fileID string = "1"

	id, err := s.Repo.AddCandidate(request.CandidateName, request.CandidateSurname, fileID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(CandidateResponse{CandidateID: id}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetRequests inputs all user requests.
func (s *Server) GetRequests(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	t := r.URL.Query().Get("type")
	userRequests, err := s.Repo.GetRequests(currentUserID, t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(UserRequests{Requests: userRequests}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DownloadCV downloads CV of a particular candidate.
func (s *Server) DownloadCV(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := s.Repo.GetCVID(id)
	if errors.Is(err, repository.ErrNoFile) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: download id from storage - storage

	if err = json.NewEncoder(rw).Encode(DownloadResponse{File: []byte{}}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CheckCandidateSendingRequest validates data after sending a candidate
func (r *CandidateRequest) CheckCandidateSendingRequest() error {
	if len(r.CandidateName) == 0 || len(r.CandidateSurname) == 0 || len(r.FileName) == 0 {
		return errParameterRequired
	}

	isRightFile, _ := regexp.MatchString("([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$", r.FileName)
	if !isRightFile {
		return errInvalidFile
	}

	isValidName, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", r.CandidateName)
	if !isValidName {
		return errInvalidName
	}

	isValidSurname, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", r.CandidateSurname)
	if !isValidSurname {
		return errInvalidName
	}

	return nil
}

// CheckSignUpRequest validates data after login/signup
func (r *AuthRequest) CheckSignUpRequest() error {
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

func (r *AuthRequest) CheckLogInRequest(user repository.User) error {
	if r.Name == "" || r.Password == "" {
		return errParameterRequired
	}

	if !hash.CheckPassowrdHash(r.Password, user.Password) {
		return errWrongPassword
	}

	return nil
}
