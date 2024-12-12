package routes

import (
	"greenenvironment/configs"
	"greenenvironment/constant/route"
	"greenenvironment/features/admin"
	"greenenvironment/features/cart"
	"greenenvironment/features/challenges"
	"greenenvironment/features/chatbot"
	"greenenvironment/features/dashboard"
	"greenenvironment/features/forum"
	"greenenvironment/features/guest"
	"greenenvironment/features/impacts"
	"greenenvironment/features/leaderboard"
	"greenenvironment/features/products"
	reviewproducts "greenenvironment/features/review_products"
	"greenenvironment/features/transactions"
	"greenenvironment/features/users"
	"greenenvironment/features/webhook"
	"greenenvironment/helper"
	"greenenvironment/utils/storages"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uh users.UserControllerInterface, cfg configs.GEConfig) {
	e.POST(route.UserRegisterOTP, uh.RequestRegisterOTP)
	e.POST(route.UserVerifyRegisterOTP, uh.VerifyRegisterOTP)
	e.POST(route.UserLogin, uh.Login)

	e.POST(route.UserForgotPassword, uh.ForgotPasswordRequest)
	e.POST(route.UserVerifyForgotOTP, uh.VerifyForgotPasswordOTP)
	e.PUT(route.UserResetPassword, uh.ResetPassword)

	e.GET(route.UserLoginGoogle, uh.GoogleLogin)
	e.GET(route.UserGoogleCallback, uh.GoogleCallback)

	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET(route.UserProfile, uh.GetUserData, echojwt.WithConfig(jwtConfig))
	e.PUT(route.UserUpdate, uh.UpdateUserInfo, echojwt.WithConfig(jwtConfig))
	e.PUT(route.UserUpdateAvatar, uh.UpdateAvatar, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.UserPath, uh.Delete, echojwt.WithConfig(jwtConfig))

	e.POST(route.UserRequestUpdateOTP, uh.RequestPasswordUpdateOTP, echojwt.WithConfig(jwtConfig))
	e.PUT(route.UserUpdatePassword, uh.UpdateUserPassword, echojwt.WithConfig(jwtConfig))

	// Admin
	e.GET(route.AdminManageUserPath, uh.GetAllUsersForAdmin, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminManageUserByID, uh.GetUserByIDForAdmin, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminManageUserByID, uh.UpdateUserForAdmin, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminManageUserByID, uh.DeleteUserForAdmin, echojwt.WithConfig(jwtConfig))
}

func RouteAdmin(e *echo.Echo, ah admin.AdminControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.AdminLogin, ah.Login)

	e.GET(route.AdminPath, ah.GetAdminData, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminPath, ah.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminPath, ah.Delete, echojwt.WithConfig(jwtConfig))
}

func RoutesProducts(e *echo.Echo, ph products.ProductControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.POST(route.ProductPath, ph.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ProductPath, ph.GetAll)
	e.GET(route.ProductByID, ph.GetById)
	e.GET(route.CategoryProduct, ph.GetByCategory)
	e.PUT(route.ProductByID, ph.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ProductByID, ph.Delete, echojwt.WithConfig(jwtConfig))
}

func RouteImpacts(e *echo.Echo, ic impacts.ImpactControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.ImpactCategoryPath, ic.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ImpactCategoryPath, ic.GetAll, echojwt.WithConfig(jwtConfig))
	e.GET(route.ImpactCategoryByID, ic.GetByID, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ImpactCategoryByID, ic.Delete, echojwt.WithConfig(jwtConfig))
}

func RouteStorage(e *echo.Echo, sc storages.StorageInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/media/upload", sc.UploadFileHandler, echojwt.WithConfig(jwtConfig))
}

func RouteGuest(e *echo.Echo, gc guest.GuestControllerInterface) {
	e.GET("/api/v1/guest/products", gc.GetGuestProduct)
}

func RouteCart(e *echo.Echo, cc cart.CartControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.GET(route.CartPath, cc.Get, echojwt.WithConfig(jwtConfig))
	e.POST(route.CartPath, cc.Create, echojwt.WithConfig(jwtConfig))
	e.PUT(route.CartPath, cc.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.CartByID, cc.Delete, echojwt.WithConfig(jwtConfig))
}

func RouteTransaction(e *echo.Echo, tc transactions.TransactionControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.TransactionPath, tc.CreateTransaction, echojwt.WithConfig(jwtConfig))
	e.GET(route.TransactionPath, tc.GetUserTransaction, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.TransactionByID, tc.DeleteTransaction, echojwt.WithConfig(jwtConfig))

	e.GET("/api/v1/admin/transactions", tc.GetAllTransaction, echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/transactions/:id", tc.GetTransactionByID, echojwt.WithConfig(jwtConfig))
}

func PaymentNotification(e *echo.Echo, wh webhook.MidtransNotificationController) {
	e.POST("/midtrans-notification", wh.HandleNotification)
}

func RouteReviewProduct(e *echo.Echo, rpc reviewproducts.ReviewProductControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.POST(route.ReviewProduct, rpc.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ReviewProductByID, rpc.GetProductReview)
}

func RouteChatbot(e *echo.Echo, ch chatbot.ChatbotControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST(route.ChatbotPath, ch.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.ChatbotPathByID, ch.GetByID, echojwt.WithConfig(jwtConfig))
}

func RouteForum(e *echo.Echo, fh forum.ForumControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET(route.ForumPath, fh.GetAllForum, echojwt.WithConfig(jwtConfig))
	e.GET(route.ForumByID, fh.GetForumByID, echojwt.WithConfig(jwtConfig))
	e.GET(route.GetForumByUserID, fh.GetForumByUserID, echojwt.WithConfig(jwtConfig))
	e.POST(route.ForumPath, fh.PostForum, echojwt.WithConfig(jwtConfig))
	e.PUT(route.ForumByID, fh.UpdateForum, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ForumByID, fh.DeleteForum, echojwt.WithConfig(jwtConfig))

	e.GET(route.ForumMessageByID, fh.GetMessageForumByID, echojwt.WithConfig(jwtConfig))
	e.POST(route.ForumMessage, fh.PostMessageForum, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.ForumMessageByID, fh.DeleteMessageForum, echojwt.WithConfig(jwtConfig))
	e.PUT(route.ForumMessageByID, fh.UpdateMessageForum, echojwt.WithConfig(jwtConfig))
}

func RouteChallenge(e *echo.Echo, cc challenges.ChallengeControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.POST(route.AdminChallengePath, cc.Create, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminChallengePath, cc.GetAll, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminChallengeByID, cc.GetByID, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminChallengeByID, cc.Update, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminChallengeByID, cc.Delete, echojwt.WithConfig(jwtConfig))

	// Challenge Task
	e.POST(route.AdminChallengeTask, cc.CreateTask, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminChallengeTaskbyChallengeID, cc.GetAllTasksByChallengeID, echojwt.WithConfig(jwtConfig))
	e.GET(route.AdminChallengeTaskByID, cc.GetTaskByID, echojwt.WithConfig(jwtConfig))
	e.PUT(route.AdminChallengeTaskByID, cc.UpdateTask, echojwt.WithConfig(jwtConfig))
	e.DELETE(route.AdminChallengeTaskByID, cc.DeleteTask, echojwt.WithConfig(jwtConfig))

	// User
	e.POST(route.TakeChallenge, cc.CreateChallengeLog, echojwt.WithConfig(jwtConfig))
	e.PUT(route.TaskConfirmationProgress, cc.UpdateChallengeConfirmationProgress, echojwt.WithConfig(jwtConfig))
	e.POST(route.ClaimRewards, cc.ClaimRewards, echojwt.WithConfig(jwtConfig))
	e.GET(route.ActiveChallenge, cc.GetActiveChallenges, echojwt.WithConfig(jwtConfig))
	e.GET(route.UnclaimedChallenge, cc.GetUnclaimedChallenges, echojwt.WithConfig(jwtConfig))
	e.GET(route.UserChallengeDetails, cc.GetChallengeDetailsWithConfirmations, echojwt.WithConfig(jwtConfig))
	e.GET(route.UserUnclaimedChallengeDetails, cc.GetChallengeDetails, echojwt.WithConfig(jwtConfig))
}

func RouteDashboard(e *echo.Echo, dc dashboard.DashboardControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.GET(route.AdminDashboard, dc.GetDashboard, echojwt.WithConfig(jwtConfig))
}

func RouteLeaderboard(e *echo.Echo, lc leaderboard.LeaderboardControllerInterface, cfg configs.GEConfig) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.GET(route.LeaderboardPath, lc.GetLeaderboard, echojwt.WithConfig(jwtConfig))
}
