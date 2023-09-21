// Code generated by ent, DO NOT EDIT.

package ent

import (
	"StoryGoAPI/ent/guestuser"
	"StoryGoAPI/ent/predicate"
	"StoryGoAPI/ent/user"
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// GuestUserQuery is the builder for querying GuestUser entities.
type GuestUserQuery struct {
	config
	ctx          *QueryContext
	order        []guestuser.OrderOption
	inters       []Interceptor
	predicates   []predicate.GuestUser
	withFollowed *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GuestUserQuery builder.
func (guq *GuestUserQuery) Where(ps ...predicate.GuestUser) *GuestUserQuery {
	guq.predicates = append(guq.predicates, ps...)
	return guq
}

// Limit the number of records to be returned by this query.
func (guq *GuestUserQuery) Limit(limit int) *GuestUserQuery {
	guq.ctx.Limit = &limit
	return guq
}

// Offset to start from.
func (guq *GuestUserQuery) Offset(offset int) *GuestUserQuery {
	guq.ctx.Offset = &offset
	return guq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (guq *GuestUserQuery) Unique(unique bool) *GuestUserQuery {
	guq.ctx.Unique = &unique
	return guq
}

// Order specifies how the records should be ordered.
func (guq *GuestUserQuery) Order(o ...guestuser.OrderOption) *GuestUserQuery {
	guq.order = append(guq.order, o...)
	return guq
}

// QueryFollowed chains the current query on the "followed" edge.
func (guq *GuestUserQuery) QueryFollowed() *UserQuery {
	query := (&UserClient{config: guq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := guq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := guq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guestuser.Table, guestuser.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, guestuser.FollowedTable, guestuser.FollowedPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(guq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first GuestUser entity from the query.
// Returns a *NotFoundError when no GuestUser was found.
func (guq *GuestUserQuery) First(ctx context.Context) (*GuestUser, error) {
	nodes, err := guq.Limit(1).All(setContextOp(ctx, guq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{guestuser.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (guq *GuestUserQuery) FirstX(ctx context.Context) *GuestUser {
	node, err := guq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first GuestUser ID from the query.
// Returns a *NotFoundError when no GuestUser ID was found.
func (guq *GuestUserQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = guq.Limit(1).IDs(setContextOp(ctx, guq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{guestuser.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (guq *GuestUserQuery) FirstIDX(ctx context.Context) int {
	id, err := guq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single GuestUser entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GuestUser entity is found.
// Returns a *NotFoundError when no GuestUser entities are found.
func (guq *GuestUserQuery) Only(ctx context.Context) (*GuestUser, error) {
	nodes, err := guq.Limit(2).All(setContextOp(ctx, guq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{guestuser.Label}
	default:
		return nil, &NotSingularError{guestuser.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (guq *GuestUserQuery) OnlyX(ctx context.Context) *GuestUser {
	node, err := guq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only GuestUser ID in the query.
// Returns a *NotSingularError when more than one GuestUser ID is found.
// Returns a *NotFoundError when no entities are found.
func (guq *GuestUserQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = guq.Limit(2).IDs(setContextOp(ctx, guq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{guestuser.Label}
	default:
		err = &NotSingularError{guestuser.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (guq *GuestUserQuery) OnlyIDX(ctx context.Context) int {
	id, err := guq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of GuestUsers.
func (guq *GuestUserQuery) All(ctx context.Context) ([]*GuestUser, error) {
	ctx = setContextOp(ctx, guq.ctx, "All")
	if err := guq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*GuestUser, *GuestUserQuery]()
	return withInterceptors[[]*GuestUser](ctx, guq, qr, guq.inters)
}

// AllX is like All, but panics if an error occurs.
func (guq *GuestUserQuery) AllX(ctx context.Context) []*GuestUser {
	nodes, err := guq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of GuestUser IDs.
func (guq *GuestUserQuery) IDs(ctx context.Context) (ids []int, err error) {
	if guq.ctx.Unique == nil && guq.path != nil {
		guq.Unique(true)
	}
	ctx = setContextOp(ctx, guq.ctx, "IDs")
	if err = guq.Select(guestuser.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (guq *GuestUserQuery) IDsX(ctx context.Context) []int {
	ids, err := guq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (guq *GuestUserQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, guq.ctx, "Count")
	if err := guq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, guq, querierCount[*GuestUserQuery](), guq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (guq *GuestUserQuery) CountX(ctx context.Context) int {
	count, err := guq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (guq *GuestUserQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, guq.ctx, "Exist")
	switch _, err := guq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (guq *GuestUserQuery) ExistX(ctx context.Context) bool {
	exist, err := guq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GuestUserQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (guq *GuestUserQuery) Clone() *GuestUserQuery {
	if guq == nil {
		return nil
	}
	return &GuestUserQuery{
		config:       guq.config,
		ctx:          guq.ctx.Clone(),
		order:        append([]guestuser.OrderOption{}, guq.order...),
		inters:       append([]Interceptor{}, guq.inters...),
		predicates:   append([]predicate.GuestUser{}, guq.predicates...),
		withFollowed: guq.withFollowed.Clone(),
		// clone intermediate query.
		sql:  guq.sql.Clone(),
		path: guq.path,
	}
}

// WithFollowed tells the query-builder to eager-load the nodes that are connected to
// the "followed" edge. The optional arguments are used to configure the query builder of the edge.
func (guq *GuestUserQuery) WithFollowed(opts ...func(*UserQuery)) *GuestUserQuery {
	query := (&UserClient{config: guq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	guq.withFollowed = query
	return guq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Token string `json:"token,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GuestUser.Query().
//		GroupBy(guestuser.FieldToken).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (guq *GuestUserQuery) GroupBy(field string, fields ...string) *GuestUserGroupBy {
	guq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GuestUserGroupBy{build: guq}
	grbuild.flds = &guq.ctx.Fields
	grbuild.label = guestuser.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Token string `json:"token,omitempty"`
//	}
//
//	client.GuestUser.Query().
//		Select(guestuser.FieldToken).
//		Scan(ctx, &v)
func (guq *GuestUserQuery) Select(fields ...string) *GuestUserSelect {
	guq.ctx.Fields = append(guq.ctx.Fields, fields...)
	sbuild := &GuestUserSelect{GuestUserQuery: guq}
	sbuild.label = guestuser.Label
	sbuild.flds, sbuild.scan = &guq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GuestUserSelect configured with the given aggregations.
func (guq *GuestUserQuery) Aggregate(fns ...AggregateFunc) *GuestUserSelect {
	return guq.Select().Aggregate(fns...)
}

func (guq *GuestUserQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range guq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, guq); err != nil {
				return err
			}
		}
	}
	for _, f := range guq.ctx.Fields {
		if !guestuser.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if guq.path != nil {
		prev, err := guq.path(ctx)
		if err != nil {
			return err
		}
		guq.sql = prev
	}
	return nil
}

func (guq *GuestUserQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*GuestUser, error) {
	var (
		nodes       = []*GuestUser{}
		_spec       = guq.querySpec()
		loadedTypes = [1]bool{
			guq.withFollowed != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*GuestUser).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &GuestUser{config: guq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, guq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := guq.withFollowed; query != nil {
		if err := guq.loadFollowed(ctx, query, nodes,
			func(n *GuestUser) { n.Edges.Followed = []*User{} },
			func(n *GuestUser, e *User) { n.Edges.Followed = append(n.Edges.Followed, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (guq *GuestUserQuery) loadFollowed(ctx context.Context, query *UserQuery, nodes []*GuestUser, init func(*GuestUser), assign func(*GuestUser, *User)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*GuestUser)
	nids := make(map[int]map[*GuestUser]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(guestuser.FollowedTable)
		s.Join(joinT).On(s.C(user.FieldID), joinT.C(guestuser.FollowedPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(guestuser.FollowedPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(guestuser.FollowedPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*GuestUser]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*User](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "followed" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (guq *GuestUserQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := guq.querySpec()
	_spec.Node.Columns = guq.ctx.Fields
	if len(guq.ctx.Fields) > 0 {
		_spec.Unique = guq.ctx.Unique != nil && *guq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, guq.driver, _spec)
}

func (guq *GuestUserQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(guestuser.Table, guestuser.Columns, sqlgraph.NewFieldSpec(guestuser.FieldID, field.TypeInt))
	_spec.From = guq.sql
	if unique := guq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if guq.path != nil {
		_spec.Unique = true
	}
	if fields := guq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guestuser.FieldID)
		for i := range fields {
			if fields[i] != guestuser.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := guq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := guq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := guq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := guq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (guq *GuestUserQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(guq.driver.Dialect())
	t1 := builder.Table(guestuser.Table)
	columns := guq.ctx.Fields
	if len(columns) == 0 {
		columns = guestuser.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if guq.sql != nil {
		selector = guq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if guq.ctx.Unique != nil && *guq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range guq.predicates {
		p(selector)
	}
	for _, p := range guq.order {
		p(selector)
	}
	if offset := guq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := guq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// GuestUserGroupBy is the group-by builder for GuestUser entities.
type GuestUserGroupBy struct {
	selector
	build *GuestUserQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (gugb *GuestUserGroupBy) Aggregate(fns ...AggregateFunc) *GuestUserGroupBy {
	gugb.fns = append(gugb.fns, fns...)
	return gugb
}

// Scan applies the selector query and scans the result into the given value.
func (gugb *GuestUserGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gugb.build.ctx, "GroupBy")
	if err := gugb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuestUserQuery, *GuestUserGroupBy](ctx, gugb.build, gugb, gugb.build.inters, v)
}

func (gugb *GuestUserGroupBy) sqlScan(ctx context.Context, root *GuestUserQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(gugb.fns))
	for _, fn := range gugb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*gugb.flds)+len(gugb.fns))
		for _, f := range *gugb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*gugb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gugb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GuestUserSelect is the builder for selecting fields of GuestUser entities.
type GuestUserSelect struct {
	*GuestUserQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (gus *GuestUserSelect) Aggregate(fns ...AggregateFunc) *GuestUserSelect {
	gus.fns = append(gus.fns, fns...)
	return gus
}

// Scan applies the selector query and scans the result into the given value.
func (gus *GuestUserSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gus.ctx, "Select")
	if err := gus.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuestUserQuery, *GuestUserSelect](ctx, gus.GuestUserQuery, gus, gus.inters, v)
}

func (gus *GuestUserSelect) sqlScan(ctx context.Context, root *GuestUserQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(gus.fns))
	for _, fn := range gus.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*gus.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gus.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}