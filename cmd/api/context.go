package main

import (
	"github.com/schlafer/EventApp/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *application) GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &database.User{}
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}

// Here we are getting the user from the context and returning it.
// If the user is not found we return an empty user.
