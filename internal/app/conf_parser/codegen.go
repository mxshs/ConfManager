package confparser

import (
	"fmt"
	"os/exec"
	"os/user"
	//"strings"
)

type CommandSet struct {
    Name string
    After []string
    Commands []*exec.Cmd
}

func GenCommands(jobs []ConfValue) ([]*CommandSet, error) {
    res := []*CommandSet{}

    for _, job := range jobs {
        command_set, err := GenCommandSet(job)
        if err != nil {
            return nil, err
        }

        res = append(res, command_set)
    }

    return res, nil
}

func GenCommandSet(job ConfValue) (*CommandSet, error) {
    block, ok := job.(*Block)
    if !ok {
        return nil, unexpectedOptionError("Block", "Job") 
    }

    res := &CommandSet{}

    for _, cmd := range block.Opts {
        opt, ok := cmd.(*Opt)
        if !ok {
            return nil, unexpectedOptionError(cmd.Type(), block.Type())
        }

        switch opt.Name.String() {
        case "name":
            res.Name = opt.String()
        case "action":
            action_set, err := genAction(opt.Value)
            if err != nil {
                return nil, err
            }

            res.Commands = append(res.Commands, action_set...)
        //case "after":
        default:
            return nil, unexpectedOptionError(opt.Name.String(), block.Type())
        }
    }

    return res, nil
}

func genAction(action_set ConfValue) ([]*exec.Cmd, error) {
    res := []*exec.Cmd{}
    switch action_set.Type() {
    case "Block":
        dir := ""
        actions, _ := action_set.(*Block)
        for _, action := range actions.Opts {
            a, _ := action.(*Opt)
            if a.Value.Type() == "Block" {
                v, _ := a.Value.(*Block)
                sub_block, err := genAction(v)
                if err != nil {
                    return nil, err
                }

                res = append(res, sub_block...)
            } else {
                if a.Name.String() == "dir" {
                    dir = expandDirectory(a.Value.String())
                } else {
                    cmd := exec.Command("bash", "-c", a.Value.String())
                    if len(dir) > 0 {
                        cmd.Dir = dir
                    }

                    res = append(res, cmd)
                }
            }
        }
    default:
        fmt.Println(action_set.Type())
    }
    return res, nil
}

func unexpectedOptionError(option, context string) error {
    return fmt.Errorf("unexpected option %s inside %s context", option, context)
}

func expandDirectory(path string) string {
    expanded := ""
    for _, ch := range path {
        if ch == '~' {
            usr, _ := user.Current()
            expanded += usr.HomeDir
        } else {
            expanded += string(ch)
        }
    }

    return expanded
}

