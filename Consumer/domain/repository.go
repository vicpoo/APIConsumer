// repository.go
package domain

import models "github/vicpoo/APIConsumer/Consumer/domain/entities"

type OrderRepository interface {
	Save(order models.Order) error
	GetAll() ([]models.Order, error)
}
