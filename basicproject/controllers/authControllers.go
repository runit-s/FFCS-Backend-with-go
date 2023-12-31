package controllers

import (
	"basicproject/models"
	u "basicproject/utils"
	"encoding/json"
	"net/http"
)

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.RegistrationNo, account.Password)
	u.Respond(w, resp)
}

var DisconnectDB = func(w http.ResponseWriter, r *http.Request) {
	models.Disconnect()
	u.Respond(w, u.Message(true, "Disconnected"))
	return
}
