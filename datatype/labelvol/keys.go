/*
	This file supports keyspaces for labelvol data type.
*/

package labelvol

import (
	"encoding/binary"

	"github.com/janelia-flyem/dvid/dvid"
	"github.com/janelia-flyem/dvid/storage"
)

const (
	// keyUnknown should never be used and is a check for corrupt or incorrectly set keys
	keyUnknown storage.TKeyClass = iota

	// keyLabelBlockRLE have keys ordered by label + block coord, and have a sparse volume
	// encoding for its value. They are also useful for returning all blocks
	// intersected by a label.
	keyLabelBlockRLE = 227

	keyLabelMax = 228

	keyRepoLabelMax = 229
)

// DescribeTKeyClass returns a string explanation of what a particular TKeyClass
// is used for.  Implements the datastore.TKeyClassDescriber interface.
func (d *Data) DescribeTKeyClass(tkc storage.TKeyClass) string {
	switch tkc {
	case keyLabelBlockRLE:
		return "labelvol label + block coord key"
	case keyLabelMax:
		return "labelvol label max key"
	case keyRepoLabelMax:
		return "labelvol repo label max key"
	default:
	}
	return "unknown labelvol key"
}

// NewTKey returns a TKey for storing a "label + spatial index", where
// the spatial index references a block that contains a voxel with the given label.
func NewTKey(label uint64, block dvid.IZYXString) storage.TKey {
	sz := len(block)
	ibytes := make([]byte, 8+sz)
	binary.BigEndian.PutUint64(ibytes[0:8], label)
	copy(ibytes[8:], []byte(block))
	return storage.NewTKey(keyLabelBlockRLE, ibytes)
}

// DecodeTKey returns a label and block index bytes from a label block RLE key.
// The block index bytes are returned because different block indices may be used (e.g., CZYX),
// and its up to caller to determine which one is used for this particular key.
func DecodeTKey(tk storage.TKey) (label uint64, block dvid.IZYXString, err error) {
	ibytes, err := tk.ClassBytes(keyLabelBlockRLE)
	if err != nil {
		return
	}
	label = binary.BigEndian.Uint64(ibytes[0:8])
	block = dvid.IZYXString(ibytes[8:])
	return
}

var (
	maxLabelTKey     = storage.NewTKey(keyLabelMax, nil)
	maxRepoLabelTKey = storage.NewTKey(keyRepoLabelMax, nil)
)
