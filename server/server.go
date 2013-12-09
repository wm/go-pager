package main

import (
	"log"
	"encoding/json"
	"os"
	"net/http"
	"strconv"
	"github.com/codegangsta/martini"
)

type Contact struct {
	Name string
	Number string
}


type Config struct {
    Contacts []Contact
}

func main() {
	reader, _ := os.Open("server/config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	m := martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	m.Map(config)

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/`, welcome)
	r.Post(`/attempt_call`, attemptCall)
	r.Post(`/screen_for_machine`, screenForMachine)
	r.Post(`/complete_call`, completeCall)

	// Add the router action
	m.Action(r.Handle)
	m.Run()

}

func welcome(config *Config) string {
    return "Hello world!"
}

func attemptCall(req *http.Request, w http.ResponseWriter, config *Config) string {
	var response string

	dialCallStatus := req.PostFormValue("DialCallStatus")
	numberIndexStr := req.PostFormValue("number_index")

	if numberIndexStr == "" {
		numberIndexStr = "0"
	}

	numberIndex, _ := strconv.Atoi(numberIndexStr)

	if dialCallStatus != "completed" && (numberIndex < len(config.Contacts)) {
		response = dialNumber(numberIndex + 1, (config.Contacts)[numberIndex].Number)
	} else {
		response = hangUp()
	}

	log.Printf("Contact: %v", config.Contacts[numberIndex].Name)
	log.Printf("response: %v", response)
	return response
}

func screenForMachine() string {
	return "<?xml version='1.0' encoding='UTF-8'?>" + 
	"<Response>" +
	"<Gather action='complete_call'>" +
	"<Say>Press any key to accept this call</Say>" +
	"</Gather>" +
	"<Hangup/>" +
	"</Response>"
}

func completeCall() string {
	return "<?xml version='1.0' encoding='UTF-8'?>" +
	"<Response>" +
	"<Say>Connecting</Say>" +
	"</Response>"
}

func hangUp() string {
	return "<Response><Hangup/></Response>"
}

func dialNumber(nextNumberIndex int, number string) string {
	return "<?xml version='1.0' encoding='UTF-8'?>" + 
    "<Response>" +
		"<Dial action='attempt_call?number_index=" + strconv.Itoa(nextNumberIndex) + "'>" +
			"<Number url='screen_for_machine'>" + number + "</Number>" +
		"</Dial>" +
    "</Response>"
}
