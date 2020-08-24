package configurer

import (
	"os"
	"io/ioutil"

	"stasher/errorer"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Stasher	StasherSection		`yaml:"stasher"`
	CouchDB	CouchDBSection	`yaml:"couchdb"`
}

type StasherSection struct {
	Address		string	`yaml:"address"`
	Port				string	`yaml:"port"`
	Salt				string	`yaml:"salt"`
	Hostname	string	`yaml:"hostname"`
}

type CouchDBSection struct {
	Protocol		string 	`yaml:"protocol"`
	Address		string 	`yaml:"address"`
	Port				string 	`yaml:"port"`
	DBName		string	`yaml:"dbname"`
	Username	string	`yaml:"username"`
	Password	string	`yaml:"password"`
	CertCheck	bool		`yaml:"cert_check"`
}

//add check for not empty values
func ParseConfig( filepath string ) Config {
	file, fileError := os.Open( filepath )
	errorer.LogError( fileError )
	config, configError := ioutil.ReadAll( file )
	errorer.LogError( configError )

	var marshaledConfig Config
	_ = yaml.Unmarshal( config, &marshaledConfig )

	return marshaledConfig
}