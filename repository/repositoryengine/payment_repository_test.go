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

type paymentRepositoryTestSuite struct {
	suite.Suite

	payment repository.Payment
}

func (s *paymentRepositoryTestSuite) SetupSuite() {
	s.payment = testutil.RepositoryPayment()
}

func (s *paymentRepositoryTestSuite) TestGetAllPaymentsFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectQuery("^SELECT id, amount").
		WillReturnError(errors.New("fail"))

	repository := newPaymentRepository(sqlxDB)
	payments, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
	s.Nil(payments)
}

func (s *paymentRepositoryTestSuite) TestGetAllPaymentsEmptySucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	allRows := sqlmock.
		NewRows([]string{"id", "amount", "payer_account_uid", "recipient_account_uid", "created_at"})

	mockSQL.
		ExpectQuery("^SELECT id, amount").
		WillReturnRows(allRows)

	repository := newPaymentRepository(sqlxDB)
	payments, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.Empty(payments)
}

func (s *paymentRepositoryTestSuite) TestGetAllPaymentsSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	allRows := sqlmock.
		NewRows([]string{"id", "amount", "payer_account_uid", "recipient_account_uid", "created_at"}).
		AddRow(s.payment.ID,
			s.payment.Amount, s.payment.PayerAccountUID, s.payment.RecipientAccountUID, s.payment.CreatedAt)

	mockSQL.
		ExpectQuery("^SELECT id, amount").
		WillReturnRows(allRows)

	repository := newPaymentRepository(sqlxDB)
	payments, err := repository.GetAll(context.Background())

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.Len(payments, 1)
	s.EqualValues(s.payment, payments[0])
}

func (s *paymentRepositoryTestSuite) TestStorePaymentFailed() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^INSERT INTO payments").
		WithArgs(s.payment.Amount, s.payment.PayerAccountUID, s.payment.RecipientAccountUID, s.payment.CreatedAt).
		WillReturnError(errors.New("fail"))

	repository := newPaymentRepository(sqlxDB)
	payment := s.payment
	err = repository.Store(context.Background(), &payment)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.Error(err)
}

func (s *paymentRepositoryTestSuite) TestStorePaymentSucceeded() {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mockSQL.
		ExpectExec("^INSERT INTO payments").
		WithArgs(s.payment.Amount, s.payment.PayerAccountUID, s.payment.RecipientAccountUID, s.payment.CreatedAt).
		WillReturnResult(sqlmock.NewResult(s.payment.ID, 1))

	repository := newPaymentRepository(sqlxDB)
	payment := s.payment
	payment.ID = 0
	err = repository.Store(context.Background(), &payment)

	s.NoError(mockSQL.ExpectationsWereMet())
	s.NoError(err)
	s.EqualValues(s.payment, payment)
}
