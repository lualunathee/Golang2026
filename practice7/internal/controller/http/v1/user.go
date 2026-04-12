package v1

import (
	"net/http"
	"practice7/internal/entity"
	"practice7/internal/usecase"
	"practice7/utils"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserInterface
}

func NewUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
    r := &userRoutes{t}
    

	handler.Use(utils.RateLimitMiddleware())

    
    h := handler.Group("/users")
    {
        h.POST("/", r.RegisterUser)
        h.POST("/login", r.LoginUser)
    }

    
    protected := handler.Group("/users")
    protected.Use(utils.JWTAuthMiddleware()) 
    {
        protected.GET("/me", r.GetMe) 
    }

	adminOnly := handler.Group("/users")
   
    adminOnly.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("admin")) 
    {
        
        adminOnly.PATCH("/promote/:id", r.PromoteUser) 
    }
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials or " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
  
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no user id in context"})
        return
    }

    
    user, err := r.t.GetUserByID(userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    
    c.JSON(http.StatusOK, gin.H{"email": user.Email})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
    id := c.Param("id") 
    
    err := r.t.UpdateUserRole(id, "admin")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin", "user_id": id})
}

