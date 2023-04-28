package main

import (
	"fmt"
	"golang-app/auth"
	"golang-app/campaign"
	"golang-app/handler"
	"golang-app/helper"
	"golang-app/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/golang-app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
  	}

	// campaignRepository := campaign.NewRepository(db)
	// campaigns, err := campaignRepository.FindAll()

	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(campaigns))
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// }

	// campaignRepository := campaign.NewRepository(db)
	// campaigns, err := campaignRepository.FindByUserID(1)

	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(campaigns))
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println("Jumlah Gambar",len(campaign.CampaignImages))
	// 		fmt.Println(campaign.CampaignImages[0].FileName)
	// 	}
	// }

	

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	// fmt.Println(authService.GenerateToken(1001))
	// userService.SaveAvatar(1, "runggu.png")
	
	campaigns, _ := campaignService.FindCampaigns(311)
	fmt.Println("campaigns")
	fmt.Println(len(campaigns))

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()

}
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//ambil data token 
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}


//ambil nilai header Authorization: bearer token
//ambil header Authorization, kita ambil nilai tokenya saja
//kita validasi token 
//kita ambil user_id
//ambil user dari db berdasarkan user_id lewat service
//kita set context isinya user

