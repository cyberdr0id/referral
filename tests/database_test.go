package tests

import (
	"fmt"
)

const (
	defaultName             = "username"
	defaultPassword         = "password"
	defaultIsAdmin          = false
	defaultUserID           = "1"
	defaultFileID           = "1"
	defaultRequestID        = "1"
	defaultStatus           = "Submitted"
	defaultPageNumber       = 1
	defaultPageSize         = 1
	defaultCandidateName    = "candidate"
	defaultCandidateSurname = "candidate"
	defaultRequestsLength   = 1

	statusAccepted = "Accepted"
)

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
	id, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	_, err = s.repo.AddCandidate(id, defaultCandidateName, defaultCandidateSurname, defaultFileID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot add candidate: %w", err).Error())
	}
	s.NoError(err)

	s.clearTables()
}

func (s *ReferralAPISuite) TestGetRequests() {
	id, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	_, err = s.repo.AddCandidate(id, defaultCandidateName, defaultCandidateSurname, defaultFileID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot add candidate: %w", err).Error())
	}
	s.NoError(err)

	requests, err := s.repo.GetRequests(id, defaultStatus, defaultPageNumber, defaultPageSize)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot get requests: %w", err).Error())
	}
	s.NoError(err)

	s.Equal(defaultRequestsLength, len(requests))

	s.clearTables()
}

func (s *ReferralAPISuite) TestUpdateRequest() {
	id, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	requestID, err := s.repo.AddCandidate(id, defaultCandidateName, defaultCandidateSurname, defaultFileID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot add candidate: %w", err).Error())
	}
	s.NoError(err)

	err = s.repo.UpdateRequest(requestID, statusAccepted)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot update request: %w", err).Error())
	}
	s.NoError(err)

	s.clearTables()
}

func (s *ReferralAPISuite) TestGetCVID() {
	id, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	requestID, err := s.repo.AddCandidate(id, defaultCandidateName, defaultCandidateSurname, defaultFileID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot add candidate: %w", err).Error())
	}
	s.NoError(err)

	fileID, err := s.repo.GetCVID(requestID)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot get id of cv: %w", err).Error())
	}
	s.NoError(err)

	s.Equal(defaultFileID, fileID)

	s.clearTables()
}
