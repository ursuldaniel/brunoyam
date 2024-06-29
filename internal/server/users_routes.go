package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

func (s *Server) handleListUsers(c *gin.Context) {
	users, err := s.store.ListUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *Server) handleGetUser(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleCreateUser(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "user successfully created"})
}

func (s *Server) handleAddUsers(c *gin.Context) {
	users := []*models.User{}
	if err := c.ShouldBindBodyWithJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	for _, user := range users {
		if err := s.validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
			return
		}
	}

	if err := s.store.AddUsers(users); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "users successfully added"})
}

func (s *Server) handleUpdateUser(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	newUser := &models.User{}
	if err := c.ShouldBindBodyWithJSON(newUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(newUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.UpdateUser(id, newUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "user successfully updated"})
}

func (s *Server) handleDeleteUser(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "user successfully deleted"})
}
