package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
)

func main() {
	p := flag.Int("p", 10000, "total production")
	k := flag.Int("k", 1000, "desired number of serial numbers")
	mean := flag.Float64("m", 1443.0, "desired mean serial number")
	stdDev := flag.Float64("s", 100.0, "desired standard deviation")
	distributionType := flag.String("t", "uniform", "s/n distribution, \"uniform\", \"normal\"")
	gnuplotOutput := flag.Bool("g", false, "gnuplot output on stdout")
	tsvOutput := flag.Bool("T", false, "TSV output on stdout")
	flag.Parse()

	var serialNumbers []int

	switch *distributionType {
	case "uniform":
		serialNumbers = uniformDistribution(*p, *k, *mean, *stdDev)
	case "normal":
		serialNumbers = normalDistribution(*p, *k, *mean, *stdDev)
	default:
		fmt.Fprintf(os.Stderr, "no distribution %q, use \"uniform\" or \"normal\"\n", *distributionType)
		return
	}

	sort.Sort(sort.IntSlice(serialNumbers))

	largestSerialNumber := serialNumbers[len(serialNumbers)-1]
	smallestSerialNumber := serialNumbers[0]
	serialNumberDifference := float64(largestSerialNumber - smallestSerialNumber)

	// find estimated total production
	totalProduction := serialNumberDifference * float64(*k+1) / float64(*k-1)

	// find continuous distribution line
	slope := float64(len(serialNumbers)-1) / serialNumberDifference
	yOffset := 0.0 - slope*float64(serialNumbers[0])

	// find largest difference from step to line
	step := 0.0
	maxDiff := -1.
	maxSN := 0 // serial number where maxDiff is taken
	maxStep := 0
	for i := range serialNumbers {
		diff := math.Abs(step - (slope*float64(serialNumbers[i]) + yOffset))
		if diff > maxDiff {
			maxDiff = diff
			maxSN = serialNumbers[i]
			maxStep = i
		}
		step += 1.0
	}

	if *tsvOutput {
		fmt.Printf("%d\t%d\t%.02f\t%.02f\t%s\t", *p, *k, *mean, *stdDev, *distributionType)
		fmt.Printf("%.0f\t%.05f\t%.05f\t", totalProduction, slope, yOffset)
		fmt.Printf("%d\t%.04f\t%d\t", maxSN, maxDiff, maxStep)
		leader := ""
		for _, sn := range serialNumbers {
			fmt.Printf("%s%d", leader, sn)
			leader = ","
		}
		fmt.Println("")
		return
	}

	if *gnuplotOutput {
		fmt.Printf("# Max difference %.03f at %d\n", maxDiff, maxSN)
		fmt.Printf("# Kolmogorov statistic %.05f\n", maxDiff/float64(*k))
		fmt.Printf("# estimated total production %.02f\n", totalProduction)
		fmt.Println("$MAXDIFF << ENDDIFF")
		fmt.Printf("%d\t%d\n", maxSN, maxStep)
		fmt.Printf("%d\t%.03f\n", maxSN, slope*float64(maxSN)+yOffset)
		fmt.Println("ENDDIFF")
		fmt.Println("$DATA << EOD")
		for i, sn := range serialNumbers {
			fmt.Printf("%d\t%d\n", sn, i)
		}
		fmt.Println("EOD")
		fmt.Printf("f(x) = %.06f * x + %.06f\n", slope, yOffset)
		fmt.Println("set key left")
		fmt.Println("set grid")
		fmt.Println("plot $DATA with steps title \"Cumulative distribution\", $MAXDIFF with line title \"Maximum difference\", f(x) with line title \"Continuos distribution\"")
		return
	}

	leader := ""
	for i := range serialNumbers {
		if i > 1 && (i%8) == 1 {
			fmt.Println("")
			leader = ""
		}
		fmt.Printf("%s%d", leader, serialNumbers[i])
		leader = ", "
	}
	fmt.Println("")
	fmt.Printf("y = %.04f * x + %.04f\n", slope, yOffset)
	fmt.Printf("max difference %.03f at %d\n", maxDiff, maxSN)
}

// uniformDistribution finds k serial numbers from a total production of p,
// uniformly distributed
func uniformDistribution(p int, k int, _, _ float64) []int {

	already := make(map[int]bool)

	outputCount := 0
	for outputCount < k {
		x := rand.Intn(int(p))
		if already[x] {
			continue
		}
		outputCount++
		already[x] = true
	}

	serialNumbers := make([]int, k)
	i := 0
	for x := range already {
		serialNumbers[i] = x
		i++
	}

	return serialNumbers
}

// normalDistribution finds k serial numbers from a total production of p,
// with specified mean and std deviation
func normalDistribution(p int, k int, mean, stdDev float64) []int {
	alreadySeen := make(map[int]bool)

	snCount := 0

	for snCount < k {
		f := rand.NormFloat64()*stdDev + mean

		if f < 0.0 || f > float64(p) {
			continue
		}

		serialNumber := int(f)

		if alreadySeen[serialNumber] {
			continue
		}

		alreadySeen[serialNumber] = true
		snCount++
	}

	serialNumbers := make([]int, k)
	i := 0
	for x := range alreadySeen {
		serialNumbers[i] = x
		i++
	}

	return serialNumbers
}
