package passgen

import (
	// "fmt"
	"math/rand"
	"time"
)

const (
	lowercase string	= "abcdefghijklmnopqrstuvwxyz"
	uppercase string	= "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers string		= "1234567890"
	brackets string		= "()[]{}<>"
	specials string		= "~`!@#$%^&*-_=+/\\?.,'\""
)

var (
	randomizator *rand.Rand = rand.New( rand.NewSource( time.Now().UnixNano() ) )
)

func getRandLowercase() string {
	return string( lowercase[randomizator.Intn(26)] )
}

func getRandUppercase() string {
	return string( uppercase[randomizator.Intn(26)] )
}

func getRandNumber() string {
	return string( numbers[randomizator.Intn(10)] )
}

func getRandBracket() string {
	return string( brackets[randomizator.Intn(8)] )
}

func getRandSpecial() string {
	return string( specials[randomizator.Intn(21)] )
}

func GenerateID( length int ) string {
	var id string

	for i := 0; i <= length; i++ {
		r := randomizator.Intn(2)

		if r == 0 {
			id += getRandLowercase()
		} else if r == 1 {
			id += getRandNumber()
		}
	} 

	return id
}

func GetPassL( length int ) string {
	password := getRandLowercase()

	for i := 0; i <= length; i++ {
		password += getRandLowercase()
	}

	return password
}

func GetPassU( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		password += getRandUppercase()
	}

	return password
}

func GetPassN( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		password += getRandNumber()
	}

	return password
}

func GetPassLU( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandLowercase()
		} else {
			password += getRandUppercase()
		}
	}

	return password
}

func GetPassLN( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandLowercase()
		} else {
			password += getRandNumber()
		}
	}

	return password
}

func GetPassLB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandLowercase()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassLS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandLowercase()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassUN( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandUppercase()
		} else {
			password += getRandNumber()
		}
	}

	return password
}

func GetPassUB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandUppercase()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassUS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandUppercase()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassNB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandNumber()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassNS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(2)

		if selector == 0 {
			password += getRandNumber()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLUN( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else {
			password += getRandNumber()
		}
	}

	return password
}

func GetPassLUB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassLUS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLNB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassLNS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassUNB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandUppercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassUNS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandUppercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassUBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandUppercase()
		} else if selector == 1 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassNBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(3)

		if selector == 0 {
			password += getRandNumber()
		} else if selector == 1 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLUNB( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(4)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else if selector == 2 {
			password += getRandNumber()
		} else {
			password += getRandBracket()
		}
	}

	return password
}

func GetPassLUNS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(4)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else if selector == 2 {
			password += getRandNumber()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLUBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(4)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else if selector == 2 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLNBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(4)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else if selector == 2 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassUNBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(4)

		if selector == 0 {
			password += getRandUppercase()
		} else if selector == 1 {
			password += getRandNumber()
		} else if selector == 2 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}

func GetPassLUNBS( length int ) string {
	var password string

	for i := 0; i <= length; i++ {
		selector := randomizator.Intn(5)

		if selector == 0 {
			password += getRandLowercase()
		} else if selector == 1 {
			password += getRandUppercase()
		} else if selector == 2 {
			password += getRandNumber()
		} else if selector == 3 {
			password += getRandBracket()
		} else {
			password += getRandSpecial()
		}
	}

	return password
}