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
	"errors"
)

func ConstraintFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("ConstraintFusion: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("ConstraintFusion: Inputs cannot be null opinions")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	har := b1*u2 + b2*u1 + b1*b2
	con := b1*d2 + b2*d1

	// 1 - con is computed as the sum of the remaining terms of (b1+d1+u1)*(b2+d2+u2) = 1, since subtracting a con close
	// to 1 from 1 loses most significant digits
	disHar := d1*u2 + d2*u1 + d1*d2
	nonConflicting := har + disHar + u1*u2

	if con == 1 || nonConflicting == 0 {
		return Opinion{}, errors.New("ConstraintFusion: mathematically possible only if input opinions are not conflicting and do not result in Con = 1")
	}

	b := har / nonConflicting
	u := u1 * u2 / nonConflicting
	d := disHar / nonConflicting

	a := -1.0
	if u1+u2 < 2 {
		a = (a1*(1-u1) + a2*(1-u2)) / (2 - u1 - u2)
	} else {
		a = (a1 + a2) / 2
	}

	return NewOpinionWithTolerance(b, d, u, a)
}
