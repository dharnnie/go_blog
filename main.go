package main 

import(
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"mycf/handlers"

)

func main() {
	serveWeb()
}

func serveWeb(){
	
	port := os.Getenv(){
		if port == ""{
			port = ":9000"
		}
	}
	
	http.HandleFunc("/styles/", handlers.ServeResource)
	myMux := mux.NewRouter()
	myMux.HandleFunc("/", handlers.Home)
	myMux.HandleFunc("/comment", handlers.Comment)
	myMux.HandleFunc("/login", handlers.LoginPage)
	myMux.HandleFunc("/signup", handlers.SignUp)
	myMux.HandleFunc("/upage", handlers.ULPage)
	myMux.HandleFunc("/say", handlers.PostComment)
	http.Handle("/", myMux)
	err := http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux))
	if err != nil{
		log.Fatal("ListenAndServe error: -> This error occured ", err)
	}
}



