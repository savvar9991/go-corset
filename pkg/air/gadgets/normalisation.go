package gadgets

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/go-corset/pkg/air"
	sc "github.com/consensys/go-corset/pkg/schema"
	"github.com/consensys/go-corset/pkg/schema/assignment"
	tr "github.com/consensys/go-corset/pkg/trace"
	"github.com/consensys/go-corset/pkg/util"
)

// Normalise constructs an expression representing the normalised value of e.
// That is, an expression which is 0 when e is 0, and 1 when e is non-zero.
// This is done by introducing a computed column to hold the (pseudo)
// mutliplicative inverse of e.
func Normalise(e air.Expr, schema *air.Schema) air.Expr {
	// Construct pseudo multiplicative inverse of e.
	ie := ApplyPseudoInverseGadget(e, schema)
	// Return e * e⁻¹.
	return e.Mul(ie)
}

// ApplyPseudoInverseGadget constructs an expression representing the
// (pseudo) multiplicative inverse of another expression.  Since this cannot be computed
// directly using arithmetic constraints, it is done by adding a new computed
// column which holds the multiplicative inverse.  Constraints are also added to
// ensure it really holds the inverted value.
func ApplyPseudoInverseGadget(e air.Expr, schema *air.Schema) air.Expr {
	// Determine enclosing module.
	ctx := e.Context(schema)
	// Sanity check
	if ctx.IsVoid() || ctx.IsConflicted() {
		panic("conflicting (or void) context")
	}
	// Construct inverse computation
	ie := &Inverse{Expr: e}
	// Determine computed column name
	name := ie.String()
	// Look up column
	index, ok := sc.ColumnIndexOf(schema, ctx.Module(), name)
	// Add new column (if it does not already exist)
	if !ok {
		// Add computed column
		index = schema.AddAssignment(assignment.NewComputedColumn(ctx, name, ie))
	}

	// Construct 1/e
	inv_e := air.NewColumnAccess(index, 0)
	// Construct e/e
	e_inv_e := e.Mul(inv_e)
	// Construct 1 == e/e
	one_e_e := air.NewConst64(1).Equate(e_inv_e)
	// Construct (e != 0) ==> (1 == e/e)
	e_implies_one_e_e := e.Mul(one_e_e)
	// Construct (1/e != 0) ==> (1 == e/e)
	inv_e_implies_one_e_e := inv_e.Mul(one_e_e)
	// Ensure (e != 0) ==> (1 == e/e)
	l_name := fmt.Sprintf("[%s <=]", ie.String())
	schema.AddVanishingConstraint(l_name, ctx, nil, e_implies_one_e_e)
	// Ensure (e/e != 0) ==> (1 == e/e)
	r_name := fmt.Sprintf("[%s =>]", ie.String())
	schema.AddVanishingConstraint(r_name, ctx, nil, inv_e_implies_one_e_e)
	// Done
	return air.NewColumnAccess(index, 0)
}

// Inverse represents a computation which computes the multiplicative
// inverse of a given AIR expression.
type Inverse struct{ Expr air.Expr }

// EvalAt computes the multiplicative inverse of a given expression at a given
// row in the table.
func (e *Inverse) EvalAt(k int, tbl tr.Trace) *fr.Element {
	inv := new(fr.Element)
	val := e.Expr.EvalAt(k, tbl)
	// Go syntax huh?
	return inv.Inverse(val)
}

// Bounds returns max shift in either the negative (left) or positive
// direction (right).
func (e *Inverse) Bounds() util.Bounds { return e.Expr.Bounds() }

// Context determines the evaluation context (i.e. enclosing module) for this
// expression.
func (e *Inverse) Context(schema sc.Schema) tr.Context {
	return e.Expr.Context(schema)
}

// RequiredColumns returns the set of columns on which this term depends.
// That is, columns whose values may be accessed when evaluating this term
// on a given trace.
func (e *Inverse) RequiredColumns() *util.SortedSet[uint] {
	return e.Expr.RequiredColumns()
}

func (e *Inverse) String() string {
	return fmt.Sprintf("(inv %s)", e.Expr)
}
