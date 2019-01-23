package TextFileRead

import (
	"encoding/json"
	"io/ioutil" 
	"log"
	"os" 
	)


type Person struct {
	FirstName string	`json: "firstname: omitempty"`
	LastName string		`json: "lastname: omitempty"`
	ID 		string 		`json: "id"`
}	

type TextFileRead struct{

	
}

func (T *TextFileRead)  Init()   {
 
}

func (T *TextFileRead)  Close()   {
 
}

func (T *TextFileRead) DatabaseName () string {
	return "Text File Database"
}

 

func (T *TextFileRead) ReadFile(parameter map[string] string) []byte {

	var _, err = os.Stat("person.txt")

	var content = "Person not found"
	
	

    // create file if not exists
    if err!=nil { 
		log.Println(err)
		return []byte(content)
		// _, err = os.Create("person.txt")
		// if err != nil {
		// 	log.Fatal(err.Error())
		// } 
    }
  
 


	data, err := ioutil.ReadFile("person.txt")	

	log.Println(string(data))
 
	for  i := 0;i<len(data );{
		ok := false
		for j :=i + 1 ;j < len(data);j++ {
			if data[j]==byte('}'){
				
				var tmp Person
				log.Println(string (data[i:j+1]))

				json.Unmarshal([]byte(string (data[i:j+1]) )  ,&tmp )
				log.Println(tmp)

				if tmp.ID == parameter["id"]{
					content = string (data[i:j+1] )
					ok = true
				}
				i = j + 1
			}			
		}

		if ok {
			break
		}

		log.Println(i)
 
	}

	return [] byte (content)
} 


func (T *TextFileRead) WriteFile(content map[string] interface{} ) {

	b,_:=json.Marshal(content)

//	a,ok:=content["id"].(string)

	// if ok==false{
	// 	log.Fatal("ID not found")
	// }

    var _, err = os.Stat("person.txt")

    // create file if not exists
    if os.IsNotExist(err) {
        _, err = os.Create("person.txt")
		if err != nil {
			log.Fatal(err.Error())
		} 
    }
 

	data, err := ioutil.ReadFile("person.txt")	


	var allPerson = string(data)


	allPerson += string(b)



	_ = ioutil.WriteFile( "person.txt" ,[]byte(allPerson),0)
 
	// file, err = os.Open("person.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// } 


	// writer := bufio.NewWriter(file)

	// for _,val:= range allPerson{
	// 	log.Println(val)
	// 	x,y := writer.Write([]byte("hii") )
	// 	log.Println(x)
	// 	log.Println(y)
	// }

	// defer file.Close()

}



func (T *TextFileRead) GetPeople() []byte{
 

    var _, err = os.Stat("person.txt")


	// create file if not exists
	var content = "Person not found"	
    if err!=nil { 
		log.Println(err)
		return []byte(content)
		// _, err = os.Create("person.txt")
		// if err != nil {
		// 	log.Fatal(err.Error())
		// } 
    }	
 

	data, err := ioutil.ReadFile("person.txt")	
		
	log.Println(string(data) )
	return []byte(string(data) )
 

} 


func (T *TextFileRead) Update(id string, person map[string]string ) []byte {

	var _, err = os.Stat("person.txt")

	var content = "Person not found"
	
	

    // create file if not exists
    if err!=nil { 
		log.Println(err)
		return []byte(content)
		// _, err = os.Create("person.txt")
		// if err != nil {
		// 	log.Fatal(err.Error())
		// } 
    }
  
	if id != person["id"]{
		return []byte ("url id and body id not same")
	}

	b,_:=json.Marshal(person) 
	



	data, err := ioutil.ReadFile("person.txt")	
	
	var allPerson string 
	for  i := 0;i<len(data );{

		for j :=i + 1 ;j < len(data);j++ {
			if data[j]==byte('}'){
				
				var tmp Person 

				json.Unmarshal([]byte(string (data[i:j+1]) )  ,&tmp ) 

				if tmp.ID == id{ 
					allPerson+=(string (b) )	
					content = string (b) 
				} else {
					allPerson+=(string (data[i:j+1]) )					
				} 
				i = j + 1
			}			
		}
 
 
 
	}

	


	_ = ioutil.WriteFile( "person.txt" ,[]byte(allPerson),0)	
	return [] byte (content)
} 