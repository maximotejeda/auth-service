package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// comparePasswd
func comparePasswd(password string, hash string) bool {
	//TODO compare 2 passwords
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err, password)
	}
	return err == nil
}

// LoginAccount
func LoginAccount(user *User) (*User, error) {
	//TODO login account
	newUser := User{}
	newUser.Username = user.Username
	newUser.Password = user.Password
	newUser.Email = user.Email

	now := time.Now()

	user, err := getUser(user)
	if err != nil {
		return nil, err
	}
	result := comparePasswd(newUser.Password, user.Password)
	if !result {
		return nil, fmt.Errorf("failed to authenticate user")
	}
	alterLogin := UserInfo{Loged: &now, UserID: user.ID}
	_, err = editUserInfo(&alterLogin)
	if err != nil {
		fmt.Println(err)
	}
	info, err := getUserInfo(user.ID)
	if err != nil {
		fmt.Println(err)
	}
	if info != nil {
		user.UserInfo = *info
	}
	user.Password = ""
	return user, nil

}

// CreateAccount
func CreateAccount(user *User) error {
	if user.Email == "" {
		return fmt.Errorf("user creation must have an email")
	}
	password, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	err = addUser(user.Username, password, user.Email, user.UserInfo.FirstName, user.UserInfo.LastName, user.Active)
	return err
}

// AccountExist
func AccountExist(data string) bool {
	user := User{}
	if strings.Contains(data, "@") {
		user.UserLogin.Email = data
	} else {
		user.UserLogin.Username = data
	}
	_, err := getUser(&user)

	return err == nil
}

// NewRecoverAccount
func NewRecoverAccount(data string) *User {
	user := &User{}
	recover := RecoverAccount{}
	if strings.Contains(data, "@") {
		user.UserLogin.Email = data
	} else {
		user.UserLogin.Username = data
	}
	user, err := getUser(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	recover.UserID = user.ID
	fmt.Println("Recovering ID: ", user.ID)
	// think about rigger the email from here or from the handler
	_, err = editRecoverAccount(&recover)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	recovery, _ := getRecoverAccount(user.ID)
	user.RecoverAccount = *recovery
	return user

}

// EndRecoverAccount
func EndRecoverAccount(id uint) {
	if id == 0 {
		fmt.Println(fmt.Errorf("id of userid on recovery cnat be: %d", id))
		return
	}
	recovery := RecoverAccount{UserID: id}
	cancelRecoverAccount(&recovery)
}

// GetRecoverInfo
func GetRecoverInfo(id uint) (*RecoverAccount, error) {
	recover, error := getRecoverAccount(id)
	return recover, error
}

// ActivateAccount
func ActivateAccount(user *User) (*User, error) {
	user, err := switchAccount(user)
	return user, err
}

// Admin Functions with special priviledges
func AdminGetUsers() ([]User, error) {
	users, err := getUsers()
	return users, err
}

// AdminBanUser
func AdminBanUser(userStr string) error {
	if userStr == "" {
		return errors.New("user cant be empty email or username needed")
	}
	user := &User{}
	if strings.Contains(userStr, "@") {
		user.Email = userStr
	} else {
		user.Username = userStr
	}
	user, err := switchAccount(user)
	if err != nil {
		fmt.Println("error adminBAnUser: ", err)
		return err
	}
	return err
}
