package oncall

import (
	"fmt"
	"log"
	"github.com/wm/go_pager"
	"net/http"
	"strconv"
	"strings"
)

func Welcome(config *go_pager.Config) string {
	list := ""

	for _, contact := range(config.Contacts) {
		list += contact.Name + ": "
		list += contact.Number + ", "
	}

	return fmt.Sprintf(`here is: %v`, list)
}

func AttemptCall(req *http.Request, w http.ResponseWriter, config *go_pager.Config) string {
	var commands string

	dialCallStatus := req.PostFormValue("DialCallStatus")
	numberIndexStr := req.FormValue("number_index")

	if numberIndexStr == "" {
		numberIndexStr = "0"
	}

	numberIndex, _ := strconv.Atoi(numberIndexStr)
	log.Printf("param[`DialCallStatus`] %v", dialCallStatus)

	if dialCallStatus == "completed"{
		commands = hangUp()
	} else if (numberIndex < len(config.Contacts)) {
		commands = dialContact(numberIndex + 1, (config.Contacts)[numberIndex])
	} else {
		commands = appologise() + hangUp()
	}

	return respond(commands)
}

func ScreenForMachine() string {
	return `
	<Gather action='/oncall/complete_call'>
	<Say>Press any key to accept this call</Say>
	</Gather>
	<Hangup/>
	`
}

func CompleteCall() string {
	return `<Say>Connecting</Say>`
}

func respond(commands string) string {
	response := fmt.Sprintf(`
	<?xml version='1.0' encoding='UTF-8'?>
	<Response>%v</Response>
	`, commands)

	log.Printf("response: %v", response)

	return strings.TrimSpace(response)
}

func hangUp() string {
	return `<Hangup/>`
}

func appologise() string {
	return `
	<Say>
	We are sorry but we are unable to connect you with an on call Engineer.
	Please try again.
	</Say>
	`
}

func dialContact(nextContactIndex int, contact go_pager.Contact) string {
	log.Printf("Forwarding to contact: %v", contact.Name)

	return fmt.Sprintf(`
	<Say>Attempting to connect you with %v, today's on call engineer.</Say>
	<Dial action='/oncall/attempt_call?number_index=%v'>
		<Number url='/oncall/screen_for_machine'>%v</Number>
	</Dial>
	`, contact.Name, strconv.Itoa(nextContactIndex), contact.Number)
}
