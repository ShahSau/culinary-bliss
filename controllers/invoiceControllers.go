package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetInvoices(c *gin.Context) {
	fmt.Println("GetInvoices")
}

func GetInvoice(c *gin.Context) {
	fmt.Println("GetInvoice")
}

func CreateInvoice(c *gin.Context) {
	fmt.Println("CreateInvoice")
}

func UpdateInvoice(c *gin.Context) {
	fmt.Println("UpdateInvoice")
}

func DeleteInvoice(c *gin.Context) {
	fmt.Println("DeleteInvoice")
}
