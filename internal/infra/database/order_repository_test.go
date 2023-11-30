package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/viniciusmgaspar/cleanArch/internal/entity"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestListOrders() {
	repo := NewOrderRepository(suite.Db)

	order, err := entity.NewOrder("0001", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	err = repo.Save(order)
	suite.NoError(err)

	order2, err := entity.NewOrder("0002", 15.0, 3.0)
	suite.NoError(err)
	suite.NoError(order2.CalculateFinalPrice())
	err = repo.Save(order2)
	suite.NoError(err)

	order3, err := entity.NewOrder("0003", 15.0, 3.0)
	suite.NoError(err)
	suite.NoError(order3.CalculateFinalPrice())
	err = repo.Save(order3)
	suite.NoError(err)

	orders, err := repo.FindAll()
	suite.NoError(err)

	suite.Equal(3, len(orders))

	suite.Equal("0001", orders[0].ID)
	suite.Equal(10.0, orders[0].Price)
	suite.Equal(2.0, orders[0].Tax)
	suite.Equal(12.0, orders[0].FinalPrice)

	suite.Equal("0002", orders[1].ID)
	suite.Equal(15.0, orders[1].Price)
	suite.Equal(3.0, orders[1].Tax)
	suite.Equal(18.0, orders[1].FinalPrice)

	suite.Equal("0003", orders[2].ID)
	suite.Equal(15.0, orders[2].Price)
	suite.Equal(3.0, orders[2].Tax)
	suite.Equal(18.0, orders[2].FinalPrice)
}
