package controllers

import (
	"basicproject/models"
	u "basicproject/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/datatypes"

	"github.com/gorilla/mux"
)

var GetFaculty = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	faculty := models.GetFacultyFromDB(id)

	if faculty.ID == "" {
		resp := u.BodyBuilder(false)
		resp["data"] = nil
		u.Respond(w, resp)
		return
	}

	resp := u.BodyBuilder(true)
	resp["data"] = faculty
	u.Respond(w, resp)
}

var GetCourse = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	course := models.GetCourseFromDB(id)

	if course.ID == "" {
		resp := u.BodyBuilder(false)
		resp["data"] = nil
		u.Respond(w, resp)
		return
	}

	resp := u.BodyBuilder(true)
	resp["data"] = course
	u.Respond(w, resp)
}

type RegisterPayLoad struct {
	Course_id  string         `json:"course_id"`
	Faculty_id string         `json:"faculty_id"`
	Slot_ids   datatypes.JSON `json:"slot_ids"`
}

var GetTimetable = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	var Account models.Account = models.GetAccountFromDB(userId)
	resp := u.BodyBuilder(true)
	resp["data"] = Account.ParseCourse()
	u.Respond(w, resp)
	return
}

var RegisterCourse = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	var Account models.Account = models.GetAccountFromDB(userId)
	var registerPayload RegisterPayLoad
	err := json.NewDecoder(r.Body).Decode(&registerPayload)

	if err != nil {
		resp := u.BodyBuilder(false)
		resp["Error"] = "Invalid request"
		u.Respond(w, resp)
		return
	}

	NewCourseId := registerPayload.Course_id
	NewCourse := models.GetCourseFromDB(NewCourseId)
	NewFaculty := registerPayload.Faculty_id

	// Check if that faculty teaches or not
	var FacultiesTeaching []string
	u.JSONToSlice(&NewCourse.Faculty_ids, &FacultiesTeaching)
	isTeaching := false

	for _, v := range FacultiesTeaching {
		if v == NewFaculty {
			isTeaching = true
		}
	}

	if !isTeaching {
		resp := u.BodyBuilder(false)
		resp["Error"] = "This faculty doesn't teach this course"
		u.Respond(w, resp)
		return
	}
	// End faculty Check

	// Start slot clash check
	// Finding New Slots
	var NewSlotsIDs []string
	u.JSONToSlice(&registerPayload.Slot_ids, &NewSlotsIDs)
	var NewSlots []models.Slot
	for _, v := range NewSlotsIDs {
		slot := models.GetSlotFromDB(v)
		if slot.ID == "" {
			resp := u.BodyBuilder(false)
			resp["Error"] = fmt.Sprintf("Slot %s is not available", v)
			u.Respond(w, resp)
			return
		}
		NewSlots = append(NewSlots, slot)
	}

	// Finding old courses to find old slots
	var OldCourseIds []string
	u.JSONToSlice(&Account.Registered_courses, &OldCourseIds)
	var OldCourses []models.Course
	for _, v := range OldCourseIds {
		course := models.GetCourseFromDB(v)

		if course.ID == "" {
			resp := u.BodyBuilder(false)
			resp["Error"] = "Internal Error"
			u.Respond(w, resp)
			return
		}

		OldCourses = append(OldCourses, course)
	}
	// fInding old slots from old courses
	var OldSlotIds []string
	for _, v := range OldCourses {
		var slotsId []string
		u.JSONToSlice(&v.Slot_ids, &slotsId)
		for _, s := range slotsId {
			OldSlotIds = append(OldSlotIds, s)
		}
	}

	var OldSlots []models.Slot
	for _, v := range OldSlotIds {
		OldSlots = append(OldSlots, models.GetSlotFromDB(v))
	}

	// Finding New Timings
	var NewTiming []models.Timings
	for _, v := range NewSlots {
		NewTiming = append(NewTiming, v.Timing)
	}
	// Finding Old Timings
	var OldTiming []models.Timings
	for _, v := range OldSlots {
		OldTiming = append(OldTiming, v.Timing)
	}

	//  Validating if new time doesn't clash with any old time
	for _, new := range NewTiming {
		for _, old := range OldTiming {
			if new.Day == old.Day && IsOverlapping(new, old) {
				resp := u.BodyBuilder(false)
				w.WriteHeader(http.StatusBadRequest)
				resp["Error"] = fmt.Sprintf("Slot %s clashes with %s", new, old)
				u.Respond(w, resp)
				return
			}
		}
	}
	// End sLot clash check

	// Now we just need to add this course to users account
	OldCourseIds = append(OldCourseIds, NewCourseId)
	registeredCourses := u.SliceToJson(&OldCourseIds)
	Account.Registered_courses = registeredCourses
	models.GetDB().Save(Account)
	resp := u.BodyBuilder(true)
	resp["data"] = Account.ParseCourse()
	u.Respond(w, resp)
	return
}
