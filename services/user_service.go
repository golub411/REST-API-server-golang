package services

import (
	"api-go/crudsql"
	"api-go/models"
	"api-go/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *crudsql.Database
}

func DatabaseInitU(db *crudsql.Database) *PostService {
	return &PostService{db: db}
}

func (p *PostService) CreateUsersTable() error {
	err := p.db.CreateTable("users", []string{"post_id INTEGER PRIMARY KEY AUTOINCREMENT", "name TEXT", "token TEXT"})
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

	return p.db.InsertValue("users", []string{"name", "token"}, []interface{}{user.Name, user.Token})
}

func (p *UserService) Login(username, password string) (string, error) {
	// Assuming SelectValueWhere returns a single map with "token" and "name" fields
	userMap, err := p.db.SelectValueWhere("users", []string{"name", "token"}, fmt.Sprintf("name = '%s'", username))
	if err != nil {
		return "", err
	}

	// Check if the result is empty
	if len(userMap) == 0 {
		return "", fmt.Errorf("user not found")
	}

	// Extract the token and name from the map
	user := userMap[0]
	token, ok := user["token"].(string)
	if !ok {
		return "", fmt.Errorf("token is not a string")
	}

	name, ok := user["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is not a string")
	}

	// Compare the password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(token), []byte(password))
	if err != nil {
		return "", err
	}

	// Generate a JWT token using the user's name
	jwtToken, err := utils.GenerateJWT(name)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (p *UserService) DeleteUsers(id int) error {
	return p.db.DeleteValue("users", map[string]interface{}{"id": id})
}
