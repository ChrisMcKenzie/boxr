package db

import (
	"database/sql"

	"github.com/Secret-Ironman/boxr/shared/types"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

func New(file string) (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", file)
	// checkErr(err, "sql.Open failed")
	if err != nil {
		return nil, err
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(types.Pallet{}, "pallets").SetKeys(false, "Name")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	// checkErr(err, "Create tables failed")
	if err != nil {
		return nil, err
	}

	return dbmap, err
}
