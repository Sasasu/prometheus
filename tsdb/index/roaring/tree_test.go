package roaring

import (
	"github.com/prometheus/prometheus/util/testutil"
	"testing"
)

func pointer(v T) *T {
	return &v
}

func TestNodeFind(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		testutil.Equals(t, _find(nil, []byte{1, 2}, false), (*node)(nil))
		testutil.Equals(t, _find(nil, []byte{1, 2}, true).key, []byte{1, 2})
	})

	t.Run("normal", func(t *testing.T) {
		var root = &node{
			key: []byte{1, 2},
			child: []node{
				{key: []byte{2}, data: pointer(1)},
				{key: []byte{3}, data: pointer(2)},
			},
		}

		testutil.Equals(t, _find(root, []byte{}, false), root)
		testutil.Equals(t, _find(root, []byte{1}, false), root)
		testutil.Equals(t, _find(root, []byte{1, 2}, false), root)
		testutil.Equals(t, _find(root, []byte{1, 2, 4}, false), (*node)(nil))

		testutil.Equals(t, _find(root, []byte{1, 2, 2}, false).data, pointer(1))
		testutil.Equals(t, _find(root, []byte{1, 2, 3}, false).data, pointer(2))
	})

	t.Run("insert", func(t *testing.T) {
		var root = &node{
			key: []byte{1, 2},
			child: []node{
				{key: []byte{2}, data: pointer(1)},
				{key: []byte{3}, data: pointer(2)},
			},
		}

		{
			var n = _find(root, []byte{1, 2, 4}, true)
			testutil.Equals(t, n.key, []byte{4})
			n.data = pointer(1)
		}
		{
			var n = _find(root, []byte{1, 2, 4, 5}, true)
			testutil.Equals(t, n.key, []byte{5})
			n.data = pointer(2)
		}
	})
}

func TestPrefixLength(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		testutil.Equals(t, _prefixLength(nil, nil), 0)
		testutil.Equals(t, _prefixLength([]byte{1, 2}, nil), 0)
		testutil.Equals(t, _prefixLength(nil, []byte{1, 2}), 0)
	})

	t.Run("normal", func(t *testing.T) {
		testutil.Equals(t, _prefixLength([]byte{1, 2}, []byte{2, 3}), 0)
		testutil.Equals(t, _prefixLength([]byte{1, 2}, []byte{1, 3}), 1)
		testutil.Equals(t, _prefixLength([]byte{1, 3}, []byte{1, 2, 4, 7}), 1)
		testutil.Equals(t, _prefixLength([]byte("鹿目 まどか"), []byte("鹿目 まどか")), 16)
		testutil.Equals(t, _prefixLength([]byte("鹿目 まどか"), []byte("鹿目 まどか鹿目")), 16)
	})
}
