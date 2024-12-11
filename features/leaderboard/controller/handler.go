package controller

import (
	"greenenvironment/constant"
	"greenenvironment/features/leaderboard"
	"greenenvironment/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LeaderboardController struct {
	service    leaderboard.LeaderboardServiceInterface
	jwt helper.JWTInterface
}

func NewLeaderboardController(service leaderboard.LeaderboardServiceInterface, jwt helper.JWTInterface) leaderboard.LeaderboardControllerInterface {
	return &LeaderboardController{
		service:    service,
		jwt: jwt,
	}
}

// Get Leaderboard
// @Summary      Retrieve leaderboard data
// @Description  Fetch the leaderboard data for users with the role "User"
// @Tags         Leaderboard
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer token"
// @Success      200            {object}  helper.Response{data=LeaderboardResponse}
// @Failure      401            {object}  helper.Response{data=string} "Unauthorized"
// @Failure      500            {object}  helper.Response{data=string} "Internal server error"
// @Router       /leaderboard [get]
func (lc *LeaderboardController) GetLeaderboard(c echo.Context) error {
	tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
	if tokenString == "" {
		return helper.UnauthorizedError(c)
	}

	data, err := lc.service.GetLeaderboard()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "failed to fetch leaderboard", nil))
	}

	response := LeaderboardResponse{}.FromEntity(data)

	return c.JSON(http.StatusOK, helper.FormatResponse(true, "success fetch leaderboard", response))
}
