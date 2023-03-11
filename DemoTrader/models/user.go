package models

import (
	// "encoding/json"
	"fmt"
	"local/tools"
	"log"
	"time"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

const (
	userTableName = "user"
	userTokenTableName = "user_table"
	secretKey = "GolangInWonderland"
)

// table : user
type User struct {
	UserId          int    `db:"id"`
	Email           string `json:"email" db:"email"`
	TimestampJoined int    `db:"timestamp"`
	Password        string `db:"password"`
}


func (user *User) Save() *User {
	data := tools.ConvertStruct2Map(user)
	delete(data, "UserId")
	delete(data, "TimestampJoined")
	data["timestamp"] = int(time.Now().Unix())
	data["Password"] = tools.Encrypt(data["Password"].(string))

	user.TimestampJoined = data["timestamp"].(int)
	user.Password = data["Password"].(string)

	query := insertQuery("user", data)
	insertDB(query)
	return user
}

func (user *User) Read(id int, args ...[]string) []User {
	conditionData := make(map[string]interface{})

	if id != 0 {
		conditionData["user_id"] = id
	}

	// fill conditionData with dynamic arguments
	for i := 0; i<len(args); i++ {
		kv := args[i]
		conditionData[kv[0]] = kv[1]
	}

	query := selectQuery([]string{"*"}, userTableName, conditionData)
	rows := selectDB(query)
	fmt.Println(rows.Columns())

	var usr []User
	for rows.Next() {
		var tmpUser User
		rows.Scan(&tmpUser.UserId, &tmpUser.Email, &tmpUser.TimestampJoined, &tmpUser.Password)
		fmt.Println(tmpUser)
		usr = append(usr, tmpUser)
	}

	fmt.Println(usr)
	rows.Close()
	return usr
}

func (user *User) Validate(email string, password string) bool {
	var usr []User 
	encryptedPassword := tools.Encrypt(password)
	usr = user.Read(0, []string{"email", email}, []string{"password", encryptedPassword})
	switch len(usr) {
	case 0:
		return false
	case 1:
		return true 
	}
	return false
}

// ----------------------------------------------

// table : user_token
type UserToken struct {
	gorm.Model
	User 		User 		`json:"user" db:"user" gorm:"embedded"`
	ExpTime 	int		`json:"expire_time" db:"expire_time"`
	Token 	string 	`json:"token" db:"token"`
}

func (userToken *UserToken) Migrate() {
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&UserToken{})
}


func (userToken *UserToken) Save() *UserToken {
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.Create(&userToken)
	return userToken
}

func (userToken *UserToken) Read(userId int, args ...[]string) UserToken {
	db, err := gorm.Open(sqlite.Open("./db.sliqte3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var targetUserToken UserToken
	db.Find(&targetUserToken, userId)
	return targetUserToken
}

func (userToken UserToken) Validate(userId int, token string) bool {
	fetchedUserToken := userToken.Read(userId, []string{"token", token})
	if fetchedUserToken.ID != 0 {
		return true
	}
	return false
}