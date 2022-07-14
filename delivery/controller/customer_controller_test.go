package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

//controller butuh usecase
type CustomerUseCaseMock struct {
	mock.Mock
}

var dummyCustomer = []model.Customer{
	{
		Id:      "C001",
		Nama:    "Name One",
		Address: "Dummy One",
	},
	{
		Id:      "C002",
		Nama:    "Name Two",
		Address: "Dummy Two",
	},
}

// buat test stuite

type CustomerControllerTestSuite struct {
	suite.Suite
	routerMock *gin.Engine
	useCaseMock *CustomerUseCaseMock
}

func (suite *CustomerControllerTestSuite) SetupTest(){
	suite.routerMock = gin.Default()

	suite.useCaseMock = new(CustomerUseCaseMock)
}



// func (r *CustomerUseCaseMock) RegisterCustomer(newCustomer model.Customer) error {
// 	args := r.Called(newCustomer)
// 	if args.Get(1) != nil {
// 		return args.Get(0).(error)
// 	}

// 	return nil
// }

// func (r *CustomerUseCaseMock) GetAllCustomer() ([]model.Customer, error) {
// 	args := r.Called()
// 	if args.Get(1) == nil {
// 		return nil, args.Error(1)
// 	}

// 	return args.Get(0).([]model.Customer), nil

// }

// func (r *CustomerUseCaseMock) FindCustomerById(id string) (model.Customer, error)  {
// 	args := r.Called(id) //ini args adan hubungan sama returnnya
// 	if args.Get(1) != nil { //get 1 karena kita abil index data
// 		return model.Customer{}, args.Get(1).(error)
// 	}

// 	return args.Get(0).(model.Customer), nil
// }

func (r *CustomerUseCaseMock) RegisterCustomer(newCustomer model.Customer) error {
	args := r.Called(newCustomer)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (r *CustomerUseCaseMock) GetAllCustomer() ([]model.Customer, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Customer), nil

}

func (r *CustomerUseCaseMock) FindCustomerById(id string) (model.Customer, error)  {
	args := r.Called(id) //ini args adan hubungan sama returnnya
	if args.Get(1) != nil { //get 1 karena kita abil index data
		return model.Customer{}, args.Error(1)
	}

	return args.Get(0).(model.Customer), nil
}


func (suite *CustomerControllerTestSuite) TestGetAllCustomerApi_Success() {

	customers := []model.Customer{
		{
			Id: "COO1",
			Nama: "Dummy Name 1",
			Address: "Dummy Address 1",
		},
	}
	
suite.useCaseMock.On("GetAllCustomer").Return(customers, nil)

NewCustomerController(suite.routerMock, suite.useCaseMock)

//ini baru kondisikan http status

r := httptest.NewRecorder()

//request test yang sesuai

request, err := http.NewRequest(http.MethodGet, "/customer", nil)

suite.routerMock.ServeHTTP(r, request)

var actualCustomers []model.Customer

response := r.Body.String()

json.Unmarshal([]byte(response), &actualCustomers)

assert.Equal(suite.T(), http.StatusOK, r.Code)
assert.Equal(suite.T(), 1, len(actualCustomers))
assert.Equal(suite.T(), customers[0].Nama, actualCustomers[0].Nama)
//gak pake nama juga gpp
assert.Nil(suite.T(), err)


}

func (suite *CustomerControllerTestSuite) TestGetAllCustomerApi_Failed() {
suite.useCaseMock.On("GetAllCustomer").Return(nil, errors.New("failed"))

NewCustomerController(suite.routerMock, suite.useCaseMock)

//ini baru kondisikan http status

r := httptest.NewRecorder()

//request test yang sesuai

request, _ := http.NewRequest(http.MethodGet, "/customer", nil)

suite.routerMock.ServeHTTP(r, request)
var errorResponse struct {Err string}

response := r.Body.String()

json.Unmarshal([]byte(response), &errorResponse)

assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
assert.Equal(suite.T(),"failed", errorResponse.Err)


}


// yang create bikin error pas binding dan error pas masuk usecase


func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_Success() {

	dummyCustomer := dummyCustomer[0]

	suite.useCaseMock.On("RegisterCustomer", dummyCustomer).Return(nil)

	NewCustomerController(suite.routerMock, suite.useCaseMock)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(reqBody))

	suite.routerMock.ServeHTTP(r, request)

	response := r.Body.String()

	var actualCustomers model.Customer

	json.Unmarshal([]byte(response), &actualCustomers)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), dummyCustomer.Nama, actualCustomers.Nama)

}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_Failed() {
	r := httptest.NewRecorder()


	NewCustomerController(suite.routerMock, suite.useCaseMock)
	request, _ := http.NewRequest(http.MethodPost, "/customer", nil)

	suite.routerMock.ServeHTTP(r, request)

	var errorResponse struct{ Err string }

	response := r.Body.String()

	json.Unmarshal([]byte(response), &errorResponse)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.NotEmpty(suite.T(), errorResponse.Err)

	
}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_FailedBinding() {
	r := httptest.NewRecorder()
	NewCustomerController(suite.routerMock, suite.useCaseMock)

	request, _ := http.NewRequest(http.MethodPost, "/customer", nil)

	suite.routerMock.ServeHTTP(r, request)

	var errorResponse struct{ Err string }

	response := r.Body.String()
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.NotEmpty(suite.T(), errorResponse.Err)
}

func (suite *CustomerControllerTestSuite) TestRegisterCustomerApi_FailedUseCase() {

	dummyCustomer := dummyCustomer[0]

	suite.useCaseMock.On("RegisterCustomer", dummyCustomer).Return(errors.New("failed"))

	NewCustomerController(suite.routerMock, suite.useCaseMock)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(reqBody))

	suite.routerMock.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)

	response := r.Body.String()

	var errorResponse struct{ Err string}

	json.Unmarshal([]byte(response), &errorResponse)

	assert.Equal(suite.T(), "failed", errorResponse.Err)

}



func TestCustomerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerControllerTestSuite))
}

