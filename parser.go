package emailparser


import (
	"os"
	"fmt"
	"bytes"
	"regexp"
	"strings"
	"net/mail"
	"encoding/json"
	
	_ "github.com/denisenkom/go-mssqldb"
	
)

//Error : Custom error struct
type Error struct {
	msg string
	err error
}

//Email : the extracted email information from the database
type Email struct {
	EnvelDate string `json:"envelDate"`
	EnvelFrom string `json:"envelFrom"`
	EnvelSubject string `json:"envelSubject"`
	EnvelID string `json:"envelID"`
	EnvelTo string `json:"envelTo"`

	DsnBody string `json:"dsnBody"`//content type text/plain
	DsnHeader string `json:"dsnHeader"`//content type message/delivery-status
	Rfc string `json:"rfc"`//content type message/rfc822 
}

func (e * Email) print() {
	fmt.Println("Email: ")
	fmt.Println("  Envelope: ")
	fmt.Println("	Date:" + e.EnvelDate)
	fmt.Println("	From:" + e.EnvelFrom)
	fmt.Println("	Subject:" + e.EnvelSubject)
	fmt.Println("	ID:" + e.EnvelID)
	fmt.Println("	To:" + e.EnvelTo)
	fmt.Println("	DSN:" + e.DsnBody)
	fmt.Println("	Header:" + e.DsnHeader)
	fmt.Println("	RFC:" + e.Rfc)
}

//Log error and descriptive message to console
func logError(err *Error, msg string) {
	if err != nil {
		fmt.Println(msg, err.msg, ": ", err.err)
		os.Exit(1);
	}
}

func formatJSON(e string) (map[string]interface{}, *Error){
	var email Email
	

	//find boundary string 
	re := regexp.MustCompile(`boundary="?([0-9A-Za-z+/=.-]+)"?`)
	res := re.FindAllStringSubmatch(e,-1)
	
	if(res != nil) {
		token := res[0][1];

		//separate email by boundary 
		content := strings.Split(e,"--"+token)

		//first section is envelope, get envelope details 
		envelope, err := mail.ReadMessage(bytes.NewReader([]byte(e)))
		if err != nil {
			return nil, &Error{"Error parsing envelope",err}
		}

		email.EnvelDate = envelope.Header.Get("Date")
		email.EnvelFrom = envelope.Header.Get("From")
		email.EnvelSubject = envelope.Header.Get("Subject")
		email.EnvelID = envelope.Header.Get("Message-Id")
		email.EnvelTo = envelope.Header.Get("To")
		

	
		//go through each section, first line is Content-Type - 
		for i := 1; i < len(content)-1; i++{

			if(strings.Contains(content[i],"text/plain") == true){
				email.DsnBody = content[i]
			} else if(strings.Contains(content[i],"message/delivery-status") == true){
				email.DsnHeader = content[i]
			} else if(strings.Contains(content[i],"message/rfc822") == true){
				email.Rfc = content[i]
			}
			
		}
	} else {
		//if no token is found, format is weird, combo DSN and envelope
    	
	}

	var parsed map[string]interface{}

	temp, err := json.Marshal(email)
	if err != nil || strings.Compare(string(temp), "{}" ) == 0 {
		return  nil, &Error{"Error Marshalling email to JSON",err}
	}

	json.Unmarshal(temp, &parsed)
	
	return parsed, nil
}



