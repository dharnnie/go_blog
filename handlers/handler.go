package handlers

import(
	"fmt"
	"html/template"
	"github.com/gorilla/sessions"
	"net/http"
	"mycf/database"
	"log"
	"time"
)

type User struct{
	Nick string
	pass string
	Talk string
	Status string
	CTime string
	OtherComments map[string]string
}
type AuthErrs struct{
	LoginFail string
	SignUpFail string
}
type FullValues struct{// we just pass a single struct to the tpl.Execute()
	NickAndTime string
	Nick string
	Words string
	OtherComments map[string]string
	Status string
	CTime string
}
// recall: -  that, this struct is not neccesary with an {{ if pipeline}}
type InitComms struct{ // struct for user who is not logged in
	OtherComments map[string]string
	Nick string
}

var comms map[string]string

var CurrentUser User
var store  = sessions.NewCookieStore([]byte("secret"))

func Home(rw http.ResponseWriter, req *http.Request){
	if req.Method == "GET"{
		t, err := template.ParseFiles("templates/index.html")
		if err != nil{
			fmt.Println("Error occured while serving home")
		}
		t.Execute(rw, nil)
	}
}

func ServeResource(rw http.ResponseWriter, req *http.Request) {
	path := "templates/" + req.URL.Path
	http.ServeFile(rw, req, path)
}

func Comment(rw http.ResponseWriter, req *http.Request){
	if req.Method == "GET"{
		// check session
		t, err := template.ParseFiles("templates/comment.html")
		if err != nil{
			fmt.Println("Error Occured Parsing comment.html")
		}
		comms = loadComms()
		t.Execute(rw, FullValues{
			OtherComments: comms,
			Nick: CurrentUser.Nick,
			})
	}
}

func PostComment(rw http.ResponseWriter, req *http.Request){
	if req.Method == "POST"{
		req.ParseForm()
		t, err := template.ParseFiles("templates/comment.html")		
		if err != nil{
			fmt.Println("Error occured parsing temp comment + pipeline", err)
		}
		u_time := GetTime()
		// SAVE COMMENT TO DATABASE
		database.SaveComment(CurrentUser.Nick, req.FormValue("talk"), u_time)	
		comms = loadComms()
		nkAndTime := CurrentUser.Nick + " " + u_time
		t_err := t.Execute(rw, FullValues{
			NickAndTime: nkAndTime,
			Words: req.FormValue("talk"),
			CTime: u_time,
			OtherComments: comms,
			Status: "Logout",
			Nick: CurrentUser.Nick,
			})
		if t_err != nil{
			fmt.Println("Full values could not be parsed!!!", t_err)
		}
	}else{
		t,err := template.ParseFiles("templates/comment.html")
		if err != nil{
			log.Println("Could not parse comment from GET req:", err)
		}
		comms = loadComms()
		t.Execute(rw, FullValues{
			OtherComments: comms,
			Nick: CurrentUser.Nick,
			})
	}
}

func LoginPage(rw http.ResponseWriter, req *http.Request){
	if req.Method == "GET"{
		t, err := template.ParseFiles("templates/login.html")
		if err != nil{
			fmt.Println("Error occured Parsing login.html")
		}
		t.Execute(rw, nil)
	}
}

func ULPage(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST"{
		
		req.ParseForm()
		//tm := time.Now()
		oComms := database.GetComments()
		CurrentUser = User{req.FormValue("nick"), req.FormValue("password"), "", "","", oComms}
		verified := database.UAuth(CurrentUser.Nick, CurrentUser.pass)// Current User

		log.Println("The boolean returned from UAuth is :")
		log.Println(verified)
		log.Println("Values from struct User is: ")
		log.Println(CurrentUser.Nick, CurrentUser.pass)

		if verified{		
			fmt.Println("User verification completed succefully: Pass is ", verified)
			session, err := store.Get(req, "new_session")
			if err != nil{
				http.Error(rw, err.Error(), 500)
				return
			}
			session.Values["User"] = CurrentUser.Nick
			//session.Values["crack"] = "crack"
			session.Save(req, rw)
			tpl, err := template.ParseFiles("templates/comment.html")
			if err != nil{
				fmt.Println("Something happened while parsing files")
			}

			tpl.Execute(rw, User{
				Nick: CurrentUser.Nick,
				Status: "Logout",
				OtherComments: oComms,
			})// &CurrentUser... if multiple templates is returned to tpl
		}else{
			t,err := template.ParseFiles("templates/login.html")
			if err != nil{
				log.Println("Could not parse Login page")
			}
			t.Execute(rw, AuthErrs{
				LoginFail: "The details do not match!! Try again",
			})
		}
	}
}

func SignUp(rw http.ResponseWriter, req *http.Request){
	if req.Method == "GET"{
		t, err := template.ParseFiles("templates/signup.html")
		handleErr(err,"Error occured parsing signup" )
		t.Execute(rw, nil)
	}else{
		req.ParseForm()
		// add if IsSignedUp here.. To avoid user duplication
		if database.IsSignedUp(req.FormValue("nick")){
			t, err := template.ParseFiles("templates/signup.html")
			handleErr(err,"Could parse a redir to SignUp:" )
			t.Execute(rw, AuthErrs{
				SignUpFail: "The Nick is already taken!!!",
				})
		}else{
			created := database.CreateUser(req.FormValue("nick"), req.FormValue("password"))
			if created {
				log.Printf("User has %s been created with pass %s: ", req.FormValue("nick"), req.FormValue("password"))
				// set new session
				t, err := template.ParseFiles("templates/login.html")	
				handleErr(err,"Can't parse login.html @ SignUp:" )	
				t.Execute(rw, nil)
			}
		}
		
	}
}

// func redirToLogin(rw http.ResponseWriter, req *http.Request){

// // 	http.Redirect(rw, req, "http://localhost:8080/login", 303)

// }

// func serveResource() {
	
// }

func handleErr(e error, i string){
	if e != nil{
		log.Println(i, e)
	}
}

func GetTime()string {
	t := time.Now()
	timeOfPost := t.Format(time.RFC1123)
	return timeOfPost
}

func loadComms()map[string]string{
	x := database.GetComments()
	return x
}
