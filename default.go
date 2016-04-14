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
	"fmt"
	"strings"
)

var defaultVariables = map[string]interface{}{
	"lower":     strings.ToLower,
	"upper":     strings.ToUpper,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"repeat":    strings.Repeat,
	"replace":   strings.Replace,
	"map":       newMap,
}

func newMap(values ...interface{}) (nmap map[string]interface{}) {
	if len(values) % 2 > 0 {
		panic("new map: invalid number of arguments on call to map")
	}
	nmap = make(map[string]interface{})

	for i := 0; i < len(values); i += 2 {
		nmap[fmt.Sprint(values[i])] = values[i + 1]
	}
	return
}
