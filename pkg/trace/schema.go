package trace

// Describes the permitted "layout" of a given trace.  That includes
// identifying the required columns and the set of constraints which
// must hold over the trace.  Columns can be either data columns, or
// computed columns.  A data column is one whose values are expected
// to be provided by the user, whilst computed columns are derivatives
// whose values can be computed from the other columns of the trace.
// A trace of data values is said to be "well-formed" with respect to
// a schema if: (1) every data column in the schema exists in the
// trace; (2) every constraint in the schema holds for the trace.
type Schema[C Column, R Constraint]  struct {
	// Column array (either data or computed).  Columns are stored
	// such that the dependencies of a column always come before
	// that column (i.e. have a lower index).  Thus, data columns
	// always precede computed columns, etc.
	columns []C
	// Constaint array.  For a trace of values to be well-formed
	// with respect to this schema, each constraint must hold.
	constraints []R
}

func EmptySchema[C Column, R Constraint]() *Schema[C,R] {
	p := new(Schema[C,R])
	// Initially empty columns
	p.columns = make([]C,0)
	// Initially empty constraints
	p.constraints = make([]R,0)
	// Initialise height as 0
	return p
}

// Construct a new Schema initialised with a given set of columns
// and constraints.
func NewSchema[C Column, R Constraint](columns []C, constraints []R) *Schema[C,R] {
	p := new(Schema[C,R])
	p.columns = columns
	p.constraints = constraints
	//
	return p
}

// Check whether a given schema has a given column.
func (p *Schema[C, R]) HasColumn(name string) bool {
	for _,c := range p.columns {
		if c.Name() == name {
			return true
		}
	}
	return false
}

// Return the set of columns (data or computed) which are required by
// this schema.
func (p *Schema[C, R]) Columns() []C {
	return p.columns
}

// Return the set of constraints required by this schema.
func (p *Schema[C, R]) Constraints() []R {
	return p.constraints
}

// Append a new constraint onto the schema.
func (p *Schema[C, R]) AddConstraint(constraint R) {
	p.constraints = append(p.constraints,constraint)
}

// Append a new column onto the schema.
func (p *Schema[C, R]) AddColumn(column C) {
	p.columns = append(p.columns,column)
}
