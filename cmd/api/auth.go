package main

import (
	"net/http"
	"time"

	"github.com/schlafer/EventApp/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {

	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

/*
The function begins by binding the incoming
JSON request body to a loginRequest struct.
This ensures that the data is properly formatted
and meets the required criteria, such as a valid email format
and a minimum password length of 8 characters.
It checks if the user exists in the database
by calling the GetByEmail method on the UserModel.
If the user is not found, a 404 Not Found response is returned.
The function uses the bcrypt library to compare the provided password
with the stored hashed password.
If the passwords do not match, a 401 Unauthorized response is returned.
Upon successful authentication, a JWT token is generated using the jwt library.
The token includes the user’s ID and an expiration time
(e.g., 72 hours from the time of issuance).
The generated token is returned to the client in a 200 OK response.
This token can then be used by the client to access protected routes.
*/
