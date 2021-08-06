// Package env_wrapper provides simplified access to environment variables and docker secrets.
// If a secret is present the environment variable will be ignored.
// Every secret needs an "ENV_" prefix.
package env_wrapper

import (
	"os"
	"strconv"
	"strings"

	secrets "github.com/ijustfool/docker-secrets"
)

type env_wrapper struct {
	secretsEnabled bool
	secretsReader  *secrets.DockerSecrets
}

// Creates a new EnvWrapper with the default secret directory.
func Default() *env_wrapper {
	return New("")
}

// Creates a new EnvWrapper with a custom secret directory.
func New(secretsDir string) *env_wrapper {
	dockerSecrets, err := secrets.NewDockerSecrets(secretsDir)
	res := &env_wrapper{
		(err != nil),
		dockerSecrets,
	}
	if _, err := os.Stat(secretsDir); os.IsNotExist(err) {
		res.secretsEnabled = false
	}
	return res
}

// Gets a string value or returns an empty string if the variable doesn't exist.
func (w *env_wrapper) GetString(name string) string {
	return w.GetStringDef(name, "")
}

// Gets a string value or returns a default value if the string is empty.
func (w *env_wrapper) GetStringDef(name, defval string) string {
	res := defval
	hasval := false
	upname := strings.ToUpper(name)
	secname := "ENV_" + upname
	if w.secretsEnabled {
		secret, err := w.secretsReader.Get(secname)
		if err != nil {
			res = strings.TrimSpace(secret)
			hasval = true
		}
	}
	if !hasval {
		envval := strings.TrimSpace(os.Getenv(upname))
		if len(envval) > 0 {
			res = envval
		}
	}

	return res
}

// Gets a boolean value or returns false if the variable doesn't exist.
func (w *env_wrapper) GetBool(name string) bool {
	return w.GetBoolDef(name, false)
}

// Gets a boolean value or returns a default value if variable doesn't exist.
func (w *env_wrapper) GetBoolDef(name string, defval bool) bool {
	strval := w.GetString(name)
	if len(strval) > 0 {
		res, err := strconv.ParseBool(strval)
		if err == nil {
			return res
		}
	}
	return defval
}

// Gets a integer value or returns 0 if the variable doesn't exist.
func (w *env_wrapper) GetInt(name string) int {
	return w.GetIntDef(name, 0)
}

// Gets a integer value or returns a default value if variable doesn't exist.
func (w *env_wrapper) GetIntDef(name string, defval int) int {
	strval := w.GetString(name)
	if len(strval) > 0 {
		res, err := strconv.Atoi(strval)
		if err == nil {
			return res
		}
	}
	return defval
}

// Gets a string array by splitting the value with the whitespace character.
func (w *env_wrapper) GetStringArray(name string) []string {
	return w.GetStringArraySep(name, " ")
}

// Gets a string array by splitting the value with a seperator.
func (w *env_wrapper) GetStringArraySep(name, seperator string) []string {
	res := []string{}
	strval := w.GetString(name)
	if len(strval) > 0 {
		strparts := strings.Split(strval, seperator)
		for _, s := range strparts {
			cleans := strings.TrimSpace(s)
			if len(cleans) > 0 {
				res = append(res, cleans)
			}
		}
	}
	return res
}
