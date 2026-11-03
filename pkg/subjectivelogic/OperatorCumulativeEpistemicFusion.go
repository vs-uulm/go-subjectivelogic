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

/*
CumulativeEpistemicFusion fuses two opinions with Epistemic Cumulative Belief Fusion as defined by Jøsang, Subjective
Logic, Springer 2016, Definition 12.6: the opinions are fused with CumulativeFusion, and the result is
uncertainty-maximised, i.e. as much belief or disbelief as possible is converted into uncertainty while the projected
probability is preserved.
*/
func CumulativeEpistemicFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// CumulativeFusion already checks for nil pointers and null opinions, so
	// those checks are inherited by delegating to it here.
	omega3, err := CumulativeFusion(opinion1, opinion2)
	if err != nil {
		return Opinion{}, err
	}
	return uncertaintyMaximised(omega3, "CumulativeEpistemicFusion")
}

// uncertaintyMaximised returns the uncertainty-maximised opinion with the same projected probability and base rate
func uncertaintyMaximised(opinion Opinion, operator string) (Opinion, error) {
	p := opinion.ProjectedProbability()
	a := opinion.BaseRate()

	// A base rate of exactly 0 or 1 makes uMax's division undefined.
	if a == 0 || a == 1 {
		return Opinion{}, errors.New(operator + ": base rate must not be equal to 0 or 1.")
	}

	uMax := math.Min(p/a, (1-p)/(1-a))
	bMax := p - a*uMax
	dMax := math.Max(0, 1-bMax-uMax)

	return NewOpinionWithTolerance(bMax, dMax, uMax, a)
}
