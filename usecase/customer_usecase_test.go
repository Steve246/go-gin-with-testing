package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCustomer = []model.Customer{
	{
		Id: "C001",
		Nama: "Name One",
		Address: "Dummy One",

	},
	{
		Id: "C002",
		Nama: "Name Two",
		Address: "Dummy Two",

	},
}



type repoMock struct {
	mock.Mock //Mock bohongan, unit testing tidak bisa ke database, atau network makanya dibikin repo mong (repo bohongan) agar dapat di test
}

type CustomerUseCaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

func (r *repoMock) Create(newCustomer model.Customer) error {
	args := r.Called(newCustomer)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (r *repoMock) RetrieveAll() ([]model.Customer, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Customer), nil

}

func (r *repoMock) FindById(id string) (model.Customer, error)  {
	args := r.Called(id) //ini args adan hubungan sama returnnya
	if args.Get(1) != nil { //get 1 karena kita abil index data
		return model.Customer{}, args.Error(1)
	}

	return args.Get(0).(model.Customer), nil
}

func(suite *CustomerUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func(suite *CustomerUseCaseTestSuite) TestCustomerFindById_Success() {
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(dummyCustomer, nil)
	//method mengacu ke repo, kalau repo findbyid maka findbyid
	
	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	customer, err := customerUsecaseTest.FindCustomerById(dummyCustomer.Id)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyCustomer.Id, customer.Id) 
	//apakah hasilnya sama atau ngak
}

func(suite *CustomerUseCaseTestSuite) TestCustomerFindById_Failed() {
	
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(model.Customer{}, errors.New("failed"))
	//method mengacu ke repo, kalau repo findbyid maka findbyid
	
	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	customer, err := customerUsecaseTest.FindCustomerById(dummyCustomer.Id)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(),"failed", err.Error()) //apakah hasilnya sama atau ngak

	assert.Equal(suite.T(), "",customer.Id)
	
}

func (suite *CustomerUseCaseTestSuite) TestCustomerRetrieveAll_Success() {

	//EXPECTED
	dummyCustomer := dummyCustomer //kita perlu semua array

	suite.repoMock.On("RetrieveAll").Return(dummyCustomer, nil)

	//ACTUAL UJI CODE USECASE
	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	customer, err := customerUsecaseTest.GetAllCustomer()

	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), dummyCustomer,customer)

}

func (suite *CustomerUseCaseTestSuite) TestCustomerRetrieveAll_Failed() {

	//EXPECTED
	// dummyCustomer := dummyCustomer 
	//kita perlu semua array

	suite.repoMock.On("RetrieveAll").Return(nil, errors.New("failed"))
	//samain sama return reponya

	//ACTUAL UJI CODE USECASE
	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	customer, err := customerUsecaseTest.GetAllCustomer()

	assert.NotNil(suite.T(), err)
	
	assert.Empty(suite.T(), nil, customer) //empty karena kalau error dia ada nill

}


func (suite *CustomerUseCaseTestSuite) TestCustomerCreate_Success(){
	//Expected
	var newDumyCustomer = model.Customer{}

	suite.repoMock.On("Create", newDumyCustomer).Return(nil) //balikin nil, tapi diinput data customer baru

	//ACTUAL UJI CODE USECASE

	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	err := customerUsecaseTest.RegisterCustomer(newDumyCustomer)
	//return nil 

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), nil, err) //bandingin nil data dummy dan error dari usecase asli

}

func (suite *CustomerUseCaseTestSuite) TestCustomerCreate_Failed(){
	//Expected
	var newDumyCustomer = model.Customer{}

	suite.repoMock.On("Create", newDumyCustomer).Return(errors.New("failed")) 

	//ACTUAL UJI CODE USECASE

	customerUsecaseTest := NewCustomerUseCase(suite.repoMock)

	err := customerUsecaseTest.RegisterCustomer(newDumyCustomer)
	

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.New("failed"), err)

}




func TestCustomerUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerUseCaseTestSuite))
}