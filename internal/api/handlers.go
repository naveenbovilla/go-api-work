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
	"time"

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
		api.POST("/addExtraHours", addExtraHours)
		api.PUT("/updateExtraHours/:id", updateExtraHours)
		api.DELETE("/deleteExtraHours/:id", deleteExtraHoursWithServiceId)
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

type AddHours struct {
	Hours      int    `json:"hours"`
	TypeOfWork string `json:"type_of_work"`
	Notes      string `json:"notes"`
}

func addExtraHours(c *gin.Context) {
	dbw, _ := c.Get("db")
	dbwc := dbw.(db.Postgres)
	db := dbwc.Con
	body := AddHours{}
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var extraHours model.ExtraHours
	extraHours.HoursRequested = uint(body.Hours)
	extraHours.Notes = body.Notes
	extraHours.TypeOfWork = body.TypeOfWork
	extraHours.WaitingDate = time.Now().Truncate(time.Minute)
	extraHours.Status = "pending"
	extraHours.CreatedAt = time.Now()
	serviceID, _ := model.GetNextServiceID(db)
	extraHours.ServiceID = serviceID
	if err := db.Create(&extraHours).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, extraHours)
}

type UpdateHours struct {
	Hours      int    `json:"hours"`
	TypeOfWork string `json:"type_of_work"`
	Notes      string `json:"notes"`
	Status     string `json:"status"`
}

func updateExtraHours(c *gin.Context) {
	dbw, _ := c.Get("db")
	dbwc := dbw.(db.Postgres)
	db := dbwc.Con
	body := UpdateHours{}
	id := c.Params.ByName("id")
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var extraHours model.ExtraHours
	if err := db.Where("service_id = ?", id).First(&extraHours).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	extraHours.HoursRequested = uint(body.Hours)
	extraHours.Notes = body.Notes
	extraHours.TypeOfWork = body.TypeOfWork
	extraHours.WaitingDate = time.Now().Truncate(time.Minute)
	extraHours.Status = body.Status
	extraHours.UpdatedAt = time.Now()
	if err := db.Save(&extraHours).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		fmt.Println(err)
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
	if err := db.Where("service_id = ?", id).First(&extraHours).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, extraHours)
}

func deleteExtraHoursWithServiceId(c *gin.Context) {
	dbw, _ := c.Get("db")
	dbwc := dbw.(db.Postgres)
	db := dbwc.Con
	id := c.Params.ByName("id")
	var extraHours model.ExtraHours
	if err := db.Where("service_id = ?", id).Delete(&extraHours).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	c.Status(http.StatusOK)
}
