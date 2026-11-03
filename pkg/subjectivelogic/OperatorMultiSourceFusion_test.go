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
	"bufio"
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"
)

type multiSourceOperator struct {
	name   string
	multi  func([]Opinion) (Opinion, error)
	binary func(*Opinion, *Opinion) (Opinion, error)
}

var multiSourceOperators = []multiSourceOperator{
	{"MultiSourceCumulativeFusion", MultiSourceCumulativeFusion, CumulativeFusion},
	{"MultiSourceAveragingFusion", MultiSourceAveragingFusion, AveragingFusion},
	{"MultiSourceWeightedFusion", MultiSourceWeightedFusion, WeightedFusion},
	{"MultiSourceCumulativeEpistemicFusion", MultiSourceCumulativeEpistemicFusion, CumulativeEpistemicFusion},
}

func TestMultiSourceFusion(t *testing.T) {
	tableI := []Opinion{{0.1, 0.3, 0.6, 0.5}, {0.4, 0.2, 0.4, 0.5}, {0.7, 0.1, 0.2, 0.5}}
	tests := []struct {
		name     string
		operator func([]Opinion) (Opinion, error)
		opinions []Opinion
		want     Opinion
		wantErr  bool
	}{
		//example of Table I in Jøsang, Wang and Zhang (2017): (0.651, 0.209, 0.140) and (0.509, 0.164, 0.327)
		{"CumulativeTableI", MultiSourceCumulativeFusion, tableI, Opinion{0.224 / 0.344, 0.072 / 0.344, 0.048 / 0.344, 0.5}, false},
		{"AveragingTableI", MultiSourceAveragingFusion, tableI, Opinion{0.224 / 0.44, 0.072 / 0.44, 0.144 / 0.44, 0.5}, false},
		//computed by subjective-logic-java
		{"WeightedTableI", MultiSourceWeightedFusion, tableI, Opinion{0.56216216216, 0.14594594595, 0.29189189189, 0.5}, false},

		//one dogmatic opinion dominates, including its base rate (Case II of "Categories of Belief Fusion")
		{"CumulativeOneDogmatic", MultiSourceCumulativeFusion, []Opinion{{0.6, 0.4, 0, 0.3}, {0.2, 0.3, 0.5, 0.7}, {0.1, 0.1, 0.8, 0.9}}, Opinion{0.6, 0.4, 0, 0.3}, false},
		{"AveragingOneDogmatic", MultiSourceAveragingFusion, []Opinion{{0.6, 0.4, 0, 0.3}, {0.2, 0.3, 0.5, 0.7}, {0.1, 0.1, 0.8, 0.9}}, Opinion{0.6, 0.4, 0, 0.3}, false},
		{"WeightedOneDogmatic", MultiSourceWeightedFusion, []Opinion{{0.6, 0.4, 0, 0.3}, {0.2, 0.3, 0.5, 0.7}, {0.1, 0.1, 0.8, 0.9}}, Opinion{0.6, 0.4, 0, 0.3}, false},

		//several dogmatic opinions are averaged with equal weights, including their base rates only
		{"CumulativeTwoDogmatic", MultiSourceCumulativeFusion, []Opinion{{0.6, 0.4, 0, 0.2}, {0.2, 0.8, 0, 0.4}, {0.3, 0.3, 0.4, 0.9}}, Opinion{0.4, 0.6, 0, 0.3}, false},
		{"AveragingTwoDogmatic", MultiSourceAveragingFusion, []Opinion{{0.6, 0.4, 0, 0.2}, {0.2, 0.8, 0, 0.4}, {0.3, 0.3, 0.4, 0.9}}, Opinion{0.4, 0.6, 0, 0.3}, false},
		{"WeightedTwoDogmatic", MultiSourceWeightedFusion, []Opinion{{0.6, 0.4, 0, 0.2}, {0.2, 0.8, 0, 0.4}, {0.3, 0.3, 0.4, 0.9}}, Opinion{0.4, 0.6, 0, 0.3}, false},

		//vacuous opinions are neutral for cumulative and weighted fusion, but not for averaging fusion
		{"CumulativeVacuous", MultiSourceCumulativeFusion, []Opinion{{0.6, 0.3, 0.1, 0.5}, {0, 0, 1, 0.5}, {0, 0, 1, 0.5}}, Opinion{0.6, 0.3, 0.1, 0.5}, false},
		{"AveragingVacuous", MultiSourceAveragingFusion, []Opinion{{0.6, 0.3, 0.1, 0.5}, {0, 0, 1, 0.5}, {0, 0, 1, 0.5}}, Opinion{0.6 / 1.2, 0.3 / 1.2, 0.3 / 1.2, 0.5}, false},
		{"WeightedVacuous", MultiSourceWeightedFusion, []Opinion{{0.6, 0.3, 0.1, 0.5}, {0, 0, 1, 0.5}, {0, 0, 1, 0.5}}, Opinion{0.6, 0.3, 0.1, 0.5}, false},

		//only vacuous opinions give the vacuous opinion with the average base rate
		{"CumulativeAllVacuous", MultiSourceCumulativeFusion, []Opinion{{0, 0, 1, 0.2}, {0, 0, 1, 0.4}, {0, 0, 1, 0.9}}, Opinion{0, 0, 1, 0.5}, false},
		{"AveragingAllVacuous", MultiSourceAveragingFusion, []Opinion{{0, 0, 1, 0.2}, {0, 0, 1, 0.4}, {0, 0, 1, 0.9}}, Opinion{0, 0, 1, 0.5}, false},
		{"WeightedAllVacuous", MultiSourceWeightedFusion, []Opinion{{0, 0, 1, 0.2}, {0, 0, 1, 0.4}, {0, 0, 1, 0.9}}, Opinion{0, 0, 1, 0.5}, false},

		//uncertainty maximisation of (0.224, 0.072, 0.048) / 0.344 with projected probability P = 0.248 / 0.344 and a = 0.5
		//converts all disbelief into uncertainty: u = 2 (1 - P), b = P - u / 2
		{"EpistemicTableI", MultiSourceCumulativeEpistemicFusion, tableI, Opinion{2*0.248/0.344 - 1, 0, 2 - 2*0.248/0.344, 0.5}, false},
		{"EpistemicBaseRate1", MultiSourceCumulativeEpistemicFusion, []Opinion{{0.6, 0.3, 0.1, 1}, {0.2, 0.3, 0.5, 1}}, Opinion{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.operator(tt.opinions)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !compareWithin(got, tt.want, 1e-10) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiSourceFusionInputs(t *testing.T) {
	for _, op := range multiSourceOperators {
		if _, err := op.multi(nil); err == nil {
			t.Errorf("%s(nil) returned no error", op.name)
		}
		if _, err := op.multi([]Opinion{}); err == nil {
			t.Errorf("%s([]) returned no error", op.name)
		}
		if _, err := op.multi([]Opinion{{0.6, 0.3, 0.1, 0.5}, {0, 0, 0, 0}}); err == nil {
			t.Errorf("%s with a null opinion returned no error", op.name)
		}
		single := Opinion{0.6, 0.3, 0.1, 0.4}
		want := single
		if op.name == "MultiSourceCumulativeEpistemicFusion" {
			// uncertainty-maximised with P = 0.64: u = min(0.64 / 0.4, 0.36 / 0.6), b = P - 0.4 u
			want = Opinion{0.4, 0, 0.6, 0.4}
		}
		if got, err := op.multi([]Opinion{single}); err != nil || !compareWithin(got, want, 1e-12) {
			t.Errorf("%s(%v) = %v, %v, want %v", op.name, single, got, err, want)
		}
		// vacuous opinions are neutral in cumulative and weighted fusion, also for a single other opinion
		if op.name != "MultiSourceAveragingFusion" {
			got, _ := op.multi([]Opinion{single})
			withVacuous, err := op.multi([]Opinion{single, {0, 0, 1, 0.4}, {0, 0, 1, 0.4}})
			if err != nil || !compareWithin(got, withVacuous, 1e-12) {
				t.Errorf("%s with vacuous opinions = %v, %v, want %v", op.name, withVacuous, err, got)
			}
		}
	}
}

// TestMultiSourceFusionJava compares MultiSourceCumulativeFusion and MultiSourceWeightedFusion with results of
// subjective-logic-java for random opinions with equal base rates, including vacuous and dogmatic opinions.
func TestMultiSourceFusionJava(t *testing.T) {
	file, err := os.Open("testdata/multi_source_fusion.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	cases := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "|")
		in := parseFloats(t, fields[0])
		opinions := make([]Opinion, int(in[0]))
		for i := range opinions {
			opinions[i] = Opinion{in[1+4*i], in[2+4*i], in[3+4*i], in[4+4*i]}
		}
		for i, op := range []multiSourceOperator{
			{"MultiSourceCumulativeFusion", MultiSourceCumulativeFusion, nil},
			{"MultiSourceWeightedFusion", MultiSourceWeightedFusion, nil},
		} {
			out := parseFloats(t, fields[1+i])
			want := Opinion{out[0], out[1], out[2], opinions[0].baseRate}
			got, err := op.multi(opinions)
			if err != nil || !compareWithin(got, want, 1e-10) {
				t.Errorf("%s(%v) got = %v, %v, want %v", op.name, opinions, got, err, want)
			}
		}
		cases++
	}
	if cases == 0 {
		t.Fatal("no test cases found")
	}
}

// TestMultiSourceFusionProperties compares the multi-source operators with the binary operators for two opinions,
// checks that the result does not depend on the order of the opinions, and compares them with the evidence they
// correspond to.
func TestMultiSourceFusionProperties(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 20000; i++ {
		opinions := randomSpecialOpinions(r, 2+r.Intn(7))

		for _, op := range multiSourceOperators {
			got, err := op.multi(opinions)
			if len(opinions) == 2 {
				want, wantErr := op.binary(&opinions[0], &opinions[1])
				// Epistemic fusion rejects a fused base rate of exactly 0 or 1, which the binary and the multi-source
				// operator may compute with different rounding errors
				if cumulative, _ := CumulativeFusion(&opinions[0], &opinions[1]); op.name == "MultiSourceCumulativeEpistemicFusion" &&
					(err != nil) != (wantErr != nil) && math.Min(cumulative.baseRate, 1-cumulative.baseRate) < 1e-9 {
					continue
				}
				// With exactly one dogmatic opinion, averaging and weighted fusion take its base rate (Case II of
				// "Categories of Belief Fusion"), whereas the binary operators average the base rates
				if dogmatic := (opinions[0].uncertainty == 0) != (opinions[1].uncertainty == 0); dogmatic &&
					(op.name == "MultiSourceAveragingFusion" || op.name == "MultiSourceWeightedFusion") {
					want.baseRate = opinions[0].baseRate
					if opinions[1].uncertainty == 0 {
						want.baseRate = opinions[1].baseRate
					}
				}
				if (err != nil) != (wantErr != nil) || !compareWithin(got, want, 1e-9) {
					t.Fatalf("%s(%v) got = %v, %v, want %v, %v (binary operator)", op.name, opinions, got, err, want, wantErr)
				}
			}
			if err != nil {
				continue
			}
			shuffled := append([]Opinion(nil), opinions...)
			r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
			if gotShuffled, _ := op.multi(shuffled); !compareWithin(got, gotShuffled, 1e-9) {
				t.Fatalf("%s() depends on order: %v vs %v", op.name, got, gotShuffled)
			}
		}

		// for non-dogmatic opinions, the operators correspond to the sum, the average and the confidence-weighted
		// average of the evidence r = W * b / u (with W = 2)
		dogmatic := false
		var sumB, sumD, avgB, avgD, confB, confD, conf float64
		for _, o := range opinions {
			dogmatic = dogmatic || o.uncertainty == 0
			if dogmatic {
				break
			}
			rb, rd := 2*o.belief/o.uncertainty, 2*o.disbelief/o.uncertainty
			sumB, sumD = sumB+rb, sumD+rd
			avgB, avgD = avgB+rb/float64(len(opinions)), avgD+rd/float64(len(opinions))
			c := 1 - o.uncertainty
			confB, confD, conf = confB+c*rb, confD+c*rd, conf+c
		}
		if dogmatic || conf == 0 {
			continue
		}
		evidence := func(rb, rd float64) Opinion {
			return Opinion{rb / (2 + rb + rd), rd / (2 + rb + rd), 2 / (2 + rb + rd), 0}
		}
		for _, check := range []struct {
			operator func([]Opinion) (Opinion, error)
			want     Opinion
		}{
			{MultiSourceCumulativeFusion, evidence(sumB, sumD)},
			{MultiSourceAveragingFusion, evidence(avgB, avgD)},
			{MultiSourceWeightedFusion, evidence(confB/conf, confD/conf)},
		} {
			got, _ := check.operator(opinions)
			got.baseRate = 0
			if !compareWithin(got, check.want, 1e-9) {
				t.Fatalf("fusion of %v got = %v, want %v (evidence)", opinions, got, check.want)
			}
		}
	}
}

// TestMultiSourceCumulativeFusionChaining checks that MultiSourceCumulativeFusion equals chaining CumulativeFusion for
// non-dogmatic opinions, for which cumulative fusion is associative, including the evidence-weighted base rate.
func TestMultiSourceCumulativeFusionChaining(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for i := 0; i < 5000; i++ {
		opinions := make([]Opinion, 2+r.Intn(7))
		for j := range opinions {
			u := 0.01 + 0.99*r.Float64()
			b := r.Float64() * (1 - u)
			opinions[j] = Opinion{b, math.Max(0, 1-u-b), u, r.Float64()}
		}
		want := opinions[0]
		for j := 1; j < len(opinions); j++ {
			want, _ = CumulativeFusion(&want, &opinions[j])
		}
		got, err := MultiSourceCumulativeFusion(opinions)
		if err != nil || !compareWithin(got, want, 1e-9) {
			t.Fatalf("MultiSourceCumulativeFusion(%v) got = %v, %v, want %v (chained)", opinions, got, err, want)
		}
	}
}

// TestMultiSourceCumulativeEpistemicFusionProperties checks that MultiSourceCumulativeEpistemicFusion has the
// projected probability of MultiSourceCumulativeFusion and the maximum uncertainty, i.e. belief or disbelief is 0.
func TestMultiSourceCumulativeEpistemicFusionProperties(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for i := 0; i < 5000; i++ {
		opinions := randomOpinions(r, 2+r.Intn(7), 0.05+0.9*r.Float64())
		cumulative, _ := MultiSourceCumulativeFusion(opinions)
		got, err := MultiSourceCumulativeEpistemicFusion(opinions)
		if err != nil {
			t.Fatalf("MultiSourceCumulativeEpistemicFusion(%v) error = %v", opinions, err)
		}
		if math.Abs(got.ProjectedProbability()-cumulative.ProjectedProbability()) > 1e-9 || math.Min(got.belief, got.disbelief) > 1e-9 {
			t.Fatalf("MultiSourceCumulativeEpistemicFusion(%v) got = %v, cumulative %v", opinions, got, cumulative)
		}
	}
}

// TestMultiSourceFusionManySources checks that the sums do not underflow for many sources: for 1000 copies each of
// (0.5, 0.1, 0.4) and (0.1, 0.5, 0.4), the products of the uncertainties are 0.4^2000, far below the smallest float64.
// The evidence of every source is (2.5, 0.5) or (0.5, 2.5).
func TestMultiSourceFusionManySources(t *testing.T) {
	var opinions []Opinion
	for i := 0; i < 1000; i++ {
		opinions = append(opinions, Opinion{0.5, 0.1, 0.4, 0.5}, Opinion{0.1, 0.5, 0.4, 0.5})
	}
	for _, check := range []struct {
		operator func([]Opinion) (Opinion, error)
		want     Opinion
	}{
		{MultiSourceCumulativeFusion, Opinion{3000.0 / 6002, 3000.0 / 6002, 2.0 / 6002, 0.5}},
		{MultiSourceAveragingFusion, Opinion{0.3, 0.3, 0.4, 0.5}},
		{MultiSourceWeightedFusion, Opinion{0.3, 0.3, 0.4, 0.5}},
	} {
		got, err := check.operator(opinions)
		if err != nil || !compareWithin(got, check.want, 1e-12) {
			t.Errorf("got = %v, %v, want %v", got, err, check.want)
		}
	}
}

func BenchmarkMultiSourceCumulativeFusion(b *testing.B) {
	opinions := randomOpinions(rand.New(rand.NewSource(1)), 1000, 0.5)
	b.ResetTimer()
	for range b.N {
		x, err := MultiSourceCumulativeFusion(opinions)
		if err != nil {
			b.Error(err)
		}
		sink = x
	}
}
