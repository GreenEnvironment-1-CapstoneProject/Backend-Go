package repository

import (
	"errors"
	"greenenvironment/features/cart"
	productData "greenenvironment/features/products/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) cart.CartRepositoryInterface {
	return &CartRepository{DB: db}
}

func (cr *CartRepository) Create(cart cart.NewCart) error {
	var dbQty int

	stock, err := cr.GetStock(cart.ProductID)
	if err != nil {
		return err
	}
	if cart.Quantity >= 0 {
		dbQty = cart.Quantity
	} else {
		dbQty = 1
	}

	if stock < dbQty {
		return errors.New("error quantity exceeds stock")
	}
	newCart := &Cart{
		ID:        uuid.New().String(),
		UserID:    cart.UserID,
		ProductID: cart.ProductID,
		Quantity:  dbQty,
	}

	err = cr.DB.Create(newCart).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *CartRepository) Update(cart cart.UpdateCart) error {
	err := cr.DB.Model(&Cart{}).Where("user_id = ?", cart.UserID).Updates(cart).Error
	if err != nil {
		return err
	}

	return nil
}
func (cr *CartRepository) Delete(userId string, productId string) error {
	err := cr.DB.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&Cart{}).Error

	if err != nil {
		return err
	}
	return nil

}
func (cr *CartRepository) Get(userId string) (cart.Cart, error) {
	var carts []Cart
	var cartData cart.Cart

	err := cr.DB.Model(&Cart{}).
		Preload("Product").
		Preload("User").
		Preload("Product.Images").
		Preload("Product.ImpactCategories").
		Preload("Product.ImpactCategories.ImpactCategory").
		Where("user_id = ?", userId).Find(&carts).Error

	if err != nil {
		return cartData, err
	}

	if len(carts) == 0 {
		return cart.Cart{}, nil
	}

	cartData.User.ID = userId
	cartData.User.Username = carts[0].User.Username
	cartData.User.Email = carts[0].User.Email
	cartData.User.Address = carts[0].User.Address
	cartData.User.Phone = carts[0].User.Phone

	for _, cart := range carts {
		cartEntity := cart.ToEntity()
		cartData.Items = append(cartData.Items, cartEntity.Items...)
	}

	return cartData, nil
}

func (cr *CartRepository) IsCartExist(userId string, productId string) (bool, error) {
	var count int64
	err := cr.DB.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userId, productId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (cr *CartRepository) InsertIncrement(userId string, productId string, qty int) error {
	var dbQty int

	if qty >= 0 {
		dbQty = qty
	} else {
		dbQty = 1
	}

	return cr.DB.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userId, productId).Update("quantity", gorm.Expr("quantity + ?", dbQty)).Error
}

func (cr *CartRepository) InsertDecrement(userId string, productId string) error {
	return cr.DB.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userId, productId).Update("quantity", gorm.Expr("quantity - 1")).Error
}

func (c *CartRepository) GetCartQty(userId string, productId string) (int, error) {
	var cart Cart
	err := c.DB.Where("user_id = ? AND product_id = ?", userId, productId).First(&cart).Error
	return cart.Quantity, err
}

func (c *CartRepository) InsertByQuantity(userId string, productId string, quantity int) error {
	stock, err := c.GetStock(productId)
	if err != nil {
		return err
	}
	if stock < quantity {
		return errors.New("error quantity exceeds stock")
	}

	return c.DB.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userId, productId).Update("quantity", quantity).Error
}

func (c *CartRepository) GetStock(productId string) (int, error) {
	var product productData.Product
	err := c.DB.Where("id = ?", productId).First(&product).Error

	return product.Stock, err
}
