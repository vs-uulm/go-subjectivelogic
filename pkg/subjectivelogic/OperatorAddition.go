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

func Addition(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("Addition: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("Addition: Inputs cannot be null opinions")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	b := -1.0
	d := -1.0
	u := -1.0
	a := -1.0

	if a1 == 0 && a2 == 0 {
		return Opinion{}, errors.New("Addition: Base rates cannot be both equal to 0")

	} else {
		b = b1 + b2
		d = (a1*(d1-b2) + a2*(d2-b1)) / (a1 + a2)
		u = (a1*u1 + a2*u2) / (a1 + a2)
		a = a1 + a2
	}

	o, err := NewOpinion(b, d, u, a)

	if err != nil {
		return Opinion{}, errors.New("Addition: Check the validity of your input values")
	}

	if b > 1 {
		return Opinion{}, errors.New("Addition: Sum of beliefs cannot exceed 1")
	} else if a > 1 {
		return Opinion{}, errors.New("Addition: Sum of base rates cannot exceed 1")
	}

	return o, err

}
