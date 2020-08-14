package hasher

import (
	"stasher/errorer"

	"golang.org/x/crypto/bcrypt"
)

func GetHash( text string ) string {
	hash, hashError := bcrypt.GenerateFromPassword( []byte( text ), 10 ) 

	errorer.LogError( hashError )

	return string( hash )
}

func IsTextCorrect( plainText, hash string) bool {
	if bcrypt.CompareHashAndPassword( []byte( hash ), []byte( plainText ) ) != nil {
		return false
	}

	return true
}