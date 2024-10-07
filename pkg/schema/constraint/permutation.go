package constraint

import (
	"fmt"

	"github.com/consensys/go-corset/pkg/schema"
	sc "github.com/consensys/go-corset/pkg/schema"
	"github.com/consensys/go-corset/pkg/sexp"
	"github.com/consensys/go-corset/pkg/trace"
	"github.com/consensys/go-corset/pkg/util"
)

// PermutationFailure provides structural information about a failing permutation constraint.
type PermutationFailure struct {
	msg string
}

// Message provides a suitable error message
func (p *PermutationFailure) Message() string {
	return p.msg
}

func (p *PermutationFailure) String() string {
	return p.msg
}

// PermutationConstraint declares a constraint that one (or more) columns are a permutation
// of another.
type PermutationConstraint struct {
	// The target columns
	targets []uint
	// The source columns
	sources []uint
}

// NewPermutationConstraint creates a new permutation
func NewPermutationConstraint(targets []uint, sources []uint) *PermutationConstraint {
	if len(targets) != len(sources) {
		panic("differeng number of target / source permutation columns")
	}

	return &PermutationConstraint{targets, sources}
}

// RequiredSpillage returns the minimum amount of spillage required to ensure
// valid traces are accepted in the presence of arbitrary padding.
func (p *PermutationConstraint) RequiredSpillage() uint {
	return uint(0)
}

// Accepts checks whether a permutation holds between the source and
// target columns.
func (p *PermutationConstraint) Accepts(tr trace.Trace) schema.Failure {
	// Slice out data
	src := sliceColumns(p.sources, tr)
	dst := sliceColumns(p.targets, tr)
	// Sanity check whether column exists
	if util.ArePermutationOf(dst, src) {
		// Success
		return nil
	}
	// Prepare suitable error message
	src_names := trace.QualifiedColumnNamesToCommaSeparatedString(p.sources, tr)
	dst_names := trace.QualifiedColumnNamesToCommaSeparatedString(p.targets, tr)
	//
	msg := fmt.Sprintf("Target columns (%s) not permutation of source columns (%s)",
		dst_names, src_names)
	// Done
	return &PermutationFailure{msg}
}

// Lisp converts this schema element into a simple S-Expression, for example
// so it can be printed.
func (p *PermutationConstraint) Lisp(schema sc.Schema) sexp.SExp {
	targets := sexp.EmptyList()
	sources := sexp.EmptyList()

	for _, tid := range p.targets {
		target := schema.Columns().Nth(tid)
		targets.Append(sexp.NewSymbol(target.QualifiedName(schema)))
	}

	for _, sid := range p.sources {
		source := schema.Columns().Nth(sid)
		sources.Append(sexp.NewSymbol(source.QualifiedName(schema)))
	}

	return sexp.NewList([]sexp.SExp{
		sexp.NewSymbol("permutation"),
		targets,
		sources,
	})
}

// Targets returns the indices of the columns composing the "left" table of the
// permutation.
func (p *PermutationConstraint) Targets() []uint {
	return p.targets
}

// Sources returns the indices of the columns composing the "right" table of the
// permutation.
func (p *PermutationConstraint) Sources() []uint {
	return p.sources
}

func sliceColumns(columns []uint, tr trace.Trace) []util.FrArray {
	// Allocate return array
	cols := make([]util.FrArray, len(columns))
	// Slice out the data
	for i, n := range columns {
		nth := tr.Column(n)
		// Copy over
		cols[i] = nth.Data()
	}
	// Done
	return cols
}
