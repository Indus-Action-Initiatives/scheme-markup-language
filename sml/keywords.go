package sml

const (
	kwAllowed = iota + 1000
	kwAllowedValue
	kwCalendarType
	kwControlType
	kwDataType
	kwDefaultSelected
	kwDefaultValue
	kwEmptyState
	kwEnd
	kwEndLabel
	kwEndType
	kwGranularity
	kwLabel
	kwLvalue
	kwMax
	kwMin
	kwMultiSelect
	kwMultiValue
	kwName
	kwParameter
	kwPrecision
	kwRvalue
	kwSource
	kwSourceType
	kwStart
	kwStartLabel
	kwStartType
	kwType
	kwValue
	kwValues
	kwVar
	kwVariables
)

var linglingKeywords = map[string]keyword{
	"allowed_value": kwAllowedValue,
	"calendar_type": kwCalendarType,
	"control_type":  kwControlType,
	"data_type":     kwDataType,
	"default_value": kwDefaultValue,
	"empty_state":   kwEmptyState,
	"end":           kwEnd,
	"end_label":     kwEndLabel,
	"end_type":      kwEndType,
	"granularity":   kwGranularity,
	"label":         kwLabel,
	"lvalue":        kwLvalue,
	"max":           kwMax,
	"min":           kwMin,
	"multi_select":  kwMultiSelect,
	"name":          kwName,
	"parameter":     kwParameter,
	"precision":     kwPrecision,
	"rvalue":        kwRvalue,
	"source":        kwSource,
	"source_type":   kwSourceType,
	"start":         kwStart,
	"start_label":   kwStartLabel,
	"start_type":    kwStartType,
	"type":          kwType,
	"value":         kwValue,
	"values":        kwValues,
	"var":           kwVar,
	"variables":     kwVariables,
	// synonyms
	"allowed":          kwAllowedValue,
	"default_selected": kwValues,
}

var linglingKeyword2string = makeKeywordToString(linglingKeywords, []string{})
