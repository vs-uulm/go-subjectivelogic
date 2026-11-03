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
	"strconv"
	"strings"
	"testing"
)

func TestMultiSourceConsensusCompromiseFusion(t *testing.T) {
	type args struct {
		opinions []Opinion
	}
	tests := []struct {
		name    string
		args    args
		want    Opinion
		wantErr bool
	}{
		//nil input
		{"TestMultiSourceConsensusCompromiseFusion1",
			args{nil},
			Opinion{},
			true,
		},
		{"TestMultiSourceConsensusCompromiseFusion2",
			args{[]Opinion{}},
			Opinion{},
			true,
		},

		//1 argument is returned unchanged
		{"TestMultiSourceConsensusCompromiseFusion3",
			args{[]Opinion{{0.6, 0.3, 0.1, 0.5}}},
			Opinion{0.6, 0.3, 0.1, 0.5},
			false,
		},

		//null opinion input
		{"TestMultiSourceConsensusCompromiseFusion4",
			args{[]Opinion{{1, 0, 0, 0.5}, {0, 0, 0, 0}, {0, 1, 0, 0.5}}},
			Opinion{},
			true,
		},

		//mismatched base rates
		{"TestMultiSourceConsensusCompromiseFusion5",
			args{[]Opinion{{0.5, 0.3, 0.2, 0.5}, {0.5, 0.3, 0.2, 0.5}, {0.5, 0.3, 0.2, 0.6}}},
			Opinion{},
			true,
		},

		//two arguments, equal to ConsensusCompromiseFusion
		{"TestMultiSourceConsensusCompromiseFusion6",
			args{[]Opinion{{0.6, 0.3, 0.1, 0.4}, {0.091, 0.604, 0.305, 0.4}}},
			Opinion{0.35484913523, 0.35166680866, 0.29348405610, 0.4},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion7",
			args{[]Opinion{{0.4, 0.0, 0.6, 0.5}, {.7, .0, .3, .5}}},
			Opinion{.82, .0, .18, .5},
			false,
		},

		//example of Table I in van der Heijden et al. (2018)
		{"TestMultiSourceConsensusCompromiseFusion8",
			args{[]Opinion{{0.1, 0.3, 0.6, 0.5}, {0.4, 0.2, 0.4, 0.5}, {0.7, 0.1, 0.2, 0.5}}},
			Opinion{0.62875, 0.18225, 0.189, 0.5},
			false,
		},

		//general tests, expected values computed by subjective-logic-java
		{"TestMultiSourceConsensusCompromiseFusion9",
			args{[]Opinion{{0.7, 0.2, 0.1, 0.5}, {0.1, 0.6, 0.3, 0.5}, {0.4, 0.3, 0.3, 0.5}}},
			Opinion{0.35018965517, 0.25956896552, 0.39024137931, 0.5},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion10",
			args{[]Opinion{{0.6, 0.3, 0.1, 0.4}, {0.091, 0.604, 0.305, 0.4}, {0.53, 0.227, 0.243, 0.4}}},
			Opinion{0.30383443692, 0.28767444946, 0.40849111361, 0.4},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion11",
			args{[]Opinion{{0.4, 0, 0.6, 0.5}, {0.7, 0, 0.3, 0.5}, {0.2, 0, 0.8, 0.5}}},
			Opinion{0.856, 0, 0.144, 0.5},
			false,
		},

		//vacuous opinions are not neutral, they reduce the consensus to zero
		{"TestMultiSourceConsensusCompromiseFusion12",
			args{[]Opinion{{0.6, 0.3, 0.1, 0.4}, {0.091, 0.604, 0.305, 0.4}, {0.53, 0.227, 0.243, 0.4}, {0, 0, 1, 0.4}}},
			Opinion{0.58473218685, 0.40785631315, 0.0074115, 0.4},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion13",
			args{[]Opinion{{0, 0, 1, 0.4}, {0.6, 0.3, 0.1, 0.4}, {0, 0, 1, 0.4}}},
			Opinion{0.6, 0.3, 0.1, 0.4},
			false,
		},

		//dogmatic opinions, where subjective-logic-java divides by zero
		{"TestMultiSourceConsensusCompromiseFusion14",
			args{[]Opinion{{0.7, 0.3, 0, 0.5}, {0.1, 0.6, 0.3, 0.5}, {0.4, 0.3, 0.3, 0.5}}},
			Opinion{0.4, 0.3, 0.3, 0.5},
			false,
		},
		//fully opposed dogmatic opinions resolve to total uncertainty
		{"TestMultiSourceConsensusCompromiseFusion15",
			args{[]Opinion{{1, 0, 0, 0.5}, {0, 1, 0, 0.5}, {0.5, 0.2, 0.3, 0.5}}},
			Opinion{0, 0, 1, 0.5},
			false,
		},

		//no compromise mass, where Eq. 7 is undefined: the consensus is returned
		{"TestMultiSourceConsensusCompromiseFusion16",
			args{[]Opinion{{0, 0, 1, 0.5}, {0, 0, 1, 0.5}, {0, 0, 1, 0.5}}},
			Opinion{0, 0, 1, 0.5},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion17",
			args{[]Opinion{{0.3, 0.5, 0.2, 0.5}, {0.3, 0.5, 0.2, 0.5}, {0.3, 0.5, 0.2, 0.5}}},
			Opinion{0.3, 0.5, 0.2, 0.5},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion18",
			args{[]Opinion{{1, 0, 0, 0.5}, {1, 0, 0, 0.5}, {1, 0, 0, 0.5}}},
			Opinion{1, 0, 0, 0.5},
			false,
		},

		//no compromise mass with several dogmatic opinions: limit for uncertainties approaching zero at the same rate,
		//which is unique where the residual mass of the dogmatic opinions is all belief or all disbelief
		{"TestMultiSourceConsensusCompromiseFusion19",
			args{[]Opinion{{1, 0, 0, 0.5}, {1, 0, 0, 0.5}, {0.5, 0, 0.5, 0.5}}},
			Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestMultiSourceConsensusCompromiseFusion20",
			args{[]Opinion{{0.7, 0.3, 0, 0.5}, {0.7, 0.3, 0, 0.5}, {0.5, 0.2, 0.3, 0.5}}},
			Opinion{0.7, 0.3, 0, 0.5},
			false,
		},
		//contradicting dogmatic opinions: the residual belief 0.3 and disbelief 0.3 are weighted equally
		{"TestMultiSourceConsensusCompromiseFusion21",
			args{[]Opinion{{0.7, 0.3, 0, 0.5}, {0.4, 0.6, 0, 0.5}, {0.4, 0.3, 0.3, 0.5}}},
			Opinion{0.55, 0.45, 0, 0.5},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MultiSourceConsensusCompromiseFusion(tt.args.opinions)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiSourceConsensusCompromiseFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareWithin(got, tt.want, 1e-10) {
				t.Errorf("MultiSourceConsensusCompromiseFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMultiSourceConsensusCompromiseFusionJava compares MultiSourceConsensusCompromiseFusion with results of subjective-logic-java
// for random opinions.
func TestMultiSourceConsensusCompromiseFusionJava(t *testing.T) {
	file, err := os.Open("testdata/multi_source_consensus_compromise_fusion.txt")
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
		inputs, output, _ := strings.Cut(line, "|")
		in, out := parseFloats(t, inputs), parseFloats(t, output)
		opinions := make([]Opinion, int(in[0]))
		for i := range opinions {
			opinions[i] = Opinion{in[1+4*i], in[2+4*i], in[3+4*i], in[4+4*i]}
		}
		want := Opinion{out[0], out[1], out[2], opinions[0].baseRate}

		got, err := MultiSourceConsensusCompromiseFusion(opinions)
		if err != nil || !compareWithin(got, want, 1e-10) {
			t.Errorf("MultiSourceConsensusCompromiseFusion(%v) got = %v, %v, want %v", opinions, got, err, want)
		}
		cases++
	}
	if cases == 0 {
		t.Fatal("no test cases found")
	}
}

// TestMultiSourceConsensusCompromiseFusionProperties compares MultiSourceConsensusCompromiseFusion on random opinions
// with a direct evaluation of Eq. 7, with ConsensusCompromiseFusion for two opinions,
// and checks that the result does not depend on the order of the opinions.
func TestMultiSourceConsensusCompromiseFusionProperties(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 2000; i++ {
		opinions := randomOpinions(r, 2+r.Intn(6), 0.5)

		got, err := MultiSourceConsensusCompromiseFusion(opinions)
		if err != nil {
			t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) error = %v", opinions, err)
		}
		if want, _ := enumerateConsensusCompromiseFusion(opinions); !compareWithin(got, want, 1e-9) {
			t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) got = %v, want %v (Eq. 7)", opinions, got, want)
		}
		if len(opinions) == 2 {
			want, err := ConsensusCompromiseFusion(&opinions[0], &opinions[1])
			if err != nil || !compareWithin(got, want, 1e-9) {
				t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) got = %v, want %v (ConsensusCompromiseFusion)", opinions, got, want)
			}
		}
		shuffled := append([]Opinion(nil), opinions...)
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if gotShuffled, _ := MultiSourceConsensusCompromiseFusion(shuffled); !compareWithin(got, gotShuffled, 1e-9) {
			t.Fatalf("MultiSourceConsensusCompromiseFusion() depends on order: %v vs %v", got, gotShuffled)
		}
	}
}

// TestMultiSourceConsensusCompromiseFusionSpecialOpinions repeats the property test with dogmatic and vacuous opinions,
// which the Java results and the other tests do not contain, and checks that the results are valid opinions.
func TestMultiSourceConsensusCompromiseFusionSpecialOpinions(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for i := 0; i < 5000; i++ {
		opinions := randomSpecialOpinions(r, 2+r.Intn(5))
		for j := range opinions {
			opinions[j].baseRate = opinions[0].baseRate
		}

		got, err := MultiSourceConsensusCompromiseFusion(opinions)
		if err != nil {
			t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) error = %v", opinions, err)
		}
		if want, defined := enumerateConsensusCompromiseFusion(opinions); defined && !compareWithin(got, want, 1e-9) {
			t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) got = %v, want %v (Eq. 7)", opinions, got, want)
		}
		if len(opinions) == 2 {
			want, err := ConsensusCompromiseFusion(&opinions[0], &opinions[1])
			if err != nil || !compareWithin(got, want, 1e-9) {
				t.Fatalf("MultiSourceConsensusCompromiseFusion(%v) got = %v, want %v (ConsensusCompromiseFusion)", opinions, got, want)
			}
		}
		shuffled := append([]Opinion(nil), opinions...)
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if gotShuffled, _ := MultiSourceConsensusCompromiseFusion(shuffled); !compareWithin(got, gotShuffled, 1e-9) {
			t.Fatalf("MultiSourceConsensusCompromiseFusion() depends on order: %v vs %v", got, gotShuffled)
		}
	}
}

func BenchmarkMultiSourceConsensusCompromiseFusion(b *testing.B) {
	opinions := randomOpinions(rand.New(rand.NewSource(1)), 8, 0.5)
	b.ResetTimer()
	for range b.N {
		x, err := MultiSourceConsensusCompromiseFusion(opinions)
		if err != nil {
			b.Error(err)
		}
		sink = x
	}
}

func compareWithin(opinion1, opinion2 Opinion, tolerance float64) bool {
	return math.Abs(opinion1.belief-opinion2.belief) < tolerance &&
		math.Abs(opinion1.disbelief-opinion2.disbelief) < tolerance &&
		math.Abs(opinion1.uncertainty-opinion2.uncertainty) < tolerance &&
		math.Abs(opinion1.baseRate-opinion2.baseRate) < tolerance
}

func parseFloats(t *testing.T, s string) []float64 {
	var values []float64
	for _, field := range strings.Fields(s) {
		value, err := strconv.ParseFloat(field, 64)
		if err != nil {
			t.Fatal(err)
		}
		values = append(values, value)
	}
	return values
}

// randomSpecialOpinions returns random opinions with random base rates, including dogmatic and vacuous opinions and
// base rates of 0 and 1
func randomSpecialOpinions(r *rand.Rand, n int) []Opinion {
	pick := func() float64 {
		switch r.Intn(6) {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return float64(r.Intn(11)) / 10
		default:
			return r.Float64()
		}
	}
	opinions := make([]Opinion, n)
	for i := range opinions {
		u := pick()
		b := pick() * (1 - u)
		opinions[i] = Opinion{b, math.Max(0, 1-u-b), u, pick()}
	}
	return opinions
}

func randomOpinions(r *rand.Rand, n int, baseRate float64) []Opinion {
	opinions := make([]Opinion, n)
	for i := range opinions {
		u := r.Float64()
		b := r.Float64() * (1 - u)
		opinions[i] = Opinion{b, math.Max(0, 1-u-b), u, baseRate}
	}
	return opinions
}

// enumerateConsensusCompromiseFusion evaluates Eq. 7 of van der Heijden et al. (2018) for binomial opinions directly:
// the residual belief of every source times the uncertainties of all other sources, plus every combination in which
// each source contributes its residual belief or its residual disbelief.
// The second result is false if the compromise mass is zero, where Eq. 7 is undefined and the consensus is returned.
func enumerateConsensusCompromiseFusion(opinions []Opinion) (Opinion, bool) {
	bCons, dCons, uPre := 1.0, 1.0, 1.0
	for _, opinion := range opinions {
		bCons = math.Min(bCons, opinion.belief)
		dCons = math.Min(dCons, opinion.disbelief)
		uPre *= opinion.uncertainty
	}
	// mass indexed by set: 1 = {x}, 2 = {not x}, 3 = whole domain
	var mass [4]float64
	for i, opinion := range opinions {
		uOthers := 1.0
		for j, other := range opinions {
			if j != i {
				uOthers *= other.uncertainty
			}
		}
		mass[1] += (opinion.belief - bCons) * uOthers
		mass[2] += (opinion.disbelief - dCons) * uOthers
	}
	for c := 0; c < 1<<len(opinions); c++ {
		p, intersection, union := 1.0, 3, 0
		for i, opinion := range opinions {
			if c&(1<<i) == 0 {
				p *= opinion.belief - bCons
				intersection &= 1
				union |= 1
			} else {
				p *= opinion.disbelief - dCons
				intersection &= 2
				union |= 2
			}
		}
		if intersection != 0 {
			mass[intersection] += p
		} else {
			mass[union] += p
		}
	}
	b, d := bCons, dCons
	compromiseMass := mass[1] + mass[2] + mass[3]
	if compromiseMass > 0 {
		eta := (1 - bCons - dCons - uPre) / compromiseMass
		b += eta * mass[1]
		d += eta * mass[2]
	}
	return Opinion{b, d, 1 - b - d, opinions[0].baseRate}, compromiseMass > 0
}
