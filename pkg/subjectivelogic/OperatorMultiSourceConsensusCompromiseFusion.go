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

// Values of a binomial domain {x, not x} as bit sets, so that intersection and union are bitwise and/or
const (
	domainEmpty    = 0b00
	domainTrue     = 0b01
	domainFalse    = 0b10
	domainComplete = 0b11
)

// maxConsensusCompromiseOpinions limits the number of opinions, as 4^n combinations are enumerated (about 40 s for
// 15 opinions)
const maxConsensusCompromiseOpinions = 15

/*
MultiSourceConsensusCompromiseFusion fuses an arbitrary number of opinions with the Consensus & Compromise Fusion
operator for multiple sources as defined by van der Heijden, Kopp and Kargl, "Multi-Source Fusion Operations in
Subjective Logic", FUSION 2018 (Definition 5, Eq. 7).
CC fusion is only semi-associative, so fusing three or more opinions by chaining ConsensusCompromiseFusion gives
order-dependent results.
The implementation follows ccCollectionFuse of subjective-logic-java and enumerates all 4^n combinations of the values
{empty set, x, not x, whole domain} of the n sources, so its runtime is O(n * 4^n).
Unlike ccCollectionFuse, it supports dogmatic opinions.
Eq. 7 is undefined if the compromise mass is zero. This happens for identical opinions, for which the consensus is
returned, as ConsensusCompromiseFusion does. It also happens with two or more dogmatic opinions if every combination of
residual belief and disbelief contains a zero; the result is then the limit for uncertainties approaching zero at the
same rate, i.e. the remaining mass is distributed according to the residual belief and disbelief of the dogmatic opinions
with equal weights. This follows the convention of Jøsang, "Categories of Belief Fusion", Journal of Advances in
Information Fusion 13(2), 2018, for dogmatic opinions in the other fusion operators (gamma = 1/N); where the dogmatic
opinions contradict each other, the limit depends on how the uncertainties approach zero.
Unlike in the binary operator, vacuous opinions are not neutral: they reduce the consensus to zero.
For two opinions the result equals that of ConsensusCompromiseFusion, and a single opinion is returned unchanged.
*/
func MultiSourceConsensusCompromiseFusion(opinions []Opinion) (Opinion, error) {
	if err := checkMultiSourceInput(opinions, "MultiSourceConsensusCompromiseFusion"); err != nil {
		return Opinion{}, err
	}
	if len(opinions) > maxConsensusCompromiseOpinions {
		return Opinion{}, errors.New("MultiSourceConsensusCompromiseFusion: Too many Opinions")
	}

	// Checking if the base rates match
	a := 0.0
	for _, opinion := range opinions {
		if math.Abs(opinion.baseRate-opinions[0].baseRate) > 1e-12 {
			return Opinion{}, errors.New("MultiSourceConsensusCompromiseFusion: opinions must have the same base rate")
		}
		a += opinion.baseRate
	}
	a /= float64(len(opinions))

	// Step 1: consensus phase, the belief and disbelief shared by all sources
	bCons, dCons := 1.0, 1.0
	for _, opinion := range opinions {
		bCons = math.Min(bCons, opinion.belief)
		dCons = math.Min(dCons, opinion.disbelief)
	}

	n := len(opinions)
	bRes := make([]float64, n)
	dRes := make([]float64, n)
	uPre := 1.0
	for i, opinion := range opinions {
		bRes[i] = math.Max(0, opinion.belief-bCons)
		dRes[i] = math.Max(0, opinion.disbelief-dCons)
		uPre *= opinion.uncertainty
	}

	// Step 2: compromise phase, consisting of 4 sub-sums
	bComp, dComp, xComp := 0.0, 0.0, 0.0 // xComp is the belief on the whole domain, which later becomes uncertainty

	// Sub-sum 1: the residual belief of each source, weighted by the uncertainty of all other sources.
	// The product is computed without dividing by the source's own uncertainty, which may be 0.
	for i := range opinions {
		uOthers := 1.0
		for j, other := range opinions {
			if j != i {
				uOthers *= other.uncertainty
			}
		}
		bComp += bRes[i] * uOthers
		dComp += dRes[i] * uOthers
	}

	// Sub-sums 2 to 4: every source contributes one of the values of the domain, encoded with two bits per source
	for combination := uint64(0); combination < uint64(1)<<(2*n); combination++ {
		// The product of the residual belief of every source on its value. The residual belief on the empty set and
		// on the whole domain is 0, so only combinations of x and not x contribute.
		intersection, union, product := domainComplete, domainEmpty, 1.0
		for i := 0; i < n; i++ {
			value := int(combination>>(2*i)) & domainComplete
			intersection &= value
			union |= value
			switch value {
			case domainTrue:
				product *= bRes[i]
			case domainFalse:
				product *= dRes[i]
			default:
				product = 0
			}
		}

		switch {
		case intersection == domainTrue:
			// Sub-sum 2: the intersection is x
			bComp += product
		case intersection == domainFalse:
			// Sub-sum 2: the intersection is not x
			dComp += product
		case intersection == domainEmpty && union == domainComplete:
			// Sub-sum 4: the values are conflicting, which becomes belief on the whole domain
			xComp += product
		}
		// Sub-sum 3 (the union is x or not x and the intersection is not empty) is always 0 for binomial opinions,
		// and if the union is x or not x but the intersection is empty, the combination contains the empty set
	}

	// Step 3: normalisation phase
	compromiseMass := bComp + dComp + xComp
	if compromiseMass == 0 {
		// With two or more dogmatic opinions, sub-sum 1 is zero. Its leading terms for uncertainties approaching zero
		// at the same rate are the residual belief and disbelief of the dogmatic opinions, weighted equally.
		bDogmatic, dDogmatic, dogmatic := 0.0, 0.0, 0
		for i, opinion := range opinions {
			if opinion.uncertainty == 0 {
				bDogmatic += bRes[i]
				dDogmatic += dRes[i]
				dogmatic++
			}
		}
		if dogmatic < 2 || bDogmatic+dDogmatic == 0 {
			return NewOpinionWithTolerance(bCons, dCons, 1-bCons-dCons, a)
		}
		free := 1 - bCons - dCons
		return NewOpinionWithTolerance(bCons+free*bDogmatic/(bDogmatic+dDogmatic), dCons+free*dDogmatic/(bDogmatic+dDogmatic), 0, a)
	}
	eta := (1 - bCons - dCons - uPre) / compromiseMass

	return NewOpinionWithTolerance(bCons+eta*bComp, dCons+eta*dComp, uPre+eta*xComp, a)
}
