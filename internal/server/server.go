package server

import (
	"net/http"
	"os"
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
	AddUsers(users []*models.User) error
	UpdateUser(id int, newUser *models.User) error
	DeleteUser(id int) error

	Login(loginUser *models.LoginUser) (int, error)

	ListBooks(id int) ([]*models.Book, error)
	AddBook(id int, book *models.Book) error
	AddBooks(id int, book []*models.Book) error
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
	usersRoutes.POST("/add", s.handleAddUsers)
	usersRoutes.PUT("/:id", s.handleUpdateUser)
	usersRoutes.DELETE("/:id", s.handleDeleteUser)

	authRoutes := app.Group("/auth")
	authRoutes.POST("/login", s.handleLogin)
	authRoutes.GET("/profile", JWTAuth(s), s.handleProfile)

	booksRoutes := app.Group("/books", JWTAuth(s))
	booksRoutes.GET("/", s.handleListBooks)
	booksRoutes.POST("/", s.handleAddBook)
	booksRoutes.POST("/add", s.handleAddBooks)

	return app.Run(s.addr)
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

	secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
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
			return []byte(os.Getenv("SECRET_KEY")), nil
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
