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

func AveragingFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("AveragingFusion: Input cannot be nil")
	}

	b1 := opinion1.belief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	b := -1.0
	d := -1.0
	u := -1.0
	a := -1.0

	if u1 != 0 || u2 != 0 {

		b = (b1*u2 + b2*u1) / (u1 + u2)
		u = 2 * u1 * u2 / (u1 + u2)
		a = (a1 + a2) / 2

	} else {

		b = 0.5 * (b1 + b2)
		u = 0
		a = 0.5 * (a1 + a2)
	}

	d = 1 - b - u

	return NewOpinion(b, d, u, a)
}
