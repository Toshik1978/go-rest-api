package repositoryengine

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type utilsTestSuite struct {
	suite.Suite
}

func (s *utilsTestSuite) TestSqlxExtDefaultSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ext := sqlxExt(context.Background(), sqlxDB)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Equal(sqlxDB, ext)
}

func (s *utilsTestSuite) TestSqlxExtSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.ExpectBegin()
	tx, _ := sqlxDB.Beginx()
	ctx := contextWithTransaction(context.Background(), tx)
	ext := sqlxExt(ctx, sqlxDB)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Equal(tx, ext)
}
