package sml

// This file is the implementation of the "SML" language. Aka Scheme Markup Language

import (
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	log "github.com/maagodata/maago-commons/logger"
)

const (
	skwAuxiliaryColumn keyword = iota
	skwColumn
	skwCombine
	skwCriteria
	skwDataset
	skwDescription
	skwEvaluation
	skwFormat
	skwGranulariy
	skwJoin
	skwLabel
	skwName
	skwOn
	skwOperator
	skwProject
	skwScheme
	skwSQL
	skwTable
	skwTerm
	skwThen
	skwTransformer
	skwTransformerName
	skwType
	skwValue
)

const (
	decoratorOneToMany = "oneToMany"
	decoratorManyToOne = "manyToOne"
)

var smlKeywords = map[string]keyword{
	"auxiliary_column": skwAuxiliaryColumn,
	"column":           skwColumn,
	"combine":          skwCombine,
	"criteria":         skwCriteria,
	"dataset":          skwDataset,
	"description":      skwDescription,
	"evaluation":       skwEvaluation,
	"format":           skwFormat,
	"granularity":      skwGranulariy,
	"join":             skwJoin,
	"label":            skwLabel,
	"name":             skwName,
	"on":               skwOn,
	"operator":         skwOperator,
	"project":          skwProject,
	"scheme":           skwScheme,
	"sql":              skwSQL,
	"table":            skwTable,
	"term":             skwTerm,
	"then":             skwThen,
	"transformer":      skwTransformer,
	"transformer_name": skwTransformerName,
	"type":             skwType,
	"value":            skwValue,
}

var smlKeyword2string = makeKeywordToString(smlKeywords, []string{})

// define SML data types
const (
	TypeBool     = "bool"
	TypeDatetime = "datetime"
	TypeFloat    = "float"
	TypeInt      = "int"
	TypeString   = "string"
	TypeTime     = "time"
	TypeVerbatim = "verbatim"
)

// DataTypes ...
var DataTypes = map[string]bool{
	TypeString:   true,
	TypeVerbatim: true,
	TypeInt:      true,
	TypeFloat:    true,
	TypeDatetime: true,
	TypeBool:     true,
	TypeTime:     true,
}

var smlKwTree = makeKwTree(smlKeywords, `
project  
  dataset
    table
    join
      on
      sql
      then
        on
        sql
    label
  table
    description
    sql
    label
    column
      sql
      label
      format
      type
      transformer
  scheme
    name
    description
    criteria
      column
      table
      operator
      value
      combine
        term
          column
          table
          operator
          value
          granularity
          combine
        column
        table
        operator
        value
        granularity
    evaluation
`)

//------------ parser -------------

// Project ...
type Project struct {
	Name       string
	TableNames []string
	Tables     map[string]Table
	Datasets   map[string]Dataset
	Schemes    map[string]Scheme
}

// TODO: Add example
type Table struct {
	Name        string
	DataSource  string
	SQL         string
	ColumnNames []string
	Columns     map[string]Column
	Label       string
	Pk          []string
	Description string
}

func makeTable(name string) Table {
	return Table{
		Name:    name,
		Columns: make(map[string]Column),
	}
}

// Dataset ...
type Dataset struct {
	Name       string
	Label      string
	TableNames []string
	Tables     map[string]Table
	Joins      []Join
}

// Column ...
type Column struct {
	Name            string
	SQL             string
	Label           string
	Format          string
	Type            string
	Description     string
	Transformer     string
	TransformerName string
}

// OperatorTableSQL ...
type OperatorTableSQL struct {
	Operator string
	Table    string
	SQL      string
}

// Join ...
type Join struct {
	OneTable   string   // != "" => oneToMany / manyToOne
	OneColumns []string // != nil => oneToMany / manyToOne
	OTS        []OperatorTableSQL
}

var joinOperators = map[string]string{
	"<->":              "<->",
	"inner":            "<->",
	"inner join":       "<->",
	"join":             "<->",
	"-->":              "-->",
	"left outer join":  "-->",
	"left join":        "-->",
	"left":             "-->",
	"<--":              "<--",
	"right outer join": "<--",
	"right join":       "<--",
	"right":            "<--",
	">-<":              ">-<",
	"full outer join":  ">-<",
	"full join":        ">-<",
	"full":             ">-<",
}

type CriteriaSimpleTerm struct {
	Column      string
	Table       string
	Operator    string
	Granularity string
	Value       []interface{}
}

type CriteriaCombinationTerm struct {
	LogicalOperator string
	Terms           []CriteriaTerm
}

type CriteriaTerm struct {
	// A criteria can either have a combination term
	// Or a simple term
	Name            string
	SimpleTerm      *CriteriaSimpleTerm
	CombinationTerm *CriteriaCombinationTerm
}

type Scheme struct {
	Name          string
	Label         string
	Description   string
	Criteria      map[string]CriteriaTerm
	CriteriaNames []string
	Evaluation    string // similar to a filter string
}

var datetimeRegex = regexp.MustCompile(`^\s*((?:19|20)[0-9]{2})\s*-\s*(0[1-9]|1[0-2])\s*-\s*(0[1-9]|[1-2][0-9]|3[0-1]).*$`)

func copyTable(dst *Table, src Table) {
	// ignore Name
	if src.SQL != "" {
		dst.SQL = src.SQL
	}
	if src.Label != "" {
		dst.Label = src.Label
	}
	for k := range src.Columns {
		dst.Columns[k] = src.Columns[k] // assumption: Column has a no field with reference type
	}
	if dst.ColumnNames == nil {
		dst.ColumnNames = make([]string, 0)
	}

	dst.ColumnNames = append(dst.ColumnNames, src.ColumnNames...)
	dst.Pk = append(dst.Pk, src.Pk...)
}

// SML parser object contains parse methods and also
// all the errors encountered during parsing the model
type smlParser struct {
	maxErrors   int
	parseErrors []ParseError
}

const maxErrors = 100

// newSMLParser returns a new instance of a SML Parser
func newSMLParser() *smlParser {
	parser := new(smlParser)
	parser.maxErrors = maxErrors
	return parser
}

func (parser *smlParser) didParseFail() bool {
	return len(parser.parseErrors) > 0
}

func (parser *smlParser) getFirstParseError() ParseError {
	if len(parser.parseErrors) > 0 {
		return parser.parseErrors[0]
	}
	return ParseError{}
}

func (parser *smlParser) appendParseError(parseErrors ...ParseError) {
	emptySlots := parser.maxErrors - len(parser.parseErrors)
	if emptySlots > 0 {
		if len(parseErrors) < emptySlots {
			parser.parseErrors = append(parser.parseErrors, parseErrors...)
		} else {
			// Too many errors
			parser.parseErrors = append(parser.parseErrors, parseErrors[:emptySlots]...)
			panic(parser.parseErrors)
		}
	} else {
		// Too many errors
		panic(parser.parseErrors)
	}
}

func (parser *smlParser) parseProject(g generic) (w Project) {
	// TODO merge fragments to same table (or dataset) together
	if g.kw != skwProject {
		// Cannot parse further
		parser.appendParseError(ParseError{"value of the keyword must be project", -1})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"\"project\" has to have a name", g.line})
	}

	w.Name = g.value
	w.TableNames = make([]string, 0)
	w.Tables = make(map[string]Table)
	w.Datasets = make(map[string]Dataset)
	w.Schemes = make(map[string]Scheme)
	schemeFound := false
	for kw, list := range g.children {
		switch kw {
		case skwTable:
			for _, child := range list {
				table := parser.parseTable(child)
				if _, isPresent := w.Tables[table.Name]; !isPresent {
					w.Tables[table.Name] = table // if table definition is repeated (case sensitive names), we take the first one
					w.TableNames = append(w.TableNames, table.Name)
				}
			}
		case skwDataset:
			for _, child := range list {
				dataset := parser.parseDataset(child)
				w.Datasets[dataset.Name] = dataset
			}
		case skwScheme:
			schemeFound = true
			// parse schemes in the second pass, after the tables are populated
			// because the column's data type may be needed while dealing with values
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"project\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
		}
	}
	if schemeFound {
		schemes := g.children[skwScheme]
		for _, s := range schemes {
			scheme := parser.parseScheme(s, w.Tables)
			w.Schemes[scheme.Name] = scheme
		}
	}
	// validate table names and joins in datasets
	for i := range w.Datasets {
		tablesInDataset := make(map[string]bool)
		for _, table := range w.Datasets[i].TableNames {
			if _, isPresent := w.Tables[table]; !isPresent {
				parser.appendParseError(ParseError{fmt.Sprintf("table \"%s\" in dataset \"%s\" not defined in project \"%s\"", table, w.Datasets[i].Name, w.Name), g.line})
			}
			tablesInDataset[table] = true
		}
		for j := range w.Datasets[i].Joins {
			for _, ots := range w.Datasets[i].Joins[j].OTS {
				if _, isPresent := tablesInDataset[ots.Table]; !isPresent {
					parser.appendParseError(ParseError{fmt.Sprintf("table \"%s\" in join of dataset \"%s\" not included in dataset", ots.Table, w.Datasets[i].Name), g.line})
				}
			}
		}
	}
	// for each table, validate Pk values
	for _, table := range w.Tables {
		columns := make(map[string]bool)
		for column := range table.Columns {
			columns[strings.ToLower(column)] = true
		}
		for _, name := range table.Pk {
			if _, isPresent := columns[strings.ToLower(name)]; !isPresent {
				parser.appendParseError(ParseError{fmt.Sprintf("pk \"%s\" in table \"%s\" is not a known column", name, table.Name), g.line})
			}
		}
	}
	return
}

func computeDerivedAttributesInProject(w *Project) (parseErrors []ParseError) {
	// for each table T in dataset N, copy contents of T to N
	for name, dataset := range w.Datasets {
		for _, table := range dataset.TableNames {
			tmp := makeTable(w.Tables[table].Name)
			// isTableHiddenInDataset := dataset.Tables[tmp.Name].IsHidden
			copyTable(&tmp, w.Tables[table])
			// tmp.IsHidden = isTableHiddenInDataset
			dataset.Tables[tmp.Name] = tmp
		}
		w.Datasets[name] = dataset
	}
	return
}

func (parser *smlParser) parseTable(g *generic) (table Table) {
	if g.kw != skwTable {
		parser.appendParseError(ParseError{"value of the keyword must be table", g.line})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"\"table\" has to have a name", g.line})
	} else if !IsValidID(g.value) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid name for table: %s", g.value), g.line})
	}
	table.Name = g.value
	table.ColumnNames = make([]string, 0)
	table.Columns = make(map[string]Column)
	columnNames := make([]string, 0)

	for kw, list := range g.children {
		switch kw {
		case skwDescription:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"descriptions\"s for a \"table\"", g.line})
			}
			child := list[0]
			table.Description = parser.parseSMLString(child)
		case skwSQL:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"sql\"s for a \"table\"", g.line})
			}
			child := list[0]
			table.SQL = parser.parseSMLString(child)
		case skwLabel:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"label\"s for a \"table\"", g.line})
			}
			child := list[0]
			table.Label = parser.parseSMLString(child)
		case skwColumn, skwAuxiliaryColumn:
			for _, child := range list {
				column := parser.parseColumn(child, table.Name)
				if _, isPresent := table.Columns[column.Name]; isPresent {
					parser.appendParseError(ParseError{fmt.Sprintf("duplicate column \"%s\" in table \"%s\"", column.Name, table.Name), child.line})
					continue
				}
				table.Columns[column.Name] = column
				columnNames = append(columnNames, column.Name)
			}
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"table\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
		}
	}
	if table.Label == "" {
		table.Label = createSmartLabel(table.Name)
	}
	if table.SQL == "" {
		parser.appendParseError(ParseError{"\"table\" must have a  \"sql\"", -1})
	}

	// Fill the column names.
	table.ColumnNames = append(table.ColumnNames, columnNames...)

	return
}

func (parser *smlParser) parseDataset(g *generic) (n Dataset) {
	if g.kw != skwDataset {
		parser.appendParseError(ParseError{"value of the keyword must be dataset", g.line})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"\"dataset\" has to have a name", g.line})
	} else if !IsValidID(g.value) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid name for dataset: %s", g.value), g.line})
	}
	n.Name = g.value
	n.TableNames = make([]string, 0)
	n.Tables = make(map[string]Table)
	n.Joins = make([]Join, 0)
	for kw, list := range g.children {
		switch kw {
		case skwTable:
			var table Table
			for _, child := range list {
				ids := parser.parseIDList(child, false)
				n.TableNames = append(n.TableNames, ids...)
				for _, id := range ids {
					n.Tables[id] = table // dummy table, later overridden in computeDerivedAttributesInProject()
				}
			}
		case skwJoin:
			for _, child := range list {
				n.Joins = append(n.Joins, parser.parseJoin(child))
			}
		case skwLabel:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"label\"s for a \"dataset\"", g.line})
			}
			child := list[0]
			n.Label = parser.parseSMLString(child)
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"dataset\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
		}
	}
	if n.Label == "" {
		n.Label = createSmartLabel(n.Name)
	}
	return
}

func (parser *smlParser) handleQualifiedColumn(table, column *string, tableName, parentType string, parentLine, columnLine int) {
	if columnLine < 0 { // there was no 'column' line
		parser.appendParseError(ParseError{fmt.Sprintf("missing \"column\" for a \"%s\"", parentType), parentLine})
		return
	}
	if *column == "" {
		return // we have already produced an error message
	}
	var t, c string
	i := strings.Index(*column, ".")
	if i < 0 {
		c = *column
	} else {
		t, c = (*column)[:i], (*column)[i+1:]
		if t == "" {
			parser.appendParseError(ParseError{fmt.Sprintf("invalid column value: %s", *column), columnLine})
			return
		}
	}
	if (t != "" && !IsValidID(t)) || !IsValidID(c) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid column value: %s", *column), columnLine})
		return
	}
	if t == "" {
		if *table == "" {
			*table = tableName
		}
	} else if *table == "" {
		*table = t
	} else {
		parser.appendParseError(ParseError{fmt.Sprintf("\"table\" exists for a \"%s\" but \"column\" is already qualified with table", parentType), parentLine})
		return
	}
	*column = c
}

func (parser *smlParser) parseColumn(g *generic, tableName string) (c Column) {
	columnTypeExists := false
	if g.kw != skwColumn && g.kw != skwAuxiliaryColumn {
		parser.appendParseError(ParseError{"value of the keyword must be column", g.line})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"\"column\" has to have a name", g.line})
	} else if !IsValidID(g.value) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid name for column: %s", g.value), g.line})
	}
	c.Name = g.value
	for kw, list := range g.children {
		switch kw {
		case skwSQL:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"sql\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.SQL = parser.parseSMLString(child)
		case skwLabel:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"label\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.Label = parser.parseSMLString(child)
		case skwFormat:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"format\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.Format = parser.parseSMLString(child)
		case skwType:
			columnTypeExists = true
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"type\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.Type = parser.parseDatatype(child)
		case skwTransformer:
			// To be implemented
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"transformers\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.Transformer = parser.parseSMLString(child)
		case skwTransformerName:
			// To be implemented
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"transformers_name\"s for a \"column\"", g.line})
			}
			child := list[0]
			c.TransformerName = parser.parseSMLString(child)
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"column\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
		}
	}
	if !columnTypeExists {
		parser.appendParseError(ParseError{fmt.Sprintf("\"type\" is mandatory for a column, missing in column %s", c.Name), g.line})
	}
	if c.Label == "" {
		c.Label = createSmartLabel(c.Name)
	}

	return
}

func (parser *smlParser) parseJoin(g *generic) (j Join) {
	if g.kw != skwJoin {
		parser.appendParseError(ParseError{"value of the keyword must be join", g.line})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"join is not specified", g.line})
		panic(parser.parseErrors)
	}
	pieces := strings.Fields(g.value) // failure case: A<->B (no spaces)
	n := len(pieces)
	if n < 3 {
		parser.appendParseError(ParseError{fmt.Sprintf("cannot parse \"join\": %s", g.value), g.line})
		panic(parser.parseErrors)
	}
	oneToMany := normalizedOneToMany(pieces[0])
	if oneToMany != "" {
		n--
		pieces = pieces[1:]
		if n < 3 {
			parser.appendParseError(ParseError{fmt.Sprintf("cannot parse \"join\": %s", g.value), g.line})
			panic(parser.parseErrors)
		}
	}
	operator := strings.Join(pieces[1:n-1], " ")
	if op, exists := joinOperators[strings.ToLower(operator)]; !exists {
		parser.appendParseError(ParseError{fmt.Sprintf("expected operator in \"join\", got: %s", operator), g.line})
		panic(parser.parseErrors)
	} else {
		operator = op
	}
	// first table
	table := pieces[0]
	if !IsValidID(table) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid table \"%s\" in \"join\"", table), g.line})
		panic(parser.parseErrors)
	}
	j.OTS = []OperatorTableSQL{{"", table, ""}}
	// second table
	table = pieces[n-1]
	if !IsValidID(table) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid table \"%s\" in \"join\"", table), g.line})
		panic(parser.parseErrors)
	}
	ots := OperatorTableSQL{operator, table, ""}
	// sql/on + then
	var then []*generic
	for kw, list := range g.children {
		switch kw {
		case skwOn:
			fallthrough
		case skwSQL:
			if ots.SQL != "" || len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"sql/on\"s for a \"join\"", g.line})
				panic(parser.parseErrors)
			}
			child := list[0]
			ots.SQL = parser.parseSMLString(child)
			sqlTokens := strings.Fields(ots.SQL)
			if len(sqlTokens) == 0 {
				parser.appendParseError(ParseError{"no expression for a \"sql\" or \"on\" clause", child.line})
				panic(parser.parseErrors)
			}
			if oneToMany := normalizedOneToMany(sqlTokens[0]); oneToMany != "" {
				parser.appendParseError(ParseError{fmt.Sprintf("not expecting \"%s\" in \"sql\" or \"ON\" line, should be on the previous \"join\" line", oneToMany), g.line})
				panic(parser.parseErrors)
			}
		case skwThen:
			then = make([]*generic, len(list))
			for i, child := range list {
				then[i] = child
			}
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"join\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
			panic(parser.parseErrors)
		}
	}
	if ots.SQL == "" {
		parser.appendParseError(ParseError{"no \"on\" clause for a \"join\"", g.line})
		panic(parser.parseErrors)
	}
	j.OTS = append(j.OTS, ots)
	sort.Sort(genericByLine(then)) // original order of 'then's in data model
	for i := range then {
		if then[i].value == "" {
			parser.appendParseError(ParseError{"\"then\" has to have value", then[i].line})
			panic(parser.parseErrors)
		}
		pieces := strings.Fields(then[i].value)
		n := len(pieces)
		if n < 2 {
			parser.appendParseError(ParseError{fmt.Sprintf("cannot parse \"then\": %s", then[i].value), then[i].line})
			panic(parser.parseErrors)
		}
		operator := strings.Join(pieces[0:n-1], " ")
		if op, exists := joinOperators[strings.ToLower(operator)]; !exists {
			parser.appendParseError(ParseError{fmt.Sprintf("expected operator in \"then\", got: %s", operator), then[i].line})
			panic(parser.parseErrors)
		} else {
			operator = op
		}
		table := pieces[n-1]
		if !IsValidID(table) {
			parser.appendParseError(ParseError{fmt.Sprintf("invalid table \"%s\" in \"then\"", table), then[i].line})
			panic(parser.parseErrors)
		}
		ots := OperatorTableSQL{operator, table, ""}
		for kw, list := range then[i].children {
			switch kw {
			case skwOn:
				fallthrough
			case skwSQL:
				if ots.SQL != "" || len(list) > 1 {
					parser.appendParseError(ParseError{"cannot have more than one \"sql/on\"s for a \"then\"", then[i].line})
					panic(parser.parseErrors)
				}
				child := list[0]
				ots.SQL = parser.parseSMLString(child)
				sqlTokens := strings.Fields(ots.SQL)
				if len(sqlTokens) == 0 {
					parser.appendParseError(ParseError{"no expression for a \"sql\" or \"on\" clause", child.line})
					panic(parser.parseErrors)
				}
				if oneToMany := normalizedOneToMany(sqlTokens[0]); oneToMany != "" {
					parser.appendParseError(ParseError{fmt.Sprintf("not expecting \"%s\" in joins involving more than two tables", oneToMany), then[i].line})
					panic(parser.parseErrors)
				}
			default:
				parser.appendParseError(ParseError{fmt.Sprintf("\"then\" cannot contain \"%s\"", smlKeyword2string[kw]), then[i].line})
				panic(parser.parseErrors)
			}
		}
		if ots.SQL == "" {
			parser.appendParseError(ParseError{"no \"on\" clause for a \"then\"", then[i].line})
			panic(parser.parseErrors)
		}
		j.OTS = append(j.OTS, ots)
	}
	// deal with constraints on manyToOne / oneToMany
	if oneToMany != "" {
		if len(j.OTS) > 2 {
			parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" JOIN cannot connect more than 2 tables", oneToMany), g.line})
			panic(parser.parseErrors)
		}
		j.OneTable = j.OTS[0].Table
		if oneToMany == decoratorManyToOne {
			j.OneTable = j.OTS[1].Table
		}
		j.OneColumns = parser.validateOneToManyJoinCondition(j.OTS[1].SQL, j.OneTable, oneToMany, g.line)
	}
	return
}

func (parser *smlParser) parseEvaluationString(g *generic) (e string) {
	return parser.parseSMLString(g)
}

func (parser *smlParser) parseSimpleCriteriaTerm(g *generic, tables map[string]Table) (s *CriteriaSimpleTerm) {
	var tableFound, columnFound, operatorFound bool
	s = &CriteriaSimpleTerm{}
	for kw, list := range g.children {
		switch kw {
		case skwColumn:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"column\"s for a criteria term", g.line})
			}
			child := list[0]
			s.Column = parser.parseSMLString(child)
			columnFound = true
		case skwTable:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"table\"s for a criteria term", g.line})
			}
			child := list[0]
			s.Table = parser.parseSMLString(child)
			tableFound = true
		case skwOperator:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"operator\"s for a criteria term", g.line})
			}
			child := list[0]
			s.Operator = parser.parseSMLString(child)
			operatorFound = true
		case skwGranulariy:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"granularity\"s for a criteria term", g.line})
			}
			child := list[0]
			s.Granularity = parser.parseSMLString(child)
			operatorFound = true
		case skwValue:
			// Taken care of later on
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"term\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})

		}
	}
	if !tableFound {
		parser.appendParseError(ParseError{"must have exactly one \"table\" for a criteria term", g.line})
	}
	if !operatorFound {
		parser.appendParseError(ParseError{"must have exactly one \"operator\" for a criteria term", g.line})
	}
	if !columnFound {
		parser.appendParseError(ParseError{"must have exactly one \"column\" for a criteria term", g.line})
	}

	for kw, list := range g.children {
		if kw == skwValue {
			valueString := parser.parseSMLString(list[0])
			// dataType is guaranteed to be a correct type since it's checked already
			// while parsing columns, and criteria are parsed in the second pass
			dataType, err := parser.findDatatype(tables, s.Table, s.Column)
			if err != nil {
				parser.appendParseError(ParseError{err.Error(), g.line})
			} else {
				var (
					typedValues []interface{}
					values      []string
				)
				if IsStringInSlice(s.Operator, []string{"age_between", "between", "in"}) {
					// Potential bug here: will deal when it happens, trims all leading and trailing characters
					valueString = strings.Trim(valueString, "[]")
					values = strings.Split(valueString, ",")
				} else {
					values = []string{valueString}
				}

				for _, v := range values {
					v = strings.TrimSpace(v)
					typedValue, err := parser.getTypedValue(v, dataType)
					if err != nil {
						parser.appendParseError(ParseError{err.Error(), g.line})
					} else {
						typedValues = append(typedValues, typedValue)
					}
				}
				s.Value = typedValues
			}
		}
	}
	if s.Granularity != "" {
		var err error
		datatype := tables[s.Table].Columns[s.Column].Type
		if datatype != TypeDatetime {
			parser.appendParseError(ParseError{fmt.Sprintf("can not have granularity with the data type %s", datatype), g.line})
		}
		for i, v := range s.Value {
			s.Value[i], err = strconv.Atoi(v.(string))
			if err != nil {
				parser.appendParseError(ParseError{err.Error(), g.line})
			}
		}
	}
	return
}

func (parser *smlParser) getTypedValue(v, datatype string) (typedValue interface{}, err error) {
	switch datatype {
	case TypeString:
		typedValue = v
	case TypeVerbatim:
		typedValue = v
	case TypeInt:
		typedValue, err = strconv.Atoi(v)
	case TypeFloat:
		typedValue, err = strconv.ParseFloat(v, 32)
	case TypeDatetime:
		// @TODO: Validate date values
		typedValue = v
	case TypeBool:
		typedValue, err = strconv.ParseBool(v)
	default:
		err = fmt.Errorf("type %s not supported", datatype)
	}

	return
}

func (parser *smlParser) findDatatype(tables map[string]Table, table, column string) (datatype string, err error) {
	assert(tables != nil)
	assert(table != "")
	assert(column != "")

	// make sure that the table is present in the map
	if t, ok := tables[table]; !ok {
		err = fmt.Errorf("table %s does not exist", table)
	} else {
		// make sure that the column is present in the table
		if IsStringInSlice(column, t.ColumnNames) {
			c := t.Columns[column]
			datatype = c.Type
		} else {
			err = fmt.Errorf("column %s does not exist in table %s", column, table)
		}
	}

	return
}

func (parser *smlParser) parseCombinationCriteriaTerm(g *generic, tables map[string]Table) (c *CriteriaCombinationTerm) {
	c = &CriteriaCombinationTerm{}
	c.LogicalOperator = g.value
	for kw, list := range g.children {
		switch kw {
		case skwTerm:
			if len(list) != 2 {
				parser.appendParseError(ParseError{"must have exactly two \"term\"s for a combination criteria", g.line})
			}
			for _, term := range list {
				c.Terms = append(c.Terms, parser.parseCriteria(term, tables))
			}
		}
	}
	return
}

func (parser *smlParser) parseCriteria(g *generic, tables map[string]Table) (t CriteriaTerm) {
	var logicalOperatorFound bool
	// figure out if this is a simple term or a complex term
	for kw, list := range g.children {
		switch kw {
		case skwCombine:
			logicalOperatorFound = true
			t.CombinationTerm = parser.parseCombinationCriteriaTerm(list[0], tables)
		}
	}
	name := g.value
	if g.value == "" {
		name = RandomString(8)
	}
	t.Name = name
	if !logicalOperatorFound {
		t.SimpleTerm = parser.parseSimpleCriteriaTerm(g, tables)
	}
	return
}

func (parser *smlParser) parseScheme(g *generic, tables map[string]Table) (s Scheme) {
	if g.kw != skwScheme {
		parser.appendParseError(ParseError{"value of the keyword must be scheme", g.line})
		panic(parser.parseErrors)
	}
	if g.value == "" {
		parser.appendParseError(ParseError{"\"scheme\" has to have an ID", g.line})
	} else if !IsValidID(g.value) {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid name for scheme: %s", g.value), g.line})
	}
	s.Name = g.value
	criteriaFound := false
	s.Criteria = make(map[string]CriteriaTerm)
	for kw, list := range g.children {
		switch kw {
		case skwDescription:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"description\"s for a \"scheme\"", g.line})
			}
			child := list[0]
			s.Description = parser.parseSMLString(child)
		case skwLabel:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"label\"s for a \"scheme\"", g.line})
			}
			child := list[0]
			s.Label = parser.parseSMLString(child)
		case skwEvaluation:
			if len(list) > 1 {
				parser.appendParseError(ParseError{"cannot have more than one \"evaluation\"s for a \"scheme\"", g.line})
			}
			child := list[0]
			s.Evaluation = parser.parseEvaluationString(child)
		case skwCriteria:
			criteriaFound = true
			for _, child := range list {
				term := parser.parseCriteria(child, tables)
				s.Criteria[term.Name] = term
				s.CriteriaNames = append(s.CriteriaNames, term.Name)
			}
		default:
			parser.appendParseError(ParseError{fmt.Sprintf("\"scheme\" cannot contain \"%s\"", smlKeyword2string[kw]), g.line})
		}
	}
	if !criteriaFound {
		parser.appendParseError(ParseError{"no criteria terms found", g.line})
	}

	return
}

func normalizedOneToMany(word string) string {
	switch strings.ToLower(word) {
	case strings.ToLower(decoratorOneToMany):
		return decoratorOneToMany
	case strings.ToLower(decoratorManyToOne):
		return decoratorManyToOne
	}
	return ""
}

func (parser *smlParser) validateOneToManyJoinCondition(sql, oneTable, oneToMany string, linenum int) (oneColumns []string) {
	tokens := TokenizeSQL(sql)
	// the intent is to support only conjunctions of equalities
	// for now, we only support:  t1.c1 = t2.c2 AND ... with optional parentheses
	// another constraint: at least one column of 'oneTable' is involved
	for _, token := range tokens {
		if token.isReserved() && (strings.ToUpper(token.Value) == "OR" || strings.ToUpper(token.Value) == "NOT") {
			parser.appendParseError(ParseError{fmt.Sprintf("expecting simple conjunctions in %s JOIN, found \"%s\"", oneToMany, token.Value), linenum})
		}
		if token.isOp() && token.Value != "=" {
			parser.appendParseError(ParseError{fmt.Sprintf("expecting only equalities (=) in %s JOIN, found \"%s\"", oneToMany, token.Value), linenum})
		}
		if token.isName() && token.Value2 != "" && strings.ToLower(token.Value) == strings.ToLower(oneTable) {
			oneColumns = append(oneColumns, token.Value2)
		}
	}
	if len(oneColumns) == 0 {
		parser.appendParseError(ParseError{fmt.Sprintf("did not find any column of table \"%s\" participating in %s JOIN", oneTable, oneToMany), linenum})
	}
	return
}

func (parser *smlParser) parseSMLString(g *generic) string {
	if g.children != nil {
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" cannot have descendents", smlKeyword2string[g.kw]), g.line})
	}
	return g.value
}

func (parser *smlParser) parseID(g *generic) string {
	v := parser.parseSMLString(g)
	if v == "" {
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has empty name", smlKeyword2string[g.kw]), g.line})
	} else if !IsValidID(v) {
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has invalid name: \"%s\"", smlKeyword2string[g.kw], v), g.line})
	}
	return v
}

func (parser *smlParser) parseIDList(g *generic, allowEmpty bool) (ids []string) {
	s := parser.parseSMLString(g)
	if s == "" {
		if allowEmpty {
			return
		}
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has no values", smlKeyword2string[g.kw]), g.line})
	}
	ids = strings.Fields(s)
	for _, id := range ids {
		if !IsValidID(id) {
			parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has invalid value: \"%s\"", smlKeyword2string[g.kw], id), g.line})
		}
	}
	return
}

func (parser *smlParser) parseDatatype(g *generic) (s string) {
	// called for column
	// string/int/float/datetime/verbatim
	s = parser.parseSMLString(g)
	if s == "" {
		parser.appendParseError(ParseError{"empty type", g.line})
	}
	if !DataTypes[s] {
		parser.appendParseError(ParseError{fmt.Sprintf("invalid type: \"%s\"", s), g.line})
	}
	return
}

func (parser *smlParser) parseBool(g *generic) (b bool) {
	s := parser.parseSMLString(g)
	switch strings.ToLower(s) {
	case "true", "1":
		b = true
	case "false", "0":
		b = false
	case "":
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has empty value", smlKeyword2string[g.kw]), g.line})
	default:
		parser.appendParseError(ParseError{fmt.Sprintf("\"%s\" has invalid value: \"%s\"", smlKeyword2string[g.kw], s), g.line})
	}
	return
}

func isColumnExpressionSpecialAtomic(expression string) bool {
	n := len(expression)
	if n > 2 {
		if expression[0] == '`' && expression[n-1] == '`' && !strings.ContainsRune(expression[1:n-1], '`') {
			return true
		}
		if expression[0] == '[' && expression[n-1] == ']' && !strings.ContainsAny(expression[1:n-1], "[]") {
			return true
		}
	}
	return false
}

// IsColumnExpressionAtomic takes an expression and returns true if the expression is atomic
func IsColumnExpressionAtomic(expression string) bool {
	if isColumnExpressionSpecialAtomic(expression) {
		return true
	}
	return len(TokenizeSQL(expression)) == 1
}

const (
	isAggUnknown int = iota
	isAggStarted
	isAggTrue
	isAggFalse
)

var isAggFunction = map[string]bool{
	"count":          true,
	"count_if":       true,
	"count_distinct": true,
	"sum":            true,
	"max":            true,
	"min":            true,
	"avg":            true,
}

func ensureUniqueColumnNames(table Table) (parseErrors []ParseError) {
	// ensure columns names are unique within a table (case insensitive)
	names := make(map[string]bool)
	for name := range table.Columns {
		l := strings.ToLower(name)
		if names[l] {
			parseErrors = append(parseErrors, ParseError{fmt.Sprintf("more than one column with the same name (%s) in table %s", l, table.Name), -1})
			continue
		}
		names[l] = true
	}
	return
}

func ensureUniqueTableAndColumnNames(project *Project) (parseErrors []ParseError) {
	// ensure table names are unique, case-insensitive, as well as columns names within a table are unique
	names := make(map[string]bool)
	for name, table := range project.Tables {
		l := strings.ToLower(name)
		if names[l] {
			parseErrors = append(parseErrors, ParseError{fmt.Sprintf("more than one table with the same name: %s", l), -1})
			continue
		}
		names[l] = true
		errs := ensureUniqueColumnNames(table)
		if len(errs) > 0 {
			parseErrors = append(parseErrors, errs...)
			continue
		}
	}
	return
}

func ensureUniqueDatasetNames(project *Project) (parseErrors []ParseError) {
	names := make(map[string]bool)
	for name := range project.Datasets {
		l := strings.ToLower(name)
		if names[l] {
			parseErrors = append(parseErrors, ParseError{fmt.Sprintf("more than one dataset with the same name: %s", l), -1})
			continue
		}
		names[l] = true
	}
	return
}

// API ParseSML ...
// given an SML model, returns a parsed project
func ParseSML(input, projectName string) (project Project, parseErrors []ParseError) {
	defer func() {
		if r := recover(); r != nil {
			pe, ok := r.([]ParseError)
			if ok {
				parseErrors = pe
			} else {
				pe, ok := r.(ParseError)
				if ok {
					parseErrors = []ParseError{pe}
				} else {
					log.Infof("panic stacktrace:\n%s\n", string(debug.Stack()))
					parseErrors = append(parseErrors, ParseError{fmt.Sprintf("unknown error: %v", r), -1})
				}
			}
		}
	}()

	l := startLex(input, smlKeywords, nil)

	token := <-l
	switch token.kw {
	case kwError:
		parseErrors = append(parseErrors, ParseError{fmt.Sprintf("error: %s", token.value), token.line})
		return
	case kwEOF:
		//Do nothing
	case skwProject:
		var g generic
		g, token, parseErrors = parseGeneric(l, token)
		if len(parseErrors) > 0 {
			return
		}
		parseErrors = findIndentationErrors(&g, smlKwTree, smlKeyword2string)
		if len(parseErrors) > 0 {
			return
		}

		parser := newSMLParser()
		project = parser.parseProject(g)
		parseErrors = append(parseErrors, parser.parseErrors...)
		if projectName != "" {
			// project name supplied in API overrides name in SML file
			project.Name = projectName
		}
	default:
		parseErrors = append(parseErrors, ParseError{fmt.Sprintf("expecting project, got: %s", smlKeyword2string[token.kw]), token.line})
		return
	}

	if token.kw != kwEOF {
		if token.kw == kwError {
			parseErrors = append(parseErrors, ParseError{"this should not happen", token.line})
		}
		parseErrors = append(parseErrors, ParseError{fmt.Sprintf("unexpected content beyond project: \"%s\"", smlKeyword2string[token.kw]), token.line})
	}

	if len(parseErrors) > 0 {
		return
	}

	var errs []ParseError
	errs = ensureUniqueTableAndColumnNames(&project)
	parseErrors = append(parseErrors, errs...)

	errs = ensureUniqueDatasetNames(&project)
	parseErrors = append(parseErrors, errs...)

	errs = computeDerivedAttributesInProject(&project)
	parseErrors = append(parseErrors, errs...)

	return
}

type smlFileType int

const (
	fileTypeUnknown smlFileType = iota
	fileTypeDataset
	fileTypeTable
	fileTypeScheme
)

var smlFileTypes = map[string]smlFileType{
	".table.sml":   fileTypeTable,
	".dataset.sml": fileTypeDataset,
	".scheme.sml":  fileTypeScheme,
}

func typeFromFilename(filename string) smlFileType {
	for k, v := range smlFileTypes {
		n := len(k)
		if len(filename) > n && strings.ToLower(filename[len(filename)-n:]) == k {
			return v
		}
	}
	return fileTypeUnknown
}

// ParseSMLPieces ...
// given an SML model, returns a parsed project
func ParseSMLPieces(inputs []string, filenames []string) (project Project, errorLine int, err error) {
	defer func() {
		if r := recover(); r != nil {
			pe, ok := r.([]ParseError)
			var msg string
			if ok {
				msg = pe[0].Msg
				errorLine = pe[0].LineNum
			} else {
				msg = fmt.Sprintf("unknown error: %v", r)
				errorLine = -1
			}
			err = errors.New(msg)
		}
	}()

	if len(inputs) != len(filenames) {
		err = fmt.Errorf("inputs length and filenames length do not match in API call")
		return
	}

	project.TableNames = make([]string, 0)
	project.Tables = make(map[string]Table)
	project.Datasets = make(map[string]Dataset)

	var tables []generic
	var datasets []generic

	for i := range inputs {
		var expectedKeyword keyword
		filename := filenames[i]
		switch typeFromFilename(filename) {
		case fileTypeUnknown:
			continue
		case fileTypeDataset:
			expectedKeyword = skwDataset
		case fileTypeTable:
			expectedKeyword = skwTable
		}

		l := startLex(inputs[i], smlKeywords, nil)

		token := <-l
		switch token.kw {
		case kwError:
			err = errors.New(token.value)
			errorLine = token.line
			return
		case kwEOF:
			err = nil
		case expectedKeyword:
			var g generic
			var errs []ParseError
			g, token, errs = parseGeneric(l, token)
			if errs != nil {
				err = errors.New(errs[0].Msg)
				errorLine = errs[0].LineNum
				return
			}
			errs = findIndentationErrors(&g, makeKwSubtree(smlKwTree, expectedKeyword), smlKeyword2string)
			if errs != nil {
				err = errors.New(errs[0].Msg)
				errorLine = errs[0].LineNum
				return
			}
			switch expectedKeyword {
			case skwDataset:
				datasets = append(datasets, g)
			case skwTable:
				tables = append(tables, g)
			}
		default:
			err = fmt.Errorf("expecting %s, got: %s in file %s line %d", smlKeyword2string[expectedKeyword],
				smlKeyword2string[token.kw], filename, token.line)
			return
		}

		if token.kw != kwEOF {
			if token.kw == kwError {
				err = fmt.Errorf("this should not happen")
				return
			}
			err = fmt.Errorf("unexpected content beyond %s: \"%s\" (file %s line %d)", smlKeyword2string[expectedKeyword],
				smlKeyword2string[token.kw], filename, token.line)
			return

		}
	}
	parser := newSMLParser()

	// parse
	for i := range datasets {
		dataset := parser.parseDataset(&datasets[i])
		if parser.didParseFail() {
			err = errors.New(parser.getFirstParseError().Msg)
			errorLine = parser.getFirstParseError().LineNum
			return
		}
		project.Datasets[dataset.Name] = dataset
	}
	for i := range tables {
		table := parser.parseTable(&tables[i])
		if parser.didParseFail() {
			err = errors.New(parser.getFirstParseError().Msg)
			errorLine = parser.getFirstParseError().LineNum
			return
		}
		project.Tables[table.Name] = table
		project.TableNames = append(project.TableNames, table.Name)
	}
	errs := ensureUniqueTableAndColumnNames(&project)
	if len(errs) > 0 {
		err = errors.New(errs[0].Msg)
		errorLine = errs[0].LineNum
		return
	}
	errs = ensureUniqueDatasetNames(&project)
	if len(errs) > 0 {
		err = errors.New(errs[0].Msg)
		errorLine = errs[0].LineNum
		return
	}
	return
}

func quotedString(s string) string {
	n := len(s)
	if n > 2 && (s[0] == '"' || s[0] == '[' || s[0] == '`') {
		endQuote := s[0]
		if endQuote == '[' {
			endQuote = ']'
		}
		if s[n-1] == endQuote && s[n-2] != '\\' {
			if possiblyID := strings.TrimSpace(s[1 : n-1]); IsValidID(possiblyID) {
				return possiblyID
			}
			if endQuote == '"' {
				if strings.IndexByte(s, '\'') < 0 {
					return "'" + s + "'"
				}
				return fmt.Sprintf("%q", s)
			}
		}
	}
	return s
}

func ConvertColumnName(uniqueNames map[string]bool, name string) string {
	name = ConvertName(name)
	if _, isPresent := isAggFunction[strings.ToLower(name)]; isPresent {
		name = "_" + name
	}
	i := 1
	base := name
	for uniqueNames[name] {
		name = fmt.Sprintf("%s%d", base, i)
		i++
	}
	uniqueNames[name] = true
	return name
}
