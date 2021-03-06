package tests

import (
	"fmt"
)

const (
	defaultName             = "username"
	defaultPassword         = "password"
	defaultIsAdmin          = false
	defaultFileID           = "1"
	defaultStatus           = "submitted"
	defaultPageNumber       = 1
	defaultPageSize         = 1
	defaultCandidateName    = "candidate"
	defaultCandidateSurname = "candidate"
	defaultRequestsLength   = 1

	statusAccepted = "accepted"
)

func makeRequest(s *ReferralAPISuite) (id string, requestID string) {
	id, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	requestID, err = s.repo.AddCandidate(id, defaultCandidateName, defaultCandidateSurname, defaultFileID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot add candidate: %w", err).Error())
	}
	s.NoError(err)

	return id, requestID
}

func (s *ReferralAPISuite) TestCreateUser() {
	_, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	s.clearTables()
}

func (s *ReferralAPISuite) TestGetUser() {
	_, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	user, err := s.repo.GetUser(defaultName)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot get user: %w", err).Error())
	}
	s.NoError(err)

	s.Equal(defaultName, user.Name)
	s.Equal(defaultPassword, user.Password)
	s.Equal(defaultIsAdmin, user.IsAdmin)

	s.clearTables()
}

func (s *ReferralAPISuite) TestAddCandidate() {
	_, _ = makeRequest(s)
	s.clearTables()
}

func (s *ReferralAPISuite) TestGetRequests() {
	id, _ := makeRequest(s)

	requests, err := s.repo.GetRequests(id, defaultStatus, defaultPageNumber, defaultPageSize)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot get requests: %w", err).Error())
	}
	s.NoError(err)

	s.Equal(defaultRequestsLength, len(requests))

	s.clearTables()
}

func (s *ReferralAPISuite) TestUpdateRequest() {
	_, requestID := makeRequest(s)

	err := s.repo.UpdateRequest(requestID, statusAccepted)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot update request: %w", err).Error())
	}
	s.NoError(err)

	s.clearTables()
}

func (s *ReferralAPISuite) TestGetCVID() {
	userID, requestID := makeRequest(s)

	fileID, err := s.repo.GetCVID(requestID, userID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot get id of cv: %w", err).Error())
	}
	s.NoError(err)

	s.Equal(defaultFileID, fileID)

	s.clearTables()
}
