package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	return g
}

/*
Routes is a function that initializes a new Gin server instance using gin.Default(),
which sets up some default middleware (like logging and recovery).
Currently, it just returns the Gin instance.
We will add some routes to this instance later.
*/
