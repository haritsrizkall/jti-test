package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	auth_http "github.com/haritsrizkall/jti-test/auth/delivery/http"
	auth_usecase "github.com/haritsrizkall/jti-test/auth/usecase"
	"github.com/haritsrizkall/jti-test/middlewares"
	phone_http "github.com/haritsrizkall/jti-test/phone/delivery/http"
	phone_mysql "github.com/haritsrizkall/jti-test/phone/repository/mysql"
	phone_usecase "github.com/haritsrizkall/jti-test/phone/usecase"
	phone_ws "github.com/haritsrizkall/jti-test/phone/websocket"
	"github.com/haritsrizkall/jti-test/pkg"
	user_mysql "github.com/haritsrizkall/jti-test/user/repository/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pkg.InitValidator()

	r := mux.NewRouter()
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	mySqlDB := pkg.MySQL{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := mySqlDB.Connect()
	if err != nil {
		log.Fatal(err)
	}

	phone_hub := phone_ws.NewHub()
	go phone_hub.Run()

	googleOauth := pkg.NewGoogleOAuth(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_REDIRECT_URL"))

	phoneRepository := phone_mysql.NewPhoneRepository(db)
	userRepository := user_mysql.NewUserRepository(db)

	phoneUsecase := phone_usecase.NewPhoneUsecase(phoneRepository, phone_hub)
	authUsecase := auth_usecase.NewAuthUsecase(userRepository, *googleOauth)

	authHandler := auth_http.NewAuthHandler(authUsecase)
	phoneHandler := phone_http.NewPhoneHandler(phoneUsecase)

	r.Use(mux.CORSMethodMiddleware(r))

	// auth
	r.HandleFunc("/api/auth/login/google", authHandler.LoginWithGoogle).Methods("GET")
	r.HandleFunc("/api/auth/login/google/callback", authHandler.LoginWithGoogleCallback).Methods("GET")

	// sub router phones
	r.HandleFunc("/api/phones/ws", func(w http.ResponseWriter, r *http.Request) {
		phoneHandler.ServeWs(phone_hub, w, r)
	})

	phoneRoute := r.PathPrefix("/api/phones").Subrouter()
	phoneRoute.Use(middlewares.AuthMiddleware)
	phoneRoute.HandleFunc("/{id}", phoneHandler.Delete).Methods("DELETE")
	phoneRoute.HandleFunc("", phoneHandler.GetAll).Methods("GET")
	phoneRoute.HandleFunc("/{id}", phoneHandler.GetByID).Methods("GET")
	phoneRoute.HandleFunc("", phoneHandler.Create).Methods("POST")
	phoneRoute.HandleFunc("/auto-generate", phoneHandler.AutoGenerate).Methods("POST")
	phoneRoute.HandleFunc("/{id}", phoneHandler.Update).Methods("PUT")

	// serve views
	privateRoute := r.Methods(http.MethodGet).Subrouter()
	privateRoute.Use(middlewares.AuthMiddleware)
	privateRoute.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/input.html")
	}).Methods("GET")
	privateRoute.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/output.html")
	}).Methods("GET")
	privateRoute.HandleFunc("/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/edit.html")
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/login.html")
	}).Methods("GET")

	// serve wav files on /sounds/notif.wav
	r.PathPrefix("/sounds/").Handler(http.StripPrefix("/sounds/", http.FileServer(http.Dir("./sounds/"))))

	log.Fatal(http.ListenAndServe(":8082", r))
}
