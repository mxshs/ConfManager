package confcodegen

import (
	confparser "confmanager/internal/app/conf_parser"
	"fmt"
	"os/exec"
	"os/user"
)

type CommandSet struct {
	Name     string
	After    []string
	Commands []*exec.Cmd
}

func GenCommands(jobs []confparser.Node) ([]*CommandSet, error) {
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

func GenCommandSet(job confparser.Node) (*CommandSet, error) {
	seq, ok := job.(*confparser.Sequence)
	if !ok {
		return nil, unexpectedOptionError(job.Type(), "Config")
	}

	res := &CommandSet{}

	for _, cmd := range seq.Members {
		opt, ok := cmd.(*confparser.Mapping)
		if !ok {
			return nil, unexpectedOptionError(cmd.Type(), seq.Type())
		}

		switch opt.Name.String() {
		case "name":
			res.Name = opt.Value.String()
		case "action":
			action_set, err := genActionSet(opt.Value)
			if err != nil {
				return nil, err
			}

			res.Commands = append(res.Commands, action_set...)
		//case "after":
		default:
			return nil, unexpectedOptionError(opt.Name.String(), seq.Type())
		}
	}

	return res, nil
}

func genActionSet(action_set confparser.Node) ([]*exec.Cmd, error) {
	res := []*exec.Cmd{}

	switch action_set.Type() {
	case "Sequence":
		dir := ""
		actions, _ := action_set.(*confparser.Sequence)
		for _, action := range actions.Members {
			if action.Type() == "Mapping" {
				a, _ := action.(*confparser.Mapping)

				// Currently directory is not passed to child sequences
				// I.e. each nested sequence has to specify its own directory if
				// it is not supposed to run commands from base (/usr/bin or /)
				if a.Value.Type() == "Sequence" {
					v, _ := a.Value.(*confparser.Sequence)
					sub_block, err := genActionSet(v)
					if err != nil {
						return nil, err
					}

					res = append(res, sub_block...)
				} else {
					if a.Name.String() == "dir" {
						dir = expandDirectory(a.Value.String())
					} else {
						cmd, err := genAction(a)
						if err != nil {
							return nil, err
						}

						if len(dir) > 0 {
							cmd.Dir = dir
						}

						res = append(res, cmd)
					}
				}
			} else {
				cmd, err := genAction(action)
				if err != nil {
					return nil, err
				}

				if len(dir) > 0 {
					cmd.Dir = dir
				}

				res = append(res, cmd)
			}
		}
	default:
		fmt.Println(action_set.Type())
	}

	return res, nil
}

func genAction(action confparser.Node) (*exec.Cmd, error) {
	switch action.(type) {
	case *confparser.Mapping:
		mapping, _ := action.(*confparser.Mapping)
		cmd := exec.Command("bash", "-c", mapping.Value.String())

		return cmd, nil
	case confparser.Scalar:
		cmd := exec.Command("bash", "-c", action.String())

		return cmd, nil
	default:
		return nil, unexpectedOptionError(action.Type(), "Action")
	}
}

func unexpectedOptionError(option, context string) error {
	return fmt.Errorf("unexpected option %s in %s context", option, context)
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
