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

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

/*
c.Param("id") and c.Param("userId") are used to extract URL parameters
for the event and user IDs, respectively.
The function checks if the event and user exist in the database.
If not, it returns a 404 Not Found response.
It verifies if the attendee already exists for the event.
If so, it returns a 409 Conflict response.
The Insert method is called to add the attendee to the database if all checks pass.
If any operation fails, appropriate HTTP error responses are returned.
*/

func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

/*
The method extracts the event ID from the URL parameters using c.Param("id")
and converts it to an integer.
If the conversion fails, it returns a 400 Bad Request response indicating
an invalid event ID.
It calls the GetAttendeesByEvent method from the Attendees model
to fetch a list of users attending the specified event.
If an error occurs during data retrieval,
a 500 Internal Server Error response is returned with the error message.
If the data retrieval is successful,
a 200 OK response is returned along with the list of attendees.
*/

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

/*
Extracts id (attendee ID) and eventId from the URL parameters.
Validates the IDs and returns a 400 Bad Request if they are invalid.
Calls the Delete method on the AttendeeModel to remove the attendee.
Returns a 204 No Content status if the operation is successful,
indicating that the request was successful but there is no content to send back.
If any errors occur during the deletion process,
an appropriate error message is returned to the client.
*/

func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	events, err := app.models.Events.GetByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

/*
This function retrieves all events an attendee is attending.
It extracts id (attendee ID) from the URL parameters.
Validates the ID and returns a 400 Bad Request if it is invalid.
Calls the GetByAttendee method on the EventModel to fetch the events.
Returns a 200 OK status with the list of events if the operation is successful.
If any errors occur during the retrieval process,
an appropriate error message is returned to the client.
*/
