package main

import (
	//"confmanager/internal/app/cli"
	confparser "confmanager/internal/app/conf_parser"
	"fmt"
	//"confmanager/internal/app/conf_fetch"
)


func main() {
    t, err := confparser.Start("test.file")
    if err != nil {
        panic(err)
    }

    t.ReadToken()
    p, err := confparser.ParseEval(t)
    if err != nil {
        panic(err)
    }

    for _, cmd := range p {
        fmt.Printf("%+v", cmd)
    }
    //tok := t.ReadToken()
    //fmt.Printf("Token: %s\n", tok.Literal)

    //for t.CurToken.Type != confparser.EOF {
    //    tok = t.ReadToken()
    //    fmt.Printf("Token: %s\n", tok.Literal)
    //}

    //cli.Read()
}
