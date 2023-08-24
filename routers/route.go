package routers

import (
	"assignment-project/controllers"
	"assignment-project/database"

	"github.com/gin-gonic/gin"
)

func init (){
	database.StartDB()

}

func StartServer() *gin.Engine{
	router := gin.Default()

	router.POST("/student", controllers.CreatedStudent)
	router.GET("/students/", controllers.GetAllStudent)
	router.PUT("/student/:ID", controllers.UpdateOrder)
	router.DELETE("/student/:ID", controllers.DeleteStudent)


	return router
}