package sml

import (
	"regexp"
	"strings"

	log "github.com/maagodata/maago-commons/logger"
)

// CreateQualifiedNames creates a qualified name using "names" by combining them
// based on the dialect.
func CreateQualifiedNames(names []string, dialect string) (qName string, err error) {

	if len(names) == 0 {
		// Nothing to combine.
		return
	}

	separator := "."

	// Table and column names in the backend namespace (i.e. in database)
	// can have spaces and other special characters in them. So, they
	// should be enclosed in dialect specific quotes.
	startQuote, endQuote, err := GetQuotes(dialect)
	if err != nil {
		log.Error(err)
		return
	}

	for _, name := range names {
		if !strings.HasPrefix(name, startQuote) {
			qName += startQuote + name + endQuote + separator
		} else {
			qName += name + separator
		}
	}
	qName = strings.TrimSuffix(qName, separator)

	return
}

// RemoveExtraSpaces removes extra spaces in a given string.
// It trims spaces from left and right ends. It also replaces
// extra spaces from between words in a string.
func RemoveExtraSpaces(s string) string {

	// Remove spaces from left and right ends.
	str := strings.Trim(s, " \t")

	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	str = string(re.ReplaceAll([]byte(str), []byte(" ")))

	return str
}

type StringSet struct {
	values map[string]struct{}
}

func NewStringSet() StringSet {
	return StringSet{values: make(map[string]struct{}, 0)}
}

func (s StringSet) Exists(v string) bool {
	_, ok := s.values[v]
	return ok
}

func (s StringSet) Insert(v string) {
	s.values[v] = struct{}{}
}

func (s StringSet) Remove(v string) {
	delete(s.values, v)
}

func ArrayToSet(values []string) StringSet {
	s := make(map[string]struct{}, len(values))
	for _, v := range values {
		s[v] = struct{}{}
	}
	return StringSet{values: s}
}

func IsAlpha(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// re-define SML data types, to avoid import cycle
const (
	typeString   = "string"
	typeVerbatim = "verbatim"
	typeInt      = "int"
	typeFloat    = "float"
	typeDatetime = "datetime"
	typeBool     = "bool"
	typeTime     = "time"
)

var UNIVERSAL_TYPE_TO_CLICKHOUSE = map[string]string{
	"TINYINT":     "Int16",
	"SMALLINT":    "Int32",
	"INT":         "Int64",
	"INTEGER":     "Int64",
	"BIGINT":      "Int64",
	"INT8":        "Int8",
	"INT16":       "Int16",
	"INT32":       "Int32",
	"INT64":       "Int64",
	"UINT8":       "UInt8",
	"UINT16":      "UInt16",
	"UINT32":      "UInt32",
	"UINT64":      "UInt64",
	"DECIMAL":     "Float64",
	"DOUBLE":      "Float64",
	"FLOAT":       "Float64",
	"FLOAT32":     "Float32",
	"FLOAT64":     "Float64",
	"NUMERIC":     "Float64",
	"REAL":        "Float64",
	"CHAR":        "String",
	"VARCHAR":     "String",
	"STRING":      "String",
	"TEXT":        "String",
	"FIXEDSTRING": "String",
	"TIME":        "DateTime",
	"TIMESTAMP":   "DateTime",
	"DATE":        "Date",
	"DATETIME":    "DateTime",
	"BOOL":        "Boolean",
	"BOOLEAN":     "Boolean",
}

func makeReservedWords() StringSet {
	keywords := []string{
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
	m := ArrayToSet(keywords)
	for k := range UNIVERSAL_TYPE_TO_CLICKHOUSE {
		m.Insert(k)
	}
	return m
}

var SQL_RESERVED_WORDS = makeReservedWords()

func IsReservedWord(w string) bool {
	return SQL_RESERVED_WORDS.Exists(strings.ToUpper(w))
}
