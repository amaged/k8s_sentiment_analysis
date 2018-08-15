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

	imin := 461
	imax := 813

	for i := imin; i < (imax + 1); i++ {
		message := GetTextFromID(i)
		analysis := model.SentimentAnalysis(message, sentiment.English) // 0

		//count #0s and #1s
		num0 := 0
		num1 := 0
		for i := 0; i < len(analysis.Words); i++ {
			if analysis.Words[i].Score == 1 {
				num1++
			} else {
				num0++
			}
		}

		//update database
		SetSentimentOfID(i, num0, num1)
	}

	database, _ := sql.Open("sqlite3", "issues.db")
	rows, _ := database.Query("select sum(num0),sum(num1) from issues where id>" + strconv.Itoa(imin-1) + " and id<" + strconv.Itoa(imax+1))
	var sumnum0 int
	var sumnum1 int
	for rows.Next() {
		rows.Scan(&sumnum0, &sumnum1)
		fmt.Println(strconv.Itoa(sumnum0) + " " + strconv.Itoa(sumnum1))
	}
}

func GetTextFromID(id int) string {

	database, _ := sql.Open("sqlite3", "issues.db")

	rows, _ := database.Query("select message from issues where id=" + strconv.Itoa(id) + " limit 1")
	var message string
	for rows.Next() {
		rows.Scan(&message)
		//fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	return message
}

func SetSentimentOfID(id int, num0 int, num1 int) {
	database, _ := sql.Open("sqlite3", "issues.db")
	statement, _ := database.Prepare("update issues set num0=" + strconv.Itoa(num0) + ", num1=" + strconv.Itoa(num1) + " where id=" + strconv.Itoa(id))
	statement.Exec()
}
