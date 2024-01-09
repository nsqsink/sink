package json

import (
	"errors"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/nsqsink/sink/config"
	"github.com/nsqsink/sink/contract"
	"github.com/nsqsink/sink/pkg/stack"
)

type (
	Module struct {
		template          string
		templateVariables []Variable
	}

	Variable struct {
		variable          string   // example: $user.id or $users[0].name
		orderedComponents []string // example: []string{user, id}, or []string{users[0], name}
	}
)

const varSymbol string = "$"

func New(cfg config.Parser) (contract.Parser, error) {
	if cfg.Template == "" {
		return nil, errors.New("not valid template")
	}

	// parse template
	templateVariables := preProcessTemplate(cfg.Template)

	m := Module{
		template:          cfg.Template,
		templateVariables: templateVariables,
	}

	return m, nil
}

// Parse parse template based on given data into new payload based on template
func (m Module) Parse(data []byte) (parsed []byte, err error) {
	payload := m.template

	// loop for each variables
	for _, variable := range m.templateVariables {
		// get variable components
		components := variable.GetComponents()

		// parse from given data based on variable
		val, valType, _, errParsed := jsonparser.Get(data, components...)
		if errParsed != nil {
			err = errParsed
			return
		}

		if valType.String() != "string" && valType.String() != "number" && valType.String() != "boolean" {
			err = fmt.Errorf("not valid template parser on data %s on variable %s", data, variable.GetVariable())
			return
		}

		// replace
		payload = strings.Replace(payload, variable.GetVariable(), string(val), 1)
	}

	return []byte(payload), nil
}

// extractTemplateVariables return list of variables in template
func extractTemplateVariables(data string) []string {
	var templates []string

	// preprocess
	data = strings.ReplaceAll(data, " ", "")

	// return if not valid
	if data == "" {
		return []string{}
	}

	temp := strings.Split(data, varSymbol)

	// normalize
	if data[0] != varSymbol[0] {
		temp = temp[1:]
	}

	for _, subStr := range temp {
		cleanSubStr := extractVariable(varSymbol + subStr)
		if cleanSubStr != "" {
			templates = append(templates, cleanSubStr)
		}
	}

	return templates
}

func extractVariable(subStr string) string {
	var cleanSubStr string

	// use to determine if the variables are in the array or is accessing array
	// example case: field: "$somethings[0].item" or "field: [$something.item]"
	stackArray := stack.New()

	for idx, r := range subStr {
		// if end of chars, add to the template the complete substring
		if idx == (len(subStr) - 1) {
			cleanSubStr = subStr

			if isEndOfVariable(r) {
				cleanSubStr = subStr[:idx]
			}
			break
		}

		if r == '[' {
			stackArray.Push('[')
		}

		if r == ']' {
			stackArray.Pop()
			continue
		}

		if isEndOfVariable(r) && stackArray.Size() == 0 {
			cleanSubStr = subStr[:idx]
			break
		}
	}

	return strings.TrimSpace(cleanSubStr)
}

// breakdownVariable into array
func breakdownVariable(variable string) (parts []string) {
	if variable == "" {
		parts = []string{}
		return
	}

	// remove dollar symbol
	variable = variable[1:]

	// breakdown by dot '.'
	tempParts := strings.Split(variable, ".")

	for _, tempPart := range tempParts {
		// split by [] to split array
		tempSplitted := strings.Split(tempPart, "[")

		for _, tempPart2 := range tempSplitted {
			if tempPart2 == "" {
				continue
			}

			// check if a bracket
			if strings.Contains(tempPart2, "]") {
				parts = append(parts, "["+tempPart2)
				continue
			}

			parts = append(parts, tempPart2)
		}
	}

	return
}

// isEndOfVariable return true if given char value is a symbol
func isEndOfVariable(val rune) bool {
	return val == ',' || val == '}' || val == ']' || val == '"'
}

// preProcessTemplate
func preProcessTemplate(template string) (templateVariables []Variable) {
	// find all of dollar symbols
	variables := extractTemplateVariables(template)

	// extract each templates
	for _, variable := range variables {
		subVars := breakdownVariable(variable)

		templateVariables = append(templateVariables, Variable{
			variable:          variable,
			orderedComponents: subVars,
		})
	}

	return
}

func (v Variable) GetVariable() string {
	return v.variable
}

func (v Variable) GetComponents() []string {
	return v.orderedComponents
}
