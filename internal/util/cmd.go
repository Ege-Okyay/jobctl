package util

import (
	"fmt"
	"sort"

	"github.com/Ege-Okyay/jobctl/internal/types"
)

// HandleUnknownCommand provides suggestions for mistyped commands.
func HandleUnknownCommand(commands map[string]types.Command, cmdName string) {
	closestCommands := findClosestCommands(commands, cmdName, 2)

	fmt.Printf("'%s' is not a jobctl command. See 'help'.\n", cmdName)

	if len(closestCommands) > 0 {
		fmt.Println("\nThe most similar commands are:")
		for _, command := range closestCommands {
			fmt.Printf("\t%s\n", command)
		}
	}
}

// IsHelpFlag checks if an argument is a help flag.
func IsHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

// findClosestCommands finds the most similar command names to a given input string
// using the Levenshtein distance algorithm.
func findClosestCommands(commands map[string]types.Command, cmdName string, maxResults int) []string {
	str1 := []rune(cmdName)
	var commandDistances []types.CommandDistance

	for name := range commands {
		str2 := []rune(name)
		score := Levenshtein(str1, str2)

		commandDistances = append(commandDistances, types.CommandDistance{Name: name, Score: score})
	}

	sort.Slice(commandDistances, func(i, j int) bool {
		return commandDistances[i].Score < commandDistances[j].Score
	})

	var closestCommands []string
	for i := 0; i < min(maxResults, len(commandDistances)); i++ {
		closestCommands = append(closestCommands, commandDistances[i].Name)
	}

	return closestCommands
}
