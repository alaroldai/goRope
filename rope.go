package rope

type Any interface{}

type LengthFunc func(Any) uint
type GetFunc func(Any, uint) Any
type SplitFunc func(Any, uint) (Any, Any)
type JoinFunc func(Any, Any) Any
type LeafLengthFunc func() uint

type RopeCallbacks struct {
	lengthFunc     LengthFunc
	getFunc        GetFunc
	splitFunc      SplitFunc
	joinFunc       JoinFunc
	leafLengthFunc LeafLengthFunc
}

func MakeCallbacks(length LengthFunc, get GetFunc, split SplitFunc, join JoinFunc, leafLength LeafLengthFunc) RopeCallbacks {
	return RopeCallbacks{length, get, split, join, leafLength}
}

type Rope interface {
	/**
	 *	Insert before offset
	 *	'abc'.Insert('de', 1) -> 'adebc'
	 *	Inserts such that the first character of `n` is inserted at `offset`
	 */
	Insert(n Any, offset uint) Rope

	/**
	 *	Remove in range [start, end] (i.e., including the characters at index 'start' and 'end')
	 *	If start == end, will remove only the character at that index
	 *	'abcd'.Remove(1,2) -> 'ad'
	 */
	Remove(start, end uint) Rope

	Each(func(item Any))

	Get(offset uint) Any

	/**
	 *	Returns a balanced version of a rope
	 */
	Balance() Rope
}

type node interface {
	insert(n Any, offset uint, callbacks RopeCallbacks) node
	remove(start, end uint, callbacks RopeCallbacks) node

	get(offset uint, callbacks RopeCallbacks) Any
	balance(callbacks RopeCallbacks) node

	split(offset uint, callbacks RopeCallbacks) (node, node)
	join(other node, callbacks RopeCallbacks) node

	length() uint // Length of the collection the node represents
	weight() uint // Number of leaves represented by the node
	depth() uint  // Maximum distance from the node to a leaf

	each(func(item Any))

	eachLeaf(func(n node))

	eachBalanced(func(n node))

	balanced() bool
}

type rope struct {
	root      node
	callbacks RopeCallbacks
}

func (r *rope) Insert(n Any, offset uint) Rope {
	if r.root == nil {
		r.root = makeLeaf(0, nil)
	}

	result := &rope{r.root.insert(n, offset, r.callbacks), r.callbacks}

	return result
}

func (r *rope) Remove(start, end uint) Rope {
	if r.root == nil {
		return r
	}

	return &rope{r.root.remove(start, end, r.callbacks), r.callbacks}
}

func (r *rope) Get(offset uint) Any {
	if r.root == nil {
		return nil
	}

	return r.root.get(offset, r.callbacks)
}

func (r *rope) Balance() Rope {
	return &rope{r.root.balance(r.callbacks), r.callbacks}
}

func (r *rope) Each(f func(item Any)) {
	if r.root != nil {
		r.root.each(f)
	}
}

func MakeRope(callbacks RopeCallbacks) Rope {
	return &rope{nil, callbacks}
}
