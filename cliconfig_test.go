package cliconfig

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
)

type testConfig struct {
	Name                string
	SomeOtherName       string
	CompletelyDifferent string `cli:"different,some usage here, with comma"`
	HiddenFlag          bool   `cli:",hidden"`
	OtherHidden         bool   `cli:"other,hidden"`
	NotHidden           bool   `cli:",nothidden"`
	ErrHidden           bool   `cli:","`
	Ignore              int    `cli:"-"`
	Nested              struct {
		Val1 int
		Val2 float64
	}
	Dur time.Duration
	SS  []string
}

func TestCLIConfig(t *testing.T) {
	cc := New(testConfig{
		SS: []string{
			"one",
			"two",
		},
	})

	app := cli.NewApp()
	app.Flags = cc.Flags()

	Convey("flags should be properly generated", t, func() {
		So(app.Flags, ShouldResemble, []cli.Flag{
			cli.StringFlag{
				Name:   "name",
				EnvVar: "NAME",
			},
			cli.StringFlag{
				Name:   "some-other-name",
				EnvVar: "SOME_OTHER_NAME",
			},
			cli.StringFlag{
				Name:   "different",
				EnvVar: "DIFFERENT",
				Usage:  "some usage here, with comma",
			},
			cli.BoolFlag{
				Name:   "hidden-flag",
				EnvVar: "HIDDEN_FLAG",
				Hidden: true,
			},
			cli.BoolFlag{
				Name:   "other",
				EnvVar: "OTHER",
				Hidden: true,
			},
			cli.BoolFlag{
				Name:   "not-hidden",
				EnvVar: "NOT_HIDDEN",
				Usage:  "nothidden",
			},
			cli.BoolFlag{
				Name:   "err-hidden",
				EnvVar: "ERR_HIDDEN",
			},
			cli.IntFlag{
				Name:   "nested-val1",
				EnvVar: "NESTED_VAL1",
			},
			cli.Float64Flag{
				Name:   "nested-val2",
				EnvVar: "NESTED_VAL2",
			},
			cli.DurationFlag{
				Name:   "dur",
				EnvVar: "DUR",
			},
			cli.StringSliceFlag{
				Name:   "ss",
				EnvVar: "SS",
				Value:  &cli.StringSlice{"one", "two"},
			},
		})
	})

	Convey("flags should be properly parsed", t, func() {
		app.Before = before(t, cc)
		// So(app.Run([]string{"testapp", "-test-ss", "three"}), ShouldBeNil)
		So(app.Run([]string{"testapp"}), ShouldBeNil)
	})
}

func before(t *testing.T, cc *CLIConfig) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var config testConfig
		So(cc.Parse(c, &config), ShouldBeNil)
		So(config, ShouldResemble, testConfig{
			SS: []string{"one", "two"},
		})
		return nil
	}
}

func testEq(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

type camelTest struct {
	Value  string
	Expect []string
}

func TestFromCamel(t *testing.T) {
	table := []camelTest{{
		Value:  "",
		Expect: []string{},
	}, {
		Value:  "lowercase",
		Expect: []string{"lowercase"},
	}, {
		Value:  "Class",
		Expect: []string{"Class"},
	}, {
		Value:  "MyClass",
		Expect: []string{"My", "Class"},
	}, {
		Value:  "MyC",
		Expect: []string{"My", "C"},
	}, {
		Value:  "HTML",
		Expect: []string{"HTML"},
	}, {
		Value:  "PDFLoader",
		Expect: []string{"PDF", "Loader"},
	}, {
		Value:  "AString",
		Expect: []string{"A", "String"},
	}, {
		Value:  "SimpleXMLParser",
		Expect: []string{"Simple", "XML", "Parser"},
	}, {
		Value:  "vimRPCPlugin",
		Expect: []string{"vim", "RPC", "Plugin"},
	}, {
		Value:  "GL11Version",
		Expect: []string{"GL11", "Version"},
	}, {
		Value:  "99Bottles",
		Expect: []string{"99Bottles"},
	}, {
		Value:  "99bottles",
		Expect: []string{"99bottles"},
	}, {
		Value:  "May5",
		Expect: []string{"May5"},
	}, {
		Value:  "BFG9000",
		Expect: []string{"BFG9000"},
	}, {
		Value:  "BöseÜberraschung",
		Expect: []string{"Böse", "Überraschung"},
	}, {
		Value:  "Two  spaces",
		Expect: []string{"Two", "  ", "spaces"},
	}, {
		Value:  "BadUTF8\xe2\xe2\xa1",
		Expect: []string{"BadUTF8\xe2\xe2\xa1"},
	}, {
		Value:  "ipv4",
		Expect: []string{"ipv4"},
	}, {
		Value:  "IPV4",
		Expect: []string{"IPV4"},
	}, {
		Value:  "99",
		Expect: []string{"99"},
	}}

	for _, c := range table {
		if !testEq(fromCamel(c.Value), c.Expect) {
			t.Errorf("fromCamel: expected %s to return %v but got %v instead", c.Value, c.Expect, fromCamel(c.Value))
		}
	}
}
