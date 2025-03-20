//Copyright 2025 Institute of Distributed Systems, Ulm University
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
	"errors"

)

// This is an experimental operator and is only for testing purposes.
// func ImpactFusion(opinion1 *Opinion, opinion2 *Opinion, impactFactor float64) (Opinion, error) {
func ImpactFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {

	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("ImpactFusion: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("ImpactFusion: Inputs cannot be null opinions")
	}

	impactFactor := 0.8

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	i1 := impactFactor
	i2 := 1 - i1

	b := -1.0
	d := -1.0
	u := -1.0
	a := -1.0

	b = b1*i1 + b2*i2
	d = d1*i1 + d2*i2
	u = u1*i1 + u2*i2
	a = a1*i1 + a2*i2 // This has yet to be tested whether it makes sense.

	return NewOpinion(b, d, u, a)
}
