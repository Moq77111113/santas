// Code generated by ent, DO NOT EDIT.

package exclusion

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the exclusion type in the database.
	Label = "exclusion"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGroupID holds the string denoting the group_id field in the database.
	FieldGroupID = "group_id"
	// FieldMemberID holds the string denoting the member_id field in the database.
	FieldMemberID = "member_id"
	// FieldExcludeID holds the string denoting the exclude_id field in the database.
	FieldExcludeID = "exclude_id"
	// EdgeGroup holds the string denoting the group edge name in mutations.
	EdgeGroup = "group"
	// EdgeMember holds the string denoting the member edge name in mutations.
	EdgeMember = "member"
	// EdgeExclude holds the string denoting the exclude edge name in mutations.
	EdgeExclude = "exclude"
	// Table holds the table name of the exclusion in the database.
	Table = "exclusions"
	// GroupTable is the table that holds the group relation/edge.
	GroupTable = "exclusions"
	// GroupInverseTable is the table name for the Group entity.
	// It exists in this package in order to avoid circular dependency with the "group" package.
	GroupInverseTable = "groups"
	// GroupColumn is the table column denoting the group relation/edge.
	GroupColumn = "group_id"
	// MemberTable is the table that holds the member relation/edge.
	MemberTable = "exclusions"
	// MemberInverseTable is the table name for the Member entity.
	// It exists in this package in order to avoid circular dependency with the "member" package.
	MemberInverseTable = "members"
	// MemberColumn is the table column denoting the member relation/edge.
	MemberColumn = "member_id"
	// ExcludeTable is the table that holds the exclude relation/edge.
	ExcludeTable = "exclusions"
	// ExcludeInverseTable is the table name for the Member entity.
	// It exists in this package in order to avoid circular dependency with the "member" package.
	ExcludeInverseTable = "members"
	// ExcludeColumn is the table column denoting the exclude relation/edge.
	ExcludeColumn = "exclude_id"
)

// Columns holds all SQL columns for exclusion fields.
var Columns = []string{
	FieldID,
	FieldGroupID,
	FieldMemberID,
	FieldExcludeID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Exclusion queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByGroupID orders the results by the group_id field.
func ByGroupID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGroupID, opts...).ToFunc()
}

// ByMemberID orders the results by the member_id field.
func ByMemberID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMemberID, opts...).ToFunc()
}

// ByExcludeID orders the results by the exclude_id field.
func ByExcludeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExcludeID, opts...).ToFunc()
}

// ByGroupField orders the results by group field.
func ByGroupField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGroupStep(), sql.OrderByField(field, opts...))
	}
}

// ByMemberField orders the results by member field.
func ByMemberField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMemberStep(), sql.OrderByField(field, opts...))
	}
}

// ByExcludeField orders the results by exclude field.
func ByExcludeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newExcludeStep(), sql.OrderByField(field, opts...))
	}
}
func newGroupStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GroupInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, GroupTable, GroupColumn),
	)
}
func newMemberStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MemberInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, MemberTable, MemberColumn),
	)
}
func newExcludeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ExcludeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ExcludeTable, ExcludeColumn),
	)
}
