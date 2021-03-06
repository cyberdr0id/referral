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
	userIDParameter       = "user_id"

	anyUserID         = ""
	defaultPageNumber = 1
	defaultPageSize   = 10
)

// ErrorResponse presents a custom error type for error response.
type ErrorResponse struct {
	Message string `json:"error"`
}

// ErrInvalidParameter presents an error when user enters invalid parameter.
var ErrInvalidParameter = errors.New("invalid parameter")

// sendResponse sends response with specified object in body.
func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SignUpRequest type that presents data for registration.
type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SignUpResponse type presents response after successful registration.
type SignUpResponse struct {
	ID string `json:"id"`
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
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, SignUpResponse{ID: id}, http.StatusCreated)
}

// LogInRequest type that presents data for authorization.
type LogInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LogInResponse type presents response after successful authorization.
type LogInResponse struct {
	Token string `json:"token"`
}

// LogIn logs in user
func (s *Server) LogIn(rw http.ResponseWriter, r *http.Request) {
	var request LogInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.ValidateLogInRequest(); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}

	token, err := s.Auth.LogIn(request.Name, request.Password)
	if errors.Is(err, service.ErrNoUser) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, LogInResponse{Token: token}, http.StatusOK)
}

// CandidateSendingResponse type that presents ID of sent candidate.
type CandidateSendingResponse struct {
	CandidateID string `json:"id"`
}

// SendCandidate sends candidate info and his cv.
func (s *Server) SendCandidate(rw http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile(filenameParam)
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := fileHeader.Open()
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
	}
	filetype := strings.Split(fileHeader.Filename, ".")
	request := service.SubmitCandidateRequest{
		File:             f,
		CandidateName:    r.FormValue(candidateNameParam),
		CandidateSurname: r.FormValue(candidateSurnameParam),
		Filetype:         filetype[len(filetype)-1],
	}

	if err := ValidateCandidateSendingRequest(request); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := s.Referral.AddCandidate(r.Context(), request)
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, CandidateSendingResponse{CandidateID: id}, http.StatusOK)
}

// GetRequests outputs all user requests.
func (s *Server) GetRequests(rw http.ResponseWriter, r *http.Request) {
	status := strings.ToLower(r.URL.Query().Get(statusParameter))
	pageNumber := r.URL.Query().Get(pageNumberParameter)
	pageSize := r.URL.Query().Get(pageSizeParameter)

	userID, ok := context.GetUserID(r.Context())
	if !ok {
		s.Logger.ErrorLogger.Println(fmt.Errorf("cannot get user id from context"))
		sendResponse(rw, fmt.Errorf("cannot get user id from context"), http.StatusInternalServerError)
		return
	}

	pageNumberInt, pageSizeInt, err := ValidateGetRequestsRequest(status, pageNumber, pageSize, userID)
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	userRequests, err := s.Referral.GetRequests(userID, status, pageNumberInt, pageSizeInt)
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, userRequests, http.StatusOK)
}

// GetAllRequests admin handler that returns list of all requests.
func (s *Server) GetAllRequests(rw http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get(statusParameter)
	pageNumber := r.URL.Query().Get(pageNumberParameter)
	pageSize := r.URL.Query().Get(pageSizeParameter)
	userID := r.URL.Query().Get(userIDParameter)

	pageNumberInt, pageSizeInt, err := ValidateGetRequestsRequest(status, pageNumber, pageSize, userID)
	if err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	userRequests, err := s.Referral.GetRequests(userID, status, pageNumberInt, pageSizeInt)
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, userRequests, http.StatusOK)
}

// DownloadResponse presents a type which contains link to file for download candidate cv.
type DownloadResponse struct {
	Link string `json:"link"`
}

// DownloadCV downloads CV of a particular candidate.
func (s *Server) DownloadCV(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(idParameter)

	if err := ValidateNumber(id); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	userID, ok := context.GetUserID(r.Context())
	if !ok {
		s.Logger.ErrorLogger.Println(fmt.Errorf("cannot get user id from context"))
		sendResponse(rw, fmt.Errorf("cannot get user id from context"), http.StatusInternalServerError)
		return
	}

	url, err := s.Referral.DownloadFile(r.Context(), id, userID)
	if errors.Is(err, service.ErrNoFile) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, DownloadResponse{Link: url}, http.StatusOK)
}

// DownloadAnyCV provides access for admin to download any CV of candidates.
func (s *Server) DownloadAnyCV(rw http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get(idParameter)

	if err := ValidateNumber(fileID); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	url, err := s.Referral.DownloadFile(r.Context(), fileID, anyUserID)
	if errors.Is(err, service.ErrNoFile) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, DownloadResponse{Link: url}, http.StatusOK)
}

// UpdateRequest type presents data for request update.
type UpdateRequest struct {
	ID        string `json:"id"`
	NewStatus string `json:"status"`
}

// UpdateResponse presents type with info about request update.
type UpdateResponse struct {
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

	if err := request.ValidateUpdateRequest(); err != nil {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	err := s.Referral.UpdateRequest(request.ID, strings.ToLower(request.NewStatus))
	if errors.Is(err, service.ErrNoResult) {
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err != nil {
		s.Logger.ErrorLogger.Println(err)
		sendResponse(rw, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(rw, UpdateResponse{Message: fmt.Sprintf("request status with %s ID has been updated", request.ID)}, http.StatusOK)
}

// ValidateUpdateRequest validates data before request update.
func (r *UpdateRequest) ValidateUpdateRequest() error {
	if !requestsStatus[strings.ToLower(r.NewStatus)] {
		return fmt.Errorf("%w: request status", ErrInvalidParameter)
	}

	idExp := "^([1-9])\\d*$"
	ok, err := regexp.MatchString(idExp, r.ID)
	if !ok {
		return fmt.Errorf("%w: id has bad format", ErrInvalidParameter)
	}
	if err != nil {
		return fmt.Errorf("cannot validate id parameter: %w", err)
	}

	return nil
}

// ValidateSignUpRequest validates registration data .
func (r *SignUpRequest) ValidateSignUpRequest() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name", ErrInvalidParameter)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", ErrInvalidParameter)
	}

	if len(r.Name) < 6 || len(r.Name) > 18 {
		return fmt.Errorf("%w: name must be between 6 and 18 symbols", ErrInvalidParameter)
	}

	if len(r.Password) < 6 || len(r.Password) > 18 {
		return fmt.Errorf("%w: password must be between 6 and 18 symbols", ErrInvalidParameter)
	}

	return nil
}

// ValidateLogInRequest validates authorization data.
func (r *LogInRequest) ValidateLogInRequest() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name", ErrInvalidParameter)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", ErrInvalidParameter)
	}

	return nil
}

// ValidateCandidateSendingRequest validates data before candidate sending.
func ValidateCandidateSendingRequest(r service.SubmitCandidateRequest) error {
	if len(r.CandidateName) == 0 || len(r.CandidateSurname) == 0 {
		return fmt.Errorf("%w: wrong length", ErrInvalidParameter)
	}

	nameSurnameExp := "^(^[A-Za-z??-????-??]{2,16})?$"
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

var requestsStatus = map[string]bool{
	"accepted":  true,
	"rejected":  true,
	"submitted": true,
	"":          true,
}

// ValidateGetRequestsRequest validates parameters of request of getting requests.
func ValidateGetRequestsRequest(status, pageNumber, pageSize, id string) (int, int, error) {
	var pn int
	var ps int

	if !requestsStatus[strings.ToLower(status)] {
		return 0, 0, fmt.Errorf("%w: request status", ErrInvalidParameter)
	}

	idExp := "^([1-9])\\d*$"

	if pageSize != "" {
		ok, err := regexp.MatchString(idExp, pageSize)
		if !ok {
			return 0, 0, fmt.Errorf("%w: page size has bad format", ErrInvalidParameter)
		}
		if err != nil {
			return 0, 0, fmt.Errorf("cannot validate page size: %w", err)
		}

		ps, err = strconv.Atoi(pageSize)
		if err != nil {
			return 0, 0, fmt.Errorf("cannot conver page size to int: %w", err)
		}
	} else {
		ps = defaultPageSize
	}

	if pageNumber != "" {
		ok, err := regexp.MatchString(idExp, pageNumber)
		if !ok {
			return 0, 0, fmt.Errorf("%w: page number has bad format", ErrInvalidParameter)
		}
		if err != nil {
			return 0, 0, fmt.Errorf("cannot validate page number: %w", err)
		}

		pn, err = strconv.Atoi(pageNumber)
		if err != nil {
			return 0, 0, fmt.Errorf("cannot convert page number to integer: %w", err)
		}
	} else {
		pn = defaultPageNumber
	}

	if id != "" {
		ok, err := regexp.MatchString(idExp, id)
		if !ok {
			return 0, 0, fmt.Errorf("%w: user id has bad format", ErrInvalidParameter)
		}
		if err != nil {
			return 0, 0, fmt.Errorf("cannot validate id parameter: %w", err)
		}
	}

	return pn, ps, nil
}

// ValidateNumber checks if parameter is number.
func ValidateNumber(id string) error {
	idExp := "^([1-9])\\d*$"
	ok, err := regexp.MatchString(idExp, id)
	if !ok {
		return fmt.Errorf("%w: id has bad format", ErrInvalidParameter)
	}
	if err != nil {
		return fmt.Errorf("cannot validate id parameter: %w", err)
	}

	return nil
}
