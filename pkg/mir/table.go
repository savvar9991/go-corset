package mir

import (
	"github.com/consensys/go-corset/pkg/air"
	"github.com/consensys/go-corset/pkg/trace"
)

// ============================================================================
// Table
// ============================================================================

type Table = trace.Table[Column,Constraint]

type Column = trace.Column

// For now, all constraints are vanishing constraints.
type Constraint = *trace.VanishingConstraint[Expr]

// Lower (or refine) an MIR table into an AIR table.  That means
// lowering all the columns and constraints, whilst adding additional
// columns / constraints as necessary to preserve the original
// semantics.
func LowerToAir(mirTbl Table, airTbl air.Table) {
	for _,col := range mirTbl.Columns() {
		airTbl.AddColumn(col)
	}
	for _,c := range mirTbl.Constraints() {
		// FIXME: this is broken because its currently
		// assuming that an AirConstraint is always a
		// VanishingConstraint.  Eventually this will not be
		// true.
		air_expr := c.Expr.LowerTo(airTbl)
		airTbl.AddConstraint(&trace.VanishingConstraint[air.Expr]{Handle: c.Handle,Expr: air_expr})
	}
}
