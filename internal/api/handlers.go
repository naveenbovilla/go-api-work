package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "gitlab.com/kleene/extra-hours/internal/auth"
	"gitlab.com/kleene/extra-hours/internal/db"
	model "gitlab.com/kleene/extra-hours/internal/models"
	"net/http"
	_ "net/http"
	_ "net/netip"

	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
)

// SetupRoutes sets up the routes for the API
func SetupRoutes(r *gin.Engine) {
	//ka := auth.KeycloakAuth()
	api := r.Group("/api")

	{
		// README
		// TODO crud apis + all the needed api for the workflow.
		// Bonus: Once accepted, a payment needs to be schedule using POST api to another STRIPE_SERVICE
		//
		api.GET("/getAll", getExtraHours)
		api.GET("/getExtraHours/:id", getExtraHoursWithServiceId)
		private := api.Group("/private")
		println(private.BasePath())
		//private.Use(ka)
		{
			//private.DELETE("customer/grouprequest/:id", deleteGroupRequest)
			//private.PATCH("customer/grouprequest/:id", patchGroupRequest)
		}
	}
}

func getExtraHours(c *gin.Context) {
	dbw, _ := c.Get("db")
	dbwc := dbw.(db.Postgres)
	db := dbwc.Con
	// get tokens
	//tc, _ := c.Get("token")
	//tcs, _ := c.Get("tokenStr")
	//_ = tcs.(string)
	//token := tc.(jwt.Token)
	//userID, _ := tc.Get("sub")
	// TODO to fill ..
	var extraHours []model.ExtraHours
	if result := db.Unscoped().Find(&extraHours); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	c.JSON(http.StatusOK, extraHours)
}

func getExtraHoursWithServiceId(c *gin.Context) {
	dbw, _ := c.Get("db")
	dbwc := dbw.(db.Postgres)
	db := dbwc.Con
	id := c.Params.ByName("id")
	var extraHours model.ExtraHours
	if err := db.Where("serviceid = ?", id).First(&extraHours).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, extraHours)
}
