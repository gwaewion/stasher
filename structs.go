package main

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
	Web			string		`json:"web"`
	ApiHint	ApiHint	`json:"api"`
}

type ApiHint struct {
	Url	string	`json:"url"`
	Id		string	`json:"id"`
}


type RootHandlerNew struct {}