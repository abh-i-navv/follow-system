package handlers

import (
	"follow-system/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowHandler struct {
	Service *services.FollowService
}

func NewFollowHandler(s *services.FollowService) *FollowHandler {
	return &FollowHandler{Service: s}
}

type followRequest struct {
	FollowerID string `json:"follower_id"`
	TargetID   string `json:"target_id"`
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	var req followRequest

	//check if request body is valid
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	//conver id to uuid
	followerUUID, err := uuid.Parse(req.FollowerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid follower id",
		})
		return
	}

	targetUUID, err := uuid.Parse(req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid target id",
		})
		return
	}

	// parsing idempotency key from header
	idempotencyKey := c.GetHeader("Idempotency-Key")

	// sending follow request
	err = h.Service.FollowUser(c.Request.Context(), followerUUID, targetUUID, idempotencyKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//200 ok
	c.JSON(http.StatusCreated, gin.H{
		"status": "followed",
	})
}
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	var req followRequest

	//check if request body is valid
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	//conver id to uuid
	followerUUID, err := uuid.Parse(req.FollowerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid follower id",
		})
		return
	}

	targetUUID, err := uuid.Parse(req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid target id",
		})
		return
	}

	// sending follow request
	err = h.Service.UnfollowUser(c, followerUUID, targetUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//200 ok
	c.JSON(http.StatusCreated, gin.H{
		"status": "unfollowed",
	})
}

func (h *FollowHandler) GetFollower(c *gin.Context) {
	userIDStr := c.Param("id")

	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	followers, err := h.Service.GetFollower(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": followers,
	})
}
