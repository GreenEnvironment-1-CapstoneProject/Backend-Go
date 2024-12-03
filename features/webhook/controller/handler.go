package controller

import (
	"greenenvironment/features/webhook"

	"github.com/labstack/echo/v4"
)

type WebhookRequest struct {
	s webhook.MidtransNotificationService
}

func NewWebhookRequest(s webhook.MidtransNotificationService) webhook.MidtransNotificationController {
	return &WebhookRequest{
		s: s,
	}
}

func (h *WebhookRequest) HandleNotification(c echo.Context) error {
	var notification webhook.PaymentNotification
	err := c.Bind(&notification)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	err = h.s.HandleNotification(notification)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return c.JSON(200, map[string]string{
		"message": "success",
	})
}
