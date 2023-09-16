package confparser

import (
	conftoken "confmanager/internal/app/conf_token"
	"fmt"
)

type Parser struct {
	t *conftoken.Tokenizer

	CurToken  conftoken.Token
	PeekToken conftoken.Token
	Depth     int

	Jobs []Node
}

type Node interface {
	String() string
	Type() string
}

type Mapping struct {
	Name  Key
	Value Node
}

func (m *Mapping) String() string {
	return fmt.Sprintf("{%s: %s}", m.Name.String(), m.Value.String())
}

func (m *Mapping) Type() string {
	return "Mapping"
}

type Key string

func (k Key) String() string {
	return string(k)
}

func (k Key) Type() string {
	return "Key"
}

type Sequence struct {
	Members []Node
}

func (s *Sequence) String() string {
	var seq string
	for _, mem := range s.Members {
		seq += fmt.Sprintf("\n\t%s", mem)
	}

	return "{" + seq + "\n}"
}

func (s *Sequence) Type() string {
	return "Sequence"
}

type Scalar string

func (sc Scalar) String() string {
	return string(sc)
}

func (sc Scalar) Type() string {
	return "Scalar"
}

type ParseError struct {
	errorString string
}

func (p ParseError) Error() string {
	return p.errorString
}

func GetParser(t *conftoken.Tokenizer) *Parser {
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
	if p.t.CurOffset < p.t.Length {
		p.PeekToken = p.t.ReadToken()
	}
}

func (p *Parser) Parse() ([]Node, error) {
	res := []Node{}

	for p.CurToken.Type != conftoken.EOF {
		if p.CurToken.Type == conftoken.NAME {
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

func (p *Parser) parseJob() (Node, error) {
	block, err := p.parseSequence()
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (p *Parser) parseSequence() (*Sequence, error) {
	seq := &Sequence{}
	curr_depth := p.Depth

	for p.CurToken.Type != conftoken.EOF && p.Depth >= curr_depth {
		line, err := p.parseLine()
		if err != nil {
			return nil, err
		}

		seq.Members = append(seq.Members, line)

		if p.CurToken.Type == conftoken.NEWL {
			p.readToken()
		}
	}

	return seq, nil
}

func (p *Parser) parseLine() (Node, error) {
	name := p.parseKey()
	if name.Type() == "Scalar" {
		return name, nil
	}

	key, _ := name.(Key)
	opt := &Mapping{Name: key}

	p.foldSpaces(0)
	p.readToken()
	p.foldSpaces(0)

	if p.CurToken.Type == conftoken.NEWL {
		p.readToken()

		value, err := p.parseSequence()
		if err != nil {
			return nil, err
		}

		opt.Value = value

		return opt, nil
	} else {
		value := p.parseValue()

		opt.Value = value

		return opt, nil
	}
}

func (p *Parser) parseKey() Node {
	if p.CurToken.Type == conftoken.DASH {
		p.readToken()
		p.foldSpaces(0)

		return p.parseValue()
	}

	var key string

	for p.CurToken.Type != conftoken.COLON {
		if p.CurToken.Type == conftoken.SPACE {
			key += p.foldSpaces(1)
		}
		key += p.CurToken.Literal
		p.readToken()
	}

	return Key(key)
}

func (p *Parser) parseValue() Scalar {
	var scalar string

	for p.CurToken.Type != conftoken.NEWL && p.CurToken.Type != conftoken.EOF {
		scalar += p.CurToken.Literal
		p.readToken()
	}

	return Scalar(scalar)
}

func unexpectedTokenError(exp, got conftoken.TokenType) error {
	return ParseError{
		errorString: fmt.Sprintf("expected token %s, got: %s", exp, got),
	}
}

func (p *Parser) foldSpaces(fold_to int) (res string) {
	for p.CurToken.Type == conftoken.SPACE && fold_to > 0 {
		fold_to -= 1
		res += " "
		p.readToken()
	}

	for p.CurToken.Type == conftoken.SPACE {
		p.readToken()
	}

	return res
}
