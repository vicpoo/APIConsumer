// consumer_usecase.go
package application

import (
	"github/vicpoo/APIConsumer/Consumer/domain"
	models "github/vicpoo/APIConsumer/Consumer/domain/entities"
)

type OrderUseCase struct {
	repo domain.OrderRepository
}

func NewOrderUseCase(repo domain.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) SaveOrder(order models.Order) error {
	return uc.repo.Save(order)
}

func (uc *OrderUseCase) GetAllOrders() ([]models.Order, error) {
	return uc.repo.GetAll()
}
