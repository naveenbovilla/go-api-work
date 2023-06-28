package db

import (
	"fmt"
	model "gitlab.com/kleene/extra-hours/internal/models"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Con *gorm.DB
}

func (d Postgres) Open() (*Postgres, error) {
	var err error
	//connstr := "postgres://" + os.Getenv("PG_USER") + ":" + os.Getenv("PG_PASS") + "@" + os.Getenv("PG_HOST") + "/" + os.Getenv("PG_DB_NAME")

	dsn := "host=" + os.Getenv("PG_HOST") + " user=" + os.Getenv("PG_USER") + " password=" + os.Getenv("PG_PASS") + " dbname=" + os.Getenv("PG_DB_NAME") + " port=5432 sslmode=disable"
	d.Con, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Connection error: %w", err)
	}
	model.Migrate(d.Con)
	return &d, nil
}

//func (d Postgres) Close() {
//	d.Con.Close(context.Background())
//}

func (d Postgres) Middleware() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		println("------------------------")
		// Set the global variable to the context
		c.Set("db", d)
		c.Next()
	}
	return fn
}
