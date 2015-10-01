
# Commander

Commander (**cmd**) is a command/task manager that builds CLI tasks and flags out of your Go code using reflection on your data structs and functions.

### Usage

To create tasks and flags to run those tasks simply define your usual struct types amd functions on them.

```
type Test struct {
  Name   string
}

func (t *Test) Hello() {
  fmt.Printf("Hello %s!", t.Name)
}

```

Then simply call Commanderize on your struct. You can set default values by initializing your struct.

```
err:= Commanderize(cmd.Default, &Test{"World"})
if err!= nil {
  // Do something with err
}
```

What ```Commanderize()``` does is, it reads all the exported fields on your struct type (only those supported by the ```flag``` pkg) and binds them to flags with their name in lower case which defaults to their initial value and uses the ```cmd``` tag as the usage string, calls internally ```flag.Parse()``` so that your fields get the parsed values if any, then parses the first argument ```Arg(0)``` which should be in the format ```<struct>:<method>``` looks for all the exported methods and executes the right one. It returns a detailed error (which you can print) if there was a problem in any of the steps. So after that you can compile your binary and run:

```
# binary test:hello
Hello World!

# binary -name=John test:hello
Hello John!
```

You can code multiple structs and pass them to ```Commanderize()``` in one call, the first parameter is an Options struct and a Default is provided (use it for now), the second one is a variadic argument which takes any number of structs so your CLI can better separate one task from another. Default usage is done via the ```cmd``` tag and you can untag fields so they aren´t checked. Any unexported member (field or function) is not checked by Commander.

```
type MyStruct struct {
  Check       bool      `cmd:"Usage string for check."`
  Price       float64   `cmd:"Usage string for price."`
  unexported  string    // Unexported fields do not generate flags
  Untagged   *SomeType  `cmd:"-"`  // Untagged fields do not generate flags
}
```

### Todo

* Better documentation
* Support variadic argument functions with extra Args[1:]
* Colorize CLI
