//Carla Groenland and Tom Johnston
//go run intersection-sizes-of-linear-subspaces-with-the-hypercube.go -h for options.

package main

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
)

//partition is a way of partitioning n elements into k non-negative numbers such that every number except the first is at most the value in capacities.
type partition struct {
	n          int
	currentSum int
	state      []int
	capacities []int
}

//MakePartition makes a partition of n into k non-negative numbers with the given capacities.
func MakePartition(n, k int, capacities []int) partition {
	state := make([]int, k)
	state[0] = n
	currentSum := 0
	return partition{n: n, currentSum: currentSum, state: state, capacities: capacities}
}

//Next gets the next partition if there is one and returns if it was successful.
func (p *partition) Next() bool {
	k := len(p.state)
	for i := 1; i < k; i++ {
		if p.currentSum == p.n || p.state[i] == p.capacities[i] {
			p.currentSum -= p.state[i]
			p.state[i] = 0
			continue
		}
		p.state[i]++
		p.currentSum++
		p.state[0] = p.n - p.currentSum
		return true
	}
	return false
}

//pattern holds the information for a shape.
type pattern struct {
	previousSize int   //The last size added. W.L.O.G. we can restrict to adding the constraints in non-increasing sizes.
	m            int   //The number of constraints.
	k            int   //The number of non-trivial rows.
	shape        []int //Holds the number of each row type in the shape matrix. Given a row type, convert it to a binary number i. The number of this row type is then in shape[i]. Note that shape[0] = 0 as all 0 rows are ignored.
	upperbound   int   //This is roughly the maximum intersection size of the shape. In order to avoid adding redundant conditions, the upper bound is the largest maximum intersection size which covers a lower proportion than its parent.
	secondBound  int   //This is the size of the largest intersecation that is smaller than the upper bound.
}

func main() {
	minPropPtr := flag.Float64("prop", 0.46875, "Only find shapes which can intersect with a proportion of 2^k strictly greater than prop.")
	maxKPtr := flag.Int("k", -1, "The largest value of k to check. -1 is unbounded.")
	maxMPtr := flag.Int("m", -1, "The largest value of m to check. -1 is unbounded.")
	minPtr := flag.Bool("minimality", false, "Whether the code should enfore minimality on the shapes or just attempt to stop redundant conditions.")

	flag.Parse()
	fmt.Printf("prop: %v k: %v m: %v minimality: %v \n", *minPropPtr, *maxKPtr, *maxMPtr, *minPtr)

	toCheck := make([]pattern, 1) //Store the patterns to check.
	//Compute the maximum number of non-zero entries in a column.
	maxSize := 1
	for true {
		if float64(BinomialCoeffSingle(maxSize+2, (maxSize+2)/2))/float64(int(1)<<uint(maxSize+1)) > *minPropPtr {
			maxSize++
		} else {
			break
		}
	}
	//Add an initial pattern
	toCheck[0] = pattern{previousSize: maxSize, m: 0, k: 0, shape: []int{0}, upperbound: 1}

	var f pattern
	previousM := 0
	for len(toCheck) > 0 {
		f, toCheck = toCheck[0], toCheck[1:] //Do a BFS.

		lenf := len(f.shape)
		newF := make([]int, 2*lenf) //This will hold the new shapes as we check them. Copies will need to be made to be stored.

		for i := 2; i <= f.previousSize; i++ { //The number of non-zero entries in the next column.

			p := MakePartition(i, lenf, f.shape) //Make the partition. This chooses where to distribute the non-zero entries in the next column.
		partitionLoop:
			for true {
				newF[lenf] = p.state[0]
				for j := 1; j < lenf; j++ {
					newF[j] = f.shape[j] - p.state[j]
					newF[j+lenf] = p.state[j]
				}

				if *minPtr {
					//Check for minimality. For minimality every column must have a row where it has the only non-zero entry. That means shape[2^i] must be non-zero for i=1,...,m.
					a := 1
					for true {
						if newF[a] == 0 {
							//Not minial so update the partition and check the next state.
							if !p.Next() {
								break partitionLoop
							}
							continue partitionLoop
						}
						a *= 2
						if a >= len(newF) {
							break
						}
					}
				}

				ub := f.upperbound
				for i := 0; i < p.state[0]; i++ {
					ub *= 2
				}

				if (f.m+1 <= *maxMPtr || *maxMPtr == -1) && (f.k+p.state[0] <= *maxKPtr || *maxKPtr == -1) { //Enforce the conditions.
					if c, d := MaxIntersection(pattern{previousSize: i, m: f.m + 1, k: f.k + p.state[0], shape: newF, upperbound: ub}); float64(c)/math.Pow(2, float64(f.k+p.state[0])) > *minPropPtr {
						tmp := make([]int, 2*lenf) //Make a copy of newF to store.
						copy(tmp, newF)
						tmpPatt := pattern{previousSize: i, m: f.m + 1, k: f.k + p.state[0], shape: tmp, upperbound: c, secondBound: d}
						toCheck = append(toCheck, tmpPatt)

						// Print the found result.
						if tmpPatt.m > previousM {
							fmt.Printf("\n******** M = %v ********\n\n", tmpPatt.m)
							previousM = tmpPatt.m
						}
						proportion := float64(c) / float64(int(1)<<uint(tmpPatt.k))
						proportion2 := float64(d) / float64(int(1)<<uint(tmpPatt.k))
						fmt.Printf("m: %2d k: %2d previousSize: %2d shape: %v upperbound: %3d prop: %f secondbound: %3d secondprop :%f \n", tmpPatt.m, tmpPatt.k, tmpPatt.previousSize, tmpPatt.shape, c, proportion, d, proportion2)
					}
				}
				//Update the partition.
				if !p.Next() {
					break partitionLoop
				}
			}
		}
	}
}

//MaxIntersection finds the largest intersection size for this shape which is strictly less than the upperbound. It also returns the second largest pattern it finds.
func MaxIntersection(patt pattern) (int, int) {
	m := patt.m
	k := patt.k
	f := patt.shape
	//This code attempts all the possible values for the non-zero elements in rows with at least 2 and all possible numbers of 1s in each columns remaining non-zero entries. It calculates the intersection size by testing every point in the hypercube over the shared rows and calculating the number of ways of using the other entries to satisfy each constraint.

	expandedF := make([]int, 0, k) //Expand out the rows of the shape which contain more than 1 non-zero element.
	optionsForSingleRows := 1      //Count the number of ways of choosing how many 1s in the other entries.
	for i, v := range f {
		if bits.OnesCount(uint(i)) != 1 {
			for j := 0; j < v; j++ {
				expandedF = append(expandedF, i)
			}
		} else {
			optionsForSingleRows *= v + 1
		}
	}

	k2 := len(expandedF)        //Number of rows with more than one non-zero entry.
	matrix := make([]int, m*k2) //Turn expandedF into a proper matrix and index the non-zero entries from 1. The actaul value of a non-zero entry i in matrix is given by options[state[matrix[i] -1]].
	numNonZero := 0             //Number of non-zero entries in the matrix.
	for i := 0; i < k2; i++ {
		for j := 0; j < m; j++ {
			if expandedF[i]&(1<<uint(j)) > 0 { //This checks if the appropiate bit is set.
				numNonZero++
				matrix[(i+1)*m-j-1] = numNonZero
			}
		}
	}
	state := make([]int, numNonZero) //The index in options of the given non-zero entry in matrix.
	state2 := make([]int, m)         //The number of -1s in each columns entries in single rows.
	options := []int{1, -1}          //The possible values.

	maximumPassingPoints := 0 //The maximum intesection size less than upper bound.
	secondMax := 0            //The second largest intersection size less than upper bound.

	for i := 0; i < Pow(len(options), numNonZero); i++ { //Iterate over all possible values of state.
		for j := 0; j < optionsForSingleRows; j++ { //Iterate over all possible values of state2.

			passingPoints := 0  //Count the number that pass with these states.
			passingPoints2 := 1 //Count the number that pass for each of the poitns in the shared hypercube.

			for j := 0; j < 1<<uint(k2); j++ { //Loop over the points in the shared hypercube.
				passingPoints2 = 1       //Reset this to one.
				for a := 0; a < m; a++ { //The column to check
					sum := 0                  //The value of the column over the shared hypercube.
					for b := 0; b < k2; b++ { //Sum over the rows of the matrix
						if j&(1<<uint(b)) > 0 && matrix[b*m+a] > 0 { //If the bit for this row is set in j (so the row is selected) and the matrix entry is non-zero.
							sum += options[state[matrix[b*m+a]-1]] //Add the value. Linear map.
						}
					}
					passingPoints2 *= BinomialCoeffSingle(f[1<<uint(m-a-1)]+1, state2[a]-sum+1) //Calulate the number of ways of choosing from the columns single rows to pass. Each column is seperate so multiply them together.
				}
				passingPoints += passingPoints2 //Add the number for this state of the shared hypercube.
			}

			//Update the max and second max.
			if passingPoints > maximumPassingPoints && passingPoints < patt.upperbound {
				if maximumPassingPoints > secondMax {
					//The old max is becoming the new second max.
					secondMax = maximumPassingPoints
				}
				maximumPassingPoints = passingPoints
			} else if passingPoints < maximumPassingPoints && passingPoints > secondMax {
				secondMax = passingPoints
			}

			//Update state2
			for c := m - 1; c >= 0; c-- {
				if state2[c] < f[1<<uint(m-i-1)] {
					state2[c]++
					break
				} else {
					state2[c] = 0
				}
			}
		}
		//Update state
		for j := 0; j < numNonZero; j++ {
			if state[j] == len(options)-1 {
				state[j] = 0
				continue
			}
			state[j] = state[j] + 1
			break
		}
	}
	return maximumPassingPoints, secondMax
}

//Sum returns the sum of the elements in a.
func Sum(a []int) int {
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

//Pow returns a**b.
func Pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

//BinomialCoeffSingle returns n choose k for 0 < k <= n and 0 outside this range.
func BinomialCoeffSingle(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	comb := 1
	for i := 1; i <= k; i++ {
		comb *= (n - k + i)
		comb /= i
	}
	return comb
}
