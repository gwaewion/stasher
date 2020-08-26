//TODO: add logs out to stderr as default
//TODO: add logging to file
//TODO: add debug mode with more verbose output
//TODO: add using secure connection to couchdb
//TODO: add encryption of every secret
//TODO: add graceful shutdown
//TODO: make variables names more logical
//TODO: add check of config variables

package main

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	// "io/ioutil"
	"fmt"
	// "bytes"
	"flag"
	"os"
	"crypto/tls"
	"strings"

	"stasher/errorer"
	"stasher/configurer"
	"stasher/hasher"
	"stasher/crypter"
	"stasher/passgen"

	"github.com/gorilla/mux"
	// "github.com/google/uuid"
	"github.com/gobuffalo/packr/v2"

)

var (
	config configurer.Config
	apiVersion = "1.0"
	apiPath = "/api/" + apiVersion + "/"
	couchDBUri string
	httpClient http.Client
	idSize = 8
)

func init() {
	configPath := flag.String( "c", "", "config file location")
	flag.Parse()

	if *configPath == "" {
		fmt.Println( "Config file not found." )
		os.Exit(1)
	}

	config = configurer.ParseConfig( *configPath )	
	//add exception check
	crypter.SetSalt( config.Stasher.Salt )

	couchDBUri = config.CouchDB.Protocol + "://" + config.CouchDB.Username + ":" + config.CouchDB.Password + "@" + config.CouchDB.Address + ":" + config.CouchDB.Port + "/" + config.CouchDB.DBName

	if config.CouchDB.CertCheck == true {
		httpClient = http.Client{}
	} else {
		httpClient = http.Client{ Transport: &http.Transport{ TLSClientConfig: &tls.Config{ InsecureSkipVerify: true }, } }
	}
}


func ApiSetSecretHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	requestBody, requestBodyError := conditionCheck( responseWriter, request )
	errorer.LogError( requestBodyError ) 

    var secret SetSecret

    secretUnmarshalError := json.Unmarshal( requestBody, &secret ) 
	if secretUnmarshalError != nil || secret.Message == "" {
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

	// id := uuid.New().String()
	id := passgen.GenerateID( idSize )
	var record DBRecord

	if secret.Phrase != "" {
		hash := hasher.GetHash( secret.Phrase )
		securedMessage := crypter.Crypt( secret.Phrase, secret.Message )
		record = DBRecord{ Id: id, Message: securedMessage, Secure: true, PhraseHash: hash, CreatedAt: time.Now().Format( time.RFC3339 ) }
	} else {
		record = DBRecord{ Id: id, Message: secret.Message, CreatedAt: time.Now().Format( time.RFC3339 ) }
	}

	marshaledRecord, marshaledRecordError := json.Marshal( record )
	errorer.LogError( marshaledRecordError )

	recordStatusCode, _ := makeRequest( httpClient, "post", couchDBUri, marshaledRecord )

	if recordStatusCode == 201 {
		web := config.Stasher.Scheme + "://" + config.Stasher.Hostname + "/" + id
		apiUrl := config.Stasher.Scheme + "://" + config.Stasher.Hostname + apiPath + "getSecret"
		sendJSON( responseWriter, Hint{ Web: web, ApiHint: ApiHint{ Url: apiUrl, Id: id } }, 201 )
	} else {
		log.Fatalf( "Response code is %v", recordStatusCode )
	}
}

func ApiGetSecretHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	requestBody, requestBodyError := conditionCheck( responseWriter, request )
	errorer.LogError( requestBodyError ) 

    var secret GetSecret

    secretUnmarshalError := json.Unmarshal( requestBody, &secret ) 
	if secretUnmarshalError != nil || secret.Id == "" { 
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

	secretStatusCode, secretBody := makeRequest( httpClient, "get", couchDBUri + "/" + secret.Id, nil )

	if secretStatusCode != 200 {
		errorer.Ooopsie( responseWriter, "secret not exists", 404 )
		return
	}

	var record DBRecord
	
	recordUnmarshalError := json.Unmarshal( secretBody, &record )
	errorer.LogError( recordUnmarshalError )

	var message string

	if secret.Phrase != "" && record.Secure {
		if hasher.IsTextCorrect( secret.Phrase, record.PhraseHash ) {
			message = crypter.Decrypt( secret.Phrase, record.Message )
		} else {
			errorer.Ooopsie( responseWriter, "wrong phrase", 400 )
			return
		}
	} else if  secret.Phrase == "" && record.Secure {
		errorer.Ooopsie( responseWriter, "no phrase", 400 )
		return	
	}

	if message == "" {
		message = record.Message
	}

	sendJSON( responseWriter, Secret{ Message: message }, 200 )

	_, _ = makeRequest( httpClient, "delete", couchDBUri + "/" + secret.Id + "?rev=" + record.Revision, nil )

}

func ( rh RootHandlerNew ) ServeHTTP( responseWriter http.ResponseWriter, request *http.Request ) {
	path := request.URL.Path
	webroot := packr.New( "webroot", "./webroot" )

	if path == "/secret.js" {
		secretjs, secretjsError := webroot.Find( "secret.js" )
		errorer.LogError( secretjsError )
		responseWriter.Write( secretjs )
	} else if path == "/favicon.ico" {
		icon, iconError := webroot.Find( "favicon.ico" )
		errorer.LogError( iconError )
		responseWriter.Write( icon )
	} else if path == "/script.js" {
		script, scriptError := webroot.Find( "script.js" )
		errorer.LogError( scriptError )
		responseWriter.Write( script )
	} else if path == "/style.css" {
		style, styleError := webroot.Find( "style.css" )
		errorer.LogError( styleError )
		responseWriter.Header().Set( "Content-Type", "text/css" )
		responseWriter.Write( style )
	} else if path == "/" {
		index, indexError := webroot.Find( "index.html" )
		errorer.LogError( indexError )
		responseWriter.Write( index )
	} else {
		trimmedPath := strings.TrimLeft( path, "/" )

		if len( trimmedPath ) == idSize {
			headResponseCode, _ := makeRequest( httpClient, "head", couchDBUri + "/" + trimmedPath, nil )

			if headResponseCode == 200 {
					secret, secretError := webroot.Find( "secret.html" )
					errorer.LogError( secretError )
					responseWriter.Write( secret )
			} else {
				responseWriter.WriteHeader( http.StatusNotFound )	
			} 
		} else {
			responseWriter.WriteHeader( http.StatusNotFound )
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc( apiPath + "setSecret", ApiSetSecretHandler ).Methods( "POST" )
	router.HandleFunc( apiPath + "getSecret", ApiGetSecretHandler ).Methods( "POST" )
	router.PathPrefix( "/" ).Handler( RootHandlerNew{} ).Methods( "GET" )

    server := &http.Server{
        Addr:         			config.Stasher.Address + ":" + config.Stasher.Port,
        WriteTimeout: 	time.Second * 15,
        ReadTimeout:  	time.Second * 15,
        IdleTimeout:  		time.Second * 60,
        Handler: 				router,
    }

    go log.Fatal( server.ListenAndServe() )
    
}