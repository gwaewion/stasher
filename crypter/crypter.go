package crypter

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
	"strconv"
	// "errors"
	// "bytes"

	"stasher/errorer"

	// "golang.org/x/crypto/twofish"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	salt = "actuallyThisIsSaltedSugar"
)

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
        log.Fatal("ciphertext too short")
    }

    nonce, ciphertext := restoredText[:aead.NonceSize()], restoredText[aead.NonceSize():]

    plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
    errorer.LogError( err )

    return string( plaintext )
}

// func pad( buf []byte ) ( []byte, error ) {
// 	bufLen := len( buf )
// 	padLen := twofish.BlockSize - ( bufLen % twofish.BlockSize )
// 	padText := bytes.Repeat( []byte{byte( padLen )}, padLen )
// 	return append( buf, padText... ), nil
// }

// func unpad( buf []byte ) ( []byte, error ) {
// 	bufLen := len(buf)

// 	if bufLen == 0 {
// 		return nil, errors.New("cryptgo/padding: invalid padding size1")
// 	}

// 	pad := buf[bufLen-1]
// 	padLen := int(pad)
// 	if padLen > bufLen || padLen > twofish.BlockSize {
// 		// return nil, errors.New("cryptgo/padding: invalid padding size2")
// 		return buf, nil
// 	}

// 	for _, v := range buf[bufLen-padLen : bufLen-1] {
// 		if v != pad {
// 			return nil, errors.New("cryptgo/padding: invalid padding3")
// 		}
// 	}

// 	return buf[:bufLen-padLen], nil
// }

// func Crypt( key, message string ) string {
// 	text := []byte( message )

// 	if len( text )%twofish.BlockSize != 0 {
// 		// is this safe?
// 		text, _ = pad( text )
// 	}

// 	block, blockErr := twofish.NewCipher( []byte( key ) )
// 	if blockErr != nil {
// 		panic( blockErr )
// 	}

// 	ciphertext := make( []byte, twofish.BlockSize+len( text ) )
// 	iv := ciphertext[:twofish.BlockSize]

// 	if _, fillErr := io.ReadFull( rand.Reader, iv ); fillErr != nil {
// 		panic( fillErr )
// 	}

// 	mode := cipher.NewCBCEncrypter( block, iv )
// 	mode.CryptBlocks( ciphertext[twofish.BlockSize:], text )

// 	return fmt.Sprintf( "%x", ciphertext) 
// }

// func Decrypt( key, cryptedMessage string ) string {
// 	var text []byte

// 	for i := 0; i < len( cryptedMessage ); i += 2 {
// 		hexChars := fmt.Sprint( cryptedMessage[i:i+2] )
// 		intChars, _ := strconv.ParseInt( hexChars, 16, 64 ) 
// 		text = append( text, byte( intChars ) )
// 	}

// 	block, blockErr := twofish.NewCipher( []byte( key ) )
// 	if blockErr != nil {
// 		panic( blockErr )
// 	}

// 	if len( text ) < twofish.BlockSize {
// 		panic( "ciphertext too short" )
// 	}

// 	iv := text[:twofish.BlockSize]
// 	message := text[twofish.BlockSize:]

// 	if len( message )%twofish.BlockSize != 0 {
// 		panic( "ciphertext is not a multiple of the block size" )
// 	}

// 	mode := cipher.NewCBCDecrypter( block, iv ) 
// 	// is this safe?
// 	mode.CryptBlocks( []byte( message ), []byte( message ) )

// 	message, _ = unpad( message )

// 	return fmt.Sprintf( "%s", message )
// }

// func main() {
// 	key := "1234567812345678"
// 	c := Crypt( key, `1234567812345678` )
// 	fmt.Printf( "%q\n", c )
// 	fmt.Printf( "%q\n", Decrypt( key, c ) )
// }
