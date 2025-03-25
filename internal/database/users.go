package database

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

/*
The UserModel struct contains a DB field, which is a pointer to a sql.DB instance.
The User struct includes four fields: Id, Email, Password, and Name.
json tags are used to define how the struct fields are converted to and from JSON, ensuring proper data serialization and deserialization.
The Password field is marked with a - in the json tag, instructing the JSON package to exclude it from JSON responses, making sure we don’t expose the password in the response.
*/
