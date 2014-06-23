package rope

import (
	"bytes"
	"strconv"
	"testing"
)

func TestINodeGet(test *testing.T) {
	l := makeLeaf(3, "abc")
	r := makeLeaf(3, "def")
	i := makeInode(l, r)

	x := i.get(0, stringRopeCallbacks())
	ex := 'a'
	if x != ex {
		test.Error("Expected l.get to return " + string(ex) + ", got " + string(x.(rune)))
	}

	x = i.get(5, stringRopeCallbacks())
	ex = 'f'
	if x != ex {
		test.Error("Expected l.get to return " + string(ex) + ", got " + string(x.(rune)))
	}
}

func TestINodeGetNil(test *testing.T) {
	i := makeInode(nil, nil)
	if i.get(0, stringRopeCallbacks()) != nil {
		test.Error("Expected i.get to return nil with nil child nodes")
	}

	i.right = makeLeaf(3, "abc")
	if i.get(0, stringRopeCallbacks()) != 'a' {
		test.Error("Expected i.get to fetch from the right node when left is nil")
	}

	i.left = i.right
	i.right = nil

	if i.get(4, stringRopeCallbacks()) != nil {
		test.Error("Expected i.get to fetch from the right node when left is nil")
	}

}

func TestInodeSplit(test *testing.T) {
	var l node = makeLeaf(3, "abc")
	var r node = makeLeaf(3, "def")
	i := makeInode(l, r)

	lt, rt := i.split(3, stringRopeCallbacks())
	if rt != i.right || lt != i.left {
		result := getString(lt) + getString(rt)
		test.Error("Expected i.split to return a copy of i, got " + result)
	}
}

func TestInodeInsert(test *testing.T) {
	l := makeLeaf(3, "abc")
	r := makeLeaf(3, "def")
	i := makeInode(l, r)

	o := i.insert("a", 1, shortStringCallbacks())

	var b bytes.Buffer
	o.each(func(item Any) {
		b.WriteString(item.(string))
	})
	if b.String() != "aabcdef" {
		test.Error("Expected i.insert to result in aabcdef, got " + b.String())
	}
}

func TestInodeRemove(test *testing.T) {
	l := makeLeaf(3, "abc")
	r := makeLeaf(3, "def")
	i := makeInode(l, r)

	o := i.remove(2, 3, stringRopeCallbacks())

	var b bytes.Buffer
	o.each(func(item Any) {
		b.WriteString(item.(string))
	})
	if b.String() != "abef" {
		test.Error("Expected i.remove to result in abef, got " + b.String())
	}
}

func TestInodeBalance(test *testing.T) {
	a := makeLeaf(3, "abc")
	b := makeLeaf(3, "def")
	c := makeLeaf(3, "ghi")
	d := makeLeaf(3, "jkl")
	i := makeInode(makeInode(makeInode(a, b), nil), makeInode(nil, makeInode(c, d)))

	result := i.balance(shortStringCallbacks())
	if result.depth() != 2 {
		test.Error("Expected i.balance().depth() to be 1, got " + strconv.Itoa(int(result.depth())))
	}

	var buf bytes.Buffer
	result.each(func(item Any) {
		buf.WriteString(item.(string))
	})
	if buf.String() != "abcdefghijkl" {
		test.Error("Expected i.balance() to not change the stored string")
	}
}
