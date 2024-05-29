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
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("ConstraintFusion: Input cannot be nil")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	con := b1*d2 + d1*b2

	if con == 1 {
		return Opinion{}, errors.New("ConstraintFusion: Invalid arguments: Con = 1")
	}

	har1 := b1*u2 + b2*u1 + b1*b2
	har2 := d1*u2 + d2*u1 + d1*d2

	b := har1 / (1 - con)
	d := har2 / (1 - con)
	u := u1 * u2 / (1 - con)

	a := -1.0
	if u1 == 1 && u2 == 1 {
		a = (a1 + a2) / 2
	} else {
		a = (a1*(1-u1) + a2*(1-u2)) / (2 - u1 - u2)
	}

	return NewOpinion(b, d, u, a)
}
