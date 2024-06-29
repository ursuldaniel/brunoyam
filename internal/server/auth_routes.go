package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

func (s *Server) handleLogin(c *gin.Context) {
	loginUser := &models.LoginUser{}
	if err := c.ShouldBindBodyWithJSON(loginUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(loginUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	id, err := s.store.Login(loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err := CreateToken(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: token})
}

func (s *Server) handleProfile(c *gin.Context) {
	id := c.MustGet("id").(int)

	user, err := s.store.GetUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
