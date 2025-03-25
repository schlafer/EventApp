package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}
type Event struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"ownerId" binding:"required"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
}

/*
The Event struct includes five fields: Id, OwnerId, Name, Description, Date, and Location.
We set binding tags and some validation rules. These will used later when creating an event and binding the request body to the Event struct. This is done by the Gin framework.
For now we set a binding tag on the OwnerId field. Later we will remove it and instead use the current logged in user.
*/

func (m EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := m.DB.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
	if err != nil {
		return err
	}

	return nil
}

/*
This function inserts a new event into the events table.
It uses QueryRowContext, which executes the query with a context
that includes a 3-second timeout, ensuring the operation doesn’t hang indefinitely. If there is no error we add the id to the event and return nil.
*/

func (m EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

/*
We retrieve all records from the events table.
We then iterate over the result set and append each event to the events slice.
If the query fails, we return an error.
*/

func (m EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events WHERE id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	var event Event

	err := row.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

/*
This function retrieves a specific record from the events table
where the id matches the provided value.
It maps the result to the Event struct fields.
We check if the event is not found and return nil if it is not found.
*/

func (m EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5"

	_, err := m.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)
	if err != nil {
		return err
	}
	return nil
}

/*
This function updates an existing record in the events table.
It uses the SET clause to specify the columns to be updated and their new values.
It ensures only the record with the specified id is updated.
If the update fails, it returns an error.
*/

func (m EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM events WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

/*
Removes a record from the events table where the id matches the provided value.
Returns an error if the deletion fails.
*/
