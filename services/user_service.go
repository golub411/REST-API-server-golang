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
	err := p.db.CreateTable("users", []string{"id INTEGER PRIMARY KEY AUTOINCREMENT", "name TEXT", "password TEXT", "role TEXT"})
	return err
}

func (p *UserService) Registration(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     username,
		Role:     "user", // По умолчанию присваиваем роль "user"
		Password: string(hashedPassword),
	}

	return p.db.InsertValue("users", []string{"name", "password", "role"}, []interface{}{user.Name, user.Password, user.Role})
}

func (p *UserService) Login(username, password string) (bool, *models.User, error) {
	// Assuming SelectValueWhere returns a single map with "id", "name", "password", and "role" fields
	userMap, err := p.db.SelectValueWhere("users", []string{"id", "name", "password", "role"}, fmt.Sprintf("name = '%s'", username))
	if err != nil {
		return false, nil, err
	}

	// Check if the result is empty
	if len(userMap) == 0 {
		return false, nil, fmt.Errorf("user not found")
	}

	// Extract the user details from the map
	user := userMap[0]
	fmt.Println(user)
	id, ok := user["id"].(int64)
	if !ok {
		return false, nil, fmt.Errorf("id is not an integer")
	}
	name, ok := user["name"].(string)
	if !ok {
		return false, nil, fmt.Errorf("name is not a string")
	}
	hashedPassword, ok := user["password"].(string)
	if !ok {
		return false, nil, fmt.Errorf("password is not a string")
	}
	role, ok := user["role"].(string)
	if !ok {
		return false, nil, fmt.Errorf("role is not a string")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil, err
	}

	return true, &models.User{ID: int(id), Name: name, Role: role}, nil
}

func (p *UserService) DeleteUser(id int) error {
	return p.db.DeleteValue("users", map[string]interface{}{"id": id})
}
