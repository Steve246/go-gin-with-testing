package usecase

import (
	"testing"

	"enigmacamp.com/golatihanlagi/model"
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
	if args.Get(0) == nil { //get o karena kita abil index data
		return model.Customer{}, args.Error(1)
	}

	return args.Get(0).(model.Customer), nil
}

func(suite *CustomerUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func(suite *CustomerUseCaseTestSuite) TestCustomerFindById_Success() {
	dummyCustomer := dummyCustomer[0]
	suite.repoMock.On("FindById", dummyCustomer).Return(dummyCustomer, nil)
	//method mengacu ke repo, kalau repo findbyid maka findbyid
}

func(suite *CustomerUseCaseTestSuite) TestCustomerFindById_Failed() {
	
}

func TestCustomerUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerUseCaseTestSuite))
}