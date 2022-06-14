package score_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/score"
)

func TestCalculate(t *testing.T) {
	type scenario struct {
		title   string
		mode    string
		players []string
	}

	star := "Milton"
	one := "• zero"
	two := "elguapo"
	null := "unknown"

	const s4on4Div0To1 = "4on4 div0-1"
	const s4on4Div1 = "4on4 div1"
	const s4on4Div1Need1 = "4on4 div1 need 1"
	const s4on4Div1Need2 = "4on4 div1 need 2"
	const s4on4Div1To2 = "4on4 div1-2"
	const s4on4Div2 = "4on4 div2"
	const s4on4Mix = "4on4 mix"
	const s3on3Div1 = "3on3 div1"
	const s2on2Div1 = "2on2 div1"
	const s2on2Div1To2 = "2on2 div1-2"
	const s1on1Div1 = "1on1 div1"
	const s1on1Div1Need1 = "1on1 div1 need 1"
	const sFFA = "FFA"
	const sCoopDiv1 = "coop div1"
	const sRaceDiv1 = "race div1"
	const sUnknownMode = "unknown mode"
	const sEmptyServer = "empty server"

	scenarios := []scenario{
		{
			s4on4Div0To1,
			"4on4",
			[]string{star, star, one, one, one, one, one, one},
		},
		{
			s4on4Div1,
			"4on4",
			[]string{one, one, one, one, one, one, one, one},
		},
		{
			s4on4Div1Need1,
			"4on4",
			[]string{one, one, one, one, one, one, one},
		},
		{
			s4on4Div1Need2,
			"4on4",
			[]string{one, one, one, one, one, one},
		},
		{
			s4on4Div1To2,
			"4on4",
			[]string{one, one, one, one, two, two, two, two},
		},
		{
			s4on4Div2,
			"4on4",
			[]string{two, two, two, two, two, two, two, two},
		},
		{
			s3on3Div1,
			"3on3",
			[]string{one, one, one, one, one, one},
		},
		{
			s2on2Div1,
			"2on2",
			[]string{one, one, one, one},
		},
		{
			s1on1Div1,
			"1on1",
			[]string{one, one},
		},
		{
			s4on4Mix,
			"4on4",
			[]string{one, one, two, two, two, null, null, null},
		},
		{
			s2on2Div1To2,
			"2on2",
			[]string{one, one, two, two},
		},
		{
			sFFA,
			"ffa",
			[]string{one, two, two, null, null},
		},
		{
			sCoopDiv1,
			"coop",
			[]string{one},
		},
		{
			sRaceDiv1,
			"race",
			[]string{one},
		},
		{
			s1on1Div1Need1,
			"1on1",
			[]string{one},
		},
		{
			sUnknownMode,
			"unknown",
			[]string{one},
		},
		{
			sEmptyServer,
			"4on4",
			[]string{},
		},
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(scenarios), func(i, j int) { scenarios[i], scenarios[j] = scenarios[j], scenarios[i] })

	type result struct {
		title string
		score int
	}

	results := make([]result, 0)

	for _, s := range scenarios {
		results = append(results, result{
			title: s.title,
			score: score.Calculate(s.mode, s.players),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	sortedTitles := make([]string, 0)

	for _, r := range results {
		sortedTitles = append(sortedTitles, r.title)
	}

	expect := []string{
		s4on4Div0To1,
		s4on4Div1,
		s4on4Div1Need1,
		s4on4Div1To2,
		s3on3Div1,
		s2on2Div1,
		s4on4Mix,
		s4on4Div2,
		s1on1Div1,
		s2on2Div1To2,
		s4on4Div1Need2,
		sCoopDiv1,
		sFFA,
		sRaceDiv1,
		s1on1Div1Need1,
		sUnknownMode,
		sEmptyServer,
	}

	assert.Equal(t, expect, sortedTitles)
}

func BenchmarkCalculate(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("4on4 missing players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			score.Calculate("4on4", []string{"xantom", "bps", "foo", "bar", "baz"})
		}
	})

	b.Run("1on1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			score.Calculate("1on1", []string{"xantom", "bps"})
		}
	})

	b.Run("ffa", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			score.Calculate("ffa", []string{"xantom", "bps", "foo", "bar", "baz", "alpha", "beta", "gamma"})
		}
	})
}
