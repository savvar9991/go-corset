package table

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/go-corset/pkg/util"
)

// Acceptable represents an element which can "accept" a trace, or either reject
// with an error or report a warning.
type Acceptable interface {
	Accepts(Trace) error
}

// Column describes an individual column of data within a trace table.
type Column interface {
	// Get the name of this column
	Name() string
	// Return the height (i.e. number of rows) of this column.
	Height() uint
	// Return the data stored in this column.
	Data() []*fr.Element
	// Return the value to use for padding this column.
	Padding() *fr.Element
}

// Trace describes a set of named columns.  Columns are not required to have the
// same height and can be either "data" columns or "computed" columns.
type Trace interface {
	// Add a new column of data
	AddColumn(name string, data []*fr.Element, padding *fr.Element)
	// Get the name of the ith column in this trace.
	ColumnName(int) string
	// ColumnByIndex returns the ith column in this trace.
	ColumnByIndex(uint) Column
	// ColumnByName returns the data of a given column in order that it can be
	// inspected.  If the given column does not exist, then nil is returned.
	ColumnByName(name string) Column
	// Check whether this trace contains data for the given column.
	HasColumn(name string) bool
	// Get the value of a given column by its name.  If the column
	// does not exist or if the index is out-of-bounds then an
	// error is returned.
	//
	// NOTE: this operation is expected to be slower than
	// GetByindex as, depending on the underlying data format,
	// this may first resolve the name into a physical column
	// index.
	GetByName(name string, row int) *fr.Element
	// Get the value of a given column by its index. If the column
	// does not exist or if the index is out-of-bounds then an
	// error is returned.
	GetByIndex(col int, row int) *fr.Element
	// Pad each column in this trace with n items at the front.  An iterator over
	// the padding values to use for each column must be given.
	Pad(n uint)
	// Determine the height of this table, which is defined as the
	// height of the largest column.
	Height() uint
	// Get the number of columns in this trace.
	Width() uint
}

// ConstraintsAcceptTrace determines whether or not one or more groups of
// constraints accept a given trace.  It returns the first error or warning
// encountered.
func ConstraintsAcceptTrace[T Acceptable](trace Trace, constraints []T) error {
	for _, c := range constraints {
		err := c.Accepts(trace)
		if err != nil {
			return err
		}
	}
	//
	return nil
}

// ===================================================================
// Array Trace
// ===================================================================

// ArrayTrace provides an implementation of Trace which stores columns as an
// array.
type ArrayTrace struct {
	// Holds the maximum height of any column in the trace
	height uint
	// Holds the name of each column
	columns []*ArrayTraceColumn
}

// EmptyArrayTrace constructs an empty array trace into which column data can be
// added.
func EmptyArrayTrace() *ArrayTrace {
	p := new(ArrayTrace)
	// Initially empty columns
	p.columns = make([]*ArrayTraceColumn, 0)
	// Initialise height as 0
	p.height = 0
	// done
	return p
}

// Width returns the number of columns in this trace.
func (p *ArrayTrace) Width() uint {
	return uint(len(p.columns))
}

// ColumnName returns the name of the ith column in this trace.
func (p *ArrayTrace) ColumnName(index int) string {
	return p.columns[index].Name()
}

// Clone creates an identical clone of this trace.
func (p *ArrayTrace) Clone() *ArrayTrace {
	clone := new(ArrayTrace)
	clone.columns = make([]*ArrayTraceColumn, len(p.columns))
	clone.height = p.height
	//
	for i, c := range p.columns {
		// TODO: can this be avoided?
		clone.columns[i] = c.Clone()
	}
	// done
	return clone
}

// HasColumn checks whether the trace has a given column or not.
func (p *ArrayTrace) HasColumn(name string) bool {
	for _, c := range p.columns {
		if c.name == name {
			return true
		}
	}

	return false
}

// AddColumn adds a new column of data to this trace.
func (p *ArrayTrace) AddColumn(name string, data []*fr.Element, padding *fr.Element) {
	// Sanity check the column does not already exist.
	if p.HasColumn(name) {
		panic("column already exists")
	}
	// Construct new column
	column := ArrayTraceColumn{name, data, padding}
	// Append it
	p.columns = append(p.columns, &column)
	// Update maximum height
	if uint(len(data)) > p.height {
		p.height = uint(len(data))
	}
}

// Columns returns the set of columns in this trace.
func (p *ArrayTrace) Columns() []*ArrayTraceColumn {
	return p.columns
}

// GetByName gets the value of a given column (as identified by its name) at a
// given row.  If the column does not exist, an error is returned.
func (p *ArrayTrace) GetByName(name string, row int) *fr.Element {
	// NOTE: Could improve performance here if names were kept in
	// sorted order.
	c := p.getColumnByName(name)
	if c != nil {
		// Matched column
		return c.Get(row)
	}
	// Precondition failure
	panic(fmt.Sprintf("Invalid column: {%s}", name))
}

// ColumnByIndex looks up a column based on its index.
func (p *ArrayTrace) ColumnByIndex(index uint) Column {
	return p.columns[index]
}

// ColumnByName looks up a column based on its name.  If the column doesn't
// exist, then nil is returned.
func (p *ArrayTrace) ColumnByName(name string) Column {
	for _, c := range p.columns {
		if name == c.name {
			// Matched column
			return c
		}
	}

	return nil
}

// GetByIndex returns the value of a given column (as identifier by its index or
// register) at a given row.  If the column is out-of-bounds an error is
// returned.
func (p *ArrayTrace) GetByIndex(col int, row int) *fr.Element {
	if col < 0 || col >= len(p.columns) {
		// Precondition failure
		panic(fmt.Sprintf("Invalid column: {%d}", col))
	}

	return p.columns[col].Get(row)
}

func (p *ArrayTrace) getColumnByName(name string) *ArrayTraceColumn {
	for _, c := range p.columns {
		if name == c.name {
			// Matched column
			return c
		}
	}

	return nil
}

// Height determines the maximum height of any column within this trace.
func (p *ArrayTrace) Height() uint {
	return p.height
}

// Pad each column in this trace with n items at the front.  An iterator over
// the padding values to use for each column must be given.
func (p *ArrayTrace) Pad(n uint) {
	for _, c := range p.columns {
		c.Pad(n)
	}
	// Increment height
	p.height += n
}

func (p *ArrayTrace) String() string {
	// Use string builder to try and make this vaguely efficient.
	var id strings.Builder

	id.WriteString("{")

	for i := 0; i < len(p.columns); i++ {
		if i != 0 {
			id.WriteString(",")
		}

		id.WriteString(p.columns[i].name)
		id.WriteString("={")

		for j := uint(0); j < p.height; j++ {
			jth := p.GetByIndex(i, int(j))

			if j != 0 {
				id.WriteString(",")
			}

			if jth == nil {
				id.WriteString("_")
			} else {
				id.WriteString(jth.String())
			}
		}
		id.WriteString("}")
	}
	id.WriteString("}")
	//
	return id.String()
}

// ===================================================================
// Array Trace Column
// ===================================================================

// ArrayTraceColumn represents a column of data within an array trace.
type ArrayTraceColumn struct {
	// Holds the name of this column
	name string
	// Holds the raw data making up this column
	data []*fr.Element
	// Value to be used when padding this column
	padding *fr.Element
}

// Name returns the name of the given column.
func (p *ArrayTraceColumn) Name() string {
	return p.name
}

// Height determines the height of this column.
func (p *ArrayTraceColumn) Height() uint {
	return uint(len(p.data))
}

// Padding returns the value which will be used for padding this column.
func (p *ArrayTraceColumn) Padding() *fr.Element {
	return p.padding
}

// Data returns the data for the given column.
func (p *ArrayTraceColumn) Data() []*fr.Element {
	return p.data
}

// Clone an ArrayTraceColumn
func (p *ArrayTraceColumn) Clone() *ArrayTraceColumn {
	clone := new(ArrayTraceColumn)
	clone.name = p.name
	clone.data = make([]*fr.Element, len(p.data))
	clone.padding = p.padding
	copy(clone.data, p.data)

	return clone
}

// Pad this column with n copies of a given value, either at the front
// (sign=true) or the back (sign=false).
func (p *ArrayTraceColumn) Pad(n uint) {
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

// Get the value at the given row of this column.
func (p *ArrayTraceColumn) Get(row int) *fr.Element {
	if row >= 0 && row < len(p.data) {
		return p.data[row]
	}
	// For out-of-bounds access, return padding value.
	return p.padding
}

// ===================================================================
// JSON Parser
// ===================================================================

// ParseJsonTrace parses a trace expressed in JSON notation.  For example, {"X":
// [0], "Y": [1]} is a trace containing one row of data each for two columns "X"
// and "Y".
func ParseJsonTrace(bytes []byte) (*ArrayTrace, error) {
	var zero fr.Element = fr.NewElement((0))

	var rawData map[string][]*big.Int
	// Unmarshall
	jsonErr := json.Unmarshal(bytes, &rawData)
	if jsonErr != nil {
		return nil, jsonErr
	}

	trace := EmptyArrayTrace()

	for name, rawInts := range rawData {
		// Translate raw bigints into raw field elements
		rawElements := util.ToFieldElements(rawInts)
		// Add new column to the trace
		trace.AddColumn(name, rawElements, &zero)
	}

	// Done.
	return trace, nil
}
