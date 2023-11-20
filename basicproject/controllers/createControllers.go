package controllers

import (
	"basicproject/models"
	u "basicproject/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var CreateStudent = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
	}
	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var CreateFaculty = func(w http.ResponseWriter, r *http.Request) {
	faculty := &models.Faculty{}
	err := json.NewDecoder(r.Body).Decode(faculty)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	if !models.IsFacultyNotAvailable(faculty.ID) {
		u.Respond(w, u.Message(false, "Faculty already present"))
		return
	}

	resp := faculty.Create()
	u.Respond(w, resp)
}

var CreateCourse = func(w http.ResponseWriter, r *http.Request) {
	course := &models.Course{}

	err := json.NewDecoder(r.Body).Decode(course)

	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	if models.IsCoursePresent(course.ID) {
		u.Respond(w, u.Message(false, "Course already present"))
		return
	}

	if course.CourseType != "THEORY" && course.CourseType != "LAB" {
		u.Respond(w, u.Message(false, "Invalid course type"))
		return
	}

	if len(course.Faculty_ids) <= 0 {
		u.Respond(w, u.Message(false, "No Faculties"))
		return
	}

	if len(course.Slot_ids) <= 0 {
		u.Respond(w, u.Message(false, "No Slots"))
		return
	}
	var facId []string
	err = u.JSONToSlice(&course.Faculty_ids, &facId)

	
	if err != nil {
		u.Respond(w, u.Message(false, fmt.Sprintf("Internal Error : %s", err.Error())))
		return
	}
	
	var slotId []string
	err = u.JSONToSlice(&course.Slot_ids, &slotId)

	if err != nil {
		u.Respond(w, u.Message(false, fmt.Sprintf("Internal Error : %s", err.Error())))
		return
	}

	for _, v := range facId {
		if models.IsFacultyNotAvailable(v) {
			u.Respond(w, u.Message(false, fmt.Sprintf("Faculty %s not present", v)))
			return
		}
	}

	for _, v := range slotId {
		if models.IsSlotNotAvailable(v) {
			u.Respond(w, u.Message(false, fmt.Sprintf("Slot %s not present", v)))
			return
		}
	}

	resp := course.Create()
	u.Respond(w, resp)
}

var CreateSlot = func(w http.ResponseWriter, r *http.Request) {
	slot := &models.Slot{}

	err := json.NewDecoder(r.Body).Decode(slot)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	day := slot.Timing.Day
	if day != "mon" && day != "tue" && day != "wed" && day != "thu" && day != "fri" {
		u.Respond(w, u.Message(false, "Slot day not allowed"))
		return
	}

	// Validating time format
	t := u.ChangeTimeFormat(slot.Timing.Start)
	if t == "err" {
		u.Respond(w, u.Message(false, "Start time not valid"))
		return
	}

	t = u.ChangeTimeFormat(slot.Timing.End)
	if t == "err" {
		u.Respond(w, u.Message(false, "End time not valid"))
		return
	}

	// find slots on the same day
	var sameDaySlot []models.Slot
	err = models.GetSameDaySlots(slot.Timing.Day, &sameDaySlot)

	if err != nil {
		u.Respond(w, u.Message(false, "Error from GetSameDaySlot function"))
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	// check if new slot does not clash with another

	for _, new := range sameDaySlot {
		if IsOverlapping(new.Timing, slot.Timing) {
			u.Respond(w, u.Message(false, fmt.Sprintf("Slot clashed with another slot: %s", new.ID)))
			return
		}
	}

	err = models.SaveSlot(slot)

	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	resp := u.Message(true, "Course Created")
	resp["data"] = slot
	u.Respond(w, resp)
}

func IsOverlapping(A models.Timings, B models.Timings) bool {
	AStart, t := ConvertToTime(A.Start)
	BStart, t := ConvertToTime(B.Start)
	AEnd, t := ConvertToTime(A.End)
	BEnd, t := ConvertToTime(B.End)

	if t != nil {
		// We have to do nothing as we'll always be validating the time at the time of adding the slot into the DB
	}

	if AEnd.After(BStart) && BEnd.After(AStart) {
		return true
	}
	return false
}

func ConvertToTime(tt string) (time.Time, error) {
	t, err := time.Parse("03:04 PM", tt)
	if err != nil {
		t, err = time.Parse("3:04 PM", tt)
		if err != nil {
			return t, err
		}
	}

	return t, nil
}
