package httphandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/mock"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/Toshik1978/go-rest-api/service/testutil"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// failResponseWriter declare ResponseWriter, which always fails
type failResponseWriter struct {
	headers http.Header
}

func (r *failResponseWriter) Header() http.Header {
	return r.headers
}

func (r *failResponseWriter) Write(body []byte) (int, error) {
	return 0, errors.New("fail")
}

func (r *failResponseWriter) WriteHeader(status int) {
}

func newFailResponseWriter() http.ResponseWriter {
	return &failResponseWriter{
		headers: make(http.Header),
	}
}

// apiHandlerTestSuite test suite
type apiHandlerTestSuite struct {
	suite.Suite
}

func (s *apiHandlerTestSuite) TestWriteResponseFailed1() {
	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil)
	apiHandler.writeResponse(newFailResponseWriter(), make(chan struct{}))

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to marshal HTTP response", zapRecorded.All()[0].Message)
}

func (s *apiHandlerTestSuite) TestWriteResponseFailed2() {
	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil)
	apiHandler.writeResponse(newFailResponseWriter(), handler.ServerStatusResponse{})

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to write HTTP response", zapRecorded.All()[0].Message)
}

func (s *apiHandlerTestSuite) TestServerStatusHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger:    zap.New(zapCore),
		BuildTime: time.Now().String(),
		Version:   "test",
	}, nil).ServerStatusHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusOK, r.Code)
}

func (s *apiHandlerTestSuite) TestCreateAccountHandlerNoBodyFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil).CreateAccountHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreateAccountHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreateAccountHandlerBadRequestFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	payload := `{"uid": 1}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil).CreateAccountHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreateAccountHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreateAccountHandlerFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	request := testutil.AccountRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}

	accountBuilder := mock.NewMockAccountBuilder(ctrl)
	accountBuilder.
		EXPECT().
		SetUID(gomock.Eq(request.UID)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		SetCurrency(gomock.Eq(request.Currency)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		SetBalance(gomock.Eq(request.Balance)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(nil, errors.New("fail"))

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AccountBuilder().
		Return(accountBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreateAccountHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreateAccountHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, r.Code)
}

func (s *apiHandlerTestSuite) TestCreateAccountHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountResponse()
	request := testutil.AccountRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}

	accountBuilder := mock.NewMockAccountBuilder(ctrl)
	accountBuilder.
		EXPECT().
		SetUID(gomock.Eq(request.UID)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		SetCurrency(gomock.Eq(request.Currency)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		SetBalance(gomock.Eq(request.Balance)).
		Return(accountBuilder)
	accountBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(&account, nil)

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AccountBuilder().
		Return(accountBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreateAccountHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	var response handler.Account
	_ = json.Unmarshal(r.Body.Bytes(), &response)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusCreated, r.Code)
	s.EqualValues(account, response)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerNoBodyFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerBadRequestFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	payload := `{"amount": "100"}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerBadURLFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountRequest()
	payload := `{"amount": "100"}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		s.T().Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"uid": account.UID,
	})

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, nil).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountRequest()
	request := testutil.PaymentRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"uid": account.UID,
	})

	paymentBuilder := mock.NewMockPaymentBuilder(ctrl)
	paymentBuilder.
		EXPECT().
		SetPayer(gomock.Eq(account.UID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetRecipient(gomock.Eq(request.RecipientUID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetAmount(gomock.Eq(request.Amount)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(nil, errors.New("fail"))

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		PaymentBuilder().
		Return(paymentBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerClientErrorFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountRequest()
	request := testutil.PaymentRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"uid": account.UID,
	})

	paymentBuilder := mock.NewMockPaymentBuilder(ctrl)
	paymentBuilder.
		EXPECT().
		SetPayer(gomock.Eq(account.UID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetRecipient(gomock.Eq(request.RecipientUID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetAmount(gomock.Eq(request.Amount)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(nil, handler.NewError("fail", handler.ClientError))

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		PaymentBuilder().
		Return(paymentBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusBadRequest, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerServerErrorFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountRequest()
	request := testutil.PaymentRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"uid": account.UID,
	})

	paymentBuilder := mock.NewMockPaymentBuilder(ctrl)
	paymentBuilder.
		EXPECT().
		SetPayer(gomock.Eq(account.UID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetRecipient(gomock.Eq(request.RecipientUID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetAmount(gomock.Eq(request.Amount)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(nil, handler.WrapError(errors.New("fail"), "fail", handler.ServerError))

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		PaymentBuilder().
		Return(paymentBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle CreatePaymentHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, r.Code)
}

func (s *apiHandlerTestSuite) TestCreatePaymentHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	account := testutil.AccountRequest()
	payment := testutil.PaymentResponse()
	request := testutil.PaymentRequest()
	payload, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		s.T().Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"uid": account.UID,
	})

	paymentBuilder := mock.NewMockPaymentBuilder(ctrl)
	paymentBuilder.
		EXPECT().
		SetPayer(gomock.Eq(account.UID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetRecipient(gomock.Eq(request.RecipientUID)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		SetAmount(gomock.Eq(request.Amount)).
		Return(paymentBuilder)
	paymentBuilder.
		EXPECT().
		Build(gomock.Any()).
		Return(&payment, nil)

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		PaymentBuilder().
		Return(paymentBuilder)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).CreatePaymentHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	var response handler.Payment
	_ = json.Unmarshal(r.Body.Bytes(), &response)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusCreated, r.Code)
	s.EqualValues(payment, response)
}

func (s *apiHandlerTestSuite) TestGetAllAccountsHandlerFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AllAccounts(gomock.Any()).
		Return(nil, errors.New("fail"))

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).GetAllAccountsHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle GetAllAccountsHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, r.Code)
}

func (s *apiHandlerTestSuite) TestGetAllAccountsHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AllAccounts(gomock.Any()).
		Return(nil, nil)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).GetAllAccountsHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusOK, r.Code)
}

func (s *apiHandlerTestSuite) TestGetAllPaymentsHandlerFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AllPayments(gomock.Any()).
		Return(nil, errors.New("fail"))

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).GetAllPaymentsHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to handle GetAllPaymentsHandler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, r.Code)
}

func (s *apiHandlerTestSuite) TestGetAllPaymentsHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	accountManager := mock.NewMockAccountManager(ctrl)
	accountManager.
		EXPECT().
		AllPayments(gomock.Any()).
		Return(nil, nil)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(server.Globals{
		Logger: zap.New(zapCore),
	}, accountManager).GetAllPaymentsHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusOK, r.Code)
}
