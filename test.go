package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// book_id, book_name, book_author, date
var WRITE_BOOK_QUERY string = "INSERT INTO books VALUES (?, ?, ?, ?)";
var READ_ALL_BOOK_QUERY string = "SELECT * FROM books";
var READ_BOOK_QUERY string = "SELECT * FROM books WHERE book_id=?"

type book_struct struct {
	book_id int
	book_name string 
	book_author string
}

func BookView(response http.ResponseWriter, request *http.Request){
	meth := request.Method;
	database, db_err := sql.Open("sqlite3", "./test.db");
	if db_err != nil{
		log.Fatal(db_err);
	}
	
	dict := make(map[string]string);

	if meth == "GET"{
		get_dict := make(map[string][]book_struct)
		all_books, query_err := database.Query(READ_ALL_BOOK_QUERY);
		if query_err != nil{
			log.Fatal(query_err)
		}

		book_arr := make([]book_struct, 0);

		for all_books.Next(){
			book := book_struct{}
			book_scan_err := all_books.Scan(&book.book_id, &book.book_name, &book.book_author);
			if book_scan_err != nil{
				log.Fatal(book_scan_err);
			}
			book_arr = append(book_arr, book);
		}

		get_dict["data"] = book_arr

		fmt.Print(get_dict);
		get_json, get_json_err := json.Marshal(get_dict);
		if get_json_err != nil{
			log.Fatal(get_json_err)
		}
		
		response.Write(get_json);
		database.Close();


	}else if meth == "POST"{
		dict["message"] = "writing book in database";
		dict["status"] = "success";
		res, json_err := json.Marshal(dict);
		if json_err != nil{
			log.Fatal(json_err);
		}

		response.Write(res);
		database.Close();

	}else{
		dict["message"] = meth;
		dict["status"] = "success";
		res, json_err := json.Marshal(dict);
		if json_err != nil{
			log.Fatal(json_err);
		}

		response.Write(res);
		database.Close();
	}


}



func main() {

	http.HandleFunc("/books", BookView);

	fmt.Print("server running\n");

	server_err := http.ListenAndServe("0.0.0.0:9000", nil);
	
	if errors.Is(server_err, http.ErrServerClosed){
		fmt.Print("Server closed\n");
	}else if server_err != nil{
		fmt.Printf("Server error: %s\n", server_err);
	}

}