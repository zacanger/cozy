package evaluator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/zacanger/cozy/object"
)

var searchPaths []string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting cwd: %s", err)
	}

	if e := os.Getenv("COZYPATH"); e != "" {
		tokens := strings.Split(e, ":")
		for _, token := range tokens {
			addPath(token) // ignore errors
		}
	} else {
		searchPaths = append(searchPaths, cwd)
	}
}

func addPath(path string) error {
	path = os.ExpandEnv(filepath.Clean(path))
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	searchPaths = append(searchPaths, absPath)
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FindModule finds a module based on name, used by the evaluator
func FindModule(name string) string {
	basename := fmt.Sprintf("%s.cz", name)
	for _, p := range searchPaths {
		filename := filepath.Join(p, basename)
		if exists(filename) {
			return filename
		}
	}
	return ""
}

// IsNumber checks to see if a value is a number
func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)

	return err == nil
}

// Interpolate (str, env)
// return input string with $vars interpolated from environment
func Interpolate(str string, env *object.Environment) string {
	// Match all strings preceded by ${
	// TODO: make this work on complex expressions, like
	// "${foo.bar.split().join(\",\")}"
	re := regexp.MustCompile(`(\\)?\$(\{)([a-zA-Z_0-9]{1,})(\})`)
	str = re.ReplaceAllStringFunc(str, func(m string) string {
		// If the string starts with a backslash, that's an escape, so we should
		// replace it with the remaining portion of the match. \${VAR} becomes
		// ${VAR}
		if string(m[0]) == "\\" {
			return m[1:]
		}

		varName := ""

		// If you type a variable wrong, forgetting the closing bracket, we
		// simply return it to you: eg "my ${variable"

		if m[len(m)-1] != '}' {
			return m
		}

		varName = m[2 : len(m)-1]

		v, ok := env.Get(varName)

		// If the variable is not found, we just dump an empty string
		if !ok {
			return ""
		}

		return v.Inspect()
	})

	return str
}
