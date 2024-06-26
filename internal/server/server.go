package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt"
	"github.com/ursuldaniel/brunoyam/internal/domain/models"
)

type Storage interface {
	ListUsers() ([]*models.User, error)
	GetUser(id int) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(id int, newUser *models.User) error
	DeleteUser(id int) error
	Login(loginUser *models.User) (int, error)
	// IsTokenValid(token string) error
	// DisableToken(token string) error
}

type Server struct {
	addr     string
	store    Storage
	validate *validator.Validate
}

func NewServer(addr string, store Storage) *Server {
	return &Server{
		addr:     addr,
		store:    store,
		validate: validator.New(),
	}
}

func (s *Server) Run() error {
	app := gin.Default()

	usersRoutes := app.Group("/users")
	usersRoutes.GET("/", s.handleListUsers)
	usersRoutes.GET("/:id", s.handleGetUser)
	usersRoutes.POST("/", s.handleCreateUser)
	usersRoutes.PUT("/:id", s.handleUpdateUser)
	usersRoutes.DELETE("/:id", s.handleDeleteUser)

	app.POST("/login", s.handleLogin)
	app.GET("/profile", JWTAuth(s), s.handleProfile)

	return app.Run(s.addr)
}

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

func (s *Server) handleLogin(c *gin.Context) {
	loginUser := &models.User{}
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

func ParseId(idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func CreateToken(id int) (string, error) {
	claims := &jwt.MapClaims{
		"id":        id,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	// secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret_secret_secret"))
}

func JWTAuth(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header["Authorization"]
		if tokenString == nil {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Authorization token is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
			return []byte( /*os.Getenv("SECRET_KEY")*/ "secret_secret_secret"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid token claims"})
			c.Abort()
			return
		}

		id, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusForbidden, models.Response{Message: "Unauthorized access to the account"})
			c.Abort()
			return
		}

		c.Set("id", int(id))
		c.Next()
	}
}
