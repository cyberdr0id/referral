package tests

import (
	"fmt"
)

const (
	defaultName     = "username"
	defaultPassword = "password"
)

func (s *ReferralAPISuite) TestCreateUser() {
	_, err := s.repo.CreateUser(defaultName, defaultPassword)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create user: %w", err).Error())
	}
	s.NoError(err)

	s.clearTables()
}
