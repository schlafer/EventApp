package main

import (
	"net/http"

	"github.com/schlafer/EventApp/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

func (app *application) registerUser(c *gin.Context) {
	var register registerRequest
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}
	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

/*
The registerUserHandler function is responsible for
handling user registration requests. Here’s a breakdown of its functionality:

Data Binding and Validation: The function begins by binding the incoming
JSON request body to a registerRequest struct.
This ensures that the data is properly formatted and meets the required criteria,
such as a valid email format and minimum password length.

Password Hashing: To enhance security, the user’s password is hashed
using the bcrypt library. This step is crucial as it ensures that
the password is not stored in plain text in the database.

User Creation: A new User instance is created with the provided email,
hashed password, and name. This instance is then inserted into the database.

Response Handling: If the user is successfully created,
the function responds with a 201 Created status and the user data.
If any errors occur during the process, appropriate error messages are returned to the client.
*/
