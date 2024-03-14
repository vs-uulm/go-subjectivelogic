//Copyright 2024 Institute of Distributed Systems, Ulm University
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package sllib

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

const nrOfRandomValues = 20
const nrOfValidOpinions = 9

var testValuesOpinions = [19][5]float64{
	//valid Opinions
	{1, 0, 0, 0.5, 1},
	{0, 1, 0, 0.5, 0},
	{0, 0, 1, 0.5, 0.5},
	{0.6, 0.3, 0.1, 0, 0.6},
	{0.091, 0.604, 0.305, 1, 0.396},
	{0.53, 0.227, 0.243, 1.000, 0.773},
	{0.000473, 0.555506, 0.444021, 0, 0.000473},
	{0.004, 0.5950000001, 0.4009999999, 0.00000000004334, 0.00400000001737934},
	{0, 0.999999999999999, 0.000000000000001, 0.5000000000000000000000000000000000000003, 0}, //PP not actually 0 but smaller than 1e-15

	//b+d+u != 1
	{0.5, 0.3, 0.7, 0.5},
	{0.3, 0.3, 0.2, 0.5},

	//b+d+u = 1 but one value x ist x<0
	{-0.1, 0.3, 0.8, 0.5},
	{0.7, -0.3, 0.6, 0},
	{0.5, 0.3, -0.2, 1},
	{0.1, 0.5, 0.4, -0.5},

	//one value x ist x>1
	{1.1, 0.3, 0.8, 0.5},
	{-0.7, 1.3, 0.4, 0.5},
	{0.5, 0.3, 1.2, 0.5},
	{0.1, 0.5, 0.4, 1.5}}

func TestNewOpinion(t *testing.T) {
	for i := 0; i < len(testValuesOpinions); i++ {
		_, err := NewOpinion(testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
		if i < nrOfValidOpinions && err != nil {
			t.Errorf("False negative on i = %d | Error: %s | Values: %f, %f, %f, %f", i, err, testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
		}
		if i >= nrOfValidOpinions && err == nil {
			t.Errorf("False positive on i = %d | Values: %f, %f, %f, %f", i, testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
		}
	}

	for i := 0; i < nrOfRandomValues; i++ {
		var values = []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		_, err := NewOpinion(values[0], values[1], values[2], values[3])

		valid := math.Abs(1-(values[0]+values[1]+values[2])) < 3*Precision
		if valid && err != nil {
			t.Errorf("False negative on random try %d | Error: %s | Values: %f, %f, %f, %f", i, err, values[0], values[1], values[2], values[3])
		}
		if !valid && err == nil {
			t.Errorf("False positive on random try %d | Values: %f, %f, %f, %f", i, values[0], values[1], values[2], values[3])
		}
	}
}

func TestOpinion_GetBelief(t *testing.T) {
	var o *Opinion

	gotPanic := false
	defer func() {
		if err := recover(); err != nil {
			gotPanic = true
		}
	}()
	_ = o.GetBelief()
	if !gotPanic {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	for i := 0; i < nrOfValidOpinions; i++ {
		o = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
		b := o.GetBelief()

		if math.Abs(b-o.belief) >= Precision {
			t.Errorf("Incorrect output at i = %d: Output %f | Expected: %f | Deviation: %v", i, b,
				o.belief, strconv.FormatFloat(math.Abs(b-o.belief), 'e', -1, 64))
		}

	}
}

func TestOpinion_GetDisbelief(t *testing.T) {
	var o *Opinion

	gotPanic := false
	defer func() {
		if err := recover(); err != nil {
			gotPanic = true
		}
	}()
	_ = o.GetDisbelief()
	if !gotPanic {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	for i := 0; i < nrOfValidOpinions; i++ {
		o = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
		d := o.GetDisbelief()

		if math.Abs(d-o.disbelief) >= Precision {
			t.Errorf("Incorrect output at i = %d: Output %f | Expected: %f | Deviation: %v", i, d, o.disbelief,
				strconv.FormatFloat(math.Abs(d-o.disbelief), 'e', -1, 64))
		}
	}
}

func TestOpinion_GetUncertainty(t *testing.T) {
	var o *Opinion

	gotPanic := false
	defer func() {
		if err := recover(); err != nil {
			gotPanic = true
		}
	}()
	_ = o.GetUncertainty()
	if !gotPanic {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	for i := 0; i < nrOfValidOpinions; i++ {
		o = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
		u := o.GetUncertainty()

		if math.Abs(u-o.uncertainty) >= Precision {
			t.Errorf("Incorrect output at i = %d: Output %f | Expected: %f | Deviation: %v", i, u, o.uncertainty,
				strconv.FormatFloat(math.Abs(u-o.uncertainty), 'e', -1, 64))
		}
	}
}

func TestOpinion_GetBaseRate(t *testing.T) {
	var o *Opinion

	gotPanic := false
	defer func() {
		if err := recover(); err != nil {
			gotPanic = true
		}
	}()
	_ = o.GetBaseRate()
	if !gotPanic {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	for i := 0; i < nrOfValidOpinions; i++ {
		o = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
		a := o.GetBaseRate()

		if math.Abs(a-o.baseRate) >= Precision {
			t.Errorf("Incorrect output at i = %d: Output %f | Expected: %f | Deviation: %v", i, a, o.baseRate,
				strconv.FormatFloat(math.Abs(a-o.baseRate), 'e', -1, 64))
		}
	}
}

func TestOpinion_Modify(t *testing.T) {

	var o *Opinion
	err := o.Modify(testValuesOpinions[0][0], testValuesOpinions[0][1], testValuesOpinions[0][2], testValuesOpinions[0][3])
	if err == nil {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	o = &Opinion{1, 0, 0, 0}
	for i := 0; i < len(testValuesOpinions); i++ {

		err = o.Modify(testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])

		if i < nrOfValidOpinions {
			expected := &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
			if !o.Compare(expected) {
				t.Errorf("Invalid output on i = %d | Output: %f, %f, %f, %f | Expected %f, %f, %f, %f", i,
					o.belief, o.disbelief, o.uncertainty, o.baseRate, testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
			} else {
				if err != nil {
					t.Errorf("False negative on i = %d | Error: %s | Values: %f, %f, %f, %f", i, err,
						testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
				}
			}
		}

		if i >= nrOfValidOpinions && err == nil {
			t.Errorf("False positive on i = %d | Values: %f, %f, %f, %f", i, testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3])
		}
	}
}

func TestOpinion_ProjProb(t *testing.T) {

	var o *Opinion

	gotPanic := false
	defer func() {
		if err := recover(); err != nil {
			gotPanic = true
		}
	}()
	_ = o.ProjProb()
	if !gotPanic {
		t.Errorf("Invalid call from \"nil\" passed undetected")
	}

	for i := 0; i < nrOfValidOpinions; i++ {
		o = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}

		pp := o.ProjProb()
		expected := testValuesOpinions[i][4]

		if math.Abs(pp-expected) >= ((1 + testValuesOpinions[i][1] + testValuesOpinions[i][2] + Precision) * Precision) {
			t.Errorf("Invalid output on i = %d: Output: %f| Expected %f", i, pp, expected)
		}
	}
}

func TestOpinion_Compare(t *testing.T) {

	var o1, o2 *Opinion
	if !o1.Compare(o2) {
		t.Errorf("Incorrect output: Output: %t, Expected: %t | Opinion1: \"nil\" | Opinion2 \"nil\" ", true, false)
	}

	o1 = &Opinion{1, 0, 0, 1}
	if o1.Compare(o2) {
		t.Errorf("Incorrect output: Output: %t, Expected: %t | Opinion1: 1, 0, 0, 1 | Opinion2 \"nil\" ", false, true)
	}
	if o2.Compare(o1) {
		t.Errorf("Incorrect output: Output: %t, Expected: %t | Opinion1: \"nil\" | Opinion2 1, 0, 0, 1 ", false, true)
	}

	for i := 0; i < len(testValuesOpinions); i++ {
		if i == 8 {
			i++
		} // Skip element 9

		for j := 0; j < len(testValuesOpinions); j++ {
			if j == 8 {
				j++
			} // Skip element 9

			o1 = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
			o2 = &Opinion{testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3]}
			passed1 := o1.Compare(o2)
			passed2 := o2.Compare(o1)

			if passed1 != passed2 {
				t.Errorf("Inconsistent output on i = %d, j = %d | Opinion1: %f, %f, %f, %f | Opinion2 %f, %f, %f, %f | o1.Compare(o2) = %t, o2.Compare(o1) = %t", i, j,
					testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3],
					testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3], passed1, passed2)
			} else {
				if (i == j && !passed1) || (i != j && passed1) {
					t.Errorf("Incorrect output on i = %d, j = %d | Output: %t, Expected: %t | Opinion1: %f, %f, %f, %f | Opinion2 %f, %f, %f, %f ", i, j, passed1, i == j,
						testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3],
						testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3])
				}
			}

		}
	}
	//Test element 2 == element 9, as they are within tolerance from each other
	i := 1
	j := 8
	o1 = &Opinion{testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3]}
	o2 = &Opinion{testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3]}
	passed1 := o1.Compare(o2)
	passed2 := o2.Compare(o1)

	if passed1 != passed2 {
		t.Errorf("Inconsistent output on i = %d, j = %d | Opinion1: %f, %f, %f, %f | Opinion2 %f, %f, %f, %f | o1.Compare(o2) = %t, o2.Compare(o1) = %t", i, j,
			testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3],
			testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3], passed1, passed2)
	} else {
		if !passed1 {
			t.Errorf("Incorrect output on i = %d, j = %d | Output: %t, Expected: %t | Opinion1: %f, %f, %f, %f | Opinion2 %f, %f, %f, %f ", i, j, passed1, true,
				testValuesOpinions[i][0], testValuesOpinions[i][1], testValuesOpinions[i][2], testValuesOpinions[i][3],
				testValuesOpinions[j][0], testValuesOpinions[j][1], testValuesOpinions[j][2], testValuesOpinions[j][3])
		}
	}
}

func TestOpinion_ToString(t *testing.T) {

	var o *Opinion
	str := o.ToString()
	if str != "nil" {
		t.Errorf("Invalid call from \"nil\" passed undetected | Output: %v", str)
	}

	o = &Opinion{testValuesOpinions[0][0], testValuesOpinions[0][1], testValuesOpinions[0][2], testValuesOpinions[0][3]}
	str = o.ToString()
	expected := "1, 0, 0, 0.5"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}

	o = &Opinion{testValuesOpinions[1][0], testValuesOpinions[1][1], testValuesOpinions[1][2], testValuesOpinions[1][3]}
	str = o.ToString()
	expected = "0, 1, 0, 0.5"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}

	o = &Opinion{testValuesOpinions[2][0], testValuesOpinions[2][1], testValuesOpinions[2][2], testValuesOpinions[2][3]}
	str = o.ToString()
	expected = "0, 0, 1, 0.5"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}

	o = &Opinion{testValuesOpinions[5][0], testValuesOpinions[5][1], testValuesOpinions[5][2], testValuesOpinions[5][3]}
	str = o.ToString()
	expected = "0.53, 0.227, 0.243, 1"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}

	o = &Opinion{testValuesOpinions[6][0], testValuesOpinions[6][1], testValuesOpinions[6][2], testValuesOpinions[6][3]}
	str = o.ToString()
	expected = "0.000473, 0.555506, 0.444021, 0"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}

	o = &Opinion{testValuesOpinions[7][0], testValuesOpinions[7][1], testValuesOpinions[7][2], testValuesOpinions[7][3]}
	str = o.ToString()
	expected = "0.004, 0.5950000001, 0.4009999999, 4.334e-11"
	if str != expected {
		t.Errorf("Icorrect output | Output: %v | Expected: %v", str, expected)
	}
}
