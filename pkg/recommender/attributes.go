package recommender

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"sort"
)

// AttributeValueSelector interface comprises attribute selection algorythm entrypoints
type AttributeValueSelector interface {
	// SelectAttributeValues selects a range of attributes from the given
	SelectAttributeValues(min float64, max float64) ([]float64, error)
}

// AttributeValues type representing a slice of attribute values
type AttributeValues []float64

// sort method for sorting this slice in ascending order
func (av AttributeValues) sort() {
	sort.Float64s(av)
}

// SelectAttributeValues selects values between the min and max values considering the focus strategy
// When the interval between min and max is "out of range" with respect to this slice the lowest or highest values are returned
func (av AttributeValues) SelectAttributeValues(min float64, max float64) ([]float64, error) {
	logrus.Debugf("selecting attributes from %f, min [%f], max [%f]", av, min, max)
	if len(av) == 0 {
		return nil, fmt.Errorf("empty attribute values")
	}
	var (
		// holds the selected values
		selectedValues []float64
		// vars representing "distances" to the max from the "left" and right
		lDist, rDist = math.MaxFloat64, math.MaxFloat64
		// indexes of the "closest" values to max in the values slice
		rIdx, lIdx = -1, -1
	)
	// sort the slice in increasing order
	av.sort()
	logrus.Debugf("sorted attributes: [%f]", av)

	for i, v := range av {
		if v < max {
			// distance to max from "left"
			if lDist > max-v {
				lDist = max - v
				lIdx = i
			}
		} else {
			// distance to max from "right"
			if rDist > v-max {
				rDist = v - max
				rIdx = i
			}
		}
		if min <= v && v <= max {
			logrus.Debugf("found value between min[%f]-max[%f]: [%f], index: [%d]", min, max, v, i)
			selectedValues = append(selectedValues, v)
		}
	}

	logrus.Debugf("lower-closest index: [%d], higher-closest index: [%d]", lIdx, rIdx)
	if len(selectedValues) == 0 {
		// there are no values between the two limits
		if rIdx == -1 {
			//there are no values higher than max, return the closest less value
			// this covers the case when the [min, max] interval is out of range of the value set
			// the left index is either 0 or len(av)-1 in the above case
			return []float64{av[lIdx]}, nil
		}
		// the right index is higher than -1 -> there are higher values than the max, return the closest to it
		return []float64{av[rIdx]}, nil
	}
	return selectedValues, nil
}
