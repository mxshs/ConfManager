package cli

import (
	"confmanager/internal/app/conf_fetch"
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
                    fmt.Println(ctx.Args().First())
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

