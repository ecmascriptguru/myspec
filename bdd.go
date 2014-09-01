package gomspec

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

type specification struct {
	T       *testing.T
	Feature string
	Context string
	When    string
	Title   string
	Fn      func(Expect)
}

func (spec *specification) run() {
	spec.Fn(func(val interface{}) *expectation {
		return &expectation{spec, val}
	})
}

// Given defines the Feature's specific context to be spec'd out.
func Given(t *testing.T, context string, scenerioWrapper func(When)) {

	pc, _, _, _ := runtime.Caller(1)
	featureDesc := func() string {
		m := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
		i := strings.LastIndex(m, ".")
		m = m[i+1 : len(m)]
		m = strings.Replace(m, "Test_", "", 1)
		m = strings.Replace(m, "Test", "", 1)
		return strings.Replace(m, "_", " ", -1)
	}

	scenerioWrapper(func(when string, testWrapper func(It)) {
		testWrapper(func(it string, fn func(Expect)) {
			spec := &specification{
				t,
				featureDesc(),
				context,
				when,
				it,
				fn,
			}
			spec.run()
		})
		fmt.Println()
	})
}

// When defines the action or event when Given a specific context.
type When func(when string, fn func(It))

// It defines the specification of when something happens.
type It func(title string, fn func(Expect))

// Setup is used to define before/after (setup/teardown) functions.
func Setup(before, after func()) func(fn func(Expect)) func(Expect) {
	return func(fn func(Expect)) func(Expect) {
		before()
		return func(expect Expect) {
			fn(expect)
			after()
		}
	}
}

// NotImplemented is used to mark a specification that needs coding out.
func NotImplemented() func(Expect) {
	return func(expect Expect) { expect(nil).notImplemented() }
}

// NA is shorthand for the NotImplemented() function.
func NA() func(Expect) {
	return NotImplemented()
}

// Desc is legacy support for existing Zen users.
func Desc(t *testing.T, desc string, wrapper func(It)) {
	wrapper(func(it string, fn func(Expect)) {
		spec := &specification{
			t,
			"<not set>",
			"",
			desc,
			it,
			fn,
		}
		spec.run()
	})
}