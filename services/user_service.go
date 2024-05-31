package services

import (
	"api-go/crudsql"
	"api-go/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *crudsql.Database
}

func DatabaseInitU(db *crudsql.Database) *UserService {
	return &UserService{db: db}
}

func (p *UserService) CreateUsersTable() error {
	err := p.db.CreateTable("users", []string{"id INTEGER PRIMARY KEY AUTOINCREMENT", "name TEXT", "password TEXT"})
	return err
}

func (p *UserService) Registration(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:  username,
		Token: string(hashedPassword),
	}

	return p.db.InsertValue("users", []string{"name", "password"}, []interface{}{user.Name, user.Token})
}

func (p *UserService) Login(username, password string) (bool, error) {
	// Assuming SelectValueWhere returns a single map with "token" and "name" fields
	userMap, err := p.db.SelectValueWhere("users", []string{"name", "password"}, fmt.Sprintf("name = '%s'", username))
	if err != nil {
		return false, err
	}

	// Check if the result is empty
	if len(userMap) == 0 {
		return false, fmt.Errorf("user not found")
	}

	// Extract the token and name from the map
	user := userMap[0]
	HeshPassword, ok := user["password"].(string)
	if !ok {
		return false, fmt.Errorf("token is not a string")
	}

	err = bcrypt.CompareHashAndPassword([]byte(HeshPassword), []byte(password))
	return err == nil, err

}

func (p *UserService) DeleteUser(id int) error {
	return p.db.DeleteValue("users", map[string]interface{}{"user_id": id})
}
