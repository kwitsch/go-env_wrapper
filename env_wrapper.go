// Package EnvWrapper provides simplified access to environment variables and docker secrets.
// If a secret is present the environment variable will be ignored.
// Every secret needs an "ENV_" prefix.
package EnvWrapper

import (
	"os"
	"strconv"
	"strings"

	secrets "github.com/ijustfool/docker-secrets"
)

type EnvWrapper struct {
	secretsEnabled bool
	secretsReader  *secrets.DockerSecrets
}

// Creates a new EnvWrapper with the default secret directory.
func Default() *EnvWrapper {
	return New("")
}

// Creates a new EnvWrapper with a custom secret directory.
func New(secretsDir string) *EnvWrapper {
	dockerSecrets, err := secrets.NewDockerSecrets(secretsDir)
	res := &EnvWrapper{
		(err != nil),
		dockerSecrets,
	}
	return res
}

// Gets a string value or returns an empty string if the variable doesn't exist.
func (w *EnvWrapper) GetString(name string) string {
	return w.GetStringDef(name, "")
}

// Gets a string value or returns a default value if the string is empty.
func (w *EnvWrapper) GetStringDef(name, defval string) string {
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
func (w *EnvWrapper) GetBool(name string) bool {
	return w.GetBoolDef(name, false)
}

// Gets a boolean value or returns a default value if variable doesn't exist.
func (w *EnvWrapper) GetBoolDef(name string, defval bool) bool {
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
func (w *EnvWrapper) GetInt(name string) int {
	return w.GetIntDef(name, 0)
}

// Gets a integer value or returns a default value if variable doesn't exist.
func (w *EnvWrapper) GetIntDef(name string, defval int) int {
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
func (w *EnvWrapper) GetStringArray(name string) []string {
	return w.GetStringArraySep(name, " ")
}

// Gets a string array by splitting the value with a seperator.
func (w *EnvWrapper) GetStringArraySep(name, seperator string) []string {
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
