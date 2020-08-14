package main

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"io/ioutil"
	// "io"
	"fmt"
	"bytes"
	"strings"
	"flag"
	"os"

	"stasher/errorer"
	"stasher/configurer"
	"stasher/hasher"
	"stasher/crypter"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/gobuffalo/packr/v2"

)

var (
	config configurer.Config
	apiVersion = "1.0"
	apiPath = "/api/" + apiVersion + "/"
	couchDBUri string
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

	couchDBUri = config.CouchDB.Protocol + "://" + config.CouchDB.Address + ":" + config.CouchDB.Port + "/" + config.CouchDB.DBName
}

type Ooops struct {
	Error	string	`json:"error"`
}

type Secret struct {
	Message	string	`json:"message"`	
}

type SetSecret struct {
	Message	string	`json:"message"`
	Phrase		string	`json:"phrase,omitempty"`
}

type GetSecret struct {
	Id				string	`json:"id"`
	Phrase		string	`json:"phrase,omitempty"`
}

type DBRecord struct {
	Id						string	`json:"_id"`
	Revision			string	`json:"_rev,omitempty"`
	Message			string	`json:"message"`
	Secure				bool		`json:"secure,omitempty"`
	PhraseHash		string	`json:"phrase,omitempty"`
	CreatedAt		string	`json:"created_at"`
}

type Hint struct {
	Url	string	`json:"url"`
}

func contains( list []string, word string ) bool {
	lowerWord := strings.ToLower( word )
	realList := strings.Split( list[0], ";")
	for _, element := range realList {
		if lowerWord == strings.ToLower( element ) {
			return true
		}
	}
	return false
}

func ApiSetSecretHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	// couchDBUri := config.CouchDB.Protocol + "://" + config.CouchDB.Address + ":" + config.CouchDB.Port + "/" + config.CouchDB.DBName

	contentTypesList := request.Header.Values( "Content-Type" )

	if !( contains( contentTypesList, "application/json" ) ) {
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusUnsupportedMediaType )
		return
	}

	requestBody, requestBodyError := ioutil.ReadAll( request.Body )
    if requestBodyError != nil {
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

    var secret SetSecret
    secretUnmarshalError := json.Unmarshal( requestBody, &secret ) 
	if secretUnmarshalError != nil || secret.Message == "" {
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

	id := uuid.New().String()
	var record DBRecord

	if secret.Phrase != "" {
		hash := hasher.GetHash( secret.Phrase )
		securedMessage := crypter.Crypt( secret.Phrase, secret.Message )
		record = DBRecord{ Id: id, Message: securedMessage, Secure: true, PhraseHash: hash, CreatedAt: time.Now().Format( time.RFC3339 ) }
		// fmt.Println( crypter.Decrypt( secret.Phrase, securedMessage ) )
	} else {
		record = DBRecord{ Id: id, Message: secret.Message, CreatedAt: time.Now().Format( time.RFC3339 ) }
	}
	marshaledRecord, marshaledRecordError := json.Marshal( record )
	errorer.LogError( marshaledRecordError )

	recordResponse, recordResponseError := http.Post( couchDBUri, "application/json", bytes.NewReader( marshaledRecord ) )
	errorer.LogError( recordResponseError )
	if recordResponse.StatusCode == 201 {
		url := "http://" + config.Stasher.Address + ":" + config.Stasher.Port + "/secret/"
		hint, hintError := json.Marshal( Hint{ Url: url + id } )
		errorer.LogError( hintError )
		responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusCreated )
		responseWriter.Write( hint )
	} else {
		log.Fatalf( "Response code is %v", recordResponse.StatusCode )
	}
}

func ApiGetSecretHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	// couchDBUri := config.CouchDB.Protocol + "://" + config.CouchDB.Address + ":" + config.CouchDB.Port + "/" + config.CouchDB.DBName

	contentTypesList := request.Header.Values( "Content-Type" )

	if !( contains( contentTypesList, "application/json" ) ) {
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusUnsupportedMediaType )
		return
	}

	requestBody, requestBodyError := ioutil.ReadAll( request.Body )
    if requestBodyError != nil {
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

    var secret GetSecret
    secretUnmarshalError := json.Unmarshal( requestBody, &secret ) 
	if secretUnmarshalError != nil || secret.Id == "" { //what if secret ID existst but wrong?
		// responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusBadRequest )
		return
    }

	dbRecord, dbRecordError := http.Get( couchDBUri + "/" + secret.Id )
	errorer.LogError( dbRecordError )

	if dbRecord.StatusCode != 200 {
		marshaledOoops, marshaledOoopsError := json.Marshal( Ooops{ Error: "secret not exists" } )
		errorer.LogError( marshaledOoopsError )

		responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusNotFound )
		responseWriter.Write( marshaledOoops)
		return
	}

	var record DBRecord
	
	recordBody, recordBodyError := ioutil.ReadAll( dbRecord.Body )
	errorer.LogError( recordBodyError )
	recordUnmarshalError := json.Unmarshal( recordBody, &record )
	errorer.LogError( recordUnmarshalError )

	var message string
	fmt.Printf( "%#v\n", secret )
	fmt.Printf( "%#v\n", record )

	if secret.Phrase != "" && record.Secure {
		if hasher.IsTextCorrect( secret.Phrase, record.PhraseHash ) {
			message = crypter.Decrypt( secret.Phrase, record.Message )
		} else {
			marshaledOoops, marshaledOoopsError := json.Marshal( Ooops{ Error: "wrong phrase" } )
			errorer.LogError( marshaledOoopsError )

			responseWriter.Header().Set( "Content-Type", "application/json" )
			responseWriter.WriteHeader( http.StatusBadRequest )
			responseWriter.Write( marshaledOoops )
			return
		}
	} else if  secret.Phrase == "" && record.Secure {
		marshaledOoops, marshaledOoopsError := json.Marshal( Ooops{ Error: "no phrase" } )
		errorer.LogError( marshaledOoopsError )

		responseWriter.Header().Set( "Content-Type", "application/json" )
		responseWriter.WriteHeader( http.StatusBadRequest )
		responseWriter.Write( marshaledOoops )
		return	
	}

	fmt.Printf( "%#v\n", message )
	if message == "" {
		message = record.Message
	}
	marshaledSecret, marshaledSecretError := json.Marshal( Secret{ Message: message } )
	errorer.LogError( marshaledSecretError )

	responseWriter.Header().Set( "Content-Type", "application/json" )
	responseWriter.Write( marshaledSecret )	

	client := &http.Client{}
	deleteRequest, deleteRequestError := http.NewRequest( "DELETE", couchDBUri + "/" + secret.Id + "?rev=" + record.Revision, nil )
	errorer.LogError( deleteRequestError )
	deleteRequestResponse, deleteRequestResponseError := client.Do( deleteRequest )
	errorer.LogError( deleteRequestResponseError )
	_, deleteRequestResponseBodyError := ioutil.ReadAll( deleteRequestResponse.Body )
	errorer.LogError( deleteRequestResponseBodyError )
}

// func SecretHandler( responseWriter http.ResponseWriter, request *http.Request ) {
// 	couchDBUri := config.CouchDB.Protocol + "://" + config.CouchDB.Address + ":" + config.CouchDB.Port + "/" + config.CouchDB.DBName

// 	muxerVars := mux.Vars( request )
// 	id := muxerVars[ "id" ]
// 	dbRecord, dbRecordError := http.Get( couchDBUri + "/" + id )
// 	errorer.LogError( dbRecordError )

// 	var secret SetSecret
// 	var record DBRecord
	
// 	recordBody, recordBodyError := ioutil.ReadAll( dbRecord.Body )
// 	errorer.LogError( recordBodyError )

// 	secretUnmarshalError := json.Unmarshal( recordBody, &secret )
// 	errorer.LogError( secretUnmarshalError )
// 	recordUnmarshalError := json.Unmarshal( recordBody, &record )
// 	errorer.LogError( recordUnmarshalError )

// 	marshaledSecret, marshaledSecretError := json.Marshal( secret )
// 	errorer.LogError( marshaledSecretError )
// 	responseWriter.Header().Set( "Content-Type", "application/json" )
// 	responseWriter.Write( marshaledSecret )

// 	client := &http.Client{}
// 	deleteRequest, deleteRequestError := http.NewRequest( "DELETE", couchDBUri + "/" + id + "?rev=" + record.Revision, nil )
// 	errorer.LogError( deleteRequestError )
// 	deleteRequestResponse, deleteRequestResponseError := client.Do( deleteRequest )
// 	errorer.LogError( deleteRequestResponseError )
// 	_, deleteRequestResponseBodyError := ioutil.ReadAll( deleteRequestResponse.Body )
// 	errorer.LogError( deleteRequestResponseBodyError )
// }

func SecretJSHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	webroot := packr.New( "webroot", "./webroot" )
	secretjs, secretjsError := webroot.Find( "secret.js" )
	errorer.LogError( secretjsError )
	responseWriter.Write( secretjs )
}

func SecretHTMLHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	webroot := packr.New( "webroot", "./webroot" )
	secret, secretError := webroot.Find( "secret.html" )
	errorer.LogError( secretError )
	responseWriter.Write( secret )
}

func ScriptHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	webroot := packr.New( "webroot", "./webroot" )
	script, scriptError := webroot.Find( "script.js" )
	errorer.LogError( scriptError )
	responseWriter.Write( script )
}

func RootHandler( responseWriter http.ResponseWriter, request *http.Request ) {
	webroot := packr.New( "webroot", "./webroot" )
	index, indexError := webroot.Find( "index.html" )
	errorer.LogError( indexError )
	responseWriter.Write( index )
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc( apiPath + "setSecret", ApiSetSecretHandler ).Methods( "POST" )
	router.HandleFunc( apiPath + "getSecret", ApiGetSecretHandler ).Methods( "POST" )
	// router.HandleFunc( "/secret/{id}", SecretHandler ).Methods( "GET" )
	
	router.HandleFunc( "/secret/{id}", SecretHTMLHandler ).Methods( "GET" ) 
	router.HandleFunc( "/secret.js", SecretJSHandler ).Methods( "GET" ) 
	router.HandleFunc( "/script.js", ScriptHandler ).Methods( "GET" )
	router.HandleFunc( "/", RootHandler ).Methods( "GET" )
	// router.PathPrefix( "/" ).HandleFunc( RootHandler ).Methods( "GET" )

    server := &http.Server{
        Addr:         			config.Stasher.Address + ":" + config.Stasher.Port,
        WriteTimeout: 	time.Second * 15,
        ReadTimeout:  	time.Second * 15,
        IdleTimeout:  		time.Second * 60,
        Handler: 				router,
    }

    go log.Fatal( server.ListenAndServe() )
    
}