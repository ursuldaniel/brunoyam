package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

func (s *Server) handleListBooks(c *gin.Context) {
	id := c.MustGet("id").(int)

	books, err := s.store.ListBooks(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (s *Server) handleAddBook(c *gin.Context) {
	id := c.MustGet("id").(int)

	book := &models.Book{}
	if err := c.ShouldBindBodyWithJSON(book); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(book); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.AddBook(id, book); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Message: "book successfully added"})
}

func (s *Server) handleAddBooks(c *gin.Context) {
	id := c.MustGet("id").(int)

	books := []*models.Book{}
	if err := c.ShouldBindBodyWithJSON(&books); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	for _, book := range books {
		if err := s.validate.Struct(book); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
			return
		}
	}

	if err := s.store.AddBooks(id, books); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Message: "books successfully added"})
}
