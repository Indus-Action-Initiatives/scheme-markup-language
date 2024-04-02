package sml

import (
	"fmt"
	"strings"
)

/*
This file is the core of the implementation of various modeling languages
including "sml".

From Urban Dictionary:
	Ling Ling: A student who overachieves, winning awards, and is ostensibly better than you in every way,
	especially when used in a direct comparison.
From Wikipedia:
	Lingling dialect: an unclassified mixed Chinese dialect

The first piece is the lexer/tokenizer.
Each line in the sml model file is broken into 4 pieces:
	1. indentation depth
	2. a keyword
	3. (optional) a value
	4. (optional) a trailing comment starting with # or //
A line can also be empty or entirely comment.

The indentations define a hierarchical structure.
The meaning of indentation is identical to that of Python except for how tabs are handled.  Tabs are
replaced with 4 spaces (const tabSize; this can be potentially parameterized).  The following Python warning applies
to this as well:  "it is unwise to use a mixture of spaces and tabs for the indentation in a single source file."
If you think of the hierarchy as a tree, all children of a node should have exactly the same indentation.

Keywords are a small set of strings, defined independently for each language like sml.

Values are strings and can be specified as:
	(a) "some string"
	(b) 'some string'
	(c) """some string"""
	(d) some string
In cases (a) and (b), embedded quotes have to be escaped by a preceding backslash.
In case (d), end of line or a trailing comment signifies the end of the string.
In case (d), leading and trailing whitespace are stripped from the value.

An optional comment can appear at the end of any line, starting with # or //.  The lexer discards the comment.

Once a file has been lexed into a sequence of tokens (indentation, keyword, value), the overall hierarchy (tree) is
first inferred from the indentations in parseGeneric().  This transforms the entire model file into an semi-typed
data structure using the "generic" type.  In the second phase, the final hierarchical data structure is constructed
using the types (e.g. Table, Column etc.) based on the language.
*/

// ParseError stores datamodel parse errors and line number
type ParseError struct {
	Msg     string
	LineNum int
}

// ------------ lexer / tokenizer -------------
const tabSize = 4

type token struct {
	indent int
	kw     keyword
	value  string
	line   int
}

type keyword int

const (
	kwError keyword = -iota - 1
	kwEOF
)

// Comment is a lingling comment
type Comment struct {
	Comment string
	Line    int
}

func makeKeywordToString(keywords map[string]keyword, synonyms []string) map[keyword]string {
	keyword2string := make(map[keyword]string)
	for k, v := range keywords {
		if !sliceContains(synonyms, k) {
			keyword2string[v] = k
		}
	}
	return keyword2string
}

// process input, line by line, identifying 3 attributes per line:
//
//	indentation, keyword and value
//	the value may be followed by a comment starting with // or #
//
// a value may be
//
//	any unquoted string (leading and trailing spaces are trimmed)
//	a quoted "string"
//	a triple-quoted """string"""
//
// a triple-quoted string is allowed to span multiple lines (new lines are replaced with space)
func runLexer(input string, string2keyword map[string]keyword, tokens chan token, comments *[]Comment) {
	lines := strings.Split(input, "\n")

	lineNum := 0
	for lineNum < len(lines) {
		line := lines[lineNum]
		pos := 0
		//--- first: indentation
		indent := 0
		for pos < len(line) {
			c := line[pos]
			if c == ' ' {
				pos++
				indent++
			} else if c == '\t' {
				pos++
				indent += tabSize
			} else {
				break
			}
		}
		if pos >= len(line) || line[pos] == '#' || pos < len(line)-1 && line[pos:pos+2] == "//" {
			// empty line
			lineNum++
			comment := "#"
			if pos < len(line) {
				comment = strings.TrimSpace(line[pos:])
			}
			if comments != nil {
				*comments = append(*comments, Comment{comment, lineNum})
			}
			continue
		}
		//--- second: keyword
		if !isAlpha(line[pos]) {
			tokens <- errorToken(lineNum, pos, "expected keyword, got %c", line[pos])
			// go to next line
			lineNum++
			continue
		}
		end := pos + 1
		for end < len(line) && isAlphaNumeric(line[end]) {
			end++
		}
		kwString := line[pos:end]
		kw, isPresent := string2keyword[kwString]
		if !isPresent {
			kw, isPresent = linglingKeywords[kwString]
			if !isPresent {
				tokens <- errorToken(lineNum, pos, "unknown keyword: %s", kwString)
				// go to next line
				lineNum++
				continue
			}
		}
		pos = end
		//--- third: value
		for pos < len(line) && isSpace(line[pos]) {
			pos++
		}
		var value string
		if pos < len(line) {
			if strings.HasPrefix(line[pos:], `"""`) {
				var linesConsumed int
				pos, linesConsumed, value = multilineString(pos+3, lines[lineNum:])
				if linesConsumed < 0 {
					tokens <- errorToken(lineNum, pos, "multiline string not terminated")
					break
				} else if linesConsumed > 0 {
					lineNum += linesConsumed
					line = lines[lineNum]
				}
			} else if line[pos] == '"' || line[pos] == '\'' {
				delimiter := line[pos : pos+1]
				i := pos + 1
				for {
					// handle embedded " escaped by \
					j := strings.IndexByte(line[i:], delimiter[0])
					if j < 0 {
						i = j
						break
					}
					if j == 0 || line[i+j-1] != '\\' {
						value += line[i : i+j]
						i += j
						break
					}
					value += line[i:i+j-1] + delimiter
					i += j + 1
				}
				if i < 0 {
					tokens <- errorToken(lineNum, pos, "string not terminated")
					break
				}
				pos = i + 1
			} else {
				// if value is not delimited by quotes, it terminated by EOLN or trailing comment
				i := strings.Index(line[pos:], "//")
				if i < 0 {
					i = strings.IndexByte(line[pos:], '#')
				}
				if i < 0 {
					i = len(line)
				} else {
					i += pos
				}
				value = strings.TrimSpace(line[pos:i]) // can be empty
				pos = i
			}
		}
		//--- fourth: optional comment
		for pos < len(line) && isSpace(line[pos]) {
			pos++
		}
		if pos < len(line) {
			if line[pos] == '#' || strings.HasPrefix(line[pos:], "//") {
				if comments != nil {
					*comments = append(*comments, Comment{strings.TrimSpace(line[pos:]), lineNum + 1})
				}
			} else {
				tokens <- errorToken(lineNum, pos, "unexpected content: %s", line[pos:])
				// go to next line
				lineNum++
				continue
			}
		}
		//--- success!  emit token
		tokens <- token{indent, kw, value, lineNum + 1}
		// go to next line
		lineNum++
	}
	if comments != nil {
		*comments = append(*comments, Comment{"", lineNum + 1})
	}
	tokens <- token{-1, kwEOF, "", lineNum + 1}
	close(tokens)
}

func errorToken(line int, pos int, format string, args ...interface{}) token {
	return token{-1, kwError,
		fmt.Sprintf("error at position %d: %s", pos+1,
			fmt.Sprintf(format, args...)),
		line + 1}
}

func multilineString(pos int, lines []string) (newPos int, linesConsumed int, value string) {
	var pieces []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		j := strings.Index(line[pos:], `"""`)
		if j >= 0 {
			pieces = append(pieces, line[pos:pos+j])
			return pos + j + 3, i, strings.Join(pieces, "")
		}
		s := line[pos:]
		if strings.HasSuffix(s, `\`) {
			s = s[:len(s)-1]
		} else {
			s = s + "\n"
		}
		pieces = append(pieces, s)
		pos = 0
	}
	return -1, -1, ""
}

// lexer top-level function
func startLex(input string, string2keyword map[string]keyword, comments *[]Comment) chan token {
	tokens := make(chan token)
	// '\r\n' is a line breaker used in windows system
	input = strings.ReplaceAll(input, "\r\n", "\n")
	go runLexer(input, string2keyword, tokens, comments)
	return tokens
}

/*
func stopLex(c chan token) {
	for range c {
	}
}
*/

//------------ generic parser -------------

type generic struct {
	kw       keyword
	value    string
	children map[keyword][]*generic
	line     int // TODO: this can be saved in parsed 'tables', 'columns' etc
}

func parseGeneric(l chan token, root token) (g generic, nextToken token, parseErrors []ParseError) {
	// consume the root
	g.kw = root.kw
	g.value = root.value
	g.line = root.line
	// consume children, if any
	nextToken = <-l
	for nextToken.kw == kwError {
		parseErrors = append(parseErrors, ParseError{nextToken.value, nextToken.line})
		// skip this token, take the next token
		nextToken = <-l
	}
	if nextToken.indent > root.indent {
		childIndent := nextToken.indent
		for nextToken.indent == childIndent {
			var child generic
			var errs []ParseError
			child, nextToken, errs = parseGeneric(l, nextToken)
			parseErrors = append(parseErrors, errs...)
			if g.children == nil {
				g.children = make(map[keyword][]*generic)
			}
			g.children[child.kw] = append(g.children[child.kw], &child)
		}
		if nextToken.indent > root.indent {
			// children should all be at the same indentation
			parseErrors = append(parseErrors, ParseError{"inconsistent indentation", nextToken.line})
			// skip this token, take the next token
			nextToken = <-l
		}
	}
	return
}

type genericByLine []*generic

func (a genericByLine) Len() int           { return len(a) }
func (a genericByLine) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a genericByLine) Less(i, j int) bool { return a[i].line < a[j].line }

// debugging function
func generic2string(g *generic, indent string, kw2string map[keyword]string) string {
	var out []string
	out = append(out, fmt.Sprintf("%s%s %s", indent, kw2string[g.kw], g.value))
	indent = indent + "  "
	for _, list := range g.children {
		//fmt.Printf("%s%s?\n", indent, kw2string[k])
		for _, child := range list {
			g1 := child
			out = append(out, generic2string(g1, indent, kw2string))
		}
	}
	return strings.Join(out, "\n")
}

func allKeyword2string(kw keyword) (kwstr string) {
	var ok bool
	if kwstr, ok = linglingKeyword2string[kw]; !ok {
		// should not happen
		kwstr = "unknown keyword"
	}
	// }
	return
}

type genericVisitor func(*generic) bool

func findIndentationErrors(g *generic, grammar *kwTree, keyword2string map[keyword]string) []ParseError {
	// first, we walk the entire tree to find any invalid parent-child relationship
	if g.kw != grammar.kw {
		panic("WTF")
	}
	line, level := firstLineMismatchWithGrammar(1, g, grammar)
	var kw keyword
	if line < 0 {
		return nil
	}
	// found it! keyword on 'line' has unexpected indentation
	// find previous line number
	previousLine := -1
	gv := func(g1 *generic) bool {
		if g1.line == line {
			kw = g1.kw
		} else if g1.line < line && g1.line > previousLine {
			previousLine = g1.line
		}
		return true
	}
	walkGeneric(g, gv, nil)
	// find path of nodes from root to previous line
	var path []*generic
	preFunc := func(g1 *generic) bool {
		path = append(path, g1)
		return g1.line != previousLine
	}
	postFunc := func(g1 *generic) bool {
		path = path[:len(path)-1]
		return true
	}
	walkGeneric(g, preFunc, postFunc)
	// the sequence of nodes in 'path' is valid
	// examine all possible positions of the offending keyword relative to this path
	keywordSequence := make([]keyword, len(path))
	for i := range path {
		keywordSequence[i] = path[i].kw
	}
	for len(keywordSequence) > 0 {
		keywordSequence = append(keywordSequence, kw)
		if isValidPath(keywordSequence, grammar) {
			possibleParent := len(keywordSequence) - 2
			newLevel := len(keywordSequence) - 1
			direction := "left"
			if newLevel > level {
				direction = "right"
			}
			parentKw := keywordSequence[possibleParent]
			parentKwName, isPresent := keyword2string[parentKw]
			if !isPresent {
				parentKwName = linglingKeyword2string[parentKw]
			}
			parentValue := strings.TrimSpace(path[possibleParent].value)
			if len(parentValue) > 8 {
				parentValue = parentValue[:5] + "..."
			}
			if parentValue != "" {
				parentValue = " (" + parentValue + ")"
			}
			kwName, isPresent := keyword2string[kw]
			if !isPresent {
				kwName = linglingKeyword2string[kw]
			}
			// special case for parentKwName = "project"
			if parentKwName == "project" {
				return []ParseError{{LineNum: line, Msg: fmt.Sprintf("indentation error for \"%s\", maybe move %s", kwName, direction)}}
			}
			return []ParseError{{LineNum: line, Msg: fmt.Sprintf("indentation error for \"%s\", maybe move %s to be under \"%s\"%s", kwName, direction, parentKwName, parentValue)}}
		}
		keywordSequence = keywordSequence[:len(keywordSequence)-2]
	}
	// if we are here, we could not find a cure for the invalid parent-child relationship
	// we do not return error and rely on the subsequent parsing to issue an error message
	return nil
}

func walkGeneric(g *generic, preFunc, postFunc genericVisitor) bool {
	if preFunc != nil && !preFunc(g) {
		return false
	}
	for kw := range g.children {
		for i := range g.children[kw] {
			if !walkGeneric(g.children[kw][i], preFunc, postFunc) {
				return false
			}
		}
	}
	if postFunc != nil {
		return postFunc(g)
	}
	return true
}

func firstLineMismatchWithGrammar(level int, g *generic, grammar *kwTree) (minLine, minLevel int) {
	// return first line of original text file where the hierarchical structure does not follow the grammar
	assert(grammar != nil)
	assert(g.kw == grammar.kw)
	minLine, minLevel = -1, -1
	for kw := range g.children {
		found := false
		for i := range grammar.children {
			if kw == grammar.children[i].kw {
				found = true
				for j := range g.children[kw] {
					line, level2 := firstLineMismatchWithGrammar(level+1, g.children[kw][j], grammar.children[i])
					if line >= 0 && (minLine < 0 || line < minLine) {
						minLine, minLevel = line, level2
					}
				}
				break
			}
		}
		if !found {
			// 'kw' is in g.children, but not in grammar.children
			for j := range g.children[kw] {
				if minLine < 0 || g.children[kw][j].line < minLine {
					minLine, minLevel = g.children[kw][j].line, level
				}
			}
		}
	}
	return
}

func isValidPath(path []keyword, tree *kwTree) bool {
	children := []*kwTree{tree}
	for i, kw := range path {
		found := -1
		for j := range children {
			if children[j].kw == kw {
				found = j
				break
			}
		}
		if found >= 0 {
			children = children[found].children
		} else {
			assert(i == len(path)-1) // we always expect the path excluding the last keyword to be valid
			return false
		}
	}
	return true
}
