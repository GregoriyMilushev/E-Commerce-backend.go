package controllers

import (
	"log"
	"net/http"
	"pharmacy-backend/models"
	"pharmacy-backend/services"
	"pharmacy-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/oauth2/v2"
)

type AuthController struct {
    userService *services.UserService
}

func NewAuthController(db *gorm.DB) *AuthController {
    return &AuthController{
        userService: services.NewUserService(db),
    }
}

func (ac *AuthController) Me(c *gin.Context) {
	user, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized Error",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user,err := ac.userService.GetUserByEmail(input.Email)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    token, err := utils.GenerateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authoriztion", token, 3600 * 24 * 30, "", "", false, true) // secure false for localhost , TODO: for prod switch to true function

    // c.JSON(http.StatusOK, gin.H{"token": token})
}


func (ac *AuthController) Register(c *gin.Context) {
    
    var input struct {
        Name     string `json:"name" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Fatalln("1111")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Could not start server: %v", input)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatalln("2222")

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Name:        input.Name,
        Email:       input.Email,
        Password: 	 string(hashedPassword),
        Role: models.RoleClient,
    }

    if err := ac.userService.CreateUser(&user).Error; err != nil {
        log.Fatalln("3333")

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    log.Fatalln("4444")

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// var googleOAuthConfig = &oauth2.Config{
//     ClientID:     config.GoogleClientID,
//     ClientSecret: config.GoogleClientSecret,
//     RedirectURL:  config.GoogleRedirectURL,
//     Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
//     Endpoint:     google.Endpoint,
// }
// func (ac *AuthController) GoogleLogin(c *gin.Context) {
//     url := googleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
//     c.Redirect(http.StatusTemporaryRedirect, url)
// }

// func GoogleCallback(c *gin.Context) {
//     code := c.Query("code")

//     token, err := googleOAuthConfig.Exchange(context.Background(), code)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
//         return
//     }

//     oauth2Service, err := oauth2.New(config.HttpClient(context.Background(), token))
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth2 service"})
//         return
//     }

//     userInfo, err := oauth2Service.Userinfo.V2.Me.Get().Do()
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
//         return
//     }

//     var user models.User
//     if err := config.DB.Where("google_id = ?", userInfo.Id).First(&user).Error; err != nil {
//         if err == gorm.ErrRecordNotFound {
//             // Create a new user
//             user = models.User{
//                 Name:              userInfo.Name,
//                 Email:             userInfo.Email,
//                 GoogleID:          userInfo.Id,
//                 ProfilePictureURL: userInfo.Picture,
//             }
//             config.DB.Create(&user)
//         } else {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
//             return
//         }
//     }

//     // Generate JWT token
//     tokenString, err := GenerateJWT(user)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }

