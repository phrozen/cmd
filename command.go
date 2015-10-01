package cmd

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Package variables
var {
	commands map[string]*Command
	Default = Options{false}

}

func init() {
	commands = make(map[string]*Command)
}

type Options struct {
	Namespace bool
}

type Command struct {
	Name      string
	Reference interface{}
}

type Commander struct {
	Commands []Command
	Config   Options
}

// Uses reflection to check if interface 'val' is of type
// reflect.Struct returns an error if it isn't and nil otherwise.
func checkIfStruct(val interface{}) error {
	cmdType := reflect.TypeOf(val).Elem()
	if cmdType.Kind() != reflect.Struct {
		return fmt.Errorf("Type %v is not a struct.", cmdType.Name())
	}
	return nil
}

// Returns a new command pointer based on struct 'cmd'
// Checks if interface is of Kind reflect.Struct
// returns error if it isn't and nil otherwise.
func NewCommand(cmd interface{}) (*Command, error) {
	err := checkIfStruct(cmd)
	if err != nil {
		return nil, err
	}
	cmdType := reflect.TypeOf(cmd).Elem()
	return &Command{cmdType.Name(), cmd}, nil
}

// Goes through all the exported fields of the struct 'cmd',
// then tries to bind each exported field to a flag of the
// supported types by the flag package, returns an error if
// it finds an exported Unsupported type.
// To ignore an exported field use `cmd:"-"` as the tag.
func (cmd *Command) ParseFlags(opt Options) error {
	cmdType := reflect.TypeOf(cmd.Reference).Elem()
	for i := 0; i < cmdType.NumField(); i++ {
		field := cmdType.Field(i)
		usage := field.Tag.Get("cmd")
		// Don't try unexported and untagged fields
		if field.PkgPath == "" && usage != "-" {
			cmdValue := reflect.ValueOf(cmd.Reference).Elem()
			// Addresable interface
			value := cmdValue.FieldByName(field.Name).Addr().Interface()

			// Namespacing for field collision
			flagName := ""
			if opt.Namespace {
				flagName = strings.ToLower(cmdType.Name() + "." + field.Name)
			} else {
				flagName = strings.ToLower(field.Name)
			}
			// Only types supported by the flag package
			switch value := value.(type) {
			case *bool:
				flag.BoolVar(value, flagName, *value, usage)
			case *int:
				flag.IntVar(value, flagName, *value, usage)
			case *int64:
				flag.Int64Var(value, flagName, *value, usage)
			case *uint:
				flag.UintVar(value, flagName, *value, usage)
			case *uint64:
				flag.Uint64Var(value, flagName, *value, usage)
			case *float64:
				flag.Float64Var(value, flagName, *value, usage)
			case *string:
				flag.StringVar(value, flagName, *value, usage)
			case *time.Duration:
				flag.DurationVar(value, flagName, *value, usage)
			default:
				return fmt.Errorf("Unsupported type: %s of type %v cannot be parsed as flag.", field.Name, field.Type.Name())
			}
		}
	}
	return nil
}

// Executes method on struct 'cmd' with the given 'name',
// this match is case insensitive. Return an error if the
// methos does not exist or is unexported.
func (cmd *Command) Exec(name string) error {
	cmdType := reflect.TypeOf(cmd.Reference).Elem()
	for i := 0; i < cmdType.NumMethod(); i++ {
		method := cmdType.Method(i)
		if method.PkgPath == "" {
			if strings.ToLower(method.Name) == strings.ToLower(name) {
				method.Func.Call([]reflect.Value{reflect.ValueOf(cmd.Reference).Elem()})
				return nil
			}
		}
	}
	return fmt.Errorf("Method <%s> not found.", name)
}

// Reflects and parses any number of 'struct' values to be called
// via CLI, throws an error if any of them cannot be parsed or
// if the called method cannot be executed.
func Commanderize(opt Options, values ...interface{}) error {
	commands := make([]*Command, 0)
	for _, val := range values {
		// Check if interface is of Kind reflect.Struct
		cmd, err := NewCommand(val)
		if err != nil {
			return err
		}
		// Create flags based on StructFields
		err = cmd.ParseFlags(opt)
		if err != nil {
			return err
		}
		// Append command to the command list.
		commands = append(commands, cmd)
	}

	// Parse cli arguments into flags
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("Usage: <struct>:<method> (No command given.)")
	}

	// Parse the Arg(0) which is the actual command
	arg := strings.Split(flag.Arg(0), ":")
	if len(arg) != 2 {
		return fmt.Errorf("Usage: <struct>:<method> (Got: %s)", flag.Arg(0))
	}

	// Search for the method and execute it
	cmdName, cmdMethod := arg[0], arg[1]
	for _, cmd := range commands {
		if strings.ToLower(cmd.Name) == strings.ToLower(cmdName) {
			return cmd.Exec(cmdMethod)
		} else {
			return fmt.Errorf("Command <%s> not found.", cmdName)
		}
	}

	return nil
}
