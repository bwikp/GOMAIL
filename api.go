package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"

	"github.com/rs/cors"
)

type APIServer struct {
	addr string
}

type MailStruct struct {
	Titre   string `json:"titre"`
	Mail    string `json:"mail"`
	Message string `json:"message"`
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ca repond")
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("/gomail/send", func(w http.ResponseWriter, r *http.Request) {

		// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Max-Age", "15")

		reqBody, bad := io.ReadAll(r.Body)
		if bad != nil {
			log.Fatal(bad)
		}

		var contentMail MailStruct

		defer r.Body.Close()
		fmt.Println(string(reqBody))

		json.Unmarshal(reqBody, &contentMail)

		fmt.Println("mail recu :", contentMail)

		auth := smtp.PlainAuth(
			"",
			"exemple@gmail.com",
			"mdp",
			"smtp.gmail.com",
		)

		msg := "Subject:" + contentMail.Titre + "\n" + contentMail.Titre + "\n" + contentMail.Message

		err := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"exemple@gmail.com",
			[]string{"exemple@gmail.com"},
			[]byte(msg),
		)
		if err != nil {
			fmt.Println(err)
		}

	})

	router.HandleFunc("GET /test", getTest)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders: []string{"*"},
		// AllowCredentials: false,
	})

	handler := c.Handler(router)

	server := http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	log.Printf("Server has started %s", s.addr)

	return server.ListenAndServe()
}
