package main

import (
	"basicproject/app"
	"basicproject/controllers"
	"basicproject/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.Authenticate).Methods("POST") // Done

	router.HandleFunc("/admin/student", controllers.CreateStudent).Methods("POST") // Done
	router.HandleFunc("/admin/faculty", controllers.CreateFaculty).Methods("POST") // Done
	router.HandleFunc("/admin/course", controllers.CreateCourse).Methods("POST")   // Done
	router.HandleFunc("/admin/slot", controllers.CreateSlot).Methods("POST")       // Done

	router.HandleFunc("/faculty/{id}", controllers.GetFaculty).Methods("GET")  // Done
	router.HandleFunc("/course/{id}", controllers.GetCourse).Methods("GET")    // Done
	router.HandleFunc("/register", controllers.RegisterCourse).Methods("POST") // Done
	router.HandleFunc("/timetable", controllers.GetTimetable).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	// Creating Admmin Account
	account := models.Account{
		RegistrationNo: "Admin",
		Name:           "Admin",
		Admin:          true,
	}
	account.Create()
	fmt.Println(account)
	fmt.Println("------Admin Info--------")
	fmt.Println(fmt.Sprintf("RegistrationNo : %s", account.RegistrationNo))
	fmt.Println(fmt.Sprintf("Password : %s", account.Password))
	fmt.Println(fmt.Sprintf("Token : %s", account.Token))
	fmt.Println("------Admin Info--------")

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8080" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
