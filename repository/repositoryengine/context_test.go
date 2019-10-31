package repositoryengine

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type contextTestSuite struct {
	suite.Suite
}

func (s *contextTestSuite) TestContextWithTransactionNil() {
	s.Nil(transactionFromContext(context.Background()))
}

func (s *contextTestSuite) TestContextWithTransactionBadType() {
	ctx := context.WithValue(context.Background(), txKey, 1234)
	s.Nil(transactionFromContext(ctx))
}

func (s *contextTestSuite) TestContextWithTransactionSucceeded() {
	tx := sqlx.Tx{}
	ctx := contextWithTransaction(context.Background(), &tx)

	s.Equal(transactionFromContext(ctx), &tx)
}
