package main

import (
	"testing"
	"time"

	"github.com/agradl/golang-exercise/testutil"
)

func TestNoOperationsReturnOutOfLoop(t *testing.T) {
	state := makeState()
	state.isShutdown()
	state.getStats("foo")
	state.getHash(1)
	state.getPendingHashCt()
	state.logResponse("foo", 1)
	state.doHash("foo", 0)
	state.isShutdown()
}

func TestCanTellIfShuttingDown(t *testing.T) {
	state := makeState()
	isShutdown := state.isShutdown()
	testutil.AssertFalse(t, isShutdown, "Is shutting down")
}

func TestGetHashReturnsNotFOund(t *testing.T) {
	state := makeState()
	testutil.AssertStringEqual(t, "not found", state.getHash(0), "Non existant hash")
	testutil.AssertStringEqual(t, "not found", state.getHash(2), "Non existant hash")
}
func TestDoHashDoesHashAfterDelay(t *testing.T) {
	state := makeState()
	index := state.doHash("foo", 0)
	time.Sleep(time.Millisecond)
	testutil.AssertStringEqual(t, "9_u6bgY2-JDlb7vzKD5STG-jIErimDgtYkdB0NxmODJuKCxBvl5CVNiCB3LFUYosWowMf37aGVlKfrU5RT4e1w==", state.getHash(index), "Hash")
	testutil.AssertIntsEqual(t, 0, state.getPendingHashCt(), "Pending Count")
}
func TestDoHashQueuesHash(t *testing.T) {
	state := makeState()
	index := state.doHash("foo", 0)
	testutil.AssertStringEqual(t, "", state.getHash(index), "Hash")
	testutil.AssertIntsEqual(t, 1, state.getPendingHashCt(), "Pending Count")
}

func TestStatsReturnZerosWhenNoRequests(t *testing.T) {
	state := makeState()
	aStats := state.getStats("A")

	assertStats(t, aStats, float64(0), 0)
}
func TestStatsLogsByPattern(t *testing.T) {
	state := makeState()
	state.logResponse("A", 5)
	state.logResponse("B", 11)

	aStats := state.getStats("A")

	assertStats(t, aStats, float64(5), 1)
}

func TestStatsLogsAverageAndTotal(t *testing.T) {
	state := makeState()
	state.logResponse("A", 5)
	state.logResponse("A", 11)

	aStats := state.getStats("A")

	assertStats(t, aStats, float64(8), 2)
}

func assertStats(t *testing.T, stats *statsObj, average float64, total int) {
	testutil.AssertFloatsEqual(t, average, stats.Average, "Average")
	testutil.AssertIntsEqual(t, total, stats.Total, "Total")
}
