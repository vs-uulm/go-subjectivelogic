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

func CumulativeEpistemicFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// CumulativeFusion already checks for nil pointers and null opinions, so
	// those checks are inherited by delegating to it here.
	omega3, err := CumulativeFusion(opinion1, opinion2)
	if err != nil {
		return Opinion{}, err
	}

	p := omega3.ProjectedProbability()
	a := omega3.BaseRate()

	// A base rate of exactly 0 or 1 makes uMax's division undefined.
	if a == 0 || a == 1 {
		return Opinion{}, errors.New("CumulativeEpistemicFusion: base rate must not be equal to 0 or 1.")
	}

	uMax := math.Min(p/a, (1-p)/(1-a))
	bMax := p - a*uMax
	dMax := math.Max(0, 1-bMax-uMax)

	return NewOpinion(bMax, dMax, uMax, a)
}
