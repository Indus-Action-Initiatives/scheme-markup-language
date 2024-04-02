package sml

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	indexNumber int = iota
	indexName
	indexOperator
	indexBackquote
	indexDoublequote
	indexSinglequote
	indexCatchall
	numPatterns
)

func makeRegex() *regexp.Regexp {
	var patterns = []string{
		"(" + `(?:(?:\d+\.\d*)|(?:\.\d+)|(?:\d+))(?:[eE][-+]?\d+)?` + ")", // number
		"(" + `[a-zA-Z_]\w*` + ")",                                        // name
		"(" + `<=|>=|<>` + ")",                                            // operator
		"(" + "`[^`]+`" + ")",                                             // backquote
		"(" + `"(?:""|\\"|[^"])*"` + ")",                                  // doublequote
		"(" + `'(?:''|\\'|[^'])*'` + ")",                                  // singlequote
		"(" + `[^\s]` + ")",                                               // catchall
	}
	return regexp.MustCompile(strings.Join(patterns, "|"))
}

var reToken = makeRegex()

// SQL-2003
// https://ronsavage.github.io/SQL/
var sqlReservedWords = makeReservedWordsMap()

func makeReservedWordsMap() (m map[string]bool) {
	var list = []string{
		// "reserved"
		"ADD", "ALL", "ALLOCATE", "ALTER", "AND", "ANY", "ARE", "ARRAY",
		"AS", "ASENSITIVE", "ASYMMETRIC", "AT", "ATOMIC", "AUTHORIZATION",
		"BEGIN", "BETWEEN", "BIGINT", "BINARY", "BLOB", "BOOLEAN", "BOTH",
		"BY", "CALL", "CALLED", "CASCADED", "CASE", "CAST", "CHAR",
		"CHARACTER", "CHECK", "CLOB", "CLOSE", "COLLATE", "COLUMN", "COMMIT",
		"CONNECT", "CONSTRAINT", "CONTINUE", "CORRESPONDING", "CREATE",
		"CROSS", "CUBE", "CURRENT", "CURRENT_DATE",
		"CURRENT_DEFAULT_TRANSFORM_GROUP", "CURRENT_PATH", "CURRENT_ROLE",
		"CURRENT_TIME", "CURRENT_TIMESTAMP", "CURRENT_TRANSFORM_GROUP_FOR_TYPE",
		"CURRENT_USER", "CURSOR", "CYCLE", "DATE", "DAY", "DEALLOCATE",
		"DEC", "DECIMAL", "DECLARE", "DEFAULT", "DELETE", "DEREF", "DESCRIBE",
		"DETERMINISTIC", "DISCONNECT", "DISTINCT", "DOUBLE", "DROP",
		"DYNAMIC", "EACH", "ELEMENT", "ELSE", "END", "END-EXEC", "ESCAPE",
		"EXCEPT", "EXEC", "EXECUTE", "EXISTS", "EXTERNAL", "FALSE", "FETCH",
		"FILTER", "FLOAT", "FOR", "FOREIGN", "FREE", "FROM", "FULL",
		"FUNCTION", "GET", "GLOBAL", "GRANT", "GROUP", "GROUPING", "HAVING",
		"HOLD", "HOUR", "IDENTITY", "IMMEDIATE", "IN", "INDICATOR", "INNER",
		"INOUT", "INPUT", "INSENSITIVE", "INSERT", "INT", "INTEGER",
		"INTERSECT", "INTERVAL", "INTO", "IS", "ISOLATION", "JOIN", "LANGUAGE",
		"LARGE", "LATERAL", "LEADING", "LEFT", "LIKE", "LOCAL", "LOCALTIME",
		"LOCALTIMESTAMP", "MATCH", "MEMBER", "MERGE", "METHOD", "MINUTE",
		"MODIFIES", "MODULE", "MONTH", "MULTISET", "NATIONAL", "NATURAL",
		"NCHAR", "NCLOB", "NEW", "NO", "NONE", "NOT", "NULL", "NUMERIC",
		"OF", "OLD", "ON", "ONLY", "OPEN", "OR", "ORDER", "OUT", "OUTER",
		"OUTPUT", "OVER", "OVERLAPS", "PARAMETER", "PARTITION", "PRECISION",
		"PREPARE", "PRIMARY", "PROCEDURE", "RANGE", "READS", "REAL",
		"RECURSIVE", "REF", "REFERENCES", "REFERENCING", "REGR_AVGX",
		"REGR_AVGY", "REGR_COUNT", "REGR_INTERCEPT", "REGR_R2", "REGR_SLOPE",
		"REGR_SXX", "REGR_SXY", "REGR_SYY", "RELEASE", "RESULT", "RETURN",
		"RETURNS", "REVOKE", "RIGHT", "ROLLBACK", "ROLLUP", "ROW", "ROWS",
		"SAVEPOINT", "SCROLL", "SEARCH", "SECOND", "SELECT", "SENSITIVE",
		"SESSION_USER", "SET", "SIMILAR", "SMALLINT", "SOME", "SPECIFIC",
		"SPECIFICTYPE", "SQL", "SQLEXCEPTION", "SQLSTATE", "SQLWARNING",
		"START", "STATIC", "SUBMULTISET", "SYMMETRIC", "SYSTEM", "SYSTEM_USER",
		"TABLE", "THEN", "TIME", "TIMESTAMP", "TIMEZONE_HOUR", "TIMEZONE_MINUTE",
		"TO", "TRAILING", "TRANSLATION", "TREAT", "TRIGGER", "TRUE",
		"UESCAPE", "UNION", "UNIQUE", "UNKNOWN", "UNNEST", "UPDATE", "UPPER",
		"USER", "USING", "VALUE", "VALUES", "VAR_POP", "VAR_SAMP", "VARCHAR",
		"VARYING", "WHEN", "WHENEVER", "WHERE", "WIDTH_BUCKET", "WINDOW",
		"WITH", "WITHIN", "WITHOUT", "YEAR",
		// some of the "non-reserved"
		"ASC", "CEIL", "CEILING", "CHARACTERS",
		"CHARACTER_LENGTH", "CHARACTER_SET_CATALOG", "CHARACTER_SET_NAME",
		"CHARACTER_SET_SCHEMA", "CHAR_LENGTH", "CONTAINS",
		"DATETIME_INTERVAL_CODE", "DATETIME_INTERVAL_PRECISION",
		"DENSE_RANK", "DESC", "EQUALS", "EXCLUDE", "EXCLUDING", "EXP",
		"FIRST", "FLOOR", "FOLLOWING", "FOUND", "GO", "GOTO", "INCLUDING",
		"INTERSECTION", "KEY", "KEY_MEMBER", "KEY_TYPE", "LAST", "LENGTH",
		"LOWER", "MAXVALUE", "MESSAGE_LENGTH",
		"MESSAGE_OCTET_LENGTH", "MESSAGE_TEXT", "MINVALUE", "MOD",
		"NESTING", "NEXT", "NULLABLE", "NULLIF", "NULLS", "PERCENTILE_CONT",
		"PERCENTILE_DISC", "PERCENT_RANK", "POWER", "PRECEDING", "PRESERVE",
		"ROW_COUNT", "SIZE", "STDDEV_POP",
		"STDDEV_SAMP", "SUBSTRING", "TRIM", "TYPE",
		// extra
		"STRING", "DATETIME", "LIMIT",
	}
	m = make(map[string]bool)
	for _, v := range list {
		m[v] = true
	}
	return
}

// SQLTokenType ...
type SQLTokenType int

// TokenTypes ...
const (
	NumberSQLToken SQLTokenType = iota
	StringSQLToken
	NameSQLToken
	OpSQLToken
	ReservedSQLToken
	UnknownSQLToken
)

// SQLToken ...
type SQLToken struct {
	Position int
	Type     SQLTokenType
	Value    string
	Value2   string
}

var patternIndexToType = map[int]SQLTokenType{
	indexNumber:      NumberSQLToken,
	indexName:        NameSQLToken,
	indexOperator:    OpSQLToken,
	indexBackquote:   NameSQLToken,
	indexDoublequote: StringSQLToken,
	indexSinglequote: StringSQLToken,
	indexCatchall:    UnknownSQLToken,
}

func makeToken(position int, index int, value string) (token SQLToken) {
	token.Position = position
	token.Type = patternIndexToType[index]
	if token.Type == UnknownSQLToken {
		// TODO handle any special case for the catchall scenario
		token.Type = OpSQLToken
	}
	// if index == indexBackquote {
	// 	token.Value = value[1 : len(value)-1]
	// } else {
	// 	token.Value = value
	// }
	token.Value = value
	token.Value2 = ""
	return
}

func (token *SQLToken) toString() string {
	names := []string{"NUMBER", "STRING", "NAME", "OP", "RESERVED"}
	var value string
	if token.Value2 != "" {
		value = fmt.Sprintf("%s . %s", token.Value, token.Value2)
	} else {
		value = token.Value
	}
	return fmt.Sprintf("%s [%s] %d", names[token.Type], value, token.Position)
}

func (token *SQLToken) isNumber() bool {
	return token.Type == NumberSQLToken
}

func (token *SQLToken) isString() bool {
	return token.Type == StringSQLToken
}

func (token *SQLToken) isName() bool {
	return token.Type == NameSQLToken
}

func (token *SQLToken) isOp() bool {
	return token.Type == OpSQLToken
}

func (token *SQLToken) isReserved() bool {
	return token.Type == ReservedSQLToken
}

func parseMatch(groups []string) int {
	// make sure exactly one of the groups is not empty
	// return the corresponding index
	if len(groups) != numPatterns {
		panic("number of groups is not equal to number of patterns")
	}
	index := -1
	for i, g := range groups {
		if g != "" {
			if index >= 0 {
				panic("more than one group are not empty")
			}
			index = i
		}
	}
	if index < 0 {
		panic("none of the groups are empty")
	}
	return index
}

// TokenizeSQL ...
func TokenizeSQL(sql string) (tokens []SQLToken) {
	sql = strings.Replace(strings.Replace(sql, "\n", " ", -1), "\r", " ", -1)
	matches := reToken.FindAllStringSubmatch(sql, -1)
	indices := reToken.FindAllStringSubmatchIndex(sql, -1)
	if len(matches) != len(indices) {
		panic("number of matches and number of indices don't match")
	}
	for i, match := range matches {
		index := parseMatch(match[1:])
		value := match[0]
		position := indices[i][0]
		tokens = append(tokens, makeToken(position, index, value))
	}
	// merge x . y into one token
	i := len(tokens) - 1
	for i > 1 {
		if tokens[i].isName() && tokens[i-1].Value == "." && tokens[i-2].isName() {
			tokens[i-2].Value2 = tokens[i].Value
			tokens = append(tokens[:i-1], tokens[i+1:]...)
			i -= 3 // we don't allow a . b . c
		} else {
			i--
		}
	}
	// look for reserved words
	for i, token := range tokens {
		if token.isName() && token.Value2 == "" && sqlReservedWords[strings.ToUpper(token.Value)] {
			tokens[i].Type = ReservedSQLToken
		}
	}
	return
}
