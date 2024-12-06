package service

import (
	"errors"
	"greenenvironment/constant"
	"greenenvironment/features/cart"
)

type CartService struct {
	cartRepo cart.CartRepositoryInterface
}

func NewCartService(cartRepo cart.CartRepositoryInterface) cart.CartServiceInterface {
	return &CartService{cartRepo: cartRepo}
}

func (cs *CartService) Create(cart cart.NewCart) error {
	isExist, err := cs.cartRepo.IsCartExist(cart.UserID, cart.ProductID)
	if err != nil {
		return err
	}
	if isExist {
		existQty, err := cs.cartRepo.GetCartQty(cart.UserID, cart.ProductID)
		if err != nil {
			return err
		}
		stock, err := cs.cartRepo.GetStock(cart.ProductID)
		if err != nil {
			return err
		}
		if stock < (existQty+1) || stock < (existQty+cart.Quantity) {
			return errors.New("error quantity exceeds stock")
		}
		return cs.cartRepo.InsertIncrement(cart.UserID, cart.ProductID, cart.Quantity)
	}

	return cs.cartRepo.Create(cart)
}
func (cs *CartService) Update(cart cart.UpdateCart) error {
	if cart.Type != "increment" && cart.Type != "decrement" && cart.Type != "qty" {
		return constant.ErrFieldType
	}

	existQty, err := cs.cartRepo.GetCartQty(cart.UserID, cart.ProductID)
	if err != nil {
		return err
	}

	if existQty == 1 && cart.Type == "decrement" {
		return cs.cartRepo.Delete(cart.UserID, cart.ProductID)
	}

	if cart.Type == "increment" {
		return cs.cartRepo.InsertIncrement(cart.UserID, cart.ProductID, 1)
	} else if cart.Type == "decrement" {
		return cs.cartRepo.InsertDecrement(cart.UserID, cart.ProductID)
	} else if cart.Type == "qty" {
		return cs.cartRepo.InsertByQuantity(cart.UserID, cart.ProductID, cart.Quantity)
	}

	return constant.ErrFieldType
}

func (cs *CartService) Delete(userId string, productId string) error {
	return cs.cartRepo.Delete(userId, productId)
}
func (cs *CartService) Get(userId string) (cart.Cart, error) {
	return cs.cartRepo.Get(userId)
}
