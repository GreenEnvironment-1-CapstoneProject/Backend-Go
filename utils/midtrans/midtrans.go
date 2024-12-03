package midtrans

import (
	"context"
	"fmt"
	"greenenvironment/configs"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var snapClient snap.Client

type PaymentGatewayInterface interface {
	InitializeClientMidtrans()
	CreateTransaction(snap CreatePaymentGateway) string
	CreateUrlTransactionWithGateway(snap CreatePaymentGateway) string
}

type CreatePaymentGateway struct {
	Email    string
	Phone    string
	Address  string
	OrderId  string
	GrossAmt int64
	Items    []midtrans.ItemDetails
}

type PaymentGateway struct {
	conf configs.MidtransConfig
}

func NewPaymentGateway(conf configs.MidtransConfig) PaymentGatewayInterface {
	return &PaymentGateway{conf: conf}
}
func (r PaymentGateway) InitializeClientMidtrans() {
	snapClient.New(r.conf.ServerKey, midtrans.Sandbox)
}

func (r PaymentGateway) CreateTransaction(req CreatePaymentGateway) string {
	snapUrl, err := snapClient.CreateTransactionToken(generateSnapReq(req))

	if err != nil {
		fmt.Printf("Midtrans error : %v", err.GetMessage())
		return err.Error()
	}

	return snapUrl
}

func (r PaymentGateway) CreateUrlTransactionWithGateway(req CreatePaymentGateway) string {
	snapClient.Options.SetContext(context.Background())

	snapUrl, err := snapClient.CreateTransactionUrl(generateSnapReq(req))

	if err != nil {
		fmt.Printf("Midtrans error : %v", err.GetMessage())
		return err.Error()
	}

	return snapUrl
}

func generateSnapReq(req CreatePaymentGateway) *snap.Request {
	reqSnap := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderId,
			GrossAmt: req.GrossAmt,
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeBNIVA,
			snap.PaymentTypePermataVA,
			snap.PaymentTypeBCAVA,
			snap.PaymentTypeBRIVA,
			snap.PaymentTypeBankTransfer,
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: req.Email,
			Phone: req.Phone,
			ShipAddr: &midtrans.CustomerAddress{
				Phone:   req.Phone,
				Address: req.Address,
			},
		},
		Items: &req.Items,
	}

	return reqSnap
}
