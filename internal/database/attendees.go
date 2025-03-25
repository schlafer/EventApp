package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

/*
The Attendee struct includes three fields: Id, UserId, and EventId.
An attendee is a user that has signed up for an event. An event can have many attendees and an attendee can attend many events.
*/
