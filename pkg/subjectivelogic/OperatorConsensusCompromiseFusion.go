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

func ConsensusCompromiseFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	const eps = 1e-12
	strictBaseRate := true

	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("ConsensusCompromiseFusion: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("ConsensusCompromiseFusion: Inputs cannot be null opinions")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	// CC fusion is defined for opinions with the same base rate. If
	// strictBaseRate is false, the base rates are averaged instead.
	if strictBaseRate && math.Abs(a1-a2) > eps {
		return Opinion{}, errors.New(
			"ConsensusCompromiseFusion: opinions must have the same base rate unless strictBaseRate is false",
		)
	}
	a := 0.5 * (a1 + a2)

	// Consensus phase: shared belief/disbelief mass both sources agree on
	bCons := math.Min(b1, b2)
	dCons := math.Min(d1, d2)

	// Residual belief/disbelief after removing consensus
	b1Res := math.Max(0, b1-bCons)
	b2Res := math.Max(0, b2-bCons)

	d1Res := math.Max(0, d1-dCons)
	d2Res := math.Max(0, d2-dCons)

	consensusMass := bCons + dCons

	// Compromise phase: residual belief/disbelief, weighted by the other source's uncertainty, plus common residual belief/disbelief
	bComp := b1Res*u2 + b2Res*u1 + b1Res*b2Res
	dComp := d1Res*u2 + d2Res*u1 + d1Res*d2Res

	// Conflicting residual belief becomes belief in the whole binary domain, which is later transferred to uncertainty
	xComp := b1Res*d2Res + d1Res*b2Res

	uPre := u1 * u2
	compromiseMass := bComp + dComp + xComp

	// Normalization phase: rescale the compromise contributions so the result is a valid opinion
	eta := 1.0
	if compromiseMass > eps {
		eta = (1 - consensusMass - uPre) / compromiseMass
	}

	b := math.Max(0, math.Min(1, bCons+eta*bComp))
	d := math.Max(0, math.Min(1, dCons+eta*dComp))
	u := math.Max(0, 1-b-d)

	return NewOpinion(b, d, u, a)
}