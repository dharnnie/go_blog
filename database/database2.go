package database


import(
	_"github.com/go-sql-driver/mysql"
 	"database/sql"
 	"log"
)
const(
	save_query = "INSERT INTO comments(nick, comment, time) VALUES(?,?,?)"
	retrieve_query = "SELECT nick, comment, time FROM comments"
)

func SaveComment(nk, words, time string){
	dbase, err := sql.Open("mysql", db_path)
	if err != nil{
		log.Println("Error occured @ SaveComment: ", err)
	}
	defer dbase.Close()

	
	prep, err := dbase.Prepare(save_query)
	if err != nil{
		log.Println("Error @ PrepareStatement: ", err)
	}

	res, err := prep.Exec(nk, words, time)
	if err != nil{
		log.Println("Error while getting results: ", err)
	}

	noOfComments, _:= res.LastInsertId()
	log.Println(noOfComments)
}

func GetComments()map[string]string{
	dbase, err := sql.Open("mysql", db_path)
	if err != nil{
		log.Println("Error @ getComments dbase: ", err)
	}
	defer dbase.Close()

	query, err := dbase.Query(retrieve_query)
	if err != nil{
		log.Println("Error @ retrieve_query: ", err)
	}

	var n string
	var c string
	var t string
	var nameAndTime string

	past_comments := make(map[string]string)

	for query.Next(){
		
		err := query.Scan(&n, &c, &t)
		if err != nil{
			log.Println("Error while scanning nick and Comm:", err)
		}
		nameAndTime = n + " at " + t
		past_comments[nameAndTime] = c // saves nick and correspondin comment in a map 
	}
	return past_comments
}