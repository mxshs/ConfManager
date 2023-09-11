package main

import (
	//"confmanager/internal/app/cli"
	confcodegen "confmanager/internal/app/conf_codegen"
	confparser "confmanager/internal/app/conf_parser"
	conftoken "confmanager/internal/app/conf_token"
	"fmt"
)


func main() {
    t, err := conftoken.Start("test.asd")
    if err != nil {
        panic(err)
    }

    p := confparser.GetParser(t)
    conf, err := p.Parse()
    if err != nil {
        panic(err)
    }

    fmt.Println(conf[0].String())
    c, err := confcodegen.GenCommands(conf)
    if err != nil {
        panic(err)
    }

   for _, cmd := range c[0].Commands {
        fmt.Println(cmd)
    }
   //     cmd.Run()
        //fmt.Printf("%s", cmd.Name)
        //fmt.Printf("%+v", cmd.Commands)
    //}
    //tok := t.ReadToken()
    //fmt.Printf("Token: %s\n", tok.Literal)

    //for t.CurToken.Type != confparser.EOF {
    //    tok = t.ReadToken()
    //    fmt.Printf("Token: %s\n", tok.Literal)
    //}

//    cli.Read()
}

