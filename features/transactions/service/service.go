package service

import (
	"errors"
	"greenenvironment/features/transactions"
	midtrasService "greenenvironment/utils/midtrans"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
)

type TransactionService struct {
	transactionRepo transactions.TransactionRepositoryInterface
	midtransService midtrasService.PaymentGatewayInterface
}

func NewTransactionService(transactionRepo transactions.TransactionRepositoryInterface, midtrans midtrasService.PaymentGatewayInterface) transactions.TransactionServiceInterface {
	return &TransactionService{transactionRepo: transactionRepo, midtransService: midtrans}
}

func (ts *TransactionService) GetUserTransaction(userId string) ([]transactions.TransactionData, error) {
	transaction, err := ts.transactionRepo.GetUserTransaction(userId)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
func (ts *TransactionService) GetTransactionByID(transactionId string) (transactions.TransactionData, error) {
	transaction, err := ts.transactionRepo.GetTransactionByID(transactionId)
	if err != nil {
		return transactions.TransactionData{}, err
	}
	return transaction, nil
}
func (ts *TransactionService) CreateTransaction(transaction transactions.CreateTransaction) (transactions.Transaction, error) {
	var transactionData transactions.Transaction

	transactionData.ID = uuid.New().String()
	transactionData.UserID = transaction.UserID
	transactionData.Status = "pending"
	userData, err := ts.transactionRepo.GetUserData(transaction.UserID)
	if err != nil {
		return transactions.Transaction{}, err
	}
	cartData, err := ts.transactionRepo.GetDataCartTransaction(transaction.CartID, transaction.UserID)
	if err != nil {
		return transactions.Transaction{}, err
	}
	transactionData.Address = userData.Address

	var totalPrice float64
	items := []midtrans.ItemDetails{}
	itemsData := []transactions.TransactionItems{}

	for _, cart := range cartData {
		totalPrice += float64(int64(cart.Product.Price)) * float64(cart.Quantity)
		item := midtrans.ItemDetails{
			ID:    cart.ID,
			Name:  cart.Product.Name,
			Price: int64(cart.Product.Price),
			Qty:   int32(cart.Quantity),
		}

		itemData := transactions.TransactionItems{
			TransactionID: transactionData.ID,
			ProductID:     cart.ProductID,
			Qty:           cart.Quantity,
		}
		err = ts.transactionRepo.UpdateStockByProductID(cart.ProductID, cart.Quantity)
		if err != nil {
			return transactions.Transaction{}, err
		}
		items = append(items, item)
		itemsData = append(itemsData, itemData)
	}

	transactionData.Total = totalPrice

	if transaction.UsingCoin {
		coin, err := ts.transactionRepo.GetUserCoin(transaction.UserID)
		if err != nil {
			return transactions.Transaction{}, err
		}
		newTotal, usedCoin, err := ts.transactionRepo.DecreaseUserCoin(transaction.UserID, coin, transactionData.Total)
		if err != nil {
			return transactions.Transaction{}, err
		}
		item := midtrans.ItemDetails{
			ID:    uuid.New().String(),
			Name:  "used-coin",
			Price: int64(-usedCoin),
			Qty:   int32(1),
		}

		items = append(items, item)

		transactionData.Coin = usedCoin
		transactionData.Total = newTotal
	}

	if transactionData.Total <= 0 {
		return transactions.Transaction{}, errors.New("error gross amount must be greater than 0")
	}

	ts.midtransService.InitializeClientMidtrans()

	snapReq := midtrasService.CreatePaymentGateway{
		OrderId:  transactionData.ID,
		Email:    userData.Email,
		Phone:    userData.Phone,
		Address:  userData.Address,
		GrossAmt: int64(transactionData.Total),
		Items:    items,
	}

	snapUrl := ts.midtransService.CreateTransaction(snapReq)

	transactionData.SnapURL = snapUrl

	error := ts.transactionRepo.CreateTransactions(transactionData)
	if error != nil {
		return transactions.Transaction{}, error
	}

	err = ts.transactionRepo.CreateTransactionItems(itemsData)
	if err != nil {
		return transactions.Transaction{}, err
	}

	return transactionData, nil
}
func (ts *TransactionService) DeleteTransaction(transactionId string) error {
	_, err := ts.GetTransactionByID(transactionId)
	if err != nil {
		return err
	}
	return ts.transactionRepo.DeleteTransaction(transactionId)
}
func (ts *TransactionService) GetAllTransaction() ([]transactions.TransactionData, error) {
	return ts.transactionRepo.GetAllTransaction()
}
