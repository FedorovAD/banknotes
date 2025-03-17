package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateCash(t *testing.T) {
	tests := []struct {
		name              string
		inputAmount       int
		inputCassettes    []Cassette
		expectedCassettes []CassetteOutput
		expectedFound     bool
	}{
		{
			name:        "Simple Test",
			inputAmount: 1000,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 10, IsWorking: true},
				{Number: 2, Denomination: 200, Count: 10, IsWorking: true},
				{Number: 3, Denomination: 500, Count: 10, IsWorking: true},
			},
			expectedCassettes: []CassetteOutput{
				{Number: 3, Denomination: 500, Count: 2},
			},
			expectedFound: true,
		},
		{
			name:        "OneCassetteFoundAns",
			inputAmount: 1000,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 10, IsWorking: true},
			},
			expectedCassettes: []CassetteOutput{
				{Number: 1, Denomination: 100, Count: 10},
			},
			expectedFound: true,
		},
		{
			name:        "OneCassetteNotFoundAns",
			inputAmount: 1150,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 20, IsWorking: true},
			},
			expectedCassettes: nil,
			expectedFound:     false,
		},
		{
			name:        "GetMinCountDenominations",
			inputAmount: 900,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 10, IsWorking: true},
				{Number: 2, Denomination: 200, Count: 10, IsWorking: true},
				{Number: 3, Denomination: 500, Count: 10, IsWorking: true},
			},
			expectedCassettes: []CassetteOutput{
				{Number: 3, Denomination: 500, Count: 1},
				{Number: 2, Denomination: 200, Count: 2},
			},
			expectedFound: true,
		},
		{
			name:        "GetMinCountDenominationsNotFound",
			inputAmount: 910,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 10, IsWorking: true},
				{Number: 2, Denomination: 200, Count: 10, IsWorking: true},
				{Number: 3, Denomination: 500, Count: 10, IsWorking: true},
			},
			expectedCassettes: nil,
			expectedFound:     false,
		},
		{
			name:        "EqualDenominationsAndCounts",
			inputAmount: 300,
			inputCassettes: []Cassette{
				{Number: 1, Denomination: 100, Count: 1, IsWorking: true},
				{Number: 2, Denomination: 100, Count: 1, IsWorking: true},
				{Number: 3, Denomination: 100, Count: 1, IsWorking: true},
			},
			expectedCassettes: []CassetteOutput{
				{Number: 3, Denomination: 100, Count: 1},
				{Number: 2, Denomination: 100, Count: 1},
				{Number: 1, Denomination: 100, Count: 1},
			},
			expectedFound: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualCassettes, actualFound := calculateCash(test.inputAmount, test.inputCassettes)
			require.Equal(t, actualCassettes, test.expectedCassettes)
			require.Equal(t, actualFound, test.expectedFound)
		})
	}
}
