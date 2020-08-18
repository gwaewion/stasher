package main

import (
	"strings"
	"io/ioutil"
	"net/http"
	"errors"
	"bytes"
	"encoding/json"

	"stasher/errorer"
)

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

func conditionCheck( responseWriter http.ResponseWriter, request *http.Request ) ([]byte, error) {
	contentTypesList := request.Header.Values( "Content-Type" )

	if !( contains( contentTypesList, "application/json" ) ) {
		responseWriter.WriteHeader( http.StatusUnsupportedMediaType )
		return nil, errors.New( "not application/json type" )
	}

	requestBody, requestBodyError := ioutil.ReadAll( request.Body )
    if requestBodyError != nil {
		responseWriter.WriteHeader( http.StatusBadRequest )
		return nil, errors.New( "can't read request" )
    }

    return requestBody, nil
}

func makeRequest( httpClient http.Client, method, uri string, payload []byte ) ( int, []byte ) {
	var request *http.Request
	var requestError error
	method = strings.ToUpper( method )

	if method == "POST" {
		request, requestError = http.NewRequest( method, uri, bytes.NewReader( payload ) )
		errorer.LogError( requestError )
		request.Header.Add( "Content-Type", "application/json" )
	} else if method == "GET" || method == "DELETE" {
		request, requestError = http.NewRequest( method, uri, nil )
		errorer.LogError( requestError )
	}
	
	response, responseError := httpClient.Do( request )
	errorer.LogError( responseError )
	defer response.Body.Close()

	body, bodyError := ioutil.ReadAll( response.Body )
	errorer.LogError( bodyError )

	return response.StatusCode, body
}

func sendJSON( responseWriter http.ResponseWriter, object interface{}, code int ) {
	marshaledObject, marshaledObjectError := json.Marshal( object )
	errorer.LogError( marshaledObjectError )

	responseWriter.Header().Set( "Content-Type", "application/json" )
	responseWriter.WriteHeader( code )
	responseWriter.Write( marshaledObject )
}