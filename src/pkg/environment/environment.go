// Package environment defines helpers accessing environment values
package environment

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Env represents environmental variable instance
type Env struct {
	envCache map[string]string
}

// LoadEnv creates a new instance of Env and returns an error if any occurs depending on the environment
func LoadEnv() (*Env, error) {
	envVal := os.Getenv("APP_ENV")
	if strings.EqualFold(envVal, "dev") {
		env, err := NewLoadFromFile("./dev.env")
		if err != nil {
			return nil, err
		}

		return env, nil
	}

	env, err := New()
	if err != nil {
		return nil, err
	}

	return env, nil
}

// New creates a new instance of Env and returns an error if any occurs
func New() (*Env, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, err
	}

	return &Env{}, nil
}

// NewLoadFromFile lets you load Env object from a file
func NewLoadFromFile(fileName string) (*Env, error) {
	err := godotenv.Load(fileName)
	if err != nil {
		return nil, err
	}

	return &Env{}, nil
}

// Get retrieves the string value of an environmental variable
func (e *Env) Get(key string) string {
	return os.Getenv(key)
}

// IsUnitTest is helper that returns true or false if the environment is executed in unit test
func (e *Env) IsUnitTest() bool {
	v := e.Get("IS_UNIT_TEST")
	return strings.EqualFold(v, "true")
}

// UseMock is helper that returns true or false if the environment should use mocks when hitting 3rd party partners
func (e *Env) UseMock() bool {
	v := e.Get("APP_MOCK")
	if len(v) == 0 {
		return false
	}

	if strings.EqualFold(v, "true") {
		return true
	}

	return false
}

// HelperForMocking [do not use in logic] designed for mocking and in test suite
func (e *Env) HelperForMocking(cache map[string]string) {
	e.envCache = cache
}
