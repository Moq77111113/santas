// Code generated by ent, DO NOT EDIT.

package group

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the group type in the database.
	Label = "group"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeMembers holds the string denoting the members edge name in mutations.
	EdgeMembers = "members"
	// Table holds the table name of the group in the database.
	Table = "groups"
	// MembersTable is the table that holds the members relation/edge. The primary key declared below.
	MembersTable = "group_members"
	// MembersInverseTable is the table name for the Member entity.
	// It exists in this package in order to avoid circular dependency with the "member" package.
	MembersInverseTable = "members"
)

// Columns holds all SQL columns for group fields.
var Columns = []string{
	FieldID,
	FieldName,
}

var (
	// MembersPrimaryKey and MembersColumn2 are the table columns denoting the
	// primary key for the members relation (M2M).
	MembersPrimaryKey = []string{"group_id", "member_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Group queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByMembersCount orders the results by members count.
func ByMembersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMembersStep(), opts...)
	}
}

// ByMembers orders the results by members terms.
func ByMembers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMembersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newMembersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, MembersTable, MembersPrimaryKey...),
	)
}
