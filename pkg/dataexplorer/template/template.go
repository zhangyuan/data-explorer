package template

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`\${([\w.]+)}`)

func SimpleCompile(template string, vars map[string]string) string {
	compiled := re.ReplaceAllStringFunc(template, func(s string) string {
		for key, value := range vars {
			if fmt.Sprintf("${%s}", key) == s {
				return fmt.Sprint(value)
			}
		}
		return s
	})
	return compiled
}
