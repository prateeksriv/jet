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
	"io/ioutil"
	"path"
	"testing"
)

var parseSet = NewSet("./testData")

type ParserTestCase struct {
	*testing.T
}

func (t ParserTestCase) ExpectPrintName(name, input, output string) {
	template, err := parseSet.parse(name, input)
	if err != nil {
		t.Errorf("%q %s", input, err.Error())
		return
	}
	expected := template.String()
	if expected != output {
		t.Errorf("Unexpected tree on %s Got:\n%s\nExpected: \n%s\n", name, expected, output)
	}
}

func (t ParserTestCase) ExpectPrint(input, output string) {
	t.ExpectPrintName("", input, output)
}

func (t ParserTestCase) TestPrintFile(file string) {
	content, err := ioutil.ReadFile(path.Join("./testData", file))
	if err != nil {
		t.Errorf("file %s not found", file)
		return
	}
	parts := bytes.Split(content, []byte("==="))
	t.ExpectPrintName(file, string(bytes.TrimSpace(parts[0])), string(bytes.TrimSpace(parts[1])))
}

func (t ParserTestCase) ExpectPrintSame(input string) {
	t.ExpectPrint(input, input)
}

func TestParseTemplateAndImport(t *testing.T) {
	p := ParserTestCase{t}
	p.TestPrintFile("extends.jet")
	p.TestPrintFile("imports.jet")
}

func TestParseTemplateControl(t *testing.T) {
	p := ParserTestCase{t}
	p.TestPrintFile("if.jet")
	p.TestPrintFile("range.jet")
}

func TestParseTemplateExpressions(t *testing.T) {
	p := ParserTestCase{t}
	p.TestPrintFile("simple_expression.jet")
	p.TestPrintFile("additive_expression.jet")
	p.TestPrintFile("multiplicative_expression.jet")
}

func TestParseTemplateBlockYield(t *testing.T) {
	p := ParserTestCase{t}
	p.TestPrintFile("block_yield.jet")
}
