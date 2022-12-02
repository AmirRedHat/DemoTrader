package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// book_id, book_name, book_author, date
var INSERT_BOOK_QUERY string = "INSERT INTO books (book_id, book_name, book_author) VALUES ($1, $2, $3)";
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

		all_books, query_err := database.Query(READ_ALL_BOOK_QUERY);
		if query_err != nil{
			log.Fatal(query_err)
		}

		book_arr := make([][]string, 0);

		for all_books.Next(){
			book := book_struct{}
			book_scan_err := all_books.Scan(&book.book_id, &book.book_name, &book.book_author);
			if book_scan_err != nil{
				log.Fatal(book_scan_err);
			}

			book_arr = append(book_arr, []string{book.book_name, book.book_author});

		}
		
		get_json, get_json_err := json.Marshal(book_arr);
		if get_json_err != nil{
			log.Fatal(get_json_err)
		}
		
		response.Write(get_json);
		database.Close();


	}else if meth == "POST"{

		post_body := make(map[string]interface{})
		data, err := ioutil.ReadAll(request.Body)
		json.Unmarshal(data, &post_body)
		if err != nil{
			log.Fatal(err)
		}
		
		database.Exec(INSERT_BOOK_QUERY, post_body["book_id"], post_body["book_name"], post_body["book_author"])
		fmt.Fprint(response, `{"data": "book added", "message": "success"}`)

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


const BinanceURL = "https://api.binance.com/api/v3/klines?"
const KernelURL = "http:///data"

func data_collector(symbol string, time_frame string, limit string){
	var url string = KernelURL
	parameters := fmt.Sprintf(`{"symbol": "%s", "time_frame": "%s", "limit": %s, "broker": "binance"}`, symbol, time_frame, limit)
	fmt.Println(parameters)
	post_data := []byte(parameters);
	post := bytes.NewReader(post_data)

	fmt.Println(url);
	req, err := http.NewRequest(http.MethodPost, url, post)

	if err != nil{
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil{
		log.Fatal(err)
	}

	type ohcl_type map[string]interface{}
	res_arr := make([]ohcl_type, 0)
	response, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Fatal(err)
	}
	unjson_err := json.Unmarshal(response, &res_arr)
	if unjson_err != nil{
		log.Fatal(unjson_err)
	}

	for index, row_data := range res_arr{
		index ++
		fmt.Println(row_data["date"], row_data["close_price"])
	}
}


func main() {


	// data_collector("ETH-USDT", "4hour", "5")
	http.HandleFunc("/books", BookView);

	fmt.Print("server running\n");

	server_err := http.ListenAndServe("0.0.0.0:9000", nil);
	
	if errors.Is(server_err, http.ErrServerClosed){
		fmt.Print("Server closed\n");
	}else if server_err != nil{
		fmt.Printf("Server error: %s\n", server_err);
	}

}