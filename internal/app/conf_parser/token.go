package confparser

import (
	"fmt"
	"os"
)

type tokenType string

type Token struct {
    Type tokenType
    Literal string
}

type Tokenizer struct {
    CurToken Token
    curOffset int
    lineOffset int
    Depth int
    conf []byte
}

const (
    COLON = ":"
    DASH = "-"
    SLASH = "/"
    NEWL = "NEWL"
    SPACE = " "
    DOT = "."
    TILDE = "~"
    AMP = "&"
    PLUS = "+"
    LBR = "("
    RBR = ")"
    DQ = "\""
    Q = "'"
    RSH = ">"
    LSH = "<"
    BCK_SLASH = "\\"

    IDENT = "IDENT"
    NUM = "NUM"

    NAME = "NAME_OPT"

    EOF = ""
)

func Start(path string) (*Tokenizer, error) {
    conf, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    t := &Tokenizer{}
    t.conf = append(conf, 0)

    return t, nil
}

func (t *Tokenizer) ReadToken() Token {
    curr := t.conf[t.curOffset]
    
    switch curr {
    case ' ':
        if t.lineOffset == 0 {
            counter := 0
            for t.conf[t.curOffset] == ' ' {
                counter += 1
                t.curOffset += 1
                t.lineOffset += 1
            }

            t.Depth = counter

            return t.ReadToken()
        } else {
            t.CurToken = Token{
                Type: SPACE,
                Literal: SPACE,
            }
        }
    case ':':
        t.CurToken = Token{
            Type: COLON,
            Literal: COLON,
        }
    case '-':
        t.CurToken = Token{
            Type: DASH,
            Literal: DASH,
        }
    case '~':
        t.CurToken = Token{
            Type: TILDE,
            Literal: TILDE,
        }
    case '&':
        t.CurToken = Token{
            Type: AMP,
            Literal: AMP,
        }
    case '+':
        t.CurToken = Token{
            Type: PLUS,
            Literal: PLUS,
        }
    case '/':
        t.CurToken = Token{
            Type: SLASH,
            Literal: SLASH,
        }
    case '\\':
        t.CurToken = Token{
            Type: BCK_SLASH,
            Literal: BCK_SLASH,
        }
    case '(':
        t.CurToken = Token{
            Type: LBR,
            Literal: LBR,
        }
    case ')':
        t.CurToken = Token{
            Type: RBR,
            Literal: RBR,
        }
    case '>':
        t.CurToken = Token{
            Type: RSH,
            Literal: RSH,
        }
    case '<':
        t.CurToken = Token{
            Type: LSH,
            Literal: LSH,
        }
    case '\'':
        t.CurToken = Token{
            Type: Q,
            Literal: Q,
        }
    case '"':
        t.CurToken = Token{
            Type: DQ,
            Literal: DQ,
        }
    case '\n':
        t.CurToken = Token{
            Type: NEWL,
            Literal: NEWL,
        }

        t.lineOffset = -1
        t.Depth = 0
    case '.':
        t.CurToken = Token{
            Type: DOT,
            Literal: DOT,
        }
    case 0:
        t.CurToken = Token{
            Type: EOF,
            Literal: EOF,
        }
    default:
        if isLetter(curr) {
            t.CurToken = Token{
                Type: NAME, 
                Literal: string(t.readWord()),
            }
            return t.CurToken
        } else if isDigit(curr) {
            t.CurToken = Token{
                Type: NUM,
                Literal: string(t.readNumber()),
            }
            
            return t.CurToken
        } else {
            panic(fmt.Sprintf("unexpected char %s", string(curr)))
        }
    }
    
    t.curOffset += 1
    t.lineOffset += 1
    return t.CurToken
}

func (t *Tokenizer) readWord() []byte {
    start := t.curOffset

    for isLetter(t.conf[t.curOffset]) {
        t.curOffset += 1
        t.lineOffset += 1
    }

    return t.conf[start:t.curOffset]
}

func (t *Tokenizer) readNumber() []byte {
    start := t.curOffset

    for isDigit(t.conf[t.curOffset]) {
        t.curOffset += 1
        t.lineOffset += 1
    }

    return t.conf[start:t.curOffset]
}

func isLetter(ch byte) bool {
    return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

