package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(c *gin.Engine) {
	c.GET("/invoice", controllers.GetInvoices) //admin
	c.GET("/invoice/:id", controllers.GetInvoice)
	c.POST("/invoice", controllers.CreateInvoice)
	c.PUT("/invoice/:id", controllers.UpdateInvoice)    //admin
	c.DELETE("/invoice/:id", controllers.DeleteInvoice) //admin
}
