package util

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/Ege-Okyay/jobctl/internal/types"
)

func PrintHelp(commands map[string]types.Command) {
	fmt.Println("Usage:")
	fmt.Println("    <command> [flags]")

	fmt.Println("Available Commands:")

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, name := range names {
		cmd := commands[name]
		fmt.Fprintf(w, "    %s\t%s\n", cmd.Name, cmd.Description)
	}
	w.Flush()

	fmt.Println("\nUse \"<command> --help\" for more information about a command")
	fmt.Print("Use \"clear\" to clear the terminal\n\n")
}

func PrintCommandHelp(cmd types.Command) {
	fmt.Printf("%s\n\n", cmd.Description)

	fmt.Println("Usage:")
	fmt.Printf("    %s\n\n", cmd.Usage)

	if len(cmd.Flags) > 0 {
		fmt.Println("Flags:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, f := range cmd.Flags {
			req := ""
			if f.Required {
				req = " (required)"
			}
			fmt.Fprintf(w, "    %s\t%s%s\n", f.Name, f.Description, req)
		}
		w.Flush()
		fmt.Println()
	}
}

func PrintBanner() {
	fmt.Print(`
     ██╗ ██████╗ ██████╗  ██████╗████████╗██╗     
     ██║██╔═══██╗██╔══██╗██╔════╝╚══██╔══╝██║     
     ██║██║   ██║██████╔╝██║        ██║   ██║     
██   ██║██║   ██║██╔══██╗██║        ██║   ██║     
╚█████╔╝╚██████╔╝██████╔╝╚██████╗   ██║   ███████╗
 ╚════╝  ╚═════╝ ╚═════╝  ╚═════╝   ╚═╝   ╚══════╝` + "\n\n")
}

func SuccessMessage(msg string) {
	fmt.Printf("\xE2\x9C\x94 %s\n", msg)
}

func ErrorMessage(msg string) {
	fmt.Printf("\xE2\x9C\x97 %s\n", msg)
}
