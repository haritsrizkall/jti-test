package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	user_http "github.com/haritsrizkall/jti-test/phone/delivery/http"
	user_mysql "github.com/haritsrizkall/jti-test/phone/repository/mysql"
	"github.com/haritsrizkall/jti-test/phone/usecase"
	phone_ws "github.com/haritsrizkall/jti-test/phone/websocket"
	"github.com/haritsrizkall/jti-test/pkg"
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

	phoneRepository := user_mysql.NewPhoneRepository(db)
	phoneUsecase := usecase.NewPhoneUsecase(phoneRepository, phone_hub)
	phoneHandler := user_http.NewPhoneHandler(phoneUsecase)

	r.Use(mux.CORSMethodMiddleware(r))

	upgrader := websocket.Upgrader{}

	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Client Connected")

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			log.Println(string(message))

			err = conn.WriteMessage(messageType, []byte("Selamat Makan"))
			if err != nil {
				log.Println(err)
				break
			}
		}
	})

	r.HandleFunc("/api/phones/ws", func(w http.ResponseWriter, r *http.Request) {
		phoneHandler.ServeWs(phone_hub, w, r)
	})
	r.HandleFunc("/api/phones/{id}", phoneHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/phones", phoneHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/phones", phoneHandler.Create).Methods("POST")
	r.HandleFunc("/api/phones/auto-generate", phoneHandler.AutoGenerate).Methods("POST")
	r.HandleFunc("/api/phones/{id}", phoneHandler.Update).Methods("PUT")

	r.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/input.html")
	}).Methods("GET")
	r.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/output.html")
	}).Methods("GET")

	// serve wav files on /sounds/notif.wav
	r.PathPrefix("/sounds/").Handler(http.StripPrefix("/sounds/", http.FileServer(http.Dir("./sounds/"))))

	// serve views/index.html
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))

	log.Fatal(http.ListenAndServe(":8082", r))
}
