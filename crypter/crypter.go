package crypter

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
	"strconv"

	"stasher/errorer"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	salt string
)

func SetSalt( s string ) {
	salt = s
}

func getKey( text string ) []byte {
	return argon2.IDKey( []byte( text ), []byte( salt ), 1, 64*1024, 4, 32 )
}

func getAEAD( key []byte ) cipher.AEAD {
	aead, aeadError := chacha20poly1305.NewX( key )

	errorer.LogError( aeadError )

	return aead
}

func Crypt( phrase, text string ) string {
	aead := getAEAD( getKey( phrase ) )

    nonce := make( []byte, aead.NonceSize(), aead.NonceSize() + len( text ) + aead.Overhead() )
    if _, err := rand.Read( nonce ); err != nil {
        errorer.LogError( err )
    }

    return fmt.Sprintf( "%x", aead.Seal( nonce, nonce, []byte( text ), nil ) )
}

func Decrypt( phrase, encryptedText string ) string {
	aead := getAEAD( getKey( phrase ) )

	var restoredText []byte

	for i := 0; i < len( encryptedText ); i += 2 {
		hexChars := fmt.Sprint( encryptedText[i:i+2] )
		intChars, _ := strconv.ParseInt( hexChars, 16, 64 ) 
		restoredText = append( restoredText, byte( intChars ) )
	}

	if len( restoredText ) < aead.NonceSize() {
        log.Fatal( "ciphertext too short" )
    }

    nonce, ciphertext := restoredText[:aead.NonceSize()], restoredText[aead.NonceSize():]

    plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
    errorer.LogError( err )

    return string( plaintext )
}