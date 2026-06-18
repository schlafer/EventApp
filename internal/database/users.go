package database

import (
	"context"
	"database/sql"
	"time"
)

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

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`
	err := m.DB.QueryRowContext(ctx, stmt, user.Email, user.Password, user.Name).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

//Here we insert the user into the database and return an error if there is one.

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	return m.getUser(query, email)
}

/*
Here we did some refactoring and created a new method called getUser,
notice the ...interface{} in the method signature.
This allows us to pass in multiple arguments to the method.
Then we have the Get and GetByEmail methods that we can use
to get a user by id or email.
This refactoring reduces code duplication and centralizes the logic
for querying and handling user data.
*/
