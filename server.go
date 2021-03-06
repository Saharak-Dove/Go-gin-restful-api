package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Customer struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type CustomerHandler struct {
	DB *gorm.DB
}

func main() {
	router := setupRouter()
	router.Run()
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	customerHandler := CustomerHandler{}
	customerHandler.Initialize()

	router.GET("/customers", customerHandler.GetAllCustomer)
	router.GET("/customers/:id", customerHandler.GetCustomer)
	router.POST("/customers", customerHandler.SaveCustomer)
	router.PUT("/customers/:id", customerHandler.UpdateCustomer)
	router.DELETE("/customers/:id", customerHandler.DeleteCustomer)

	return router
}

func (h *CustomerHandler) Initialize() {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/gin_restful_api_dev?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})

	h.DB = db
}

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	customers := []Customer{}

	h.DB.Find(&customers)

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
