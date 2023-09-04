package cli

import (
	"bufio"
	"confmanager/internal/app/conf_fetch"
	"io"
	"os"
	"strings"
)

type EofError interface {
    getWord() []byte
    Error() string
}

type EOF struct {
    word []byte
}

func (e EOF) getWord() []byte {
    return e.word
}

func (e EOF) Error() string {
    return "reached the end of user input"
}

func Read() {
    reader := bufio.NewReader(os.Stdin)
    var opts []string
    var args []string

    for {
        currByte, _ := reader.ReadByte()
        if currByte == '-' {
            word, err := readWord(reader)
            if err != nil {
                if err, ok := err.(EofError); ok {
                    opts = append(opts, "-" + string(err.getWord()))
                    break
                } else {
                    panic(err)
                }
            }

            opts = append(opts, "-" + string(word))
        } else {
            word, err := readWord(reader)
            if err != nil {
                if err, ok := err.(EofError); ok {
                    args = append(args, string(currByte) + string(err.getWord()))
                    break
                } else {
                    panic(err)
                }
            }

            args = append(args, string(currByte) + string(word))
        }
    }

    io.WriteString(os.Stdout, "opts: ")
    for _, opt := range opts {
        io.WriteString(os.Stdout, string(opt))
    }
    io.WriteString(os.Stdout, "args: ")
    for _, arg := range args {
        io.WriteString(os.Stdout, string(arg))
    }
}

func readWord(reader *bufio.Reader) ([]byte, error) {
    buf := []byte{}

    currByte, _ := reader.ReadByte()
    for currByte != ' ' && currByte != '\n' {
        if currByte == '\t' {
            res, err := conf_fetch.GetNamesForAutocomp()
            if err != nil {
                io.WriteString(os.Stdout, err.Error())
            }

            io.WriteString(os.Stdout, strings.Join(res, "\t"))
        }
        buf = append(buf, currByte)
        currByte, _ = reader.ReadByte()
    }

    if currByte == '\n' {
        return nil, EOF{word: buf}
    }

    return buf, nil
}

