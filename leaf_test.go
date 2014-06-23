package rope

import (
	"reflect"
	"testing"
)

func TestLeafInterface(test *testing.T) {
	l := reflect.TypeOf(makeLeaf(0, nil))
	i := reflect.TypeOf((*node)(nil)).Elem()

	typeConformityTest(test, l, i)
}

func TestLeafInsertCenter(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("def", 1, stringRopeCallbacks())

	expected := "adefbc"

	result := getString(after)
	if result != expected {
		test.Error("expected insert result to be " + expected + ", got " + result)
	}
}

func TestLeafInsertStart(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("def", 0, stringRopeCallbacks())

	expected := "defabc"

	result := getString(after)
	if result != expected {
		test.Error("expected insert result to be " + expected + ", got " + result)
	}
}

func TestLeafInsertEnd(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("def", 3, stringRopeCallbacks())

	expected := "abcdef"

	result := getString(after)
	if result != expected {
		test.Error("expected insert result to be " + expected + ", got " + result)
	}
}

func TestLeafInsertSplitLeft(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("d", 0, shortStringCallbacks())

	a, e := after.(*inode)

	if !e {
		test.Error("Expected insert result to be an inode")
	}

	left, e := a.left.(*leaf)
	if !e {
		test.Error("Expected left branch to be a leaf")
	}
	if left.data != "d" {
		test.Error("Expected left branch to contain 'd', got " + left.data.(string))
	}

	right, e := a.right.(*leaf)
	if !e {
		test.Error("Expected right branch to be a leaf")
	}
	if right.data != "abc" {
		test.Error("Expected right branch to contain 'abc', got " + right.data.(string))
	}

}

func TestLeafInsertSplitRight(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("d", 3, shortStringCallbacks())

	a, e := after.(*inode)

	if !e {
		test.Error("Expected insert result to be an inode")
	}

	left, e := a.left.(*leaf)
	if !e {
		test.Error("Expected left branch to be a leaf")
	}
	if left.data != "abc" {
		test.Error("Expected left branch to contain 'abc', got " + left.data.(string))
	}

	right, e := a.right.(*leaf)
	if !e {
		test.Error("Expected right branch to be a leaf")
	}
	if right.data != "d" {
		test.Error("Expected right branch to contain 'd', got " + right.data.(string))
	}

}

func TestLeafInsertSplitCentreLeft(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("d", 1, shortStringCallbacks())

	a, e := after.(*inode)

	if !e {
		test.Error("Expected insert result to be an inode")
	}

	left, e := a.left.(*leaf)
	if !e {
		test.Error("Expected left branch to be a leaf")
	}
	if left.data != "ad" {
		test.Error("Expected left branch to contain 'ad', got " + left.data.(string))
	}

	right, e := a.right.(*leaf)
	if !e {
		test.Error("Expected right branch to be a leaf")
	}
	if right.data != "bc" {
		test.Error("Expected right branch to contain 'bc', got " + right.data.(string))
	}
}

func TestLeafInsertSplitCentreRight(test *testing.T) {
	l := makeLeaf(3, "abc")

	after := l.insert("d", 2, shortStringCallbacks())

	a, e := after.(*inode)

	if !e {
		test.Error("Expected insert result to be an inode")
	}

	ar, e := a.right.(*leaf)
	if !e {
		test.Error("Expected right branch to be a leaf")
	}
	if ar.data != "dc" {
		test.Error("Expected right branch to contain 'dc', got " + ar.data.(string))
	}

	al, e := a.left.(*leaf)
	if !e {
		test.Error("Expected left branch to be a leaf")
	}
	if al.data != "ab" {
		test.Error("Expected right branch to contain 'ab', got " + al.data.(string))
	}
}

func TestLeafRemoveSingle(test *testing.T) {
	l := makeLeaf(3, "abc")
	after := l.remove(1, 1, stringRopeCallbacks()).(*leaf)

	expected := "ac"
	if after.data != expected {
		test.Error("Expected remove result to be " + expected + ", got " + after.data.(string))
	}
}

func TestLeafRemoveRangeEnd(test *testing.T) {
	l := makeLeaf(3, "abc")
	after := l.remove(1, 2, stringRopeCallbacks()).(*leaf)

	expected := "a"
	if after.data != expected {
		test.Error("Expected remove result to be " + expected + ", got " + after.data.(string))
	}
}

func TestLeafRemoveRangeStart(test *testing.T) {
	l := makeLeaf(3, "abc")
	after := l.remove(0, 1, stringRopeCallbacks()).(*leaf)

	expected := "c"
	if after.data != expected {
		test.Error("Expected remove result to be " + expected + ", got " + after.data.(string))
	}
}

func TestLeafRemoveAll(test *testing.T) {
	l := makeLeaf(3, "abc")
	after := l.remove(0, 2, stringRopeCallbacks())

	if after != nil {
		test.Error("Expected remove result to be nil, got " + after.(*leaf).data.(string))
	}
}

func TestLeafGet(test *testing.T) {
	l := makeLeaf(3, "abc")
	r := l.get(0, stringRopeCallbacks())
	expected := 'a'
	if r != expected {
		test.Error("Expected l.get to return " + string(expected) + ", got " + string(r.(rune)))
	}
}
