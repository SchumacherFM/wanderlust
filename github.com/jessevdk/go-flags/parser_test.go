package flags

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type defaultOptions struct {
	Int        int `long:"i"`
	IntDefault int `long:"id" default:"1"`

	String            string `long:"str"`
	StringDefault     string `long:"strd" default:"abc"`
	StringNotUnquoted string `long:"strnot" unquote:"false"`

	Time        time.Duration `long:"t"`
	TimeDefault time.Duration `long:"td" default:"1m"`

	Map        map[string]int `long:"m"`
	MapDefault map[string]int `long:"md" default:"a:1"`

	Slice        []int `long:"s"`
	SliceDefault []int `long:"sd" default:"1" default:"2"`
}

func TestDefaults(t *testing.T) {
	var tests = []struct {
		msg      string
		args     []string
		expected defaultOptions
	}{
		{
			msg:  "no arguments, expecting default values",
			args: []string{},
			expected: defaultOptions{
				Int:        0,
				IntDefault: 1,

				String:        "",
				StringDefault: "abc",

				Time:        0,
				TimeDefault: time.Minute,

				Map:        map[string]int{},
				MapDefault: map[string]int{"a": 1},

				Slice:        []int{},
				SliceDefault: []int{1, 2},
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--i=3", "--id=3", "--str=def", "--strd=def", "--t=3ms", "--td=3ms", "--m=c:3", "--md=c:3", "--s=3", "--sd=3"},
			expected: defaultOptions{
				Int:        3,
				IntDefault: 3,

				String:        "def",
				StringDefault: "def",

				Time:        3 * time.Millisecond,
				TimeDefault: 3 * time.Millisecond,

				Map:        map[string]int{"c": 3},
				MapDefault: map[string]int{"c": 3},

				Slice:        []int{3},
				SliceDefault: []int{3},
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--i=0", "--id=0", "--str=\"\"", "--strd=\"\"", "--t=0ms", "--td=0s", "--m=:0", "--md=:0", "--s=0", "--sd=0"},
			expected: defaultOptions{
				Int:        0,
				IntDefault: 0,

				String:        "",
				StringDefault: "",

				Time:        0,
				TimeDefault: 0,

				Map:        map[string]int{"": 0},
				MapDefault: map[string]int{"": 0},

				Slice:        []int{0},
				SliceDefault: []int{0},
			},
		},
	}

	for _, test := range tests {
		var opts defaultOptions

		_, err := ParseArgs(&opts, test.args)
		if err != nil {
			t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
		}

		if opts.Slice == nil {
			opts.Slice = []int{}
		}

		if !reflect.DeepEqual(opts, test.expected) {
			t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, opts)
		}
	}
}

func TestUnquoting(t *testing.T) {
	var tests = []struct {
		arg   string
		err   error
		value string
	}{
		{
			arg:   "\"abc",
			err:   strconv.ErrSyntax,
			value: "",
		},
		{
			arg:   "\"\"abc\"",
			err:   strconv.ErrSyntax,
			value: "",
		},
		{
			arg:   "\"abc\"",
			err:   nil,
			value: "abc",
		},
		{
			arg:   "\"\\\"abc\\\"\"",
			err:   nil,
			value: "\"abc\"",
		},
		{
			arg:   "\"\\\"abc\"",
			err:   nil,
			value: "\"abc",
		},
	}

	for _, test := range tests {
		var opts defaultOptions

		for _, delimiter := range []bool{false, true} {
			p := NewParser(&opts, None)

			var err error
			if delimiter {
				_, err = p.ParseArgs([]string{"--str=" + test.arg, "--strnot=" + test.arg})
			} else {
				_, err = p.ParseArgs([]string{"--str", test.arg, "--strnot", test.arg})
			}

			if test.err == nil {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}

				if test.value != opts.String {
					t.Fatalf("Expected String to be %q but got %q", test.value, opts.String)
				}
				if q := strconv.Quote(test.value); q != opts.StringNotUnquoted {
					t.Fatalf("Expected StringDefault to be %q but got %q", q, opts.StringNotUnquoted)
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error")
				} else if e, ok := err.(*Error); ok {
					if strings.HasPrefix(e.Message, test.err.Error()) {
						t.Fatalf("Expected error message to end with %q but got %v", test.err.Error(), e.Message)
					}
				}
			}
		}
	}
}

// envRestorer keeps a copy of a set of env variables and can restore the env from them
type envRestorer struct {
	env map[string]string
}

func (r *envRestorer) Restore() {
	os.Clearenv()
	for k, v := range r.env {
		os.Setenv(k, v)
	}
}

// EnvSnapshot returns a snapshot of the currently set env variables
func EnvSnapshot() *envRestorer {
	r := envRestorer{make(map[string]string)}
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			panic("got a weird env variable: " + kv)
		}
		r.env[parts[0]] = parts[1]
	}
	return &r
}

type envDefaultOptions struct {
	Int   int            `long:"i" default:"1" env:"TEST_I"`
	Time  time.Duration  `long:"t" default:"1m" env:"TEST_T"`
	Map   map[string]int `long:"m" default:"a:1" env:"TEST_M" env-delim:";"`
	Slice []int          `long:"s" default:"1" default:"2" env:"TEST_S"  env-delim:","`
}

func TestEnvDefaults(t *testing.T) {
	var tests = []struct {
		msg      string
		args     []string
		expected envDefaultOptions
		env      map[string]string
	}{
		{
			msg:  "no arguments, no env, expecting default values",
			args: []string{},
			expected: envDefaultOptions{
				Int:   1,
				Time:  time.Minute,
				Map:   map[string]int{"a": 1},
				Slice: []int{1, 2},
			},
		},
		{
			msg:  "no arguments, env defaults, expecting env default values",
			args: []string{},
			expected: envDefaultOptions{
				Int:   2,
				Time:  2 * time.Minute,
				Map:   map[string]int{"a": 2, "b": 3},
				Slice: []int{4, 5, 6},
			},
			env: map[string]string{
				"TEST_I": "2",
				"TEST_T": "2m",
				"TEST_M": "a:2;b:3",
				"TEST_S": "4,5,6",
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--i=3", "--t=3ms", "--m=c:3", "--s=3"},
			expected: envDefaultOptions{
				Int:   3,
				Time:  3 * time.Millisecond,
				Map:   map[string]int{"c": 3},
				Slice: []int{3},
			},
			env: map[string]string{
				"TEST_I": "2",
				"TEST_T": "2m",
				"TEST_M": "a:2;b:3",
				"TEST_S": "4,5,6",
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--i=0", "--t=0ms", "--m=:0", "--s=0"},
			expected: envDefaultOptions{
				Int:   0,
				Time:  0,
				Map:   map[string]int{"": 0},
				Slice: []int{0},
			},
			env: map[string]string{
				"TEST_I": "2",
				"TEST_T": "2m",
				"TEST_M": "a:2;b:3",
				"TEST_S": "4,5,6",
			},
		},
	}

	oldEnv := EnvSnapshot()
	defer oldEnv.Restore()

	for _, test := range tests {
		var opts envDefaultOptions
		oldEnv.Restore()
		for envKey, envValue := range test.env {
			os.Setenv(envKey, envValue)
		}
		_, err := ParseArgs(&opts, test.args)
		if err != nil {
			t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
		}

		if opts.Slice == nil {
			opts.Slice = []int{}
		}

		if !reflect.DeepEqual(opts, test.expected) {
			t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, opts)
		}
	}
}
