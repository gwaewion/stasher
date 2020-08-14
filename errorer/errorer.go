package errorer

import (
	"log"
)

func LogError( err error ) {
	if err != nil {
		log.Fatalf( "Follow error was encountered: ", err )
	}
}