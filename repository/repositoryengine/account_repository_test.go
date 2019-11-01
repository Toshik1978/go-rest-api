package repositoryengine

import (
	"context"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/testutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type accountRepositoryTestSuite struct {
	suite.Suite

	account repository.Account
}

func (s *accountRepositoryTestSuite) SetupSuite() {
	s.account = testutil.RepositoryAccount()
}

func (s *accountRepositoryTestSuite) TestGetAllAccountsFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectQuery("^SELECT id, uid").
		WillReturnError(errors.New("fail"))

	repository := newAccountRepository(sqlxDB)
	accounts, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
	s.Nil(accounts)
}

func (s *accountRepositoryTestSuite) TestGetAllAccountsEmptySucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	allRows := sqlmock.
		NewRows([]string{"id", "uid", "currency", "balance", "created_at"})

	mockSQL.
		ExpectQuery("^SELECT id, uid").
		WillReturnRows(allRows)

	repository := newAccountRepository(sqlxDB)
	accounts, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.Empty(accounts)
}

func (s *accountRepositoryTestSuite) TestGetAllAccountsSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	allRows := sqlmock.
		NewRows([]string{"id", "uid", "currency", "balance", "created_at"}).
		AddRow(s.account.ID, s.account.UID, s.account.Currency, s.account.Balance, s.account.CreatedAt)

	mockSQL.
		ExpectQuery("^SELECT id, uid").
		WillReturnRows(allRows)

	repository := newAccountRepository(sqlxDB)
	accounts, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.Len(accounts, 1)
	s.EqualValues(s.account, accounts[0])
}

func (s *accountRepositoryTestSuite) TestStoreAccountFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^INSERT INTO accounts").
		WithArgs(s.account.UID, s.account.Currency, s.account.Balance, s.account.CreatedAt).
		WillReturnError(errors.New("fail"))

	repository := newAccountRepository(sqlxDB)
	account := s.account
	err = repository.Store(context.Background(), &account)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
}

func (s *accountRepositoryTestSuite) TestStoreAccountSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^INSERT INTO accounts").
		WithArgs(s.account.UID, s.account.Currency, s.account.Balance, s.account.CreatedAt).
		WillReturnResult(sqlmock.NewResult(s.account.ID, 1))

	repository := newAccountRepository(sqlxDB)
	account := s.account
	account.ID = 0
	err = repository.Store(context.Background(), &account)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.EqualValues(s.account, account)
}

func (s *accountRepositoryTestSuite) TestUpdateAccountBalanceFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^UPDATE accounts").
		WithArgs(s.account.UID, s.account.Balance).
		WillReturnError(errors.New("fail"))

	repository := newAccountRepository(sqlxDB)
	err = repository.UpdateBalance(context.Background(), s.account.UID, s.account.Balance)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
}

func (s *accountRepositoryTestSuite) TestUpdateAccountBalanceSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^UPDATE accounts").
		WithArgs(s.account.UID, s.account.Balance).
		WillReturnResult(sqlmock.NewResult(s.account.ID, 1))

	repository := newAccountRepository(sqlxDB)
	err = repository.UpdateBalance(context.Background(), s.account.UID, s.account.Balance)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
}
