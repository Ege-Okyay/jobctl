package util

import (
	"fmt"
	"strings"

	"github.com/Ege-Okyay/jobctl/internal/types"
)

// ParseFlags parses command-line arguments into a map of flags and their values.
// It supports flags with and without values, and values that are quoted.
func ParseFlags(args []string) map[string]string {
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		tok := args[i]

		if !strings.HasPrefix(tok, "--") {
			continue
		}

		if i+1 >= len(args) {
			flags[tok] = ""
			continue
		}

		val := args[i+1]
		i++

		// Handle quoted values that may contain spaces.
		if len(val) > 0 && (val[0] == '"' || val[0] == '\'') {
			quote := val[0]

			if val[len(val)-1] != quote || len(val) == 1 {
				parts := []string{val}

				for i+1 < len(args) {
					next := args[i+1]
					parts = append(parts, next)
					i++

					if len(next) > 0 && next[len(next)-1] == quote {
						break
					}
				}

				val = strings.Join(parts, " ")
			}

			if len(val) >= 2 && val[0] == quote && val[len(val)-1] == quote {
				val = val[1 : len(val)-1]
			}
		}

		flags[tok] = val
	}

	return flags
}

// ValidateFlags checks if the provided flags are valid for a given command.
func ValidateFlags(provided map[string]string, defs []types.Flag) error {
	allowed := make(map[string]bool, len(defs))

	for _, f := range defs {
		allowed[f.Name] = true
	}

	for name := range provided {
		if !allowed[name] {
			return fmt.Errorf("unknown flag: %s", name)
		}
	}

	for _, f := range defs {
		if f.Required {
			if _, ok := provided[f.Name]; !ok {
				return fmt.Errorf("missing required flag: %s", f.Name)
			}
		}
	}

	return nil
}
