package handler

import (
	"net/http"

	"github.com/KazakMaxim/auth-jwt-go/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTokens(c *gin.Context) {
	userGuid := c.Query("guid")

	tokens, tokensErr := h.services.GetTokensByGuid(userGuid)
	if tokensErr != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"Error": tokensErr.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"access_token":  tokens[1],
		"refresh_token": tokens[2],
	})
}

func (h *Handler) refreshTonkens(c *gin.Context) {
	var input models.Tokens

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"Error": err.Error(),
		})

		return
	}

	if input.Access == "" || input.Refresh == "" {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"Error": "Access token or refresh token are empty",
		})

		return
	}

	newTokens, userErr := h.services.NewTokens(input)
	if userErr != nil {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"Error": userErr.Error(),
		})

		return
	}

	c.JSON(http.StatusAccepted, map[string]interface{}{
		"access_token":  newTokens[1],
		"refresh_token": newTokens[2],
	})
}
