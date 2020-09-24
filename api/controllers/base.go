package controllers

import (
	"fmt"
	"log"
	"net/http"
	"scotch-go-lang-rest-api/api/middlewares"
	"scotch-go-lang-rest-api/api/models"
	"scotch-go-lang-rest-api/api/responses"

	"gorm.io/driver/postgres"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	fmt.Println("DBURI", DBURI)

	a.DB, err = gorm.Open(postgres.Open(DBURI), &gorm.Config{})

	if err != nil {
		fmt.Printf("\n Cannot Connect to database %s", DbName)
		log.Fatal("This is the error", err)
	} else {
		fmt.Printf("We are connected to the database %s", DbName)
	}

	a.DB.Debug().AutoMigrate(&models.User{}) // database migration

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRouter()
}

func (a *App) initializeRouter() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
	a.Router.HandleFunc("/login", a.Login).Methods("POST")

}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to home")
}
