package handler

import (
	"golang-app/helper"
	"golang-app/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context){
	//tangkap input dari user 
	//map input dari user ke struct RegisterUser
	//struct input diatas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed",http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity,response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

		if err != nil {
		response := helper.APIResponse("Register account failed",http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	formatter := user.FormatUser(newUser,"tokentoken")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	//user memasukkan input email dan password
	//input ditangkap handler
	//mapping dari input user ke input struct
	//input struct dipassing kedalam service 
	//dalam service mencari dengan bantuan repository user dengan email yang sudah dimasukkan tadi
	//mencocokkan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed",http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed",http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentoken")

	response := helper.APIResponse("Succesfully loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

