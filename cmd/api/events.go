package main

import (
	"net/http"
	"strconv"

	"github.com/schlafer/EventApp/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

/*
This handler manages the creation of a new event.
It binds the incoming JSON request body to an Event struct,
validates the data, and calls the Insert method on the EventModel
to add the event to the database.
If successful, it returns a 201 Created status with the created event data.
*/

func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

/*
This handler retrieves all events.
It calls the GetAll method on the EventModel to fetch all events from the database.
If successful, it returns a 200 OK status with the list of events.
*/

func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusOK, event)
}

/*
This handler retrieves a specific event by its ID.
It extracts the event ID from the URL parameters, validates it,
and calls the Get method on the EventModel to fetch the event from the database.
If the event is found, it returns a 200 OK status with the event data
else it returns a 404 Not Found status.
*/

func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	updateEvent := &database.Event{}

	if err := c.ShouldBindJSON(&updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateEvent.Id = id

	if err := app.models.Events.Update(updateEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updateEvent)
}

/*
This handler updates an existing event.
It extracts and validates the event ID from the URL parameters,
checks if the event exists,
binds the incoming JSON request body to an Event struct,
and calls the Update method on the EventModel to update the event in the database.
If successful, it returns a 200 OK status with the updated event data.
*/

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

/*
This handler deletes a specific event by its ID.
It extracts and validates the event ID from the URL parameters
and calls the Delete method on the EventModel
to remove the event from the database.
If successful, it returns a 204 No Content status.
*/
