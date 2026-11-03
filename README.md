# go-subjectivelogic



This Go library implements a number of Subjective Logic operators. Subjective Logic is a type of probabilistic logic that uses subjective opinions instead of probabilities. To learn more about Subjective Logic please refer to "Subjective Logic, A Formalism for Reasoning under Uncertainty" book by Audun Jøsang. 

## Table of Contents

- [Usage](#usage)
	* [Opinion](#opinion)
	* [Addition](#addition)
	* [Complement](#complement)
	* [Binomial Multiplication](#binomial-multiplication)
	* [Binomial Comultiplication](#binomial-comultiplication)
	* [Belief Constraint Fusion](#belief-constraint-fusion)
	* [Cumulative Fusion](#cumulative-fusion)
	* [Averaging Fusion](#averaging-fusion)
	* [Weighted Fusion](#weighted-fusion)
	* [Multi-Source Fusion](#multi-source-fusion)
	* [Trust Discounting](#trust-discounting)
	* [Multi-Edge Trust Discounting](#trust-discounting-for-multi-edge-path)
	* [Opposite-Belief Trust Discounting](#opposite-belief-trust-discounting)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Usage

### Opinion

Subjective Logic is a type of probabilistic logic that uses subjective opinions instead of probabilities. A binomial opinion about the truth of the proposition $x$ is $\omega_x = (b_x, d_x, u_x, a_x)$, where:

- $b_x$: belief mass in support of the proposition $x$ being true.
- $d_x$: disbelief mass in support of the proposition $x$ being false.
- $u_x$: uncertainty mass representing lack of evidence.
- $a_x$: base rate i.e. the prior probability of $x$ being true without any evidence.

```go
type Opinion struct {
	belief      float64
	disbelief   float64
	uncertainty float64
	baseRate    float64
}
```

The following methods have been implemented for the opinion:
- [NewOpinion()](#new-opinion)
- [GetBelief()]()
- [GetDisbelief()]()
- [GetUncertainty()]()
- [GetBaseRate()]()
- [Modify()]()
- [ProjProb()]()
- [Compare()]()
- [ToString]()
- [ToStringE()]()


#### NewOpinion() method

The `NewOpinion()` method takes in four `float64` values and outputs an `Opinion` as well as an Error. In case a valid `Opinion` can be formed, it will be returned and the error will be `nil`. If the input values violate the requirements for a valid `Opinion`, the `Opinion` would be initialized with zeroes and an error will be returned.
For a valid `Opinion`, each input value $i$ must fulfill $0 \leq i \leq 1$ and the following statement $b + d + u = 1$ must hold.
For values computed with floating-point arithmetic, `NewOpinionWithTolerance()` corrects rounding errors: values that lie at most `ResultTolerance` ($10^{-9}$) outside of $[0, 1]$ are clamped, and $b$, $d$, $u$ are renormalised to sum to $1$.
Larger deviations still result in an error.
All operators create their results with `NewOpinionWithTolerance()`.

```go
opinion, err := subjectivelogic.NewOpinion(.5, .25, .25, .5)

if err != nil {
	println(err.Error())
} else {
	belief := opinion.GetBelief()
	disbelief := opinion.GetDisbelief()
	uncertainty := opinion.GetUncertainty()
	baseRate := opinion.GetBaseRate()
	
	println(fmt.Sprintf("Opinion: %.2f, %.2f, %.2f, %.2f", belief, disbelief, uncertainty, baseRate))
}
```

This code generates the following output:

```go
Opinion: 0.50, 0.25, 0.25, 0.50
```
---

### Addition
This implements the Addition Operator as defined in Subjective Logic:

```math
\omega_{(x\cup y)} = 
\begin{cases}
b_{x\cup y} = b_x + b_y \\
d_{x\cup y} = \frac{a_x(d_x - b_y) + a_y(d_y - b_x)}{a_x+a_y} \\
u_{x\cup y} = \frac{a_xu_x + a_yu_y}{a_x+a_y} \\
a_{x\cup y} = 0
\end{cases}
```

#### API Reference

```go
func Addition(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
Inputs $\omega_{x} = (b_x, d_x, u_x, a_x)$ and $\omega_{y} = (b_y, d_y, u_y, a_y)$ are problematic and will result in an error, if:

```math
\begin{split}
b_x + b_y &> 1 \text{, or} \\
a_x + a_y &> 1 
\end{split}
```

Moreover, the Addition Operator will attempt to divide through $0$, if $a_x = a_y = 0$, hence this is not allowed and an error will be returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5) 
	opinion3, _ := subjectivelogic.NewOpinion(1, 0, 0, 0.5) 
	
	out1, err1 := subjectivelogic.Addition(&opinion1, &opinion2) //Case 1 
	out2, err2 := subjectivelogic.Addition(&opinion1, &opinion3) //Case 2 
	
	fmt.Println("Case 1:", "Opinion = ", out1.ToString(), "Error:", err1) 
	fmt.Println("Case 2:", "Opinion = ", out2.ToString(), "Error:", err2) 

}
```

The above code snippet shows the usage of the Addition operator. Case 1 uses two Opinions that are not problematic for the Addition operator, resulting in a valid output Opinion and no error:

```go
Case 1: Opinion =  0.6000000000000001, 0.1, 0.3, 1 Error: <nil>
```

 Case 2 uses two valid Opinions that are problematic for the Addition operator as the sum of their belief masses exceeds $1$. This results in the Addition operator returning a zeroed Opinion and an error. 

```go
Case 2: Opinion =  0, 0, 0, 0 Error: Addition: Check the validity of your input values
```
---

### Complement
This implements the Complement Operator as defined in Subjective Logic:
```math
	\omega_{\overline{x}}  :
	\begin{cases}
		b_{\overline{x}} = d_x \\
		d_{\overline{x}} = b_x \\
		u_{\overline{x}} = u_x \\
		a_{\overline{x}} = 1 - a_x
	\end{cases}
```

#### API Reference

```go
func Complement(opinion *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	
	out, err := subjectivelogic.Complement(&opinion) 

	if err != nil { 
		fmt.Println("Error:", err) 
	} else { 
		fmt.Println("Output:", out, err) 
	}
}
```
The code snippet above shows the usage of the Complement operator. This specific example will result in the following output:

```go
Output: {0.8 0.2 0 0.5} <nil>
```
---

### Binomial Multiplication
This implements the Binomial Multiplication Operator as defined in Subjective Logic:
```math
	\omega_{x \wedge y}  :
	\begin{cases}
		b_{x \wedge y} = b_x b_y + \frac{(1-a_x) a_y b_x u_y + a_x (1-a_y) b_y u_x}{1 - a_x a_y} \\
		d_{x \wedge y} = d_x + d_y - d_x d_y \\
		u_{x \wedge y} = u_x u_y + \frac{(1-a_y) b_x u_y + (1-a_x) u_x b_y}{1 - a_x a_y} \\
		a_{x \wedge y} = a_x a_y
	\end{cases}       
```

#### API Reference

```go
func Multiplication(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Binomial Multiplication Operator will try to divide through $0$, if $a_x = a_y = 1$, hence this is not allowed and an error is returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5) 

	out, err := subjectivelogic.Multiplication(&opinion1, &opinion2)

	if err != nil { 
		fmt.PrintlnImplementation("Error:", err)
	} else {
		fmt.Println("Output:", out, err) 
	}	
}
``` 

The code snippet above shows the usage of the Multiplication operator. This specific example will result in the following output:

```go
Output: {0.12000000000000002 0.8 0.08 0.25} <nil>
```
---

### Binomial Comultiplication
This implements the Comultiplication Operator as defined in Subjective Logic:

```math
	\omega_{x \wedge y}  :
	\begin{cases}
		b_{x \wedge y} = b_x + b_y - b_x b_y \\
		d_{x \wedge y} = d_x d_y + \frac{a_x (1-a_y) d_x u_y + (1-a_x) a_y u_x d_y}{a_x + a_y - a_x a_y} \\
		u_{x \wedge y} = u_x u_y + \frac{a_y d_x u_y + a_x u_x d_y}{a_x + a_y - a_x a_y} \\
		a_{x \wedge y} = a_x + a_y - a_x a_y
	\end{cases}       
```

#### API Reference

```go
func Comultiplication(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Binomial Comultiplication Operator will try to divide through $0$, if $a_x = a_y = 0$, hence this is not allowed and an error will be returned.

#### Example 

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.Comultiplication(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Comultiplication operator. This specific example will result in the following output:

```go
Output: {0.52 0.16 0.32 0.75} <nil>
```
---

### Belief Constraint Fusion
==This implements the Belief Constraint Fusion Operator as defined in Subjective Logic:==


NEEDS TO BE EDITED
```math
	\omega_{X}^{(A\&B)}  :
	\begin{cases}
		b_{X}^{(A\&B)} = \frac{b_{X}^{A}u_{X}^{B} + b_{X}^{B}u_{X}^{A} + b_{X}^{A}b_{X}^{B}}{1 - Con} \\
		d_{X}^{(A\&B)} = \frac{d_{X}^{A}u_{X}^{B} + d_{X}^{B}u_{X}^{A} + d_{X}^{A}d_{X}^{B}}{1 - Con} \\
		u_{X}^{(A\&B)} = \frac{u_{X}^{A}u_{X}^{B}}{1-Con} \\
		a_{X}^{(A\&B)} = \frac{a_{X}^{A}(1-u_{X}^{A})+a_{X}^{B}(1-u_{X}^{B})}{2-u_{X}^{A}-u_{X}^{B}} & \text{for} \hspace{2mm} u_{X}^{A} + u_{X}^{B} < 2 \\
		a_{X}^{(A\&B)} = \frac{a_{X}^{A}+a_{X}^{B}}{2} & \text{for} \hspace{2mm}u_{X}^{A} = u_{X}^{B} = 1
	\end{cases}       
```

```math
	Con = b_{X}^{A}d_{X}^{B} + d_{X}^{A}b_{X}^{B}
```

#### API Reference

```go
func ConstraintFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Belief Constraint Fusion Operator will try to divide through $0$, if the conflict variable $Con = 1$, and an error will be returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.ConstraintFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Belief Constraint operator. This specific example will result in the following output:

```go
Output: {0.2941176470588236 0.7058823529411764 0 0.5} <nil>
```
---

### Cumulative Fusion
This implements the Aleatory Cumulative Fusion Operator as defined in Subjective Logic:

Case I: For $u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0$:
```math
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \frac{b_{X}^{A}(x)u_{X}^{B} + b_{X}^{B}(x)u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - u_{X}^{A} u_{X}^{B}} \\
		u_{X}^{(A\diamond B)} = \frac{u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B} - u_{X}^{A} u_{X}^{B}} \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)u_{X}^{B} + a_{X}^{B}(x)u_{X}^{A} - (a_{X}^{A}(x) +  a_{X}^{B}(x))u_{X}^{B}u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - 2 u_{X}^{A} u_{X}^{B}} & \text{if} \hspace{2mm} u_{X}^{A} \neq 1 \vee u_{X}^{B} \neq 1 \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2} & \text{if} \hspace{2mm}u_{X}^{A} = u_{X}^{B} = 1
	\end{cases}       
```


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
```math
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\diamond B)} = 0 \\
		a_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
```

```math
	 \text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
```

#### API Reference

```go
func CumulativeFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.CumulativeFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Cumulative fusion operator. This specific example will result in the following output:

```go
Output: {0.2 0.8 0 0.5} <nil>
```
---

### Averaging Fusion
This implements the Averaging Fusion Operator as defined in Subjective Logic:

Case I: For $u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0$:
```math
	\omega_{X}^{(A\underline{\diamond} B)}  :
	\begin{cases}
		b_{X}^{(A\underline{\diamond} B)}(x) = \frac{b_{X}^{A}(x)u_{X}^{B} + b_{X}^{B}(x)u_{X}^{A}}{u_{X}^{A} + u_{X}^{B}} \\
		u_{X}^{(A\underline{\diamond} B)} = \frac{2 u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B}} \\
		a_{X}^{(A\underline{\diamond} B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2}
	\end{cases}       
```


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
```math
	\omega_{X}^{(A\underline{\diamond} B)}  :
	\begin{cases}
		b_{X}^{(A\underline{\diamond} B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\underline{\diamond} B)} = 0 \\
		a_{X}^{(A\underline{\diamond} B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
```

```math
	\text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
```

#### API Reference

```go
func AveragingFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.AveragingFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Averaging fusion operator. This specific example will result in the following output:

```go
Output: {0.2 0.8 0 0.5} <nil>
```

---

### Weighted Fusion
This implements the Weighted Fusion Operator as defined in Subjective Logic:

Case I: For $(u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0) \wedge (u_{X}^{A} \neq 1 \vee u_{X}^{B} \neq 1)$:
```math
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \frac{b_{X}^{A}(x)(1-u_{X}^{A})u_{X}^{B} + b_{X}^{B}(x)(1-u_{X}^{B})u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - 2u_{X}^{A} u_{X}^{B}} \\
		u_{X}^{(A\diamond B)} = \frac{(2 - u_{X}^{A} - u_{X}^{B})u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B} - 2u_{X}^{A} u_{X}^{B}} \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)(1-u_{X}^{A}) + a_{X}^{B}(x)(1-u_{X}^{B})}{2 - u_{X}^{A} - u_{X}^{B}}
	\end{cases}       
```


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
```math
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\diamond B)} = 0 \\
		a_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
```

```math
	\text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
```


Case III: $u_{X}^{A} = u_{X}^{B} = 1$:
```math
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = 0 \\
		u_{X}^{(A\diamond B)} = 1 \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2}
	\end{cases}       
```

#### API Reference
```go
func WeightedFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.


#### Example
```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.WeightedFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Weighted fusion operator. This specific example will result in the following output:

```go
Output: {0.2 0.8 0 0.5} <nil>
```

---

### Multi-Source Fusion
Averaging, weighted, consensus & compromise and epistemic cumulative fusion are not associative, and cumulative fusion is not associative for dogmatic opinions.
Fusing three or more opinions by chaining the binary operators therefore gives results that depend on the order of the opinions.
The multi-source operators fuse all opinions in a single operation.

For $N$ sources $C \in \mathbb{C}$, the formulas below use the product of the uncertainties of all other sources:
```math
	P^{C} = \prod_{C_j \neq C} u_{X}^{C_j}
```

All multi-source operators take a slice of opinions.
A single opinion is returned unchanged (by epistemic cumulative fusion, it is uncertainty-maximised), and an empty slice or a null opinion results in an error.
If any opinion is dogmatic, cumulative, averaging and weighted fusion return the average of the dogmatic opinions, including their base rates, with equal weights ($\gamma_{X}^{C} = 1/k$ for $k$ dogmatic opinions), as in Case II of Jøsang, "Categories of Belief Fusion", Journal of Advances in Information Fusion 13(2), 2018 (Eq. 24, 39 and 56).
For two opinions, the results equal those of the binary operators, except for the base rate of averaging and weighted fusion when exactly one opinion is dogmatic: the binary operators average the base rates of both opinions.

All examples fuse the three opinions of Table I in Jøsang, Wang and Zhang (2017) and van der Heijden, Kopp and Kargl (2018):
```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.1, 0.3, 0.6, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0.2, 0.4, 0.5)
	opinion3, _ := subjectivelogic.NewOpinion(0.7, 0.1, 0.2, 0.5)

	out, err := subjectivelogic.MultiSourceCumulativeFusion([]subjectivelogic.Opinion{opinion1, opinion2, opinion3})

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

#### Multi-Source Cumulative Fusion
This implements cumulative fusion for multiple sources as defined by Jøsang, Wang and Zhang, "Multi-Source Fusion in Subjective Logic", FUSION 2017 (Eq. 16).
It corresponds to summing the evidence of all sources, and the fused base rate is the evidence-weighted average as given by Jøsang, "Categories of Belief Fusion", Journal of Advances in Information Fusion 13(2), 2018 (Eq. 23):
```math
	\omega_{X}^{\diamond(\mathbb{C})}  :
	\begin{cases}
		b_{X}^{\diamond(\mathbb{C})}(x) = \frac{\sum_{C} b_{X}^{C}(x) P^{C}}{\sum_{C} P^{C} - (N-1) \prod_{C} u_{X}^{C}} \\
		u_{X}^{\diamond(\mathbb{C})} = \frac{\prod_{C} u_{X}^{C}}{\sum_{C} P^{C} - (N-1) \prod_{C} u_{X}^{C}} \\
		a_{X}^{\diamond(\mathbb{C})}(x) = \frac{\sum_{C} a_{X}^{C}(x) (1 - u_{X}^{C}) P^{C}}{\sum_{C} (1 - u_{X}^{C}) P^{C}}
	\end{cases}
```
If all opinions are vacuous, the base rate is the average base rate.

```go
func MultiSourceCumulativeFusion(opinions []Opinion) (Opinion, error)
```
The example above results in `Output: {0.6511627906976744 0.2093023255813954 0.13953488372093026 0.5} <nil>`.

#### Multi-Source Averaging Fusion
This implements averaging fusion for multiple sources as defined by Jøsang, Wang and Zhang (2017, Eq. 32).
It corresponds to averaging the evidence of all sources, and the fused base rate is the average base rate:
```math
	\omega_{X}^{\underline{\diamond}(\mathbb{C})}  :
	\begin{cases}
		b_{X}^{\underline{\diamond}(\mathbb{C})}(x) = \frac{\sum_{C} b_{X}^{C}(x) P^{C}}{\sum_{C} P^{C}} \\
		u_{X}^{\underline{\diamond}(\mathbb{C})} = \frac{N \prod_{C} u_{X}^{C}}{\sum_{C} P^{C}} \\
		a_{X}^{\underline{\diamond}(\mathbb{C})}(x) = \frac{1}{N} \sum_{C} a_{X}^{C}(x)
	\end{cases}
```

```go
func MultiSourceAveragingFusion(opinions []Opinion) (Opinion, error)
```
The example above results in `Output: {0.509090909090909 0.16363636363636366 0.3272727272727273 0.5} <nil>`.

#### Multi-Source Weighted Fusion
This implements weighted fusion for multiple sources as defined by van der Heijden, Kopp and Kargl, "Multi-Source Fusion Operations in Subjective Logic", FUSION 2018.
It corresponds to the confidence-weighted average of the evidence of all sources, where the confidence of a source is $1 - u_{X}^{C}$:
```math
	\omega_{X}^{\widehat{\diamond}(\mathbb{C})}  :
	\begin{cases}
		b_{X}^{\widehat{\diamond}(\mathbb{C})}(x) = \frac{\sum_{C} b_{X}^{C}(x) (1 - u_{X}^{C}) P^{C}}{\sum_{C} (1 - u_{X}^{C}) P^{C}} \\
		u_{X}^{\widehat{\diamond}(\mathbb{C})} = \frac{(N - \sum_{C} u_{X}^{C}) \prod_{C} u_{X}^{C}}{\sum_{C} (1 - u_{X}^{C}) P^{C}} \\
		a_{X}^{\widehat{\diamond}(\mathbb{C})}(x) = \frac{\sum_{C} a_{X}^{C}(x) (1 - u_{X}^{C})}{N - \sum_{C} u_{X}^{C}}
	\end{cases}
```
If all opinions are vacuous, the result is the vacuous opinion with the average base rate.

```go
func MultiSourceWeightedFusion(opinions []Opinion) (Opinion, error)
```
The example above results in `Output: {0.5621621621621622 0.14594594594594595 0.2918918918918919 0.5} <nil>`.

#### Multi-Source Epistemic Cumulative Fusion
This fuses the opinions with multi-source cumulative fusion and then maximises the uncertainty of the result once, i.e. it converts as much belief or disbelief as possible into uncertainty while preserving the projected probability.
This is CBF-UM as described by Jøsang, "Categories of Belief Fusion", Journal of Advances in Information Fusion 13(2), 2018 (Sections 6.2 and 6.5), and as argued for by van der Heijden, Kopp and Kargl (2018, Section III-A).
The binary `CumulativeEpistemicFusion` follows Definition 12.6 of the Subjective Logic book; chaining it maximises the uncertainty after every pairwise fusion instead.

```go
func MultiSourceCumulativeEpistemicFusion(opinions []Opinion) (Opinion, error)
```
As for the binary operator, a fused base rate of 0 or 1 results in an error.
The example above results in `Output: {0.44186046511627897 0 0.558139534883721 0.5} <nil>`.

#### Multi-Source Consensus & Compromise Fusion
This implements consensus & compromise fusion for multiple sources as defined by van der Heijden, Kopp and Kargl (2018, Definition 5, Eq. 7).
The consensus is the belief and disbelief shared by all sources, and the compromise distributes the remaining belief and disbelief; conflicting combinations become uncertainty.
Unlike in the binary operator, vacuous opinions are not neutral: they reduce the consensus to zero.

The implementation follows `ccCollectionFuse` of [subjective-logic-java](https://github.com/vs-uulm/subjective-logic-java) and enumerates all $4^N$ combinations of values of the sources, so its runtime grows exponentially (about 30 ms for 10 opinions, 0.5 s for 12 opinions and 40 s for 15 opinions); at most 15 opinions are accepted.
Unlike `ccCollectionFuse`, it supports dogmatic opinions.
If the compromise mass is zero, Eq. 7 is undefined.
For identical opinions, the consensus is returned, as `ConsensusCompromiseFusion` does.
With two or more dogmatic opinions, the result is the limit for uncertainties approaching zero at the same rate, which distributes the remaining mass according to the residual belief and disbelief of the dogmatic opinions with equal weights.
This limit is unique if the residual mass of the dogmatic opinions is all belief or all disbelief; if the dogmatic opinions contradict each other, equal weights are a convention, as for the other operators.

```go
func MultiSourceConsensusCompromiseFusion(opinions []Opinion) (Opinion, error)
```
All opinions must have the same base rate.
The example above results in `Output: {0.6287499999999999 0.18225000000000002 0.189 0.5} <nil>`.

---

### Trust Discounting
This implements the Trust Discounting Operator as defined in Subjective Logic:

```math
	\omega_{X}^{[A;B]}  :
	\begin{cases}
		\boldsymbol{b}_{X}^{[A;B]}(x) & = \boldsymbol{P}_B^{A}*\boldsymbol{b}_{X}^B(x) \\
		u_{X}^{[A;B]} & = 1 - \boldsymbol{P}_B^{A}*\sum_{x \in \mathscr{R}(\mathbb{X})}\boldsymbol{b}_{X}^B(x) \\
		\boldsymbol{a}_{X}^{[A;B]}(x) & = \boldsymbol{a}_{X}^B(x) \\
	\end{cases}       
```

#### API Reference
```go
func TrustDiscounting(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```
#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example
```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.TrustDiscounting(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Trust discounting operator. This specific example will result in the following output:
```go
Output: {0.08000000000000002 0 0.9199999999999999 0.5} <nil>
``` 

### Trust Discounting for Multi-edge Path
This implements the Trust Discounting Operator for Multi-edge Paths as defined in Subjective Logic.

If $[A_1, ..., A_n]$ denotes the referral trust path and $[A_n, X]$ the functional trust, the Trust Discounting Operator for multi-edge paths is defined as:

```math
	\omega_{X}^{A_1}  :
	\begin{cases}
		b_{X}^{A_1}(x) & = P_{A_n}^{A_1}*b_{X}^{A_n}(x) \\
		u_{X}^{A_1} & = 1 - P_{A_n}^{A_1}*\sum_{x \in X}b_{X}^{A_n}(x) \\
		a_{X}^{A_1}(x) & = a_{X}^{A_n}(x)
	\end{cases}       
```
and:
```math
	P_{A_n}^{A_1} =  \prod_{i=1}^{n-1} P_{A_{i+1}}^{A_i}     
```

#### API Reference
TBD

#### Problematic Inputs
TBD

#### Example
TBD 


### Opposite-Belief Trust Discounting
This implements the Opposite-Belief Trust Discounting Operator as defined in Subjective Logic:

```math
	\omega_{X}^{[A;B]}  :
	\begin{cases}
		b_X^{[A;B]} & = b^A_Bb^B_X + d^A_Bd^B_X  \\
		d_X^{[A;B]} & = b^A_Bd^B_X + d^A_Bb^B_X \\
		u_X^{[A;B]} & = u^A_B  + (b^A_B + d^A_B)u^B_X \\
		a_X^{[A;B]} & = a^B_X \\
	\end{cases}       
```

#### API Reference
```go
func TrustDiscountingOppositeBelief(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```
#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example
```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(1, 0, 0, 0.5)

	out, err := subjectivelogic.TrustDiscountingOppositeBelief(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Trust discounting operator. This specific example will result in the following output:
```go
Output: {0.2 0.8 0 0.5} <nil>
``` 


## Contributing
Contributions are very welcome! Please let us know if you find an issue and have ideas for improvement. Alternately, open an issue or submit a pull request on GitHub. 

## License
This project is licensed under the Apache License, Verion 2.0 - see the LICENSE file for details.

## Contact
In case of problems or issues, please open an issue on GitHub.   
