package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-rpg/2505-server/models"
)

// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags rewards
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/profile [get]
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"coins":         user.Coins,
			"reward_streak": user.RewardStreak,
			"last_reward":   user.LastReward,
		},
	})
}

// @Summary Get daily reward status
// @Description Check if user can claim daily reward
// @Tags rewards
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/daily-reward [get]
func GetDailyReward(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user can claim reward
	canClaim := false
	nextReward := time.Time{}

	if user.LastReward.IsZero() {
		canClaim = true
	} else {
		// Check if 24 hours have passed since last reward
		nextReward = user.LastReward.Add(24 * time.Hour)
		canClaim = time.Now().After(nextReward)
	}

	c.JSON(http.StatusOK, gin.H{
		"can_claim":   canClaim,
		"next_reward": nextReward,
		"streak":      user.RewardStreak,
	})
}

// @Summary Claim daily reward
// @Description Claim the daily reward and update streak
// @Tags rewards
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/daily-reward/claim [post]
func ClaimDailyReward(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user can claim reward
	if !user.LastReward.IsZero() {
		nextReward := user.LastReward.Add(24 * time.Hour)
		if time.Now().Before(nextReward) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Cannot claim reward yet",
				"next_reward": nextReward,
			})
			return
		}

		// Check if streak should be reset (more than 48 hours since last reward)
		if time.Now().After(user.LastReward.Add(48 * time.Hour)) {
			user.RewardStreak = 0
		}
	}

	// Calculate reward (base reward + streak bonus)
	baseReward := 100
	streakBonus := user.RewardStreak * 10
	totalReward := baseReward + streakBonus

	// Update user
	user.Coins += totalReward
	user.RewardStreak++
	user.LastReward = time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daily reward claimed successfully",
		"reward": gin.H{
			"base_reward":  baseReward,
			"streak_bonus": streakBonus,
			"total_reward": totalReward,
			"new_balance":  user.Coins,
			"new_streak":   user.RewardStreak,
			"next_reward":  user.LastReward.Add(24 * time.Hour),
		},
	})
}
