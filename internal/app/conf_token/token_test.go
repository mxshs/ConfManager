package conftoken

import (
    "testing"
)

func TestTokens(t *testing.T) {
    input := `name: test
    action:
        cmd: "test" \
        cmd 2: test 123
    `

    tests := []struct {
        expectedType TokenType
        expectedLiteral string
    }{
        {NAME, "name"},
        {COLON, COLON},
        {SPACE, SPACE},
        {NAME, "test"},
        {NEWL, NEWL},
        {NAME, "action"},
        {COLON, COLON},
        {NEWL, NEWL},
        {NAME, "cmd"},
        {COLON, COLON},
        {SPACE, SPACE},
        {DQ, DQ},
        {NAME, "test"},
        {DQ, DQ},
        {SPACE, SPACE},
        {BCK_SLASH, BCK_SLASH},
        {NEWL, NEWL},
        {NAME, "cmd"},
        {SPACE, SPACE},
        {NUM, "2"},
        {COLON, COLON},
        {SPACE, SPACE},
        {NAME, "test"},
        {SPACE, SPACE},
        {NUM, "123"},
    }

    tokenizer := Tokenizer{conf: []byte(input)}

    for i, tt := range tests {
        tok := tokenizer.ReadToken()

        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - wrong token type: expected %q, got %q",
                i, tt.expectedType, tok.Type)
        }

        if tok.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - wrong token literal: expected %q, got %q",
                i, tt.expectedLiteral, tok.Literal)
        }
    }
}

