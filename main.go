package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sbanka1/lenslocked/controllers"
	"github.com/sbanka1/lenslocked/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sbanka"
	password = "sbanka"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)
	defer services.Close()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	//Gallery routes
	r.Handle("/galleries/new", galleriesC.New).Methods("GET")
	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

	fmt.Println("Starting new server at :3000...")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
