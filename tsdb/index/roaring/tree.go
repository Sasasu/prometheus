package roaring

import "sort"

type T = int

type tree struct {
	head *node
}

type node struct {
	key []byte

	child []node
	data  *T
}

type iterator struct {
	key     []byte
	stack   []*node
	current *node
}

func newTree(t *tree) *tree {
	if t == nil {
		t = &tree{}
	}

	return t
}

func (t *tree) insert(key []byte, val *T) {
	if t.head == nil {
		t.head = _find(t.head, key, true)
		t.head.data = val
		return
	}

	var node = _find(t.head, key, true)
	node.data = val
	return
}

func (t *tree) find(key []byte) *T {
	var val = _find(t.head, key, false)

	if val == nil {
		return nil
	}
	return val.data
}

func (*tree) insertBatch(key [][]byte, val []*T)            {}
func (*tree) iterFrom(key []byte, iter *iterator) *iterator { return iter }

func _find(root *node, key []byte, create bool) *node {
	// the key suffix on sub node
	var keyOnSubNode = key

	for root != nil && root.data == nil && len(key) >= len(root.key) {
		// check if the new root contains this key
		keyOnSubNode = key[_prefixLength(root.key, key):]
		if len(keyOnSubNode) == 0 {
			// return the new root
			goto END
		}

		for _, i := range root.child {
			// check if the new root's child contains this key
			if _prefixLength(i.key, keyOnSubNode) == 0 {
				// not contain, check next
				continue
			}

			// key contained on the child, move to the next node
			root = &i
			key = keyOnSubNode
			break
		}

		// all child not contained this key, return the new root
		goto END
	}

END:
	if root == nil {
		if create {
			return &node{key: keyOnSubNode}
		} else {
			return nil
		}
	}

	var prefix = _prefixLength(root.key, keyOnSubNode)

	if prefix == len(keyOnSubNode) {
		return root
	}

	if !create {
		return nil
	}

	// split the root
	root.child = _appendClientOrdered(root.child, node{key: keyOnSubNode[prefix:]})
	return &root.child[len(root.child)-1]
}

func _prefixLength(a []byte, b []byte) (i int) {
	var _min = func(a, b int) int {
		if a > b {
			return b
		} else {
			return a
		}
	}

	for min := _min(len(a), len(b)); i < min; i++ {
		if a[i] == b[i] {
			continue
		}
		return i
	}
	return
}

func _appendClientOrdered(list []node, val node) []node {
	var i = sort.Search(len(list), func(i int) bool {
		return string(list[i].key) >= string(val.key)
	})

	if i >= len(list) {
		return append(list, val)
	}
	return append(list[:i], append([]node{val}, list[i:]...)...)
}
