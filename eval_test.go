// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package jet

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"text/template"
)

var evalTemplateSet = NewSet()

func evalTestCase(t *testing.T, variables map[string]interface{}, context interface{}, testName, testContent, testExpected string) {
	buff := bytes.NewBuffer(nil)
	tt, err := evalTemplateSet.loadTemplate(testName, testContent)
	if err != nil {
		t.Errorf("Parsing error: %s %s %s", err.Error(), testName, testContent)
		return
	}
	err = tt.Execute(buff, variables, context)
	if err != nil {
		t.Errorf("Eval error: %q executing %s", err.Error(), testName)
		return
	}
	result := buff.String()
	if result != testExpected {
		t.Errorf("Result error expected %q got %q on %s", testExpected, result, testName)
	}
}

func TestEvalTextNode(t *testing.T) {
	evalTestCase(t, nil, nil, "textNode", `hello {*Buddy*} World`, `hello  World`)
}

type User struct {
	Name, Email string
}

func (user *User) Format(str string) string {
	return fmt.Sprintf(str, user.Name, user.Email)
}

func (user *User) GetName() string {
	return user.Name
}

func TestEvalActionNode(t *testing.T) {
	var data = make(map[string]interface{})
	data["lower"] = strings.ToLower
	data["upper"] = strings.ToUpper
	data["repeat"] = strings.Repeat

	data["user"] = &User{
		"José Santos", "email@example.com",
	}

	evalTestCase(t, nil, nil, "actionNode", `hello {{"world"}}`, `hello world`)
	evalTestCase(t, data, nil, "actionNode_func", `hello {{lower: "WORLD"}}`, `hello world`)
	evalTestCase(t, data, nil, "actionNode_funcPipe", `hello {{lower: "WORLD" |upper}}`, `hello WORLD`)
	evalTestCase(t, data, nil, "actionNode_funcPipeArg", `hello {{lower: "WORLD-" |upper|repeat: 2}}`, `hello WORLD-WORLD-`)
	evalTestCase(t, data, nil, "actionNode_Field", `Oi {{ user.Name }}`, `Oi José Santos`)
	evalTestCase(t, data, nil, "actionNode_Field2", `Oi {{ user.Name }}<{{ user.Email }}>`, `Oi José Santos<email@example.com>`)
	evalTestCase(t, data, nil, "actionNode_Method", `Oi {{ user.Format: "%s<%s>" }}`, `Oi José Santos<email@example.com>`)

	evalTestCase(t, data, nil, "actionNode_Add", `{{ 2+1 }}`, fmt.Sprint(2+1))
	evalTestCase(t, data, nil, "actionNode_Add3", `{{ 2+1+4 }}`, fmt.Sprint(2+1+4))
	evalTestCase(t, data, nil, "actionNode_Add3Minus", `{{ 2+1+4-3 }}`, fmt.Sprint(2+1+4-3))
	evalTestCase(t, data, nil, "actionNode_Mult", `{{ 4*4 }}`, fmt.Sprint(4*4))
	evalTestCase(t, data, nil, "actionNode_MultAdd", `{{ 2+4*4 }}`, fmt.Sprint(2+4*4))
	evalTestCase(t, data, nil, "actionNode_MultAdd1", `{{ 4*2+4 }}`, fmt.Sprint(4*2+4))
	evalTestCase(t, data, nil, "actionNode_MultAdd2", `{{ 2+4*2+4 }}`, fmt.Sprint(2+4*2+4))
	evalTestCase(t, data, nil, "actionNode_MultFloat", `{{ 1*1.23 }}`, fmt.Sprint(1*1.23))
	evalTestCase(t, data, nil, "actionNode_Mod", `{{ 3%2 }}`, fmt.Sprint(3%2))
	evalTestCase(t, data, nil, "actionNode_MultMod", `{{ (1*3)%2 }}`, fmt.Sprint((1*3)%2))
	evalTestCase(t, data, nil, "actionNode_MultDivMod", `{{ (2*5)/ 3 %1 }}`, fmt.Sprint((2*5)/3%1))

	evalTestCase(t, data, nil, "actionNode_Comparation", `{{ (2*5)==10 }}`, fmt.Sprint((2*5) == 10))
	evalTestCase(t, data, nil, "actionNode_Comparatation2", `{{ (2*5)==5 }}`, fmt.Sprint((2*5) == 5))
	evalTestCase(t, data, nil, "actionNode_Logical", `{{ (2*5)==5 || true }}`, fmt.Sprint((2*5) == 5 || true))
	evalTestCase(t, data, nil, "actionNode_Logical2", `{{ (2*5)==5 || false }}`, fmt.Sprint((2*5) == 5 || false))

	evalTestCase(t, data, nil, "actionNode_NumericCmp", `{{ 5*5 > 2*12.5 }}`, fmt.Sprint(5*5 > 2*12.5))
	evalTestCase(t, data, nil, "actionNode_NumericCmp1", `{{ 5*5 >= 2*12.5 }}`, fmt.Sprint(5*5 >= 2*12.5))
	evalTestCase(t, data, nil, "actionNode_NumericCmp1", `{{ 5 * 5 > 2 * 12.5 == 5 * 5 > 2 * 12.5 }}`, fmt.Sprint((5*5 > 2*12.5) == (5*5 > 2*12.5)))
}

func TestEvalIfNode(t *testing.T) {
	var data = make(map[string]interface{})
	data["lower"] = strings.ToLower
	data["upper"] = strings.ToUpper
	data["repeat"] = strings.Repeat

	data["user"] = &User{
		"José Santos", "email@example.com",
	}

	evalTestCase(t, data, nil, "ifNode_simples", `{{if true}}hello{{end}}`, `hello`)
	evalTestCase(t, data, nil, "ifNode_else", `{{if false}}hello{{else}}world{{end}}`, `world`)
	evalTestCase(t, data, nil, "ifNode_elseif", `{{if false}}hello{{else if true}}world{{end}}`, `world`)
	evalTestCase(t, data, nil, "ifNode_elseif_else", `{{if false}}hello{{else if false}}world{{else}}buddy{{end}}`, `buddy`)
}

func TestEvalBlockYieldIncludeNode(t *testing.T) {
	var data = make(map[string]interface{})

	data["user"] = &User{
		"José Santos", "email@example.com",
	}

	evalTestCase(t, data, nil, "Block_simple", `{{block hello "Buddy" }}Hello {{ . }}{{end}},{{yield hello user.Name}}`, `Hello Buddy,Hello José Santos`)
	evalTestCase(t, data, nil, "Block_Extends", `{{extends "Block_simple"}}{{block hello "Buddy" }}Hey {{ . }}{{end}}`, `Hey Buddy,Hey José Santos`)
	evalTestCase(t, data, nil, "Block_Import", `{{import "Block_simple"}}{{yield hello "Buddy"}}`, `Hello Buddy`)
	evalTestCase(t, data, nil, "Block_Import", `{{import "Block_simple"}}{{yield hello "Buddy"}}`, `Hello Buddy`)

	evalTemplateSet.LoadTemplate("Block_ImportInclude1", `{{yield hello "Buddy"}}`)
	evalTestCase(t, data, nil, "Block_ImportInclude", `{{ import "Block_simple"}}{{include "Block_ImportInclude1"}}`, `Hello Buddy`)
}

func TestEvalRangeNode(t *testing.T) {

	var data = make(map[string]interface{})

	data["users"] = []User{
		{"Mario Santos", "mario@gmail.com"},
		{"Joel Silva", "joelsilva@gmail.com"},
		{"Luis Santana", "luis.santana@gmail.com"},
	}
	const resultString = `<h1>Mario Santos<small>mario@gmail.com</small></h1><h1>Joel Silva<small>joelsilva@gmail.com</small></h1><h1>Luis Santana<small>luis.santana@gmail.com</small></h1>`
	evalTestCase(t, data, nil, "Range_Expression", `{{range users}}<h1>{{.Name}}<small>{{.Email}}</small></h1>{{end}}`, resultString)
	evalTestCase(t, data, nil, "Range_ExpressionValue", `{{range user:=users}}<h1>{{user.Name}}<small>{{user.Email}}</small></h1>{{end}}`, resultString)
}

func TestEvalDefaultFuncs(t *testing.T) {
	evalTestCase(t, nil, nil, "DefaultFuncs_map", `{{ map("name","José Santos").name }}`, `José Santos`)
	evalTestCase(t, nil, nil, "DefaultFuncs_map1", `{{ map(@name,"José Santos", @email,"email@example.pt").email }}`, `email@example.pt`)
}

type devNull struct{}

func (*devNull) Write(_ []byte) (int, error) {
	return 0, nil
}

var stdSet = template.New("base")

func dummy(a string) string {
	return a
}
func init() {
	stdSet.Funcs(template.FuncMap{"dummy": dummy})
	_, err := stdSet.Parse(`
		{{define "actionNode_dummy"}}hello {{dummy "WORLD"}}{{end}}
		{{define "noAllocFn"}}hello {{ "José" }} {{1}} {{ "José" }} {{end}}
	`)
	if err != nil {
		println(err.Error())
	}
	evalTemplateSet.AddGlobal("dummy", dummy)
	evalTemplateSet.LoadTemplate("actionNode_dummy", `hello {{dummy("WORLD")}}`)
	evalTemplateSet.LoadTemplate("noAllocFn", `hello {{ "José" }} {{1}} {{ "José" }}`)
}

var ww io.Writer = (*devNull)(nil)

func BenchmarkSimpleAction(b *testing.B) {
	t, _ := evalTemplateSet.GetTemplate("actionNode_dummy")
	for i := 0; i < b.N; i++ {
		err := t.Execute(ww, nil, nil)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

func BenchmarkSimpleActionNoAlloc(b *testing.B) {
	t, _ := evalTemplateSet.GetTemplate("noAllocFn")
	for i := 0; i < b.N; i++ {
		t.Execute(ww, nil, nil)
	}
}

func BenchmarkSimpleActionStd(b *testing.B) {
	t := stdSet.Lookup("actionNode_dummy")
	for i := 0; i < b.N; i++ {
		t.Execute(ww, nil)
	}
}

func BenchmarkSimpleActionStdNoAlloc(b *testing.B) {
	t := stdSet.Lookup("noAllocFn")
	for i := 0; i < b.N; i++ {
		t.Execute(ww, nil)
	}
}
