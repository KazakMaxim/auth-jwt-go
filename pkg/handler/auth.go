package handler

import (
	"net/http"

	"github.com/KazakMaxim/auth-jwt-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": err.Error(),
		})

		return
	}

	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": "Username or password are empty",
		})

		return
	}

	userGuid := uuid.New().String()

	input.User_guid = userGuid

	createErr := h.services.CreateUser(input)
	if createErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": createErr.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"user_guid": userGuid,
	})
}

type signInput struct {
	Username string `json:"username" binding "required"`
	Password string `json:"password" binding "required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": err.Error(),
		})

		return
	}

	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": "Username or password are empty",
		})

		return
	}

	userGuid, userErr := h.services.AuthUser(input.Username, input.Password)
	if userErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": userErr.Error(),
		})

		return
	}

	tokens, tokensErr := h.services.Tokens.GenerateTokens(userGuid)
	if tokensErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error": tokensErr.Error(),
		})
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_guid":     tokens[0],
		"access_token":  tokens[1],
		"refresh_token": tokens[2],
	})
}
