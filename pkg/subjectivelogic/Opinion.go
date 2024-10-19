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

package subjectivelogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

/*
Precision defines the maximum deviance each value of an Opinion can have for the Opinion to still be regarded as a valid Binomial Opinion.
*/
const Precision float64 = 0.000000000001

/*
Opinion represents a Binomial Opinion from Subjective Logic.
It is recommended to only generate new opinions using the NewOpinions function, as this will ensure the generated Opinion to be a valid Binomial Opinion,
which the Operators in this library are designed to work with.
*/
type Opinion struct {
	belief      float64
	disbelief   float64
	uncertainty float64
	baseRate    float64
}

/*
An interface to query binomial opinions from subjective logic.
Provided to increase compatibility with other SL implementations.
*/
type QueryableOpinion interface {
	Belief() float64
	Disbelief() float64
	Uncertainty() float64
	BaseRate() float64
	String() string
}

/*
NewOpinion takes four float64 values and outputs an *Opinion as well as an Error.
In case a valid Opinion can be formed, it will be returned and the error will be nil.
If the input Values violate the requirements for a valid Opinion, the *Opinion will be nil and an error will be returned.
For a valid Opinion, all input values i must fulfill 0 <= i <= 1 and for the first tree inputs b, d, u, the statement b+d+u = 1 must hold.
*/
func NewOpinion(belief, disbelief, uncertainty, baseRate float64) (Opinion, error) {
	if !checkInput(belief, disbelief, uncertainty, baseRate) {
		return Opinion{}, errors.New("NewOpinion: Invalid Input")
	}
	op := Opinion{belief: belief, disbelief: disbelief, uncertainty: uncertainty, baseRate: baseRate}
	return op, nil
}

/*
Belief is called onto an *Opinion o and returns o.belief.
*/
func (opinion *Opinion) Belief() float64 {
	return opinion.belief
}

/*
Disbelief is called onto an *Opinion o and returns o.disbelief.
*/
func (opinion *Opinion) Disbelief() float64 {
	return opinion.disbelief
}

/*
Uncertainty is called onto an *Opinion o and returns o.uncertainty.
*/
func (opinion *Opinion) Uncertainty() float64 {
	return opinion.uncertainty
}

/*
BaseRate is called onto an *Opinion o and returns o.baseRate.
*/
func (opinion *Opinion) BaseRate() float64 {
	if opinion == nil {
		panic("BaseRate(): method call from nil pointer")
	}
	return opinion.baseRate
}

/*
Modify is called onto an *Opinion o and requires four float64 values as input.
If the input values form a valid Opinion, the values of o will be changed to the input values.
If o is nil or the input values do not form a valid Opinion, o is left unchanged and an error is returned.
For a valid Opinion, all input values i must fulfill 0 <= i <= 1 and for the first tree inputs b, d, u, the statement b+d+u = 1 must hold.
*/
func (opinion *Opinion) Modify(belief, disbelief, uncertainty, baseRate float64) error {
	if !checkInput(belief, disbelief, uncertainty, baseRate) {
		return errors.New("Modify: Invalid Input")
	}
	if opinion == nil {
		return errors.New("Modify: opinion is nil")
	}
	opinion.belief = belief
	opinion.disbelief = disbelief
	opinion.uncertainty = uncertainty
	opinion.baseRate = baseRate

	return nil
}

/*
ProjectedProbability is called onto an *Opinion o and calculates the projected probability of o.
*/
func (opinion *Opinion) ProjectedProbability() float64 {
	if opinion == nil {
		panic("ProjectedProbability(): method call from nil pointer")
	}
	return opinion.belief + opinion.uncertainty*opinion.baseRate
}

/*
Compare is called onto an Opinion o1 and compares it with the input Opinion o2.
If the values of o1 and o2 each match with a maximum difference of Precision, true is returned.
Otherwise, false is returned.
*/
func (opinion1 Opinion) Compare(opinion2 Opinion) bool {
	return math.Abs(opinion1.belief-opinion2.belief) < Precision &&
		math.Abs(opinion1.disbelief-opinion2.disbelief) < Precision &&
		math.Abs(opinion1.uncertainty-opinion2.uncertainty) < Precision &&
		math.Abs(opinion1.baseRate-opinion2.baseRate) < Precision
}

/*
Copy is called onto an *Opinion o1 and returns a new *Opinion o2 that has the same values as o1.
*/
func (opinion1 *Opinion) Copy() *Opinion {

	return &Opinion{opinion1.belief, opinion1.disbelief, opinion1.uncertainty, opinion1.baseRate}
}

/*
String is called onto an *Opinion o and returns a string containing the values of o.
If o is nil, "nil" is returned.
*/
func (opinion *Opinion) String() string {
	if opinion == nil {
		return "nil"
	}
	return fmt.Sprint(opinion.belief) + ", " + fmt.Sprint(opinion.disbelief) + ", " +
		fmt.Sprint(opinion.uncertainty) + ", " + fmt.Sprint(opinion.baseRate)
}

/*
checkInput takes four float64 values as input and returns true, if they form a valid Opinion.
Otherwise, false is returned.
*/
func checkInput(b, d, u, a float64) bool {
	if math.Abs(1-(b+d+u)) < 3*Precision &&
		0 <= b && b <= 1 &&
		0 <= d && d <= 1 &&
		0 <= u && u <= 1 &&
		0 <= a && a <= 1 {
		return true
	}
	return false
}

func (opinion *Opinion) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Belief      float64 `json:"belief"`
		Disbelief   float64 `json:"disbelief"`
		Uncertainty float64 `json:"uncertainty"`
		BaseRate    float64 `json:"base_rate"`
	}{
		Belief:      opinion.Belief(),
		Disbelief:   opinion.Disbelief(),
		Uncertainty: opinion.Uncertainty(),
		BaseRate:    opinion.BaseRate(),
	})
}
