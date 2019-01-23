package Mongodb

import (

	"encoding/json"
	"log"

	"github.com/globalsign/mgo" 
	"gopkg.in/mgo.v2/bson"
	"github.com/go-redis/redis"
)

const (
	url = "localhost"
)

type Person struct {
	FirstName string	`bson: "firstname: omitempty"`
	LastName string		`bson: "lastname: omitempty"`
	ID 		string 		`bson: "id"`
}

type Mongodb struct { 
	session *mgo.Session
	db 		*mgo.Database 
	client	*redis.Client
}

func (T *Mongodb)  Init()   {
	var err interface{}
	T.session, err = mgo.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}	
	T.session.SetMode(mgo.Monotonic, true)
	T.db = T.session.DB("People")
	if T.db == nil {
		log.Println("db People not found, exiting ")
		return
	}


	T.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (T *Mongodb)  Close()   {
	T.session.Close()
}

func (T *Mongodb) DatabaseName () string {
	return "Mongo DB"
}




func (T *Mongodb) ReadFile(parameter map[string] string) []byte {

	val, err := T.client.Get(parameter["id"]).Result()
	if err == redis.Nil {
		log.Println("Not found in Redis")
		c := T.db.C("info")
		var result Person
		count, err := c.Find(bson.M{"id": parameter["id"] }).Count()
	
		if err!=nil{
			log.Fatal(err)
		}
	
		if count==0{
			return []byte ("Person not found")
		}
	
		c.Find(bson.M{"id": parameter["id"] }).One(&result)
	
		content,_ := json.Marshal(result)

		err = T.client.Set(parameter["id"], content, 0).Err()
		if err != nil {
			log.Fatal(err)
		}		
	
	 
	} else if err != nil {
		log.Fatal(err)

	} 
	
	val, err = T.client.Get(parameter["id"]).Result()
	if err != nil {
		log.Fatal(err)
	}		
	
	return []byte(val)

} 

func (T *Mongodb) WriteFile(content map[string] interface{} ) {
 
	c := T.db.C("info")
	id,ok := content["id"].(string)

	if  ok{
		count, _ := c.Find(bson.M{"id": id}).Count()
		if count==0{

			contentToreturn,err := json.Marshal(content)
			if err != nil {
				log.Fatal(err)
			}
			err = T.client.Set(id, string (contentToreturn), 0).Err()		
			if err != nil {
				log.Fatal(err)
			}

	
			var tmp Person
			json.Unmarshal(contentToreturn,&tmp)
			err = c.Insert(tmp)
			if err != nil {
				log.Fatal(err)
			}
 
			T.client.Del("allPeople")
		}  
	}
	 
}

func (T *Mongodb) GetPeople() []byte{

	val, err := T.client.Get("allPeople").Result()
	if err == redis.Nil {
		log.Println("Not found in Redis")	
		c := T.db.C("info")

		count, _ := c.Find(nil).Count()
		if count==0{
			return []byte ("Person not found")
		}
	
		var results []map[string] interface{}
		_ = c.Find(nil).All(&results)
	
		content,_ := json.Marshal(results)

		err = T.client.Set("allPeople", string (content), 0).Err()		
		if err != nil {
			log.Fatal(err)
		}		

		return content
	

	} else {

		return [] byte(val)
	} 


} 

func (T *Mongodb) Update(id string, person map[string]string ) []byte {

	c := T.db.C("info")
 	var content = "Person not found"


	if id != person["id"]{
		return []byte ("url id and body id not same")
	}

	count,_ := c.Find(bson.M{"id": id}).Count()
 
	if count!=0{
		Bcontent,err := json.Marshal(person)
		if err != nil {
			log.Fatal(err)
		}
		err = T.client.Set(id, string (Bcontent), 0).Err()		
		if err != nil {
			log.Fatal(err)
		}

		var tmp Person
		json.Unmarshal(Bcontent,&tmp)
		err = c.Update(bson.M{"id": id} ,tmp   )
		if err != nil {
			log.Fatal(err)
		}

		T.client.Del("allPeople")		

		return Bcontent
	}

	return []byte (content)
} 