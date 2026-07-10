package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/notification/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/notification/internal/usecase"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Notify usecase.NotificationUseCase
}

func NewNotificationHandler(notify usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{Notify: notify}
}

// GetNotifications godoc
//
// @Summary List Notifications
// @Tags Notification
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {array} response.NotificationResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user ID not found in context",
		})
		return
	}

	userID, ok := userIDStr.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	notifications, err := h.Notify.GetByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
