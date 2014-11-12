package api

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Secret-Ironman/boxr/shared/db"
	"github.com/Secret-Ironman/boxr/shared/types"
	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	Message interface{}   `json:"message" binding:"required"`
	Took    time.Duration `json:"took" binding:"required"`
	Success bool          `json:"success"`
}

type Api struct {
	db   *gorp.DbMap
	Port int
}

func NewApi(dbFile string, port int) (a *Api, err error) {
	a = new(Api)
	a.db, err = db.New(dbFile)
	checkErr(err, "Create tables failed")
	a.Port = port
	return a, nil
}

func (c *Api) Run() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/commits/:id", c.CommitsHookGet)
		api.GET("/pallets", c.PalletGetAll)
		api.GET("/pallets/:name", c.PalletGetOne)
		api.GET("/pallets/:name/build", c.PalletBuild)
		api.POST("/pallets", c.PalletCreate)
	}

	r.Run(fmt.Sprintf(":%d", c.Port))
}

func (c *Api) initDb(dbFile string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", dbFile)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(types.Pallet{}, "pallets").SetKeys(false, "Name")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
