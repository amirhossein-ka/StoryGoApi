// Code generated by ent, DO NOT EDIT.

package ent

import (
	"StoryGoAPI/ent/predicate"
	"StoryGoAPI/ent/story"
	"StoryGoAPI/ent/user"
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// StoryQuery is the builder for querying Story entities.
type StoryQuery struct {
	config
	ctx          *QueryContext
	order        []story.OrderOption
	inters       []Interceptor
	predicates   []predicate.Story
	withPostedby *UserQuery
	withFKs      bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the StoryQuery builder.
func (sq *StoryQuery) Where(ps ...predicate.Story) *StoryQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit the number of records to be returned by this query.
func (sq *StoryQuery) Limit(limit int) *StoryQuery {
	sq.ctx.Limit = &limit
	return sq
}

// Offset to start from.
func (sq *StoryQuery) Offset(offset int) *StoryQuery {
	sq.ctx.Offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *StoryQuery) Unique(unique bool) *StoryQuery {
	sq.ctx.Unique = &unique
	return sq
}

// Order specifies how the records should be ordered.
func (sq *StoryQuery) Order(o ...story.OrderOption) *StoryQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// QueryPostedby chains the current query on the "postedby" edge.
func (sq *StoryQuery) QueryPostedby() *UserQuery {
	query := (&UserClient{config: sq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(story.Table, story.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, story.PostedbyTable, story.PostedbyColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Story entity from the query.
// Returns a *NotFoundError when no Story was found.
func (sq *StoryQuery) First(ctx context.Context) (*Story, error) {
	nodes, err := sq.Limit(1).All(setContextOp(ctx, sq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{story.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *StoryQuery) FirstX(ctx context.Context) *Story {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Story ID from the query.
// Returns a *NotFoundError when no Story ID was found.
func (sq *StoryQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(1).IDs(setContextOp(ctx, sq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{story.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *StoryQuery) FirstIDX(ctx context.Context) int {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Story entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Story entity is found.
// Returns a *NotFoundError when no Story entities are found.
func (sq *StoryQuery) Only(ctx context.Context) (*Story, error) {
	nodes, err := sq.Limit(2).All(setContextOp(ctx, sq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{story.Label}
	default:
		return nil, &NotSingularError{story.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *StoryQuery) OnlyX(ctx context.Context) *Story {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Story ID in the query.
// Returns a *NotSingularError when more than one Story ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *StoryQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(2).IDs(setContextOp(ctx, sq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{story.Label}
	default:
		err = &NotSingularError{story.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *StoryQuery) OnlyIDX(ctx context.Context) int {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Stories.
func (sq *StoryQuery) All(ctx context.Context) ([]*Story, error) {
	ctx = setContextOp(ctx, sq.ctx, "All")
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Story, *StoryQuery]()
	return withInterceptors[[]*Story](ctx, sq, qr, sq.inters)
}

// AllX is like All, but panics if an error occurs.
func (sq *StoryQuery) AllX(ctx context.Context) []*Story {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Story IDs.
func (sq *StoryQuery) IDs(ctx context.Context) (ids []int, err error) {
	if sq.ctx.Unique == nil && sq.path != nil {
		sq.Unique(true)
	}
	ctx = setContextOp(ctx, sq.ctx, "IDs")
	if err = sq.Select(story.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *StoryQuery) IDsX(ctx context.Context) []int {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *StoryQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, sq.ctx, "Count")
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, sq, querierCount[*StoryQuery](), sq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (sq *StoryQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *StoryQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, sq.ctx, "Exist")
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *StoryQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the StoryQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *StoryQuery) Clone() *StoryQuery {
	if sq == nil {
		return nil
	}
	return &StoryQuery{
		config:       sq.config,
		ctx:          sq.ctx.Clone(),
		order:        append([]story.OrderOption{}, sq.order...),
		inters:       append([]Interceptor{}, sq.inters...),
		predicates:   append([]predicate.Story{}, sq.predicates...),
		withPostedby: sq.withPostedby.Clone(),
		// clone intermediate query.
		sql:  sq.sql.Clone(),
		path: sq.path,
	}
}

// WithPostedby tells the query-builder to eager-load the nodes that are connected to
// the "postedby" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *StoryQuery) WithPostedby(opts ...func(*UserQuery)) *StoryQuery {
	query := (&UserClient{config: sq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	sq.withPostedby = query
	return sq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		StoryName string `json:"storyName,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Story.Query().
//		GroupBy(story.FieldStoryName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (sq *StoryQuery) GroupBy(field string, fields ...string) *StoryGroupBy {
	sq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &StoryGroupBy{build: sq}
	grbuild.flds = &sq.ctx.Fields
	grbuild.label = story.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		StoryName string `json:"storyName,omitempty"`
//	}
//
//	client.Story.Query().
//		Select(story.FieldStoryName).
//		Scan(ctx, &v)
func (sq *StoryQuery) Select(fields ...string) *StorySelect {
	sq.ctx.Fields = append(sq.ctx.Fields, fields...)
	sbuild := &StorySelect{StoryQuery: sq}
	sbuild.label = story.Label
	sbuild.flds, sbuild.scan = &sq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a StorySelect configured with the given aggregations.
func (sq *StoryQuery) Aggregate(fns ...AggregateFunc) *StorySelect {
	return sq.Select().Aggregate(fns...)
}

func (sq *StoryQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range sq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, sq); err != nil {
				return err
			}
		}
	}
	for _, f := range sq.ctx.Fields {
		if !story.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *StoryQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Story, error) {
	var (
		nodes       = []*Story{}
		withFKs     = sq.withFKs
		_spec       = sq.querySpec()
		loadedTypes = [1]bool{
			sq.withPostedby != nil,
		}
	)
	if sq.withPostedby != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, story.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Story).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Story{config: sq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := sq.withPostedby; query != nil {
		if err := sq.loadPostedby(ctx, query, nodes, nil,
			func(n *Story, e *User) { n.Edges.Postedby = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (sq *StoryQuery) loadPostedby(ctx context.Context, query *UserQuery, nodes []*Story, init func(*Story), assign func(*Story, *User)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Story)
	for i := range nodes {
		if nodes[i].user_posted == nil {
			continue
		}
		fk := *nodes[i].user_posted
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_posted" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (sq *StoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	_spec.Node.Columns = sq.ctx.Fields
	if len(sq.ctx.Fields) > 0 {
		_spec.Unique = sq.ctx.Unique != nil && *sq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *StoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(story.Table, story.Columns, sqlgraph.NewFieldSpec(story.FieldID, field.TypeInt))
	_spec.From = sq.sql
	if unique := sq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if sq.path != nil {
		_spec.Unique = true
	}
	if fields := sq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, story.FieldID)
		for i := range fields {
			if fields[i] != story.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *StoryQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(story.Table)
	columns := sq.ctx.Fields
	if len(columns) == 0 {
		columns = story.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.ctx.Unique != nil && *sq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// StoryGroupBy is the group-by builder for Story entities.
type StoryGroupBy struct {
	selector
	build *StoryQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *StoryGroupBy) Aggregate(fns ...AggregateFunc) *StoryGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the selector query and scans the result into the given value.
func (sgb *StoryGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sgb.build.ctx, "GroupBy")
	if err := sgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*StoryQuery, *StoryGroupBy](ctx, sgb.build, sgb, sgb.build.inters, v)
}

func (sgb *StoryGroupBy) sqlScan(ctx context.Context, root *StoryQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sgb.flds)+len(sgb.fns))
		for _, f := range *sgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// StorySelect is the builder for selecting fields of Story entities.
type StorySelect struct {
	*StoryQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ss *StorySelect) Aggregate(fns ...AggregateFunc) *StorySelect {
	ss.fns = append(ss.fns, fns...)
	return ss
}

// Scan applies the selector query and scans the result into the given value.
func (ss *StorySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ss.ctx, "Select")
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*StoryQuery, *StorySelect](ctx, ss.StoryQuery, ss, ss.inters, v)
}

func (ss *StorySelect) sqlScan(ctx context.Context, root *StoryQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ss.fns))
	for _, fn := range ss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}