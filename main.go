package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	user_http "github.com/haritsrizkall/jti-test/phone/delivery/http"
	user_mysql "github.com/haritsrizkall/jti-test/phone/repository/mysql"
	"github.com/haritsrizkall/jti-test/phone/usecase"
	phone_ws "github.com/haritsrizkall/jti-test/phone/websocket"
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

	// serve views/index.html
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))

	log.Fatal(http.ListenAndServe(":8082", r))
}
