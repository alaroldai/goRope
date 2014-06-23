package rope

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func methodsMissingFromType(inter, typ reflect.Type) []string {
	missingMethods := make([]string, 0)
	for n := 0; n < inter.NumMethod(); n++ {
		_, present := typ.MethodByName(inter.Method(n).Name)
		if !present {
			fmt.Println(inter.Method(n).Name)
			missingMethods = append(missingMethods, inter.Method(n).Name)
		}
	}
	return missingMethods
}

func typeConformityTest(test *testing.T, stype, itype reflect.Type) {
	if !stype.Implements(itype) {
		missingMethods := methodsMissingFromType(itype, stype)
		for stype.Kind() == reflect.Ptr {
			stype = stype.Elem()
		}
		for itype.Kind() == reflect.Ptr {
			itype = itype.Elem()
		}
		test.Error("struct '" + stype.Name() + "' does not implement interface '" + itype.Name() + "' (missing methods: " + strings.Join(missingMethods, ", ") + ")")
	}
}

func stringRopeCallbacks() RopeCallbacks {
	callbacks := MakeCallbacks(
		//lengthFunc
		func(i Any) uint {
			s := i.(string)
			return uint(len(s))
		},
		//getFunc
		func(i Any, offset uint) Any {
			s := i.(string)
			for i, a := range s {
				if uint(i) == offset {
					return a
				}
			}
			return '\u0000'
		},
		//splitFunc
		func(i Any, offset uint) (Any, Any) {
			s := i.(string)
			return s[:offset], s[offset:]
		},
		//joinFunc
		func(i Any, t Any) Any {
			var b bytes.Buffer
			b.WriteString(i.(string))
			b.WriteString(t.(string))
			return b.String()
		},
		func() uint {
			return 32
		},
	)

	return callbacks
}

func shortStringCallbacks() RopeCallbacks {
	ret := stringRopeCallbacks()
	ret.leafLengthFunc = func() uint {
		return 3
	}
	return ret
}

func getString(n node) string {
	var b bytes.Buffer
	n.each(func(i Any) {
		b.WriteString(i.(string))
	})
	return b.String()
}

func TestINodeInterface(test *testing.T) {
	n := reflect.TypeOf(makeInode(nil, nil))
	i := reflect.TypeOf((*node)(nil)).Elem()

	typeConformityTest(test, n, i)
}

func TestRopeInterface(test *testing.T) {
	r := reflect.TypeOf(&rope{})
	i := reflect.TypeOf((*Rope)(nil)).Elem()

	typeConformityTest(test, r, i)
}

func TestFibIndex(test *testing.T) {
	result := fibIndex(13)
	if result != 7 {
		test.Error("Expected fibIndex(13) to be 7, got " + strconv.Itoa(int(result)))
	}
}

func TestFib(test *testing.T) {
	result := fib(7)
	if result != 13 {
		test.Error("expected fib(7) to be 13, got " + strconv.Itoa(int(result)))
	}
}
