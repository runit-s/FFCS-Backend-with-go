package models

import "fmt"

func GetSameDaySlots(day string, sameDaySlot *[]Slot) error {
	err := db.Table("slots").Where("day = ?", day).Find(&sameDaySlot).Error
	return err
}

func SaveSlot(slot *Slot) error {
	err := db.Create(slot)
	return err.Error
}

func IsFacultyNotAvailable(facId string) bool {
	var fac Faculty

	db.Table("faculties").Where("id = ?", facId).First(&fac)
	println(".....")
	fmt.Println(fac)
	fmt.Println("...")
	fmt.Println(fac.ID)
	fmt.Println(fac.Name)
	return fac.ID == ""
}

func GetFacultyFromDB(facId string) Faculty {
	var fac Faculty

	db.Table("faculties").Where("id = ?", facId).First(&fac)

	return fac
}

func GetSlotFromDB(slotID string) Slot {
	var slot Slot
	db.Table("slots").Where("id = ?", slotID).First(&slot)
	return slot
}
func GetCourseFromDB(courseId string) Course {
	var course Course
	db.Table("courses").Where("id = ?", courseId).First(&course)
	return course
}

func IsSlotNotAvailable(slotId string) bool {
	var slot Slot
	db.Table("slots").Where("id = ?", slotId).First(&slot)
	return slot.ID == ""
}

func IsCoursePresent(courseId string) bool {
	var course Course
	db.Table("courses").Where("id = ?", courseId).First(&course)
	return course.ID != ""
}

func GetAccountFromDB(id uint) Account {
	var acc Account
	db.Table("accounts").Where("id = ?", id).First(&acc)
	return acc
}
