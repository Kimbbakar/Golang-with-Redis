package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kimbbakar/Golang-with-Redis/TextFileRead"
	"github.com/kimbbakar/Golang-with-Redis/InMemoryfile"
	"github.com/kimbbakar/Golang-with-Redis/Mongodb"
	"flag"
	)

type Person struct {
	FirstName string	`json: "firstname: omitempty"`
	LastName string		`json: "lastname: omitempty"`
	ID 		string 		`json: "id"`
}

type IO interface {
	Init()     
	Close()
	DatabaseName()	string
	ReadFile(map[string] string)	[]byte
	WriteFile(map[string]interface{} )
	GetPeople()                      []byte
	Update(string,map[string] string) []byte		
 
}

var db IO

//var db IO  = &InMemoryfile.InMemoryfile{}
//var db IO = &TextFileRead.TextFileRead{}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	
	parameter := mux.Vars(r)

	var person Person
	content:= db.ReadFile(parameter)

	if string (content) == "Person not found"{
		json.NewEncoder(w).Encode("Person not found")	

		return
	}

	err := json.Unmarshal(content,&person)

	if err != nil {
		log.Fatal(err)
	}
 

	json.NewEncoder(w).Encode(person)	
 
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
 
    var person map[string]interface{}

	_ = json.NewDecoder(r.Body).Decode(&person)
  
	db.WriteFile(person)	
 
 
	json.NewEncoder(w).Encode(person) 

}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	log.Println("collecting list")

    var person map[string]interface{}

	_ = json.NewDecoder(r.Body).Decode(&person)
  
	list := db.GetPeople()	
	
 
	json.NewEncoder(w).Encode(string(list)) 

}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	
	parameter := mux.Vars(r)

    var person map[string]string 

	_ = json.NewDecoder(r.Body).Decode(&person)


	content := db.Update(parameter["id"],person )

	if (string(content)=="Person not found" || string(content)=="url id and body id not same"){
		log.Println(string (content))
		json.NewEncoder(w).Encode( string (content) )	
	} else{
		var tmp Person
		log.Println(string(content) )
		json.Unmarshal(content,&tmp)
	
		json.NewEncoder(w).Encode(tmp)	
	
	} 
}


func DatabaseProcess(){
	log.Println(db.DatabaseName() )


	router := mux.NewRouter() 
	router.HandleFunc("/person/{id}",GetPerson ).Methods("GET")
	router.HandleFunc("/person/{id}",UpdatePerson ).Methods("POST")
	router.HandleFunc("/person",CreatePerson ).Methods("POST")
	router.HandleFunc("/person",GetPeople ).Methods("GET")
	log.Println("Listening...")
	http.ListenAndServe(":8080", router)	
}


func main() {

	var database =flag.String("db", "mongo", "Default data bse is file")
	flag.Parse()
	log.Println(*database)
	
	if *database=="file"{
		db = &TextFileRead.TextFileRead{}
	}else if *database=="mongo"  {
		db = &Mongodb.Mongodb{} 
		 
	}else{
		db = &InMemoryfile.InMemoryfile{}		
	}
	
	db.Init()	 
	DatabaseProcess() 
	defer db.Close() 
}