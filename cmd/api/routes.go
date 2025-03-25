package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)
		v1.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent)
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)
		v1.POST("/auth/login", app.login)
	}

	return g
}

/*
Routes is a function that initializes a new Gin server instance using gin.Default(),
which sets up some default middleware (like logging and recovery).
We define a route group /api/v1 to version our API.
Within this group, we map HTTP methods and paths to the corresponding handler functions for event operations.
This structure helps organize routes and makes it easier to manage API versions.
*/
