package database

import "database/sql"

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Events:    EventModel{DB: db},
		Attendees: AttendeeModel{DB: db},
	}
}

/*
Here we are creating a Models struct with 3 fields: Users, Events, and Attendees.
We are also creating a NewModels function that takes a *sql.DB instance as an argument and passes it to the UserModel, EventModel, and AttendeeModel structs.
*/
