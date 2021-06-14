package slice

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"fmt"
)

func ExampleExists01() {
	b, i, e := Exists([]int{1, 2, 3}, 3)
	fmt.Println(b, i, e)
	// Output: true 2 <nil>
}

func ExampleExists02() {
	b, i, e := Exists([]int{1, 2, 3}, 4)
	fmt.Println(b, i, e)
	// Output: false -1 <nil>
}

func ExampleExists03() {
	b, i, e := Exists("2", '2')
	fmt.Println(b, i, e)
	// Output:
	// false -1 method 'exists' is designed for a slice and can not operate on string
}
