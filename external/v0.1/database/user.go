package database

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type UserLogin struct {
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required" gorm:"unique"`
	Email    string `json:"email,omitempty" gorm:"unique"`
}

type Register struct {
	UserLogin
	UserInfo
}

type User struct {
	UserLogin
	Active         bool           `json:"active,omitempty"`
	Rol            string         `gorm:"default:user"`
	RecoverAccount RecoverAccount `json:"recovery,omitempty"`
	UserInfo       UserInfo       `json:"userinfo,omitempty"`
	gorm.Model
}

type UserInfo struct {
	FirstName string `json:"first_name,omitempty" binding:"required"`
	LastName  string `json:"last_name,omitempty" binding:"required"`
	Loged     *time.Time
	UserID    uint `gorm:"unique"`
	gorm.Model
}

type RecoverAccount struct {
	UserID       uint   `gorm:"unique"`
	Type         string `gorm:"default:'email'"`                                               // in the future could be whatsapp phone others
	RecoverCount uint   `gorm:"default:0"`                                                     // how many times the account was recovered
	Recover      int    `json:"recover,omitempty" binding:"required" gorm:"default:0"`         //is account waiting for recover
	RecoverCode  string `json:"recover_code,omitempty" binding:"required" gorm:"default:null"` // is a code issued for recover?
	gorm.Model
}

// getUser will have 2 ways to get record email and username
func getUser(user *User) (*User, error) {
	var newUser *User

	if user.UserLogin.Username != "" {
		tx := DB.Where("username = ?", user.UserLogin.Username).First(&newUser)
		return newUser, tx.Error
	} else if user.UserLogin.Email != "" {
		tx := DB.Where("email = ?", user.UserLogin.Username).First(&newUser)
		return newUser, tx.Error
	}
	return newUser, fmt.Errorf("identification needed to get record")
}

// getID
func GetID(user *User) (uint, error) {
	us, err := getUser(user)
	return us.ID, err
}

func getUsers() ([]User, error) {
	//TODO get a set of users for adm
	users := []User{}
	tx := DB.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func addUser(username, password, email, firstname, lastname string, active bool) error {
	//TODO inser an User in the table
	if username == "" || password == "" || email == "" || firstname == "" || lastname == "" {
		return fmt.Errorf("insuficient info to create user")
	}
	tx := DB.Create(&User{
		Active: active,
		UserLogin: UserLogin{
			Username: username,
			Password: password,
			Email:    email,
		},
		RecoverAccount: RecoverAccount{
			RecoverCount: 0,
			Type:         "email",
			Recover:      0,
		},
		UserInfo: UserInfo{
			FirstName: firstname,
			LastName:  lastname,
		},
	})
	return tx.Error
}

func editUser(user *User) (bool, error) {
	//TODO edit a user record from the table
	var tx *gorm.DB
	switch {
	case user.UserLogin.Email != "":
		tx = DB.Model(&user).Where("username", user.Username).Update("email", user.Email)
	}

	return tx.Error != nil, tx.Error
}

func getUserInfo(userID uint) (*UserInfo, error) {
	//TODO get the user info
	var userInfo *UserInfo = &UserInfo{}
	if userID == 0 {
		return nil, fmt.Errorf("need id to look for info")
	}
	tx := DB.Where("user_id = ?", userID).First(userInfo)
	return userInfo, tx.Error
}

func editUserInfo(user *UserInfo) (bool, error) {
	//TODO edit user info
	var tx *gorm.DB
	switch {
	case user.FirstName != "":
		tx = DB.Model(&user).Where("user_id", user.UserID).Update("first_name", user.FirstName)
	case user.LastName != "":
		tx = DB.Model(&user).Where("user_id", user.UserID).Update("last_name", user.LastName)
	case user.Loged != nil:
		tx = DB.Model(&user).Where("user_id", user.UserID).Update("loged", user.Loged)
	}

	return tx.Error != nil, tx.Error
}

func getRecoverAccount(userID uint) (*RecoverAccount, error) {
	//TODO recover account
	var userRecovery *RecoverAccount = &RecoverAccount{}
	if userID == 0 {
		return nil, fmt.Errorf("need id to look for info: %v", *userRecovery)
	}
	tx := DB.Where("user_id = ?", userID).First(userRecovery)
	return userRecovery, tx.Error
}

func editRecoverAccount(recoveryAcc *RecoverAccount) (*RecoverAccount, error) {
	//TODO edit counter of recover account and code and edited
	max := big.NewInt(999999)
	randInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Println("error generating random number:", err)
		return nil, err
	}
	//user id identifyes record
	recoveryAcc, err = getRecoverAccount(recoveryAcc.UserID)
	if err != nil {
		fmt.Println("error edit recovery", err)
		return nil, err
	}
	recoveryAcc.Recover = 1
	recoveryAcc.RecoverCode = fmt.Sprintf("%d", randInt)
	recoveryAcc.RecoverCount += 1
	DB.Save(recoveryAcc)
	fmt.Println("Random Number", randInt)
	return recoveryAcc, nil
}

// function will be called after some time the recovery is issued
func cancelRecoverAccount(recovery *RecoverAccount) (bool, error) {
	recovery, err := getRecoverAccount(recovery.UserID)
	if err != nil {
		fmt.Println("error edit recovery", err)
		return false, err
	}
	recovery.Recover = 0
	recovery.RecoverCode = ""
	tx := DB.Save(recovery)
	return tx.Error != nil, tx.Error
}

func switchAccount(user *User) (*User, error) {
	user, err := getUser(user)
	if err != nil {
		return nil, err
	}
	user.Active = !user.Active
	tx := DB.Save(user)
	return user, tx.Error
}
