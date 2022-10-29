package pgdatabase

import (
	"context"
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

	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	if err != nil {
		log.Fatal(err)
	}

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	DBURL = fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb?sslmode=disable", host, port.Port())

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

func Test_GetCurrency(t *testing.T) {

	BeforeTest()
	tests := []struct {
		name     string
		prepareF func()
		err      error
		ctype    model.Currency
	}{
		{
			name:     "get successfully",
			prepareF: func() {},
			ctype:    model.Currency{Code: "rub", Ratio: decimal.NewFromInt(1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			ctype, err := storage.GetCurrentCurrency(context.Background())
			assert.ErrorIs(t, err, tt.err)

			assert.Equal(t, ctype.Code, tt.ctype.Code)
			assert.True(t, ctype.Ratio.Equal(tt.ctype.Ratio))
		})
	}
}

func Test_dbCurrencyStorage_UpdateCurrentCurrency(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		prepareF func()
		checkF   func(err error)
	}{
		{
			"updated successfully",
			"eur",
			func() {
				DB.MustExec("insert into currencies(code, ratio) values('eur', 1)")
			},
			func(err error) {
				checkIsExist(t, "select count(1) from state where current_currency_code = 'eur'", 1)
			},
		},
		{
			name: "wrong currency type",

			code:     "undefined",
			prepareF: func() {},
			checkF: func(err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			tt.checkF(storage.UpdateCurrentCurrency(tt.code))
		})
	}
}

func Test_UpdateCurrencies(t *testing.T) {
	tests := []struct {
		name     string
		data     []model.Currency
		prepareF func()
		checkF   func(err error)
	}{
		{
			name: "upinsert successfully",
			data: []model.Currency{
				{Code: "cny", Ratio: decimal.NewFromInt(1)},
				{Code: "eur", Ratio: decimal.NewFromInt(2)},
			},
			prepareF: func() {
				DB.MustExec("insert into currencies values('cny', 0)")
			},
			checkF: func(err error) {
				assert.ErrorIs(t, err, nil)
				checkIsExist(t, "select count(1) from currencies where (code = 'cny' and ratio = 1) or (code = 'eur' and ratio = 2)", 2)
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
