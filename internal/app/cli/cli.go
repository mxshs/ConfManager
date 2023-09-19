package cli

import (
	"confmanager/internal/app/conf_autocomplete"
	confcodegen "confmanager/internal/app/conf_codegen"
	"confmanager/internal/app/conf_fetch"
	confparser "confmanager/internal/app/conf_parser"
	conftoken "confmanager/internal/app/conf_token"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var completion = &cobra.Command{
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

var configure = &cobra.Command{
    Use: "configure",
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        confName, confFileName := args[0], args[1]

        conf_fetch.W.FetchConf(confName)

        t, err := conftoken.Start(confName + "/" + confFileName)
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

        wg := &sync.WaitGroup{}
        bars := mpb.New(mpb.WithWaitGroup(wg))

        for _, job := range jobs {
            bar := bars.AddBar(
                int64(len(job.Commands)),
                mpb.PrependDecorators(
                    decor.Name(
                        job.Name,
                        decor.WC{
                            W: len(job.Name) + 1,
                            C: decor.DidentRight,
                        },
                    ),
                    decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
                ),
                mpb.AppendDecorators(
                    decor.Percentage(),
                ),
            )

            for _, cmd := range job.Commands {
                bar.Increment()
                cmd.Run()
            }
        }

        bars.Wait()
    },
    ValidArgsFunction: func(cmd *cobra.Command, args []string,
        toComplete string) ([]string, cobra.ShellCompDirective) {

        if len(args) == 0 {
            res, err := confautocomplete.GetRepoNames()
            if err != nil {
                panic(err)
                //return nil, cobra.ShellCompDirectiveError
            }

            return res, cobra.ShellCompDirectiveDefault
        } else if len(args) == 1 {
            res, err := confautocomplete.GetFileNames(args[0])
            if err != nil {
                panic(err)
            }

            return res, cobra.ShellCompDirectiveDefault
        }

        return nil, cobra.ShellCompDirectiveDefault
    },
}

func Read() {
	app := &cobra.Command{
		Use: "confmanager",
    }

	app.AddCommand(completion)
    app.AddCommand(configure)

	if err := app.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

