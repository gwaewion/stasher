package errorer

import (
	"log"
	"net/http"
	"encoding/json"
)

type Ooops struct {
	Error	string	`json:"error"`
}

func LogError( err error ) {
	if err != nil {
		log.Printf( "Follow error was encountered: ", err )
	}
}

func Ooopsie( responseWriter http.ResponseWriter, errorMessage string, code int ) {
	marshaledOoops, marshaledOoopsError := json.Marshal( Ooops{ Error: errorMessage } )
	LogError( marshaledOoopsError )

	responseWriter.Header().Set( "Content-Type", "application/json" )
	responseWriter.WriteHeader( code )
	responseWriter.Write( marshaledOoops )
}