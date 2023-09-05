package confparser

import "fmt"

type Job struct {
    Name string
    Opts []Opt
}

type Opt struct {
    Name string
    Content []string
}

type ParseError struct {
    errorString string
}

func (p ParseError) Error() string {
    return p.errorString
}

func ParseEval(t *Tokenizer) ([]Job, error) {
    res := []Job{}

    for t.CurToken.Type != EOF {
        if t.CurToken.Type == NAME {
            job, err := parseJob(t)
            if err != nil {
                return nil, err
            }

            res = append(res, *job)
        } else {
            t.ReadToken()
        }
    }

    return res, nil
}

func parseJob(t *Tokenizer) (*Job, error) {
    t.ReadToken()
    t.ReadToken()

    name, _ := parseIdent(t)
    job := &Job{Name: name}

    t.ReadToken()
    tok := t.CurToken

    for t.CurToken.Type != EOF {
        if tok.Type != NAME {
            unexpectedTokenError(NAME, t.CurToken.Type)
        }
       
        opt := Opt{Name: tok.Literal}

        tok = *t.ReadToken()
        if tok.Type != COLON {
            unexpectedTokenError(COLON, t.CurToken.Type)
        }

        tok = *t.ReadToken()
        if tok.Type == NAME {
            ident, err := parseIdent(t)
            if err != nil {
                return nil, err
            }

            opt.Content = append(opt.Content, ident)
        } else if tok.Type == NEWL {
            t.ReadToken()
            opts, err := parseOpts(t)
            if err != nil {
                return nil, err
            }

            opt.Content = opts
        }

        job.Opts = append(job.Opts, opt)

        if t.CurToken.Type != EOF {
            tok = *t.ReadToken()
        }
    }

    return job, nil
}

func parseIdent(t *Tokenizer) (string, error) {
    res := ""

    for t.CurToken.Type != NEWL && t.CurToken.Type != EOF {
        res += t.CurToken.Literal
        t.ReadToken()
    }

    return res, nil
}

func parseOpts(t *Tokenizer) ([]string, error) {
    if t.CurToken.Type != DASH {
        return nil, unexpectedTokenError(DASH, t.CurToken.Type)
    }

    res := []string{}

    for t.CurToken.Type == DASH {
        cmd := ""
        
        t.ReadToken()
        
        for t.CurToken.Type != NEWL && t.CurToken.Type != EOF {
            cmd += t.CurToken.Literal
            t.ReadToken()
        }

        res = append(res, cmd)
        
        if t.CurToken.Type != EOF {
            t.ReadToken()
        }
    }

    return res, nil
}

func unexpectedTokenError(exp, got tokenType) error {
    return ParseError{
        errorString: fmt.Sprintf("expected token %s, got: %s", exp, got),
    }
}

