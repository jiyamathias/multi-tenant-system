package storage

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	DB           *gorm.DB
	mock         sqlmock.Sqlmock
	userDatabase UserDatabase
}

func (s *Suite) SetupSuite() {
	var store *Storage
	s.mock, store = GetStorage(s.Suite.T())
	s.userDatabase = *NewUser(store)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.Suite.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_Storage_Dummy() {
	// Ha ha, dont mind me with this test. SMH
	testString := "Ninja skills are in progress. 80% loading..."
	require.Equal(s.Suite.T(), 44, len(testString))
}
