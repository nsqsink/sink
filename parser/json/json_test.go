package json

import (
	"reflect"
	"testing"

	"github.com/nsqsink/sink/config"
	"github.com/nsqsink/sink/contract"
)

func Test_extractVariable(t *testing.T) {
	type args struct {
		subStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1: array case with string jason value",
			args: args{
				subStr: `$somethings[0].item","field":"value"`,
			},
			want: "$somethings[0].item",
		},
		{
			name: "Test 2: array case with integer json value",
			args: args{
				subStr: `$somethings[0].item,"field":"value"`,
			},
			want: "$somethings[0].item",
		},
		{
			name: "Test 3: integer json value",
			args: args{
				subStr: `$something.item.data,"field":"value"`,
			},
			want: "$something.item.data",
		},
		{
			name: "Test 4: array value string",
			args: args{
				subStr: `$something[0]","field":"value"`,
			},
			want: "$something[0]",
		},
		{
			name: "Test 5: array value integer",
			args: args{
				subStr: `$something[0],"field":"value"`,
			},
			want: "$something[0]",
		},
		{
			name: "Test 6: bracket array closing",
			args: args{
				subStr: `$something[0],"field":{"subfield":"value"}`,
			},
			want: "$something[0]",
		},
		{
			name: "Test 7: bracket closing }",
			args: args{
				subStr: `$something}`,
			},
			want: "$something",
		},
		{
			name: "Test 8: bracket closing array",
			args: args{
				subStr: `$something]`,
			},
			want: "$something",
		},
		{
			name: "Test 9: direct array",
			args: args{
				subStr: `$[0].something.item.data,"field":"value"`,
			},
			want: "$[0].something.item.data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVariable(tt.args.subStr); got != tt.want {
				t.Errorf("extractVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_breakdownVariable(t *testing.T) {
	type args struct {
		variable string
	}
	tests := []struct {
		name     string
		args     args
		wantPart []string
	}{
		{
			name: "test empty string",
			args: args{
				variable: "",
			},
			wantPart: []string{},
		},
		{
			name: "test no dot",
			args: args{
				variable: "$first",
			},
			wantPart: []string{"first"},
		},
		{
			name: "test 1 dot",
			args: args{
				variable: "$first.second",
			},
			wantPart: []string{"first", "second"},
		},
		{
			name: "test 2 dot",
			args: args{
				variable: "$first.second.third",
			},
			wantPart: []string{"first", "second", "third"},
		},
		{
			name: "test 2 dot with array",
			args: args{
				variable: "$first.second[1].third",
			},
			wantPart: []string{"first", "second", "[1]", "third"},
		},
		{
			name: "test 2 dot with array 2D",
			args: args{
				variable: "$first.second[1][2].third",
			},
			wantPart: []string{"first", "second", "[1]", "[2]", "third"},
		},
		{
			name: "test direct array",
			args: args{
				variable: "$[0].field",
			},
			wantPart: []string{"[0]", "field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPart := breakdownVariable(tt.args.variable); !reflect.DeepEqual(gotPart, tt.wantPart) {
				t.Errorf("breakdownVariable() = %v, want %v", gotPart, tt.wantPart)
			}
		})
	}
}

func TestModule_Parse(t *testing.T) {
	type fields struct {
		template          string
		templateVariables []Variable
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantParsed []byte
		wantErr    bool
	}{
		{
			name: "Test simple string",
			fields: fields{
				template: `{"field":"$data.field"}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field",
						orderedComponents: []string{"data", "field"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":"value"}}`),
			},
			wantParsed: []byte(`{"field":"value"}`),
		},
		{
			name: "Test simple string array",
			fields: fields{
				template: `{"field":"$data.field[0]"}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0]",
						orderedComponents: []string{"data", "field", "[0]"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":["value"]}`),
			},
			wantParsed: []byte(`{"field":"value"}`),
		},
		{
			name: "Test simple string array 3D with subfield",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]"}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}]}}`),
			},
			wantParsed: []byte(`{"field":"value"}`),
		},
		{
			name: "Test simple string array 3D with subfield + simple integer",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":$data.number}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}],"number":1}}`),
			},
			wantParsed: []byte(`{"field":"value","second_field":1}`),
		},
		{
			name: "Test simple string array 3D with subfield + simple float",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":$data.number}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}],"number":1.12}}`),
			},
			wantParsed: []byte(`{"field":"value","second_field":1.12}`),
		},
		{
			name: "Test simple string array 3D with subfield + simple float (mapped as string)",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":"$data.number"}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}],"number":1.12}}`),
			},
			wantParsed: []byte(`{"field":"value","second_field":"1.12"}`),
		},
		{
			name: "Test simple string array 3D with subfield + simple float (mapped to array string)",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}],"number":1.12}}`),
			},
			wantParsed: []byte(`{"field":"value","second_field":["1.12"]}`),
		},
		{
			name: "Test error parsing - not valid source data (json)",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`1`),
			},
			wantParsed: nil,
			wantErr:    true,
		},
		{
			name: "Test error parsing, not valid value type",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}],"number":[1.12]}}`),
			},
			wantParsed: nil,
			wantErr:    true,
		},
		{
			name: "Test error parsing, not valid value not found",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}]}}`),
			},
			wantParsed: nil,
			wantErr:    true,
		},
		{
			name: "Test error parsing, not valid value not found",
			fields: fields{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`{"data":{"field":[{"subfield": [[["value"]]]}]}}`),
			},
			wantParsed: nil,
			wantErr:    true,
		},
		{
			name: "Test direct array",
			fields: fields{
				template: `{"field":"$[0].data.field[0].subfield[0][0][0]","second_field":["$[0].data.number"]}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$[0].data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"[0]", "data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$[0].data.number",
						orderedComponents: []string{"[0]", "data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`[{"data":{"field":[{"subfield": [[["value"]]]}],"number":1.12}}]`),
			},
			wantParsed: []byte(`{"field":"value","second_field":["1.12"]}`),
		},
		{
			name: "Test direct array v2",
			fields: fields{
				template: `$[0].data.field[0].subfield[0][0][0],$[0].data.number`,
				templateVariables: []Variable{
					Variable{
						variable:          "$[0].data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"[0]", "data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$[0].data.number",
						orderedComponents: []string{"[0]", "data", "number"},
					},
				},
			},
			args: args{
				data: []byte(`[{"data":{"field":[{"subfield": [[["value"]]]}],"number":1.12}}]`),
			},
			wantParsed: []byte(`value,1.12`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Module{
				template:          tt.fields.template,
				templateVariables: tt.fields.templateVariables,
			}
			gotParsed, err := m.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParsed, tt.wantParsed) {
				t.Errorf("Module.Parse() = %v, want %v", string(gotParsed), string(tt.wantParsed))
			}
		})
	}
}

func Test_extractTemplateVariables(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test success",
			args: args{
				data: `{"field":"$data.field[0].subfield[0][0][0]","second_field":["$data.number"]}`,
			},
			want: []string{"$data.field[0].subfield[0][0][0]", "$data.number"},
		},
		{
			name: "Test empty data",
			args: args{
				data: ``,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTemplateVariables(tt.args.data)

			if len(got) != len(tt.want) {
				t.Errorf("extractTemplateVariables() len result = %+v, want %+v", len(got), len(tt.want))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTemplateVariables() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_preProcessTemplate(t *testing.T) {
	type args struct {
		template string
	}
	tests := []struct {
		name                  string
		args                  args
		wantTemplateVariables []Variable
	}{
		{
			name: "Test simple string array 3D with subfield + simple integer",
			args: args{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":$data.number}`,
			},
			wantTemplateVariables: []Variable{
				Variable{
					variable:          "$data.field[0].subfield[0][0][0]",
					orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
				},
				Variable{
					variable:          "$data.number",
					orderedComponents: []string{"data", "number"},
				},
			},
		},
		{
			name: "Test simple string array 3D with subfield + simple float (mapped as string)",
			args: args{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":"$data.number"}`,
			},
			wantTemplateVariables: []Variable{
				Variable{
					variable:          "$data.field[0].subfield[0][0][0]",
					orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
				},
				Variable{
					variable:          "$data.number",
					orderedComponents: []string{"data", "number"},
				},
			},
		},
		{
			name: "Test simple string array",
			args: args{
				template: `{"field":"$data.field[0]"}`,
			},
			wantTemplateVariables: []Variable{
				Variable{
					variable:          "$data.field[0]",
					orderedComponents: []string{"data", "field", "[0]"},
				},
			},
		},
		{
			name: "Test simple string",
			args: args{
				template: `{"field":"$data.field"}`,
			},
			wantTemplateVariables: []Variable{
				Variable{
					variable:          "$data.field",
					orderedComponents: []string{"data", "field"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTemplateVariables := preProcessTemplate(tt.args.template); !reflect.DeepEqual(gotTemplateVariables, tt.wantTemplateVariables) {
				t.Errorf("preProcessTemplate() = %v, want %v", gotTemplateVariables, tt.wantTemplateVariables)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		cfg config.Parser
	}
	tests := []struct {
		name    string
		args    args
		want    contract.Parser
		wantErr bool
	}{
		{
			name: "Test error empty template",
			args: args{
				cfg: config.Parser{},
			},
			wantErr: true,
		},
		{
			name: "Test happy flow",
			args: args{
				cfg: config.Parser{
					Type:     "json",
					Template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":$data.number}`,
				},
			},
			want: Module{
				template: `{"field":"$data.field[0].subfield[0][0][0]","second_field":$data.number}`,
				templateVariables: []Variable{
					Variable{
						variable:          "$data.field[0].subfield[0][0][0]",
						orderedComponents: []string{"data", "field", "[0]", "subfield", "[0]", "[0]", "[0]"},
					},
					Variable{
						variable:          "$data.number",
						orderedComponents: []string{"data", "number"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
