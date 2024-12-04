package databases

import (
	DataAdmin "greenenvironment/features/admin/repository"
	DataCart "greenenvironment/features/cart/repository"
	DataChatbot "greenenvironment/features/chatbot/repository"
	DataImpact "greenenvironment/features/impacts/repository"
	DataProduct "greenenvironment/features/products/repository"
	DataReview "greenenvironment/features/review_products/repository"
	DataTransaction "greenenvironment/features/transactions/repository"
	DataUser "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&DataUser.User{})
	db.AutoMigrate(&DataAdmin.Admin{})
	db.AutoMigrate(&DataImpact.ImpactCategory{})
	db.AutoMigrate(&DataProduct.Product{})
	db.AutoMigrate(&DataProduct.ProductImage{})
	db.AutoMigrate(&DataProduct.ProductImpactCategory{})
	db.AutoMigrate(&DataProduct.ProductLog{})
	db.AutoMigrate(&DataCart.Cart{})
	db.AutoMigrate(&DataTransaction.Transaction{})
	db.AutoMigrate(&DataTransaction.TransactionItem{})
	db.AutoMigrate(&DataReview.ReviewProduct{})
	db.AutoMigrate(&DataChatbot.Chatbot{})

	return nil
}
