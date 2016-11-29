package database

import(
	"fmt"
	_"github.com/go-sql-driver/mysql"
 	"database/sql"
 	"log"
)
const(
	prepQuery = "INSERT INTO users(nick, password) VALUES(?,?)"
	db_path = "root:mysqlrootpassword@/cftrial"
	presenceQuery = "SELECT COUNT(*) FROM users WHERE nick = ?"
	retPasswordQry = "SELECT password FROM users WHERE nick = ?"
)

func UAuth(nick, pass string )(val bool){
	is := IsSignedUp(nick)
	uPass := getPass(nick)
	if is == true && uPass == pass{
		val = true
	}else{
		val = false
	}
	return val
}

func CreateUser(userNickname, userPassword string)bool{
	var done bool
  	dbase, err := sql.Open("mysql", db_path)
  	if err != nil{
    	fmt.Printf("Error occured with database", db_path)
  	}

  	defer dbase.Close()
  	prepStatement, err := dbase.Prepare(prepQuery)
  	
  	if err != nil{
    	fmt.Printf("An error occured while preparing statement")
  	}
    result, err := prepStatement.Exec(userNickname, userPassword)
    
    if err != nil {
      	fmt.Println("Error occured while getting result")
    	fmt.Printf("Error: ", err)
    }
    rowCount, err := result.LastInsertId()
    if err != nil{
      log.Println("Could not retrieve last row")
    }
    log.Println("Last user is on number: ", rowCount)
    	if done = IsSignedUp(userNickname); done{// lil issues
    		return done
    	}
    	return false
    	
}

func IsSignedUp(nk string)bool{
	//var nickVal string
	//var exists bool
	db, err := sql.Open("mysql", db_path)
	if err != nil{
		fmt.Println("Error occured while preparing abstrct DB:", err)
	}
	defer db.Close()

	stmt, err := db.Query(presenceQuery, nk)//here
	if err != nil{
		fmt.Println("Error at db.Prepare :IsSignedUp", err)
	}
   	noOfRows := checkCount(stmt)
   	if noOfRows == 0{
   		return false
   	}else{
   		return true
   	}
}

func checkCount(rows *sql.Rows)(count int){
	for rows.Next(){
		err := rows.Scan(&count)
		if err != nil{
			panic(err)
		}
	}
	return count
}

func getPass(nk string)string{
	// gets user password: User must have existed so i'm most likely not getting an error
	var thePass string
	conn, err := sql.Open("mysql", db_path)
	if err != nil{
		fmt.Println("Error occured while preparing abstrct DB:", err)
	}
	defer conn.Close()

	rows, err := conn.Query(retPasswordQry, nk)
	if err != nil{
		fmt.Println("Error occured while trying to get password:", err)
	}
	for rows.Next(){
		err := rows.Scan(&thePass)
		if err != nil{
			fmt.Println("Error occured at .Next:", err)
		}
	}
	defer rows.Close()
	return thePass
}