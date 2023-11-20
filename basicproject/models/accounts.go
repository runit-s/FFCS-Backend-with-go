package models

import (
	u "basicproject/utils"
	"os"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

// a struct to rep user account
type Account struct {
	gorm.Model
	RegistrationNo     string         `json:"registrationNo"`
	Password           string         `json:"password"`
	Name               string         `json:"name"`
	Token              string         `json:"token";sql:"-"`
	Admin              bool           `json:"admin"`
	Registered_courses datatypes.JSON `json:"registered_courses"`
}

type Course struct {
	ID          string         `gorm:"primary key" json:"id"`
	Name        string         `json: "name"`
	Faculty_ids datatypes.JSON `json: "faculty_ids"`
	CourseType  string         `json: "CourseType"`
	Slot_ids    datatypes.JSON `json: "slot_ids"`
}

type Faculty struct {
	ID   string `gorm:"primary key" json:"id"`
	Name string `json: "name"`
}

type Slot struct {
	ID     string  `gorm:"primary key" json:"id"`
	Timing Timings `gorm:"embedded" json: "timings"`
}

type Timings struct {
	Day   string `gorm:"column:day" json: "day"`
	Start string `gorm:"column:start" json: "start"`
	End   string `gorm:"column:end" json: "end"`
}

func (account *Account) Create() map[string]interface{} {

	account.Password = u.RandStringRunes(5)
	tempPass := account.Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}
	if !account.Admin {
		account.RegistrationNo = "20BCE" + strconv.Itoa(int(account.ID))
		GetDB().Save(account)
	}
	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = tempPass
	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(registrationNo string, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("registration_no = ?", registrationNo).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Registration No not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

func (fac *Faculty) Create() map[string]interface{} {

	GetDB().Create(fac)

	if fac.ID == "" {
		return u.Message(false, "Failed to create Faculty")
	}

	resp := u.Message(true, "Faculty Created")
	resp["data"] = fac
	return resp
}

func (course *Course) Create() map[string]interface{} {

	GetDB().Create(course)

	if course.ID == "" {
		return u.Message(false, "Failed to create Course")
	}

	resp := u.BodyBuilder(true)
	resp["data"] = course.ParseJSON()
	return resp
}

func (course *Course) ParseJSON() map[string]interface{} {

	var facId []string
	u.JSONToSlice(&course.Faculty_ids, &facId)
	var faculties []Faculty
	for _, v := range facId {
		faculties = append(faculties, GetFacultyFromDB(v))
	}
	var slots []Slot = GetSlots(&course.Slot_ids)

	courseMap := map[string]interface{}{
		"id":            course.ID,
		"name":          course.Name,
		"faculties":     faculties,
		"course_type":   course.CourseType,
		"allowed_slots": slots,
	}

	return courseMap
}

func (account *Account) ParseCourse() map[string]interface{} {

	var registered_coursesID []string
	var registered_courses []Course
	u.JSONToSlice(&account.Registered_courses, &registered_coursesID)
	for _, v := range registered_coursesID {
		registered_courses = append(registered_courses, GetCourseFromDB(v))
	}

	CoursesMap := []map[string]interface{}{}
	for _, v := range registered_courses {
		courseElement := v.ParseJSON()
		currrentSlots := GetSlots(&v.Slot_ids)
		currentMap := map[string]interface{}{
			"course": courseElement,
			"slots":  currrentSlots,
		}
		CoursesMap = append(CoursesMap, currentMap)
	}

	AccountMap := map[string]interface{}{
		"id":                 account.RegistrationNo,
		"name":               account.Name,
		"registered_courses": CoursesMap,
	}

	return AccountMap
}

func GetSlots(jsonData *datatypes.JSON) []Slot {
	var slotID []string
	u.JSONToSlice(jsonData, &slotID)
	var slots []Slot
	for _, v := range slotID {
		slots = append(slots, GetSlotFromDB(v))
	}
	return slots
}
