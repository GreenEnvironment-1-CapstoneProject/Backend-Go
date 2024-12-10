package repository

import (
	cart "greenenvironment/features/cart/repository"
	"greenenvironment/features/impacts"
	"greenenvironment/features/products"
	productRepo "greenenvironment/features/products/repository"
	"greenenvironment/features/transactions"
	users "greenenvironment/features/users/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) transactions.TransactionRepositoryInterface {
	return &TransactionRepository{DB: db}
}

func (tr *TransactionRepository) GetUserTransaction(userId string) ([]transactions.TransactionData, error) {
	var transactionsModel []Transaction
	err := tr.DB.Model(&Transaction{}).Preload("User").Preload("TransactionItems").Preload("TransactionItems.Product").
		Preload("TransactionItems.Product.Images").
		Preload("TransactionItems.Product.ImpactCategories").
		Preload("TransactionItems.Product.ImpactCategories.ImpactCategory").Where("user_id = ?", userId).
		Find(&transactionsModel).Error

	if err != nil {
		return nil, err
	}

	var result []transactions.TransactionData
	for _, txn := range transactionsModel {
		var txnItems []transactions.TransactionItems
		var images []products.ProductImage
		var impactCategories []products.ProductImpactCategory
		for _, item := range txn.TransactionItems {
			for _, img := range item.Product.Images {
				images = append(images, products.ProductImage{
					ID:        img.ID,
					ProductID: img.ProductID,
					AlbumsURL: img.AlbumsURL,
				})
			}
			for _, impact := range item.Product.ImpactCategories {
				impactCategories = append(impactCategories, products.ProductImpactCategory{
					ID:               impact.ID,
					ProductID:        impact.ProductID,
					ImpactCategoryID: impact.ImpactCategoryID,
					ImpactCategory: impacts.ImpactCategory{
						ID:          impact.ImpactCategory.ID,
						Name:        impact.ImpactCategory.Name,
						ImpactPoint: impact.ImpactCategory.ImpactPoint,
						Description: impact.ImpactCategory.Description,
					},
				})
			}
			txnItems = append(txnItems, transactions.TransactionItems{
				ID:            item.ID,
				TransactionID: item.TransactionID,
				ProductID:     item.ProductID,
				Qty:           item.Quantity,
				Product: products.Product{
					ID:               item.Product.ID,
					Name:             item.Product.Name,
					Description:      item.Product.Description,
					Price:            item.Product.Price,
					Coin:             item.Product.Coin,
					Stock:            item.Product.Stock,
					Images:           images,
					ImpactCategories: impactCategories,
				},
			})
		}

		result = append(result, transactions.TransactionData{
			ID:      txn.ID,
			Status:  txn.Status,
			Total:   txn.Total,
			Coin:    txn.Coin,
			SnapURL: txn.SnapURL,
			User: users.User{
				ID:        txn.User.ID,
				Name:      txn.User.Name,
				Email:     txn.User.Email,
				Username:  txn.User.Username,
				Password:  txn.User.Password,
				Address:   txn.User.Address,
				Gender:    txn.User.Gender,
				Phone:     txn.User.Phone,
				Coin:      txn.User.Coin,
				Exp:       txn.User.Exp,
				AvatarURL: txn.User.AvatarURL,
			},
			TransactionItems: txnItems,
			CreatedAt:        txn.CreatedAt,
			UpdatedAt:        txn.UpdatedAt,
		})
	}

	return result, nil
}
func (tr *TransactionRepository) GetTransactionByID(transactionId string) (transactions.TransactionData, error) {
	var transactionsData Transaction
	err := tr.DB.Model(&Transaction{}).Preload("User").Preload("TransactionItems").Preload("TransactionItems.Product").
		Preload("TransactionItems.Product.Images").
		Preload("TransactionItems.Product.ImpactCategories").
		Preload("TransactionItems.Product.ImpactCategories.ImpactCategory").Where("id = ?", transactionId).
		Take(&transactionsData).Error

	if err != nil {
		return transactions.TransactionData{}, err
	}

	var result transactions.TransactionData
	var txnItems []transactions.TransactionItems
	var images []products.ProductImage
	var impactCategories []products.ProductImpactCategory
	for _, item := range transactionsData.TransactionItems {
		for _, img := range item.Product.Images {
			images = append(images, products.ProductImage{
				ID:        img.ID,
				ProductID: img.ProductID,
				AlbumsURL: img.AlbumsURL,
			})
		}
		for _, impact := range item.Product.ImpactCategories {
			impactCategories = append(impactCategories, products.ProductImpactCategory{
				ID:               impact.ID,
				ProductID:        impact.ProductID,
				ImpactCategoryID: impact.ImpactCategoryID,
				ImpactCategory: impacts.ImpactCategory{
					ID:          impact.ImpactCategory.ID,
					Name:        impact.ImpactCategory.Name,
					ImpactPoint: impact.ImpactCategory.ImpactPoint,
					Description: impact.ImpactCategory.Description,
				},
			})
		}
		txnItems = append(txnItems, transactions.TransactionItems{
			ID:            item.ID,
			TransactionID: item.TransactionID,
			ProductID:     item.ProductID,
			Qty:           item.Quantity,
			Product: products.Product{
				ID:               item.Product.ID,
				Name:             item.Product.Name,
				Description:      item.Product.Description,
				Price:            item.Product.Price,
				Coin:             item.Product.Coin,
				Stock:            item.Product.Stock,
				Images:           images,
				ImpactCategories: impactCategories,
			},
		})
	}

	result = transactions.TransactionData{
		ID:      transactionsData.ID,
		Status:  transactionsData.Status,
		Total:   transactionsData.Total,
		Coin:    transactionsData.Coin,
		SnapURL: transactionsData.SnapURL,
		User: users.User{
			ID:        transactionsData.User.ID,
			Name:      transactionsData.User.Name,
			Email:     transactionsData.User.Email,
			Username:  transactionsData.User.Username,
			Password:  transactionsData.User.Password,
			Address:   transactionsData.User.Address,
			Gender:    transactionsData.User.Gender,
			Phone:     transactionsData.User.Phone,
			Coin:      transactionsData.User.Coin,
			Exp:       transactionsData.User.Exp,
			AvatarURL: transactionsData.User.AvatarURL,
		},
		TransactionItems: txnItems,
		CreatedAt:        transactionsData.CreatedAt,
		UpdatedAt:        transactionsData.UpdatedAt,
	}
	return result, nil
}

func (tr *TransactionRepository) CreateTransactions(transaction transactions.Transaction) error {
	transactionData := Transaction{
		ID:            transaction.ID,
		Address:       transaction.Address,
		UserID:        transaction.UserID,
		Total:         transaction.Total,
		Status:        transaction.Status,
		PaymentMethod: transaction.PaymentMethod,
		Coin:          transaction.Coin,
		SnapURL:       transaction.SnapURL,
	}

	err := tr.DB.Create(&transactionData).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransactionRepository) DeleteTransaction(transactionId string) error {
	err := tr.DB.Where("id = ?", transactionId).Delete(&Transaction{}).Error

	if err != nil {
		return err
	}
	return nil
}
func (tr *TransactionRepository) GetUserData(userId string) (users.User, error) {
	var user users.User
	err := tr.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
func (tr *TransactionRepository) GetUserCoin(userId string) (int, error) {
	var user users.User
	err := tr.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.Coin, nil
}
func (tr *TransactionRepository) DecreaseUserCoin(userId string, coin int, total float64) (float64, int, error) {
	var user users.User
	err := tr.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return total, 0, err
	}

	maxCoin := int(total * 0.80)

	usedCoin := coin
	if usedCoin > maxCoin {
		usedCoin = maxCoin
	}

	if user.Coin < usedCoin {
		usedCoin = user.Coin
	}

	newTotal := total - float64(usedCoin)

	user.Coin -= usedCoin
	err = tr.DB.Save(&user).Error
	if err != nil {
		return total, 0, err
	}

	return newTotal, usedCoin, nil
}
func (tr *TransactionRepository) CreateTransactionItems(tansactionItems []transactions.TransactionItems) error {

	tx := tr.DB.Begin()
	for _, tansactionItem := range tansactionItems {
		transactionItemId := uuid.New().String()
		transactionItem := TransactionItem{
			ID:            transactionItemId,
			TransactionID: tansactionItem.TransactionID,
			ProductID:     tansactionItem.ProductID,
			Quantity:      tansactionItem.Qty,
		}

		err := tr.DB.Create(&transactionItem).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	}

	return nil
}

func (tr *TransactionRepository) GetAllTransaction() ([]transactions.TransactionData, error) {
	var transactionsData []Transaction

	err := tr.DB.Model(&Transaction{}).Preload("User").
		Preload("TransactionItems").
		Preload("TransactionItems.Product").
		Preload("TransactionItems.Product.Images").
		Preload("TransactionItems.Product.ImpactCategories").
		Preload("TransactionItems.Product.ImpactCategories.ImpactCategory").
		Order("created_at DESC").
		Find(&transactionsData).Error

	if err != nil {
		return nil, err
	}

	var result []transactions.TransactionData
	for _, txn := range transactionsData {
		var txnItems []transactions.TransactionItems
		var images []products.ProductImage
		var impactCategories []products.ProductImpactCategory
		for _, item := range txn.TransactionItems {
			for _, img := range item.Product.Images {
				images = append(images, products.ProductImage{
					ID:        img.ID,
					ProductID: img.ProductID,
					AlbumsURL: img.AlbumsURL,
				})
			}
			for _, impact := range item.Product.ImpactCategories {
				impactCategories = append(impactCategories, products.ProductImpactCategory{
					ID:               impact.ID,
					ProductID:        impact.ProductID,
					ImpactCategoryID: impact.ImpactCategoryID,
					ImpactCategory: impacts.ImpactCategory{
						ID:          impact.ImpactCategory.ID,
						Name:        impact.ImpactCategory.Name,
						ImpactPoint: impact.ImpactCategory.ImpactPoint,
						Description: impact.ImpactCategory.Description,
					},
				})
			}
			txnItems = append(txnItems, transactions.TransactionItems{
				ID:            item.ID,
				TransactionID: item.TransactionID,
				ProductID:     item.ProductID,
				Qty:           item.Quantity,
				Product: products.Product{
					ID:               item.Product.ID,
					Name:             item.Product.Name,
					Description:      item.Product.Description,
					Price:            item.Product.Price,
					Coin:             item.Product.Coin,
					Stock:            item.Product.Stock,
					Images:           images,
					ImpactCategories: impactCategories,
				},
			})
		}

		result = append(result, transactions.TransactionData{
			ID:      txn.ID,
			Status:  txn.Status,
			Total:   txn.Total,
			Coin:    txn.Coin,
			SnapURL: txn.SnapURL,
			User: users.User{
				ID:        txn.User.ID,
				Name:      txn.User.Name,
				Email:     txn.User.Email,
				Username:  txn.User.Username,
				Password:  txn.User.Password,
				Address:   txn.User.Address,
				Gender:    txn.User.Gender,
				Phone:     txn.User.Phone,
				Coin:      txn.User.Coin,
				Exp:       txn.User.Exp,
				AvatarURL: txn.User.AvatarURL,
			},
			TransactionItems: txnItems,
			CreatedAt:        txn.CreatedAt,
			UpdatedAt:        txn.UpdatedAt,
		})
	}

	return result, nil
}

func (tr *TransactionRepository) GetDataCartTransaction(cartIds []string, userId string) ([]cart.Cart, error) {

	var result []cart.Cart
	for _, cartId := range cartIds {
		var cart cart.Cart
		err := tr.DB.Preload("Product").Where("id = ? AND user_id = ?", cartId, userId).Find(&cart).Error
		if err != nil {
			return nil, err
		}

		result = append(result, cart)
	}

	return result, nil
}

func (tr *TransactionRepository) UpdateStockByProductID(productId string, stock int) error {
	return tr.DB.Model(productRepo.Product{}).Where("id = ?", productId).Update("stock", gorm.Expr("stock - ?", stock)).Error
}
