package cli

import (
	"confmanager/internal/app/conf_fetch"
	confparser "confmanager/internal/app/conf_parser"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Read() {
    app := &cli.App{
        EnableBashCompletion: true, 
        Commands: []*cli.Command{
            {
                Name: "confmanager",
                Action: func(ctx *cli.Context) error {
                    name := ctx.Args().First()

                    conf_fetch.FetchRepo(name)

                    t, err := confparser.Start(name + "/test.conf")
                    if err != nil {
                        return err
                    }

                    p := confparser.GetParser(t)

                    parsed, err := p.Parse()
                    if err != nil {
                        return err
                    }

                    jobs, err := confparser.GenCommands(parsed)
                    if err != nil {
                        return err
                    }

                    for _, job := range jobs {
                        for _, cmd := range job.Commands {
                            cmd.Run()
                        }
                    }
                    return nil
                },
                BashComplete: func(ctx *cli.Context) {
                    res, err := conf_fetch.GetNamesForAutocomp()
                    if err != nil {
                        fmt.Println(err)
                        return
                    }

                    for _, name := range res {
                        fmt.Println(name)
                    }
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        fmt.Println(err)
        return
    }
}

