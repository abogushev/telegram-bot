package pgdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase/migrations"
)

var storage *dbCurrencyStorage
var DBURL string
var DB *sqlx.DB

func SetupTestDatabase() testcontainers.Container {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	DBURL = fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb?sslmode=disable", host, port.Port())

	var err error

	DB, err = InitDB(context.Background(), DBURL)

	if err != nil {
		log.Fatalf("error in connection to db %v", err)
	}
	return dbContainer
}

func TestMain(m *testing.M) {
	container := SetupTestDatabase()
	defer func() { _ = container.Terminate(context.Background()) }()
	migrations.Up(DBURL, "migrations")
	storage = NewCurrencyStorage(context.Background(), DB)
	code := m.Run()

	os.Exit(code)

}

func BeforeTest() {
	migrations.Down(DBURL, "migrations")
	migrations.Up(DBURL, "migrations")
}

func checkIsExist(t *testing.T, q string, expectedCount int) {
	r := 0
	assert.NoError(t, DB.Get(&r, q))
	assert.Equal(t, expectedCount, r)
}

func Test_GetCurrencyType(t *testing.T) {

	BeforeTest()
	tests := []struct {
		name     string
		prepareF func()
		err      error
		ctype    model.CurrencyType
	}{
		{
			"get successfully",
			func() {},
			nil,
			model.RUB,
		},
		{
			"wrong persisted currency: cant convert",
			func() {
				DB.MustExec("insert into currencies(code, ratio) values('unexpected', 1)")
				DB.MustExec("update state set current_currency_code = 'unexpected'")
			},
			model.ErrWrongCurrencyType,
			model.Undefined,
		},
		{
			"empty data",
			func() {
				DB.MustExec("truncate table state")
			},
			sql.ErrNoRows,
			model.Undefined,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			ctype, err := storage.GetCurrencyType()
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, ctype, tt.ctype)
		})
	}
}

func Test_dbCurrencyStorage_UpdateCurrentType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		data     model.CurrencyType
		prepareF func()
		checkF   func()
	}{
		{
			"updated successfully",
			nil,
			model.EUR,
			func() {
				DB.MustExec("insert into currencies(code, ratio) values('eur', 1)")
			},
			func() {
				checkIsExist(t, "select count(1) from state where current_currency_code = 'eur'", 1)
			},
		},
		{
			"wrong currency type",
			model.ErrWrongCurrencyType,
			model.CurrencyType(99),
			func() {},
			func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			err := storage.UpdateCurrentType(tt.data)
			assert.ErrorIs(t, err, tt.err)
			tt.checkF()
		})
	}
}

func Test_UpdateCurrencies(t *testing.T) {
	tests := []struct {
		name     string
		data     map[model.CurrencyType]decimal.Decimal
		prepareF func()
		checkF   func(err error)
	}{
		{
			"upinsert successfully",
			map[model.CurrencyType]decimal.Decimal{
				model.CNY: decimal.NewFromInt(1),
				model.EUR: decimal.NewFromInt(2),
			},
			func() {
				DB.MustExec("insert into currencies values('cny', 0) ")
			},
			func(err error) {
				assert.ErrorIs(t, err, nil)
				checkIsExist(t, "select count(1) from currencies where code = 'cny' or code = 'eur'", 2)
			},
		},
		{
			"failed on wrong type",
			map[model.CurrencyType]decimal.Decimal{
				model.CNY:              decimal.NewFromInt(1),
				model.CurrencyType(99): decimal.NewFromInt(2),
			},
			func() {},
			func(err error) {
				assert.ErrorContains(t, err, "wrong сurrency type: 99")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			tt.checkF(storage.UpdateCurrencies(tt.data))
		})
	}
}
