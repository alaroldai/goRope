package rope

type leaf struct {
	dataLength uint
	data       Any
}

func makeLeaf(dataLength uint, data Any) *leaf {
	return &leaf{dataLength, data}
}

func (l *leaf) each(f func(item Any)) {
	if l != nil && l.data != nil {
		f(l.data)
	}
}

func (l *leaf) eachLeaf(f func(n node)) {
	f(l)
}

func (l *leaf) eachBalanced(f func(n node)) {
	f(l)
}

func (l *leaf) length() uint {
	return l.dataLength
}

func (l *leaf) weight() uint {
	return 1
}

func (l *leaf) depth() uint {
	return 0
}

func (l *leaf) split(offset uint, callbacks RopeCallbacks) (node, node) {
	if l == nil || l.data == nil {
		return nil, nil
	}
	left, right := callbacks.splitFunc(l.data, offset)
	return makeLeaf(callbacks.lengthFunc(left), left), makeLeaf(callbacks.lengthFunc(right), right)
}

func (l *leaf) join(other node, callbacks RopeCallbacks) node {
	ol, succ := other.(*leaf)
	if succ && l.length()+ol.length() <= callbacks.leafLengthFunc() {
		n := callbacks.joinFunc(l.data, ol.data)
		return makeLeaf(callbacks.lengthFunc(n), n)
	}
	return makeInode(l, other)
}

func (l *leaf) insert(insert Any, offset uint, callbacks RopeCallbacks) node {
	insertLength := callbacks.lengthFunc(insert)

	left, right := l.split(offset, callbacks)
	if left == nil || left.length() == 0 {
		left = makeLeaf(insertLength, insert)
	} else if right == nil || right.length() == 0 {
		right = makeLeaf(insertLength, insert)
	} else if left.length() < right.length() {
		left = left.join(makeLeaf(insertLength, insert), callbacks)
	} else {
		right = (makeLeaf(insertLength, insert)).join(right, callbacks)
	}

	return left.join(right, callbacks)
}

func (l *leaf) remove(start, end uint, callbacks RopeCallbacks) node {
	if end-start+1 == l.dataLength {
		return nil
	}

	leftc, rightc := callbacks.splitFunc(l.data, start)
	_, rightc = callbacks.splitFunc(rightc, end+1-start)
	retc := callbacks.joinFunc(leftc, rightc)
	return makeLeaf(callbacks.lengthFunc(retc), retc)
}

func (l *leaf) get(offset uint, callbacks RopeCallbacks) Any {
	return callbacks.getFunc(l.data, offset)
}

func (l *leaf) balance(callbacks RopeCallbacks) node {
	return l
}

func (l *leaf) balanced() bool {
	return true
}
