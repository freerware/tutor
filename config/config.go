package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

func _New(filePath string, configuration interface{}) error {

	// read in the contents of the configuration file.
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// expand environment variables.
	fileContents = []byte(os.ExpandEnv(string(fileContents)))
	return yaml.Unmarshal(fileContents, configuration)
}

func New(filePath string) (c Configuration, err error) {
	err = _New(filePath, &c)
	return
}

func NewDefault() (Configuration, error) {
	return New("./configuration.yaml")
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

type ServerConfiguration struct {
	Host string
	Port int
}

type DatabaseConfiguration struct {
	Host      string
	Port      int
	Name      string
	User      string
	Password  string
	ParseTime bool `yaml:"parseTime"`
	Charset   string
}

func (c DatabaseConfiguration) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=%t&charset=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.ParseTime, c.Charset)
}
