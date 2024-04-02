package sml

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode"
)

const Debug = false

func getStackTrace() string {
	// stolen from https://www.komu.engineer/blogs/golang-stacktrace/golang-stacktrace
	const maxStackLength = 50

	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return trace
}

func assert(b bool) {
	if !b {
		fmt.Printf("%s\n", getStackTrace())
		panic("assertion failure")
	}
}

var Assert = func() func(bool) {
	if Debug {
		return assert
	}
	return func(bool) {}
}()

func isAlpha(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || (c >= '0' && c <= '9')
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// IsValidID checks if a name is a valid id
// id definition:
//
//	minimum length has to be 1
//	1st character can only be undersore or alphabets
//	2nd character onwards can be alphanumeric and undersores
func IsValidID(value string) bool {
	if len(value) < 1 {
		return false
	}
	if !isAlpha(value[0]) {
		return false
	}
	for i := 1; i < len(value); i++ {
		if !isAlphaNumeric(value[i]) {
			return false
		}
	}
	return true
}

// convertName removes all non alphanumeric chracters by underscores
// 2017# becomes 2017_
// 2017#@$Year becomes 2017_Year
// 2017#Year@NYC becomes 2017_Year_NYC
func _ConvertName(name string) string {
	if IsValidID(name) {
		return name
	}

	reg, err := regexp.Compile("[^\\w]+")
	if err != nil {
		return "SMLError_" + name
	}

	name = reg.ReplaceAllString(strings.TrimSpace(name), "_")
	name = strings.Trim(name, "_")
	if name == "" || !isAlpha(name[0]) {
		name = "_" + name
	}
	return name
}

func ConvertName(name string) string {
	name = _ConvertName(name)
	if sqlReservedWords[strings.ToUpper(name)] {
		name = "_" + name
	}
	return name
}

func dedupStringSlice(a []string) []string {
	// case insensitive
	if a == nil {
		return nil
	}
	m := make(map[string]bool)
	src := 0
	dst := 0
	for ; src < len(a); src++ {
		l := strings.ToLower(a[src])
		if !m[l] {
			m[l] = true
			a[dst] = a[src]
			dst++
		}
	}
	return a[:dst]
}

func createSmartLabel(name string) string {
	tmp := strings.Trim(name, "_")
	if tmp == "" {
		return name
	}
	parts := strings.Split(tmp, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, " ")
}

func sliceContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isTableAtomic(sql string) bool {
	var tmp = strings.TrimSpace(sql)
	tmp = strings.TrimLeft(tmp, "(")
	words := strings.Fields(tmp)
	if len(words) > 1 && strings.ToLower(words[0]) == "select" {
		return false
	}
	return true
}

func sortedKeys(m map[string]bool) (a []string) {
	a = make([]string, len(m))
	i := 0
	for k := range m {
		a[i] = k
		i++
	}
	sort.Strings(a)
	return
}

// smart diff of JSON strings (used in testing)

// JType ...
type JType int

// Jt constants...
const (
	JTInt JType = iota
	JTFloat
	JTString
	JTBool
	JTMap
	JTArray
	JTNull
)

var jt2string = map[JType]string{
	JTInt:    "int",
	JTFloat:  "float",
	JTString: "string",
	JTBool:   "bool",
	JTMap:    "map",
	JTArray:  "array",
	JTNull:   "nil",
}

// JInGo ...
type JInGo struct {
	Type  JType
	Atom  string
	Map   map[string]*JInGo
	Array []JInGo
}

// J2Go ...
func J2Go(j string) (g JInGo, e error) {
	var i interface{}
	e = json.Unmarshal([]byte(j), &i)
	if e != nil {
		return
	}
	g, e = recursiveJ2Go(i)
	return
}

func recursiveJ2Go(i interface{}) (g JInGo, e error) {
	if i == nil {
		g.Type = JTNull
		return
	}
	switch val := i.(type) {
	case string:
		g.Type = JTString
		g.Atom = val
		return
	case int: // this will never happen!  https://golang.org/pkg/encoding/json/#Unmarshal
		g.Type = JTInt
		g.Atom = fmt.Sprintf("%d", val)
		return
	case float64:
		g.Type = JTFloat
		g.Atom = fmt.Sprintf("%v", val)
		return
	case bool:
		g.Type = JTBool
		g.Atom = fmt.Sprintf("%v", val)
		return
	case []interface{}:
		g.Type = JTArray
		for _, v := range val {
			var g1 JInGo
			g1, e = recursiveJ2Go(v)
			if e != nil {
				return
			}
			g.Array = append(g.Array, g1)
		}
		return
	case map[string]interface{}:
		g.Type = JTMap
		g.Map = make(map[string]*JInGo)
		for k, v := range val {
			var g1 JInGo
			g1, e = recursiveJ2Go(v)
			if e != nil {
				return
			}
			g.Map[k] = &g1
		}
		return
	default:
		e = errors.New("unknown type")
		return
	}
}

// IsStringInSlice returns true if the string is present in the slice
func IsStringInSlice(str string, list []string) bool {
	str = strings.ToLower(str)
	for _, s := range list {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}

// DiffJSON ...
func DiffJSON(a, b string) string {
	a1, e := J2Go(a)
	if e != nil {
		return fmt.Sprintf("could not parse first JSON: %s", e.Error())
	}
	var b1 JInGo
	b1, e = J2Go(b)
	if e != nil {
		return fmt.Sprintf("could not parse second JSON: %s", e.Error())
	}
	return diffJInGo(a1, b1)
}

func diffJInGo(a, b JInGo) string {
	if a.Type != b.Type {
		return fmt.Sprintf("types: %s and %s", jt2string[a.Type], jt2string[b.Type])
	}
	switch a.Type {
	case JTInt, JTString, JTBool, JTFloat:
		if a.Atom != b.Atom {
			return fmt.Sprintf("\"%s\" v \"%s\"", a.Atom, b.Atom)
		}
	case JTMap:
		var aKeys, bKeys []string
		for k := range a.Map {
			aKeys = append(aKeys, k)
		}
		sort.Strings(aKeys)
		for k := range b.Map {
			bKeys = append(bKeys, k)
		}
		sort.Strings(bKeys)
		if len(aKeys) != len(bKeys) {
			return fmt.Sprintf("keys: %s v %s", strings.Join(aKeys, ","), strings.Join(bKeys, ","))
		}
		for i, k := range aKeys {
			if bKeys[i] != k {
				return fmt.Sprintf("keys: %s v %s", strings.Join(aKeys, ","), strings.Join(bKeys, ","))
			}
			diff := diffJInGo(*a.Map[k], *b.Map[k])
			if diff != "" {
				return fmt.Sprintf("[%s] %s", k, diff)
			}
		}
	case JTArray:
		if len(a.Array) != len(b.Array) {
			return fmt.Sprintf("array lengths: %d v %d", len(a.Array), len(b.Array))
		}
		for i := 0; i < len(a.Array); i++ {
			diff := diffJInGo(a.Array[i], b.Array[i])
			if diff != "" {
				return fmt.Sprintf("[%d] %s", i, diff)
			}
		}
	}
	return ""
}

type stringPair struct {
	first  string
	second string
}

type kwTree struct {
	kw       keyword
	children []*kwTree
}

// debugging function
func printKwTree(tree *kwTree, indent string, kw2string map[keyword]string) string {
	assert(tree != nil)
	kwName, isPresent := kw2string[tree.kw]
	if !isPresent {
		kwName = linglingKeyword2string[tree.kw]
	}
	a := []string{fmt.Sprintf("%s%s", indent, kwName)}
	for _, child := range tree.children {
		a = append(a, printKwTree(child, indent+"  ", kw2string))
	}
	return strings.Join(a, "\n")
}

func makeKwTree(keywords map[string]keyword, hierarchy string) *kwTree {
	// assumptions:
	// 1. hierarchy corresponds to exactly one top-level object (e.g. project)
	// 2. every child is indented by exactly two spaces
	// 3. all chidlren of a node are indented (correctly) at the same level
	var stack []*kwTree
	n := 0
	lastIndent := 0
	for _, line := range strings.Split(strings.TrimSpace(hierarchy), "\n") {
		v1 := strings.TrimRightFunc(line, unicode.IsSpace)
		v2 := strings.TrimLeftFunc(v1, unicode.IsSpace)
		indent := len(v1) - len(v2)
		assert(indent%2 == 0)
		kw, isPresent := keywords[v2]
		if !isPresent {
			kw, isPresent = linglingKeywords[v2]
		}
		assert(isPresent)
		if indent > lastIndent {
			assert(indent == lastIndent+2)
			stack = append(stack, &kwTree{kw: kw})
			stack[n-1].children = append(stack[n-1].children, stack[n])
			n++
		} else {
			if indent < lastIndent {
				delta := (lastIndent - indent) / 2
				stack = stack[:n-delta]
				n -= delta
			}
			if n >= 2 {
				stack[n-1] = &kwTree{kw: kw}
				stack[n-2].children = append(stack[n-2].children, stack[n-1])
			} else {
				assert(n == 0)
				stack = append(stack, &kwTree{kw: kw})
				n++
			}
		}
		lastIndent = indent
	}
	return stack[0]
}

func makeKwSubtree(tree *kwTree, kw keyword) *kwTree {
	if tree == nil {
		return nil
	}
	for _, t := range tree.children {
		if t.kw == kw {
			return t
		}
	}
	return tree
}

// RandomString returns a random string of specified length
func RandomString(n int) string {
	// generate random string of specified length -
	// Ref: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
