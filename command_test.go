package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type Test struct {
	Bool bool          `cmd:"Test (bool) var."`
	Int  int           `cmd:"Test (int) var."`
	I64  int64         `cmd:"Test (int64) var."`
	Uint uint          `cmd:"Test (uint) var."`
	U64  uint64        `cmd:"Test (uint64) var."`
	F64  float64       `cmd:"Test (float64) var."`
	Str  string        `cmd:"Test (string) var."`
	Dur  time.Duration `cmd:"Test (duration) var."`

	unexported string //unexported field
	Untagged   string `cmd:"-"`
}

func (t *Test) Try() {

}

func TestStruct(t *testing.T) {
	test := new(Test)
	ok(t, checkIfStruct(test))
}

func TestNewCommand(t *testing.T) {
	test := new(Test)
	cmd, err := NewCommand(test)
	ok(t, err)
	equals(t, "Test", cmd.Name)
	equals(t, test, cmd.Reference)
}

func TestParseFlags(t *testing.T) {
	test := &Test{false, 0, 0, 0, 0, 0.0, "", 1, "", ""}
	cmd, err := NewCommand(test)
	ok(t, err)

	flags := map[string]string{
		"bool": "true",
		"int":  "-65536",
		"i64":  "-10000000000",
		"uint": "65536",
		"u64":  "10000000000",
		"f64":  "10.5",
		"str":  "Commander",
		"dur":  "5s",
	}

	for k, v := range flags {
		os.Args = append(os.Args, "-"+k+"="+v)
	}

	//fmt.Println(os.Args)
	ok(t, cmd.ParseFlags(Default))
	flag.Parse()

	equals(t, true, test.Bool)
	equals(t, int(-65536), test.Int)
	equals(t, int64(-10000000000), test.I64)
	equals(t, uint(65536), test.Uint)
	equals(t, uint64(10000000000), test.U64)
	equals(t, 10.5, test.F64)
	equals(t, "Commander", test.Str)
	equals(t, 5*time.Second, test.Dur)

}
