package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ToDo
/*
* We are handling data storage through these data structures, since we are to be handling these in-memory.
* Using these allows us to fulfill the OpenAPI contracts that we have to implement.

* Worth noting that I have assumed no constraints on the fields used (name, description of To-Do items), meaning not much validation is
* done on the values of the fields themselves.
 */
type ToDo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ToDoList struct {
	TodoList []ToDo
}

var todoList = ToDoList{
	TodoList: make([]ToDo, 0),
}

func main() {
	//? Don't mind this. Some basic testing data to ensure things are being done right when the API is launched.
	todoList.TodoList = append(todoList.TodoList, ToDo{Title: "test1", Description: "test"})
	todoList.TodoList = append(todoList.TodoList, ToDo{Title: "test2", Description: "test"})
	todoList.TodoList = append(todoList.TodoList, ToDo{Title: "test3", Description: "test"})
	todoList.TodoList = append(todoList.TodoList, ToDo{Title: "test34", Description: "test"})
	fmt.Println(todoList.TodoList)
	http.HandleFunc("/", ToDoListHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
TODO IMPLEMENT A UNIT TEST FOR POTENTIAL TEST CASES? WOULD BE GOOD TO HAVE FOR THE CONTROLLER?
TODO IMPLEMENT A DTO AND WORK USING A REPOSITORY-SERVICE-CONTROLLER PATTERN? WOULD BE MUCH NICER TO WORK WITH WHEN IMPLEMENTING A RESTFUL API?
  THE USE OF A DTO HERE CAN BE ARGUED, BUT FOLLOWING THE REPOSITORY-SERVICE-CONTROLLER PATTERN WOULD BE MUCH NICER!
TODO IMPLEMENT TIGHTER INPUT VALIDATION ON REQUESTS TO ENSURE THAT THE FRONTEND CANNOT SEND ANYTHING MALICIOUS. FRONTEND REQUESTS CAN BE MANIPULATED!
TODO VALIDATION FRAMEWORKS ARE OFTEN AVAILABLE TO WORK WITH, AND WITH SOME RESEARCH, GOLANG DOES HAVE SOME WE COULD HAVE USED HERE!	THIS WOULD SORT OUT THE VALIDATION ISSUE
TODO IMPLEMENTING ROBUST AND MORE PERMANENT LOGGING IS ALSO SOMETHING THAT SHOULD BE CONSIDERED...

TODO SEE IF WE CAN SPLIT TODOLISTHANDLER UP?? CURRENTLY NO NEED BUT HANDLING A REQUEST SHOULD IDEALLY BE HANDLED BY A SINGULAR METHOD...
TODO THERE ARE ROUTING LIBRARIES WHICH HANDLE THIS SUPPOSEDLY, JUST A MATTER OF FINDING THE RIGHT ONE AND IMPLEMENTING! WOULD ALSO LET US SPLIT THINGS UP!
TODO SINCE WE HAVE NOT SPLIT THE HANDLING UP, WE HAVE NOT MET THE OPERATIONID OPENAPI SPEC, BUT THIS IS AN OPTIONAL ATTRIBUTE ANYWAYS, AND GENERALLY ONLY SERVES AS A METHOD REFERENCE WHEN ON THE BACKEND.
  IF THERE ARE NO METHODS TO REFER TO, THEN WE CAN ARGUE THAT WE DON'T NEED THE OPERATIONID?
*/

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(todoList.TodoList)
		//? Should not be really landing on this case for a GET request anyways. Still good to log the error.l
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
			log.Println(err)
			return
		}
	case "POST":
		var created ToDo
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&created)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		if created.Title == "" || created.Description == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Please include a name and a description!"})
			return
		}
		w.WriteHeader(http.StatusOK)
		todoList.TodoList = append(todoList.TodoList, created)
		json.NewEncoder(w).Encode(created)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request method!"})
	}
	// Your code here
}
