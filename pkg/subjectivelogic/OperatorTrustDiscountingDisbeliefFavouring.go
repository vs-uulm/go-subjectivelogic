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
	"math"
)

func TrustDiscountingDisbeliefFavouring(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("OpTrustDiscDisbelief: Input cannot be nil")
	}
	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("OpTrustDiscDisbelief: Inputs cannot be null opinions")
	}
	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate
	
	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	beliefMass := b1 + b2
	uncertaintyMass := u1 + u2

	var b, d, u, a float64

	if beliefMass == 0.0 && uncertaintyMass == 0.0 {
		b = 0
		d = 1
		u = 0
		a = a1
	} else {
		d = math.Max(d1, d2)
		b = (1 - d) * beliefMass / (beliefMass + uncertaintyMass)
		u = (1 - d) * (uncertaintyMass / (beliefMass + uncertaintyMass))

		if d1 >= d2 {
			a = a1
		} else {
			a = a2
		}
	}


	return NewOpinion(b, d, u, a)
}