package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/cyberdr0id/referral/internal/context"
	"github.com/cyberdr0id/referral/internal/service"
)

const (
	filenameParam         = "fileName"
	candidateNameParam    = "candidateName"
	candidateSurnameParam = "candidateSurname"
	statusParameter       = "status"
	idParameter           = "id"
	pageNumberParameter   = "page"
	pageSizeParameter     = "size"
	allParameter          = "all"
	userIDParameter       = "user_id"

	defaultPageNumber = 1
	defaultPageSize   = 10
)

// LogInRequest presents request for login.
type LogInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SignUpRequest presents request for signup.
type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CandidateSendingResponse type presents candidate sending response.
type CandidateSendingResponse struct {
	CandidateID string `json:"candidateid"`
}

// DownloadResponse type presents a type which contains downloaded candidate cv.
type DownloadResponse struct {
	FileLink string `json:"filelink"`
}

// LogInResponse type presents structure of the log in response.
type LogInResponse struct {
	AccessToken string `json:"accessToken"`
}

// SignUpResponse type presents structure of the sign up response.
type SignUpResponse struct {
	ID string `json:"id"`
}

// ErrorResponse presents a custom error type.
type ErrorResponse struct {
	Message string `json:"error"`
}

// ErrInvalidParameter presetns an error when user input invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SignUp registers user.
func (s *Server) SignUp(rw http.ResponseWriter, r *http.Request) {
	var request SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.ValidateSignUpRequest(); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := s.Auth.SignUp(request.Name, request.Password)
	if errors.Is(err, service.ErrUserAlreadyExists) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusConflict)
		return
	}
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, SignUpResponse{ID: id}, http.StatusCreated)
}

// LogIn logs in user
func (s *Server) LogIn(rw http.ResponseWriter, r *http.Request) {
	var request LogInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := request.ValidateLogInRequest()
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}

	accessToken, err := s.Auth.LogIn(request.Name, request.Password)
	if errors.Is(err, service.ErrNoUser) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, LogInResponse{AccessToken: accessToken}, http.StatusOK)
}

// SendCandidate sends candidate info and his cv.
func (s *Server) SendCandidate(rw http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(filenameParam)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	request := service.SubmitCandidateRequest{
		File:             file,
		CandidateName:    r.FormValue(candidateNameParam),
		CandidateSurname: r.FormValue(candidateSurnameParam),
	}

	if err := ValidateCandidateSendingRequest(request); err != nil {
		sendResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.Referral.AddCandidate(r.Context(), request)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(rw, CandidateSendingResponse{CandidateID: id}, http.StatusOK)
}

// GetRequests inputs all user requests.
func (s *Server) GetRequests(rw http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get(statusParameter)
	pageNumber := r.URL.Query().Get(pageNumberParameter)
	pageSize := r.URL.Query().Get(pageSizeParameter)

	userID, ok := context.GetUserID(r.Context())
	if !ok {
		sendResponse(rw, fmt.Errorf("cannot get user id from context"), http.StatusInternalServerError)
		return
	}

	pn, ps, err := ValidateGetRequestsRequest(t, pageNumber, pageSize, userID)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	userRequests, err := s.Referral.GetRequests(userID, t, pn, ps)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(rw, userRequests, http.StatusOK)
}

// GetAllRequests admin handler that returns list of all requests.
func (s *Server) GetAllRequests(rw http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get(statusParameter)
	pageNumber := r.URL.Query().Get(pageNumberParameter)
	pageSize := r.URL.Query().Get(pageSizeParameter)
	userID := r.URL.Query().Get(userIDParameter)

	pn, ps, err := ValidateGetRequestsRequest(t, pageNumber, pageSize, userID)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	userRequests, err := s.Referral.GetRequests(userID, t, pn, ps)
	if err != nil {
		sendResponse(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(rw, userRequests, http.StatusOK)
}

// DownloadCV downloads CV of a particular candidate.
func (s *Server) DownloadCV(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(idParameter)

	if err := ValidateNumber(id); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	link, err := s.Referral.DownloadFile(id)
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, DownloadResponse{FileLink: link}, http.StatusOK)
}

// UpdateRequest type presents data for request update.
type UpdateRequest struct {
	ID        string `json:"id"`
	NewStatus string `json:"status"`
}

// UpdateRespone prsents type with info about request update.
type UpdateRespone struct {
	Message string `json:"message"`
}

// UpdateRequest updated status of request by id.
func (s *Server) UpdateRequest(rw http.ResponseWriter, r *http.Request) {
	var request UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := ValidateRequestState(request.NewStatus); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	if err := ValidateNumber(request.ID); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	err := s.Referral.UpdateRequest(request.ID, request.NewStatus)
	if errors.Is(err, service.ErrNoResult) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, UpdateRespone{Message: fmt.Sprintf("request status updated to '%s'", request.NewStatus)}, http.StatusOK)
}

// ValidateSignUpRequest validates data after signup.
func (r *SignUpRequest) ValidateSignUpRequest() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name", ErrInvalidParameter)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", ErrInvalidParameter)
	}

	if len(r.Name) < 6 || len(r.Name) > 18 {
		return fmt.Errorf("%w: wrong length", ErrInvalidParameter)
	}

	if len(r.Password) < 6 || len(r.Password) > 18 {
		return fmt.Errorf("%w: wrong length", ErrInvalidParameter)
	}

	return nil
}

// ValidateLogInRequest validates data after login.
func (r *LogInRequest) ValidateLogInRequest() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name", ErrInvalidParameter)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", ErrInvalidParameter)
	}

	return nil
}

// ValidateCandidateSendingRequest validates data after sending a candidate.
func ValidateCandidateSendingRequest(r service.SubmitCandidateRequest) error {
	if len(r.CandidateName) == 0 || len(r.CandidateSurname) == 0 {
		return fmt.Errorf("%w: wrong length", ErrInvalidParameter)
	}

	nameSurnameExp := "(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?"
	isValid, _ := regexp.MatchString(nameSurnameExp, r.CandidateName)
	if !isValid {
		return fmt.Errorf("%w: name has invalid format", ErrInvalidParameter)
	}

	isValid, _ = regexp.MatchString(nameSurnameExp, r.CandidateSurname)
	if !isValid {
		return fmt.Errorf("%w: surname has invalid format", ErrInvalidParameter)
	}

	return nil
}

var requestsState = map[string]bool{
	"accepted":  true,
	"rejected":  true,
	"submitted": true,
	"":          true,
}

// ValidateGetRequestsRequest validates parameters of request of getting requests.
func ValidateGetRequestsRequest(state, pageNumber, pageSize, id string) (int, int, error) {
	var pn int
	var ps int

	if !requestsState[strings.ToLower(state)] {
		return 0, 0, fmt.Errorf("%w: request state", ErrInvalidParameter)
	}

	idExp := "^([1-9])\\d*$"

	if pageSize != "" {
		ok, err := regexp.MatchString(idExp, pageSize)
		if !ok {
			return 0, 0, fmt.Errorf("%w: numeric parameter has bad format", ErrInvalidParameter)
		}
		if err != nil {
			return 0, 0, fmt.Errorf("cannot validate input parameter: %w", err)
		}

		ps, err = strconv.Atoi(pageSize)
		if err != nil {
			return 0, 0, fmt.Errorf("cannot conver page number to int: %w", err)
		}
	} else {
		ps = defaultPageSize
	}

	if pageNumber != "" {
		ok, err := regexp.MatchString(idExp, pageNumber)
		if !ok {
			return 0, 0, fmt.Errorf("%w: numeric parameter has bad format", ErrInvalidParameter)
		}
		if err != nil {
			return 0, 0, fmt.Errorf("cannot validate input parameter: %w", err)
		}

		pn, err = strconv.Atoi(pageNumber)
		if err != nil {
			return 0, 0, fmt.Errorf("cannot convert page size to integer: %w", err)
		}
	} else {
		pn = defaultPageNumber
	}

	ok, err := regexp.MatchString(idExp, id)
	if !ok {
		return 0, 0, fmt.Errorf("%w: numeric parameter has bad format", ErrInvalidParameter)
	}
	if err != nil {
		return 0, 0, fmt.Errorf("cannot validate input parameter: %w", err)
	}

	return pn, ps, nil
}

// ValidateRequestState validates data for request filtering.
func ValidateRequestState(state string) error {
	if !requestsState[strings.ToLower(state)] {
		return fmt.Errorf("%w: request state", ErrInvalidParameter)
	}

	return nil
}

// ValidateNumber checks if parameter is number.
func ValidateNumber(id string) error {
	idExp := "^[1-9]\\d*"
	ok, err := regexp.MatchString(idExp, id)
	if !ok {
		return fmt.Errorf("%w: numeric parameter has bad format", ErrInvalidParameter)
	}
	if err != nil {
		return fmt.Errorf("cannot validate input parameter: %w", err)
	}

	return nil
}
