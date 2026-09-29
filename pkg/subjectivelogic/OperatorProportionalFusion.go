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

func ProportionalFusion(opinion1 *Opinion, opinion2 *Opinion, proportion1, proportion2 float64) (Opinion, error) {
	const eps = 1e-6

	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("ProportionalFusion: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("ProportionalFusion: Inputs cannot be null opinions")
	}

	// Checking if the proportions sum to 1
	if math.Abs(proportion1+proportion2-1.0) > eps {
		return Opinion{}, errors.New("ProportionalFusion: proportions must sum to 1")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	b := proportion1*b1 + proportion2*b2
	d := proportion1*d1 + proportion2*d2
	u := proportion1*u1 + proportion2*u2
	a := proportion1*a1 + proportion2*a2

	// Numerical cleanup: clamp into [0, 1]
	b = math.Max(0, math.Min(1, b))
	d = math.Max(0, math.Min(1, d))
	u = math.Max(0, math.Min(1, u))
	a = math.Max(0, math.Min(1, a))

	return NewOpinion(b, d, u, a)
}