package trace

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
)

// FieldColumn represents a column of data within a trace where each row is
// stored directly as a field element.  This is the simplest form of column,
// which provides the fastest Get operation (i.e. because it just reads the
// field element out directly).  However, at the same time, it can potentially
// use quite a lot of memory.  In particular, when there are many different
// field elements which have smallish values then this requires excess data.
type FieldColumn struct {
	// Holds the name of this column
	name string
	// Holds the raw data making up this column
	data []*fr.Element
	// Value to be used when padding this column
	padding *fr.Element
}

// Name returns the name of the given column.
func (p *FieldColumn) Name() string {
	return p.name
}

// Width determines the number of bytes per element for this column (which, in
// this case, is always 32).
func (p *FieldColumn) Width() uint {
	return 32
}

// Height determines the height of this column.
func (p *FieldColumn) Height() uint {
	return uint(len(p.data))
}

// Padding returns the value which will be used for padding this column.
func (p *FieldColumn) Padding() *fr.Element {
	return p.padding
}

// Data returns the data for the given column.
func (p *FieldColumn) Data() []*fr.Element {
	return p.data
}

// Get the value at a given row in this column.  If the row is
// out-of-bounds, then the column's padding value is returned instead.
// Thus, this function always succeeds.
func (p *FieldColumn) Get(row int) *fr.Element {
	if row < 0 || row >= len(p.data) {
		// out-of-bounds access
		return p.padding
	}
	// in-bounds access
	return p.data[row]
}

// Clone an FieldColumn
func (p *FieldColumn) Clone() Column {
	clone := new(FieldColumn)
	clone.name = p.name
	clone.padding = p.padding
	// NOTE: the following is as we never actually mutate the underlying bytes
	// array.
	clone.data = p.data

	return clone
}

// Pad this column with n copies of the column's padding value.
func (p *FieldColumn) Pad(n uint) {
	// Allocate sufficient memory
	ndata := make([]*fr.Element, uint(len(p.data))+n)
	// Copy over the data
	copy(ndata[n:], p.data)
	// Go padding!
	for i := uint(0); i < n; i++ {
		ndata[i] = p.padding
	}
	// Copy over
	p.data = ndata
}

// Write the raw bytes of this column to a given writer, returning an error
// if this failed (for some reason).
func (p *FieldColumn) Write(w io.Writer) error {
	panic("TODO")
}
