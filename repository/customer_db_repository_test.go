package repository

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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


type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockDb *sql.DB //mau mock DB connectionnya

	mockSql sqlmock.Sqlmock //github.com/DATA-DOG/go-sqlmock

}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()

	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}

	suite.mockSql = mockSql
	suite.mockDb = mockDb
}

func (suite *CustomerRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}


func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Success() {
	//siapkan column, sama seperti di field table customer
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})

	for _, v := range dummyCustomer {
		rows.AddRow(v.Id, v.Nama, v.Address)
	}

	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnRows(rows)
	//panggil querry mocknya menggunakan regex --> (.+)

	//panggil repository asli
	repo := NewCustomerDbRepository(suite.mockDb)

	//panggil method yang mau di test
	actual, err := repo.RetrieveAll()

	//test assertion

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2 , len(actual))
	assert.Equal(suite.T(), "C001", actual[0].Id)

}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Failed() {
	//siapkan column, sama seperti di field table customer
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})

	for _, v := range dummyCustomer {
		rows.AddRow(v.Id, v.Nama, v.Address)

	}

	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnError(errors.New("failed"))
	//panggil querry mocknya menggunakan regex --> (.+)

	//panggil repository asli
	repo := NewCustomerDbRepository(suite.mockDb)

	//panggil method yang mau di test
	actual, err := repo.RetrieveAll()

	//test assertion

	assert.Nil(suite.T(), actual)
	assert.Error(suite.T(), err)
	

}

func(suite *CustomerRepositoryTestSuite) TestCustomerCreate_Success() {

	dummyCustomer := dummyCustomer[0]

	suite.mockSql.ExpectExec("insert into customer values").WithArgs(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address).WillReturnResult(sqlmock.NewResult(1,1))

	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)

	assert.Nil(suite.T(), err)

}

func(suite *CustomerRepositoryTestSuite) TestCustomerCreate_Failed() {

	suite.mockSql.ExpectExec("insert into customer values").WithArgs(dummyCustomer[0].Id, dummyCustomer[0].Nama, dummyCustomer[0].Address).WillReturnResult(sqlmock.NewResult(1,1))

	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer[0])

	assert.Nil(suite.T(), err)

}

func(suite *CustomerRepositoryTestSuite) TestCustomerFindById_Success() {
	dummyCustomer := dummyCustomer[0]
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)

	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnRows(rows)
	//panggil querry mocknya menggunakan regex --> (.+)

	//panggil repository asli
	repo := NewCustomerDbRepository(suite.mockDb)

	//panggil method yang mau di test
	actual, err := repo.FindById(dummyCustomer.Id)

	//test assertion

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), actual)


}

func(suite *CustomerRepositoryTestSuite) TestCustomerFindById_Failed() {

	dummyCustomer := dummyCustomer[0]

	rows := sqlmock.NewRows([]string{"ids", "namaaaaa", "addresssss"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)

	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnError(errors.New("failed"))

	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.FindById(dummyCustomer.Id)

	assert.NotEqual(suite.T(),dummyCustomer, actual)
	assert.Error(suite.T(), err)

}




func TestCustomerRepositoryTestSuite(t *testing.T){
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
