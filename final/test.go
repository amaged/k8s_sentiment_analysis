package main

import "github.com/cdipaolo/sentiment"
import "fmt"
import "database/sql"
import "strconv"

import _ "github.com/mattn/go-sqlite3"

func main() {
	model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
	}
	analysis := model.SentimentAnalysis("issues where id", sentiment.English) // 0
	fmt.Println(analysis.Words)
	for i := 0; i < len(analysis.Words); i++ {
		fmt.Println(analysis.Words[i].Score)
		/*message := GetTextFromID(i)

		  var num0 := 0
		  var num1 := 0
		  //count #0s and #1s

		  //update database
		  SetSentimentOfID(i, num0, num1)*/
	}
}

func GetTextFromID(id int) string {

	database, _ := sql.Open("sqlite3", "issues.db")

	rows, _ := database.Query("select from issues where id=" + strconv.Itoa(id) + " limit 1")
	var message string
	for rows.Next() {
		rows.Scan(&message)
		//fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	return message
}

func SetSentimentOfID(id int, num0 int, num1 int) {
	database, _ := sql.Open("sqlite3", "issues.db")
	statement, _ := database.Prepare("update issues set num0=" + strconv.Itoa(num0) + ", num1=" + strconv.Itoa(num1))
	statement.Exec()
}
