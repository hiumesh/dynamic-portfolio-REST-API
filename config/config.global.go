package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

type ServerConfiguration struct {
	Id                   string
	MaxPerUserConnection string `envconfig:"GO_SOCKET_MAX_PER_USER_CONNECTION" default:"2"`
}

type APIConfiguration struct {
	Host string
	Port string `envconfig:"PORT" default:"8080"`
}

func (c *APIConfiguration) Validate() error {
	return nil
}

type DBConfiguration struct {
	URL string `envconfig:"DATABASE_URL"`
}

type CookieConfiguration struct {
	Key      string `json:"key" envconfig:"COOKIE_KEY"`
	Domain   string `json:"domain" envconfig:"COOKIE_KEY"`
	Duration int    `json:"duration"`
}

type JWTConfiguration struct {
	Secret string `json:"secret" required:"true"`
}

func (c *DBConfiguration) Validate() error {
	return nil
}

type CORSConfiguration struct {
	AllowedHeaders []string `json:"allowed_headers" split_words:"true"`
}

func (c *CORSConfiguration) AllAllowedHeaders(defaults []string) []string {
	set := make(map[string]bool)
	for _, header := range defaults {
		set[header] = true
	}

	var result []string
	result = append(result, defaults...)

	for _, header := range c.AllowedHeaders {
		if !set[header] {
			result = append(result, header)
		}

		set[header] = true
	}

	return result
}

type GlobalConfiguration struct {
	SERVER  ServerConfiguration
	API     APIConfiguration
	DB      DBConfiguration
	CORS    CORSConfiguration   `json:"cors"`
	JWT     JWTConfiguration    `json:"jwt"`
	COOKIE  CookieConfiguration `json:"cookies"`
	LOGGING LoggingConfig       `envconfig:"LOG"`
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func LoadGlobal(filename string) (*GlobalConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(GlobalConfiguration)
	if err := envconfig.InitWithPrefix(config, "GO"); err != nil {
		return nil, err
	}

	// config.SERVER.Id = utils.GenerateUniqueServerId()

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *GlobalConfiguration) Validate() error {
	validatables := []interface {
		Validate() error
	}{
		&c.API,
		&c.DB,
	}

	for _, validatable := range validatables {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}
