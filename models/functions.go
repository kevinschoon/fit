package models

import "sort"

type Function func(*Series) *Series

// Sum adds together all of the values in the series
func Sum(input *Series) *Series {
	out := Copy(input)
	container := make([]Value, len(input.Keys))
	for _, key := range input.Keys {
		for _, value := range input.Values(key) {
			container[int(key)] += value
		}
	}
	out.Add(input.Start(), container[1:])
	return out
}

// Avg averages together all the values in the series
func Avg(input *Series) *Series {
	out := Copy(input)
	container := Sum(input).values[0]
	for _, key := range input.Keys {
		container[int(key)] = container[int(key)] / Value(input.Len())
	}
	out.Add(input.Start(), container[1:])
	return out
}

// Min returns the minimum of all the values in the series
func Min(input *Series) *Series {
	out := Copy(input)
	container := make([]Value, len(input.Keys))
	for _, key := range input.Keys {
		values := Sortable(input.Values(key))
		sort.Sort(values)
		container[int(key)] = values[0]
	}
	out.Add(input.Start(), container[1:])
	return out
}

// Max returns the maximum of all the values in the series
func Max(input *Series) *Series {
	out := Copy(input)
	container := make([]Value, len(input.Keys))
	for _, key := range input.Keys {
		values := Sortable(input.Values(key))
		sort.Sort(values)
		container[int(key)] = values[len(values)-1]
	}
	out.Add(input.Start(), container[1:])
	return out
}

func Apply(input []*Series, fn Function) (output []*Series) {
	for _, series := range input {
		output = append(output, fn(series))
	}
	return output
}
