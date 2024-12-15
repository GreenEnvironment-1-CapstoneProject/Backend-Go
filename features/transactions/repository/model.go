package repository

import (
	products "greenenvironment/features/products/repository"
	users "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

type Transaction struct {
	*gorm.Model
	ID               string            `gorm:"primary_key;type:varchar(50);not null;column:id"`
	UserID           string            `gorm:"type:varchar(50);not null;column:user_id"`
	Address          string            `gorm:"type:varchar(255);not null;column:address"`
	Total            float64           `gorm:"type:decimal(10,2);not null;column:total"`
	Status           string            `gorm:"type:varchar(50);not null;column:status"`
	PaymentMethod    string            `gorm:"type:varchar(50);column:payment_method"`
	SnapURL          string            `gorm:"type:varchar(255);not null;column:snap_url"`
	Coin             int               `gorm:"type:int;column:coin"`
	User             users.User        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionItems []TransactionItem `gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransactionItem struct {
	*gorm.Model
	ID            string           `gorm:"primary_key;type:varchar(50);not null;column:id"`
	TransactionID string           `gorm:"type:varchar(50);not null;column:transaction_id"`
	ProductID     string           `gorm:"type:varchar(50);not null;column:product_id"`
	Quantity      int              `gorm:"type:int;not null;column:quantity"`
	Transaction   Transaction      `gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product       products.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Transaction) TableName() string {
	return "transactions"
}
func (TransactionItem) TableName() string {
	return "transaction_items"
}
