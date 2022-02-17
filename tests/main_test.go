package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/stretchr/testify/suite"
)

type ReferralAPISuite struct {
	suite.Suite

	db   *sql.DB
	repo *repository.Repository
}

func TestReferralAPISuite(t *testing.T) {
	suite.Run(t, new(ReferralAPISuite))
}

func (s *ReferralAPISuite) SetupSuite() {
	config := repository.DatabaseConfig{
		Host:         os.Getenv("DB_HOST"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         os.Getenv("DB_PORT"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
	}

	db, err := repository.NewConnection(config)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot create database connection: %w", err).Error())
	}
	s.db = db

	repo := repository.NewRepository(db)
	s.repo = repo
}

func (s *ReferralAPISuite) TearDownSuite() {
	if err := s.db.Close(); err != nil {
		s.FailNow(fmt.Errorf("cannot close database connection: %w", err).Error())
	}
}

func (s *ReferralAPISuite) SetupTest() {
	s.clearTables()
}

func (s *ReferralAPISuite) TearDownTest() {
	s.clearTables()
}

func (s *ReferralAPISuite) clearTables() {
	clearUsersQuery := `TRUNCATE TABLE users CASCADE`
	clearRequestsQuery := `TRUNCATE TABLE requests`

	_, err := s.db.Exec(clearUsersQuery)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot clear USERS table: %w", err).Error())
	}

	_, err = s.db.Exec(clearRequestsQuery)
	if err != nil {
		s.FailNow(fmt.Errorf("cannot clear REQUESTS table: %w", err).Error())
	}
}
