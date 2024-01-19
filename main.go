package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	user_http "github.com/haritsrizkall/jti-test/phone/delivery/http"
	user_mysql "github.com/haritsrizkall/jti-test/phone/repository/mysql"
	"github.com/haritsrizkall/jti-test/phone/usecase"
	"github.com/haritsrizkall/jti-test/pkg"
)

func main() {
	pkg.InitValidator()

	r := mux.NewRouter()
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	mySqlDB := pkg.MySQL{
		Host:     "localhost",
		Port:     3306,
		Database: "jti_test",
		Username: "root",
		Password: "florist123",
	}

	db, err := mySqlDB.Connect()
	if err != nil {
		log.Fatal(err)
	}

	phoneRepository := user_mysql.NewPhoneRepository(db)
	phoneUsecase := usecase.NewPhoneUsecase(phoneRepository)
	phoneHandler := user_http.NewPhoneHandler(phoneUsecase)

	r.HandleFunc("/api/phones/{id}", phoneHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/phones", phoneHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/phones", phoneHandler.Create).Methods("POST")
	r.HandleFunc("/api/phones/auto-generate", phoneHandler.AutoGenerate).Methods("POST")
	r.HandleFunc("/api/phones/{id}", phoneHandler.Update).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8082", r))
}
