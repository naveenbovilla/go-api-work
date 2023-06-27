// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/

package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.com/kleene/extra-hours/docs"
	"gitlab.com/kleene/extra-hours/internal/api"
	"gitlab.com/kleene/extra-hours/internal/db"
)

func main() {
	// Create Gin instance and configure routes
	r := gin.Default()
	r.Use(cors.Default())

	postgres, err := db.Postgres{}.Open()
	if err != nil {
		// Handle the error appropriately
		fmt.Println("Failed to open the database connection:", err)
		return
	}
	//defer postgres.Con.Close()
	r.Use(postgres.Middleware())

	api.SetupRoutes(r)

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Run the server
	r.Run()
}
