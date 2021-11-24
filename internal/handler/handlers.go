package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
)

var currentUserID string

// AuthRequest presents request for login.
type LogInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AuthRequest presents request for signup.
type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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
type DownloadResponse struct {
	FileLink string `json:"filelink"`
}

// LogInResponse type presents structure of the log in response.
type LogInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// SignUpResponse type presents structure of the sign up response.
type SignUpResponse struct {
	ID string `json:"id"`
}

// UpdateResponse type presents message about success of updating.
type UpdateCandidateResponse struct {
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, fmt.Errorf("can't write HTTP reply: %w", err).Error(), http.StatusInternalServerError)
	}
}

// SignUp registers user.
func (s *Server) SignUp(rw http.ResponseWriter, r *http.Request) {
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

	// pass, err := hash.HashPassword(request.Password)
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	id, err := s.Auth.SignUp(request.Name, request.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(rw, SignUpResponse{ID: id}, http.StatusOK)
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

	err := request.ValidateLogInRequest()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, refreshToken, err := s.Auth.LogIn(request.Name, request.Password)
	if errors.Is(err, service.ErrNoUser) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(rw, LogInResponse{AccessToken: accessToken, RefreshToken: refreshToken}, http.StatusOK)
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

	id, err := s.Referral.AddCandidate(request.CandidateName, request.CandidateSurname, fileID)
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
	typeExp := "^(id|userid|candidateid|created|updated|status)$"
	ok, err := regexp.MatchString(typeExp, t)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	userRequests, err := s.Referral.GetRequests(currentUserID, t)
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

	ok, err := ValidateID(id)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.Referral.GetCVID(id)
	if errors.Is(err, service.ErrNoFile) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: download id from storage - storage

	if err = json.NewEncoder(rw).Encode(DownloadResponse{FileLink: "example.com/path/to/file.extension"}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateRequest updated status of request by id.
func (s *Server) UpdateRequest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	state := r.URL.Query().Get("state")
	ok, err := ValidateFilterType(state)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	requestId := r.URL.Query().Get("id")
	ok, err = ValidateID(requestId)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Referral.UpdateRequest(requestId, state)
	if errors.Is(err, service.ErrNoResult) {
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
		return fmt.Errorf("%w: wrong length", service.ErrInvalidParameter)
	}
	fileExp := "([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$"
	isRightFile, _ := regexp.MatchString(fileExp, r.FileName)
	if !isRightFile {
		return fmt.Errorf("%w: invalid filename or filetype", service.ErrInvalidParameter)
	}

	nameSurnameExp := "(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?"
	isValid, _ := regexp.MatchString(nameSurnameExp, r.CandidateName)
	if !isValid {
		return fmt.Errorf("%w: name has invalid format", service.ErrInvalidParameter)
	}

	isValid, _ = regexp.MatchString(nameSurnameExp, r.CandidateSurname)
	if !isValid {
		return fmt.Errorf("%w: surname has invalid format", service.ErrInvalidParameter)
	}

	return nil
}

// ValidateSignUpRequest validates data after signup.
func (r *SignUpRequest) ValidateSignUpRequest() error {
	if r.Name == "" || r.Password == "" {
		return fmt.Errorf("%w: one of parameters is empty", service.ErrInvalidParameter)
	}

	if len(r.Name) < 6 || len(r.Name) > 18 {
		return fmt.Errorf("%w: wrong length", service.ErrInvalidParameter)
	}

	if len(r.Password) < 6 || len(r.Password) > 18 {
		return fmt.Errorf("%w: wrong length", service.ErrInvalidParameter)
	}

	return nil
}

// ValidateLogInRequest validates data after login.
func (r *LogInRequest) ValidateLogInRequest() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name", service.ErrInvalidParameter)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", service.ErrInvalidParameter)
	}

	return nil
}

func ValidateFilterType(state string) (bool, error) {
	stateExp := "^([Aa]ccepted|[Rr]ejected|[Ss]ubmitted)$"
	ok, err := regexp.MatchString(stateExp, state)
	if !ok {
		return ok, fmt.Errorf("%w: id has bad format", service.ErrInvalidParameter)
	}
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func ValidateID(id string) (bool, error) {
	idExp := "^[1-9]\\d*"
	ok, err := regexp.MatchString(idExp, id)
	if !ok {
		return ok, fmt.Errorf("%w: id has bad format", service.ErrInvalidParameter)
	}
	if err != nil {
		return ok, err
	}

	return ok, nil
}
