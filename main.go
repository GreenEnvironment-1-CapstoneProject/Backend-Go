package main

import (
	"greenenvironment/configs"
	_ "greenenvironment/docs"
	"greenenvironment/helper"

	AdminContoller "greenenvironment/features/admin/controller"
	AdminRepository "greenenvironment/features/admin/repository"
	AdminService "greenenvironment/features/admin/service"
	CartController "greenenvironment/features/cart/controller"
	CartRepository "greenenvironment/features/cart/repository"
	CartService "greenenvironment/features/cart/service"
	GuestController "greenenvironment/features/guest/controller"
	GuestRepository "greenenvironment/features/guest/repository"
	guestService "greenenvironment/features/guest/service"
	ImpactController "greenenvironment/features/impacts/controller"
	ImpactRepository "greenenvironment/features/impacts/repository"
	ImpactService "greenenvironment/features/impacts/service"
	ProductController "greenenvironment/features/products/controller"
	ProductRepository "greenenvironment/features/products/repository"
	ProductService "greenenvironment/features/products/service"
	ReviewController "greenenvironment/features/review_products/controller"
	ReviewRepository "greenenvironment/features/review_products/repository"
	ReviewService "greenenvironment/features/review_products/service"
	TransactionController "greenenvironment/features/transactions/controller"
	TransactionRepository "greenenvironment/features/transactions/repository"
	TransactionService "greenenvironment/features/transactions/service"
	UserController "greenenvironment/features/users/controller"
	UserRepository "greenenvironment/features/users/repository"
	UserService "greenenvironment/features/users/service"
	WebhookController "greenenvironment/features/webhook/controller"
	WebHookRepository "greenenvironment/features/webhook/repository"
	WebhookService "greenenvironment/features/webhook/service"

	ChatbotController "greenenvironment/features/chatbot/controller"
	ChatbotRepository "greenenvironment/features/chatbot/repository"
	ChatbotService "greenenvironment/features/chatbot/service"

	"greenenvironment/routes"
	"greenenvironment/utils/databases"
	"greenenvironment/utils/midtrans"
	OpenAIservice "greenenvironment/utils/openai"
	"greenenvironment/utils/storages"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title capstone project green environment
// @version 1.0
// @description This is a sample server Swagger server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://greenenvironment.my.id
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := configs.InitConfig()
	db, err := databases.InitDB(*cfg)
	if err != nil {
		logrus.Error("terjadi kesalahan pada database, error:", err.Error())
	}

	databases.Migrate(db)
	jwt := helper.NewJWT(cfg.JWT_Secret)
	storage := storages.NewStorage(cfg.Cloudinary)
	midtransService := midtrans.NewPaymentGateway(cfg.Midtrans)
	openAIservice := OpenAIservice.NewOpenAIService(cfg.OpenAi.ApiKey)

	e := echo.New()
	e.Validator = &helper.CustomValidator{Validator: validator.New()}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	userRepo := UserRepository.NewUserRepository(db)
	userService := UserService.NewUserService(userRepo, jwt)
	userController := UserController.NewUserController(userService, jwt, storage)

	adminRepo := AdminRepository.NewAdminRepository(db)
	adminService := AdminService.NewAdminService(adminRepo, jwt)
	adminController := AdminContoller.NewAdminController(adminService, jwt)

	impactRepo := ImpactRepository.NewImpactRepository(db)
	impactService := ImpactService.NewNewImpactService(impactRepo)
	impactController := ImpactController.NewImpactController(impactService, jwt)

	productRepo := ProductRepository.NewProductRepository(db)
	productService := ProductService.NewProductService(productRepo, impactRepo)
	productController := ProductController.NewProductController(productService, jwt)

	guestRepo := GuestRepository.NewGuestRepository(db)
	guestService := guestService.NewGuestService(guestRepo)
	guestController := GuestController.NewGuestController(guestService)

	cartRepo := CartRepository.NewCartRepository(db)
	cartService := CartService.NewCartService(cartRepo)
	cartController := CartController.NewCartController(cartService, jwt)

	transactionRepo := TransactionRepository.NewTransactionRepository(db)
	transactionService := TransactionService.NewTransactionService(transactionRepo, midtransService)
	transactionController := TransactionController.NewTransactionController(transactionService, jwt)

	webhookRepo := WebHookRepository.NewWebhookRepository(db)
	webhookService := WebhookService.NewWebhookService(webhookRepo)
	webhookController := WebhookController.NewWebhookRequest(webhookService)

	reviewRepo := ReviewRepository.NewReviewProductRepository(db)
	reviewService := ReviewService.NewReviewProductService(reviewRepo)
	reviewController := ReviewController.NewReviewProductController(reviewService, jwt)

	chatbotRepo := ChatbotRepository.NewChatbotRepository(db)
	chatbotService := ChatbotService.NewChatbotService(chatbotRepo, openAIservice)
	chatbotController := ChatbotController.NewChatbotController(chatbotService, jwt)

	routes.RouteUser(e, userController, *cfg)
	routes.RouteAdmin(e, adminController, *cfg)
	routes.RoutesProducts(e, productController, *cfg)
	routes.RouteImpacts(e, impactController, *cfg)
	routes.RouteStorage(e, storage, *cfg)
	routes.RouteGuest(e, guestController)
	routes.RouteCart(e, cartController, *cfg)
	routes.RouteTransaction(e, transactionController, *cfg)
	routes.PaymentNotification(e, webhookController)
	routes.RouteReviewProduct(e, reviewController, *cfg)
	routes.RouteChatbot(e, chatbotController, *cfg)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(cfg.APP_PORT))
}
