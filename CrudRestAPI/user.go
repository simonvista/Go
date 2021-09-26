package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
const DSN="root:1234@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=true&loc=Local"

type User struct{
	gorm.Model
	FirstName 	string 	`json:"firstname"`
	LastName 	string 	`json:"lastname"`
	Email 		string 	`json:"email"`
}

func InitialMigration()  {
	DB,err=gorm.Open(mysql.Open(DSN),&gorm.Config{})
	if err!=nil {
		fmt.Println(err.Error())
		panic("Can't connect to DB")
	}
	DB.AutoMigrate(&User{})
}
// 1
func CreateUser(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-type","application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}
// 2
func GetUsers(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-type","application/json")
	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}
// 3
func GetUser(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-type","application/json")
	params :=mux.Vars(r)
	var user User
	DB.First(&user,params["id"])
	json.NewEncoder(w).Encode(user)
}
// 4
func UpdateUser(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-type","application/json")
	params :=mux.Vars(r)
	var user User
	DB.First(&user,params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}
// 5
func DeleteUser(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-type","application/json")
	params :=mux.Vars(r)
	var user User
	DB.Delete(&user,params["id"])
	json.NewEncoder(w).Encode("The user is deleted successfully")
}