package DBManager

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"

	//"github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/lib/driver"
)

const (
	POSTGRES = iota
	CLICK_HOUSE
	MIXED
)

const (
	ch_host     = "localhost"
	ch_port     = 9000
	ch_user     = "default"
	ch_password = ""
	ch_dbname   = "tutorial"
)

const (
	pg_host     = "localhost"
	pg_port     = 5432
	pg_user     = "dimplom"
	pg_password = ""
	pg_dbname   = "monitoring"
)

type DBManger struct {
	postgres   *sql.DB
	clickHouse driver.Conn
	T          string // test
}

func (dbm *DBManger) ConnectPostgres() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
	var err error
	var db *sql.DB

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Postgres: Successfully connected!")
	return db
}

func (dbm *DBManger) connectClickHouse() {
	var err error
	dbm.clickHouse, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", ch_host, ch_port)},
		Auth: clickhouse.Auth{
			Database: ch_dbname,
			Username: ch_user,
			Password: ch_password,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// TEST THAT CONNECTION IS OK
	v, err := dbm.clickHouse.ServerVersion()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Server2:", v)
	}
}

func (dbm *DBManger) Init(mode int) *DBManger {
	switch mode {
	case POSTGRES:
		dbm.postgres = dbm.ConnectPostgres()
		if dbm.postgres == nil {
			fmt.Println("Postgres not connected")
		}
		break
	case CLICK_HOUSE:
		dbm.connectClickHouse()
		break
	case MIXED:
		break
	default:
		fmt.Println("Chose correct mode")
		break
	}

	return dbm
}
