// controller.go
package infrastructure

import (
	"net/http"

	"github/vicpoo/APIConsumer/Consumer/application"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderUseCase *application.OrderUseCase
}

func NewOrderController(orderUseCase *application.OrderUseCase) *OrderController {
	return &OrderController{
		orderUseCase: orderUseCase,
	}
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	orders, err := oc.orderUseCase.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
