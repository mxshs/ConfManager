package confparser

import "fmt"

type Parser struct {
    t *Tokenizer

    CurToken Token
    PeekToken Token
    Depth int

    Jobs []ConfValue
}

type ConfValue interface {
    String() string
    Type() string
}

type Opt struct {
    Name StandaloneValue 
    Value ConfValue 
}

func (o *Opt) String() string {
    return fmt.Sprintf("{%s: %s}", o.Name.String(), o.Value.String())
}

func (o *Opt) Type() string {
    return "Opt"
}

type Block struct {
    Opts []ConfValue
}

func (b *Block) String() string {
    var block string
    for _, opt := range b.Opts {
        block += fmt.Sprintf("\n\t%s", opt)
    }

    return "{" + block + "\n}"
}

func (b *Block) Type() string {
    return "Block"
}

type StandaloneValue struct {
    Value string
}

func (sv *StandaloneValue) String() string {
    return sv.Value
}

func (sv *StandaloneValue) Type() string {
    return "Value"
}

type ParseError struct {
    errorString string
}

func (p ParseError) Error() string {
    return p.errorString
}

func GetParser(t *Tokenizer) *Parser {
    p := &Parser{
        t: t,
    }

    p.readToken()
    p.readToken()

    return p
}

func (p *Parser) readToken() {
    p.Depth = p.t.Depth
    p.CurToken = p.PeekToken
    if p.t.curOffset < len(p.t.conf) {
        p.PeekToken = p.t.ReadToken()
    }
}

func (p *Parser) Parse() ([]ConfValue, error) {
    res := []ConfValue{}

    for p.CurToken.Type != EOF {
        if p.CurToken.Type == NAME {
            job, err := p.parseJob()
            if err != nil {
                return nil, err
            }

            res = append(res, job)
        } else {
            p.readToken()
        }
    }

    return res, nil
}

func (p *Parser) parseJob() (ConfValue, error) {
    block, err := p.parseBlock()
    if err != nil {
        return nil, err
    }

    return block, nil
}

func (p *Parser) parseBlock() (*Block, error) {
    block := &Block{}
    curr_depth := p.Depth

    for p.CurToken.Type != EOF && p.Depth >= curr_depth {
        opt, err := p.parseLine()
        if err != nil {
            return nil, err
        }

        block.Opts = append(block.Opts, opt)

        if p.CurToken.Type == NEWL {
            p.readToken()
        }
    }

    return block, nil
}

func (p *Parser) parseLine() (*Opt, error) {
    opt := &Opt{Name: p.parseKey()}

    p.foldSpaces(0)
    p.readToken()
    p.foldSpaces(0)

    if p.CurToken.Type == NEWL {
        p.readToken()

        value, err := p.parseBlock()
        if err != nil {
            return nil, err
        }

        opt.Value = value

        return opt, nil
    } else {
        value, err := p.parseValue()
        if err != nil {
            return nil, err
        }

        opt.Value = value

        return opt, nil
    }
}

func (p *Parser) parseKey() StandaloneValue {
    var key string

    if p.CurToken.Type == DASH {
        p.readToken()
        p.foldSpaces(0)
    }

    for p.CurToken.Type != COLON && p.CurToken.Type != NEWL && p.CurToken.Type != EOF {
        if p.CurToken.Type == SPACE {
            key += p.foldSpaces(1)
        }
        key += p.CurToken.Literal
        p.readToken()
    }

    return StandaloneValue{
        Value: key,
    }
}

func (p *Parser) parseValue() (ConfValue, error) {
    var key string

    for p.CurToken.Type != NEWL && p.CurToken.Type != EOF {
        key += p.CurToken.Literal
        p.readToken()
    }

    return &StandaloneValue{
        Value: key,
    }, nil
}

func unexpectedTokenError(exp, got tokenType) error {
    return ParseError{
        errorString: fmt.Sprintf("expected token %s, got: %s", exp, got),
    }
}

func (p *Parser) foldSpaces(fold_to int) (res string) {
    counter := 0
    for p.CurToken.Type == SPACE {
        counter += 1
        res += " "
        p.readToken()
    }

    if counter > fold_to {
        return res[:fold_to]
    }

    return res
}

