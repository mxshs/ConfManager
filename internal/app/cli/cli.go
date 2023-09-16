package cli

import (
	"confmanager/internal/app/conf_autocomplete"
	confcodegen "confmanager/internal/app/conf_codegen"
	"confmanager/internal/app/conf_fetch"
	confparser "confmanager/internal/app/conf_parser"
	conftoken "confmanager/internal/app/conf_token"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  "To load completions",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		cobra.OnlyValidArgs,
	),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func Read() {
	app := &cobra.Command{
		Use: "confmanager",
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			conf_fetch.FetchRepo(name)

			t, err := conftoken.Start(name + "/test.conf")
			if err != nil {
				panic(err)
			}

			p := confparser.GetParser(t)

			parsed, err := p.Parse()
			if err != nil {
				panic(err)
			}

			jobs, err := confcodegen.GenCommands(parsed)
			if err != nil {
				panic(err)
			}

			for _, job := range jobs {
				for _, cmd := range job.Commands {
					cmd.Run()
				}
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string,
			toComplete string) ([]string, cobra.ShellCompDirective) {

			res, err := autoComplete(&confautocomplete.FileCache{})
			if err != nil {
				panic(err)
			}

			return res, cobra.ShellCompDirectiveDefault
		},
	}

	app.AddCommand(CompletionCmd)

	if err := app.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func autoComplete(cache confautocomplete.Cache) ([]string, error) {
	c, err := cache.Open("names.cache")
	if err != nil {
		return nil, err
	}

	return c.ReadCache()
}
