package repositoryengine

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type scopeTestSuite struct {
	suite.Suite
}

func (s *scopeTestSuite) TestScopeWithContextFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin().
		WillReturnError(errors.New("fail"))

	scope := newScope(sqlxDB)
	ctx, err := scope.WithContext(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
	s.Nil(ctx)
}

func (s *scopeTestSuite) TestScopeWithContextSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.ExpectBegin()

	scope := newScope(sqlxDB)
	ctx, err := scope.WithContext(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.NotNil(ctx)
}

func (s *scopeTestSuite) TestScopeCompleteNoTransactionSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	scope := newScope(sqlxDB)
	err = scope.Complete(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}

func (s *scopeTestSuite) TestScopeCompleteFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectCommit().
		WillReturnError(errors.New("fail"))

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Complete(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
}

func (s *scopeTestSuite) TestScopeCompleteSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectCommit()

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Complete(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}

func (s *scopeTestSuite) TestScopeCompleteTxDoneSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectCommit().
		WillReturnError(sql.ErrTxDone)

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Complete(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}

func (s *scopeTestSuite) TestScopeCancelNoTransactionSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	scope := newScope(sqlxDB)
	err = scope.Cancel(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}

func (s *scopeTestSuite) TestScopeCancelFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectRollback().
		WillReturnError(errors.New("fail"))

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Cancel(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
}

func (s *scopeTestSuite) TestScopeCancelSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectRollback()

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Cancel(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}

func (s *scopeTestSuite) TestScopeCancelTxDoneSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectBegin()
	mockSQL.
		ExpectRollback().
		WillReturnError(sql.ErrTxDone)

	tx, _ := sqlxDB.Beginx()
	scope := newScope(sqlxDB)
	err = scope.Cancel(contextWithTransaction(context.Background(), tx))

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}
