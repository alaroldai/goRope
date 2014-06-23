package rope

import "math"

type inode struct {
	left    node
	right   node
	_length uint
	_weight uint
	_depth  uint
}

func makeInode(left, right node) *inode {
	var ldepth, rdepth float64
	var llength, lweight, rlength, rweight uint
	if left != nil {
		ldepth = float64(left.depth())
		llength = left.length()
		lweight = left.weight()
	} else {
		ldepth = 0
		llength = 0
		lweight = 0
	}
	if right != nil {
		rdepth = float64(right.depth())
		rlength = right.length()
		rweight = right.weight()
	} else {
		rdepth = 0
		rlength = 0
		rweight = 0
	}
	depth := uint(math.Max(ldepth, rdepth)) + 1

	return &inode{left, right, llength + rlength, lweight + rweight, depth}
}

func (n *inode) split(offset uint, callbacks RopeCallbacks) (node, node) {
	if offset < n.left.length() {
		l, m := n.left.split(offset, callbacks)
		return l, m.join(n.right, callbacks)
	}
	if offset > n.left.length() {
		m, r := n.right.split(offset-n.left.length(), callbacks)
		return n.left.join(m, callbacks), r
	}
	return n.left, n.right
}

func (n *inode) join(other node, callbacks RopeCallbacks) node {
	return makeInode(n, other)
}

func (n *inode) each(f func(item Any)) {
	if n != nil && n.left != nil {
		n.left.each(f)
	}
	if n != nil && n.right != nil {
		n.right.each(f)
	}
}

func (n *inode) eachLeaf(f func(n node)) {
	if n.left != nil {
		n.left.eachLeaf(f)
	}
	if n.right != nil {
		n.right.eachLeaf(f)
	}
}

func (n *inode) eachBalanced(f func(n node)) {
	if n.balanced() {
		f(n)
	} else {
		if n.left != nil {
			n.left.eachBalanced(f)
		}
		if n.right != nil {
			n.right.eachBalanced(f)
		}
	}
}

func (n *inode) length() uint {
	return n._length
}

func (n *inode) weight() uint {
	return n._weight
}

func (n *inode) depth() uint {
	return n._depth
}

func (n *inode) insert(insert Any, offset uint, callbacks RopeCallbacks) node {

	if callbacks.lengthFunc(insert) > callbacks.leafLengthFunc() {
		head, tail := callbacks.splitFunc(insert, callbacks.leafLengthFunc())
		n.insert(head, offset, callbacks)
		n.insert(tail, offset+callbacks.leafLengthFunc(), callbacks)
	}

	newNode := makeLeaf(callbacks.lengthFunc(insert), insert)

	sl, sr := n.split(offset, callbacks)
	join := func(l, r node) node {
		return l.join(r, callbacks)
	}

	return join(join(sl, newNode), sr)
}

func (n *inode) remove(start, end uint, callbacks RopeCallbacks) node {
	l, m := n.split(start, callbacks)
	m, r := m.split(end+1-start, callbacks)
	return l.join(r, callbacks)
}

func (n *inode) get(offset uint, callbacks RopeCallbacks) Any {
	if n.left != nil {
		if offset < n.left.length() {
			return n.left.get(offset, callbacks)
		}
		if n.right != nil {
			return n.right.get(offset-n.left.length(), callbacks)
		}
		return nil
	}
	if n.right != nil {
		return n.right.get(offset, callbacks)
	}
	return nil
}

func (n *inode) balance(callbacks RopeCallbacks) node {
	nslots := fib(n.depth()) + 3
	slots := make([]node, nslots, nslots)
	n.eachBalanced(func(l node) {
		insert := l
		inserted := false

		for !inserted {
			targetIndex := fibIndex(insert.weight())

			lcons := (node)(nil)
			for i := targetIndex; i <= targetIndex; i-- {
				idx := targetIndex - i

				if slots[idx] != nil {
					if lcons != nil {
						lcons = slots[idx].join(lcons, callbacks)
					} else {
						lcons = slots[idx]
					}

					slots[idx] = nil
				}
			}

			if lcons != nil {
				insert = lcons.join(insert, callbacks)
			} else {
				slots[targetIndex] = insert
				inserted = true
			}
		}
	})

	balanced := (node)(nil)
	for i, e := range slots {
		if e != nil {
			if balanced != nil {
				balanced = e.join(balanced, callbacks)
			} else {
				balanced = e
			}
			slots[i] = nil
		}
	}

	return balanced
}

func (n *inode) balanced() bool {
	return n.weight() >= fib(n.depth()+2)
}
