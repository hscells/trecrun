// Package trecrun deals with the deserialization of output from trec_eval-style run files. This package is inspired by
// the companion go library https://github.com/TimothyJones/trecresults.
package trecrun

import (
	"io"
	"bufio"
	"strings"
	"github.com/pkg/errors"
	"strconv"
)

// Run is the measurements computed for a topic.
type Run struct {
	Topic       int64
	Measurement map[string]float64
}

// Runs are the measurements completed for an experiment.
type Runs map[int64]Run

// Result is the summary of the runs.
type Result struct {
	RunId       string
	Measurement map[string]float64
}

// RunFile is the union of runs and a result.
type RunFile struct {
	Runs
	Result
}

// NewRun creates a new run with a topic.
func NewRun(topic int64) *Run {
	return &Run{
		Topic:       topic,
		Measurement: make(map[string]float64),
	}
}

// NewResult creates a new result.
func NewResult() *Result {
	return &Result{
		Measurement: make(map[string]float64),
	}
}

// Add adds a measurement and its associated value to the run.
func (r *Run) Add(measurement string, value float64) {
	r.Measurement[measurement] = value
}

// Add adds a measurement and its associated value to the results.
func (r *Result) Add(measurement string, value float64) {
	r.Measurement[measurement] = value
}

// readLine reads a single line from a trec_eval-style run file. The measurement is always a string, however for the
// case of the results, the topic and value are interfaces, as the topic may be `all`, and the value may be the runid.
//
// The topic can only be cast to a string or an int64 and the value can only be cast to a string or a float64.
func readLine(line string) (measurement string, topic interface{}, value interface{}, err error) {
	l := strings.Split(line, "\t")
	if len(l) != 3 {
		err = errors.New("run files must contain three columns")
		return
	}
	measurement = strings.TrimSpace(l[0])

	if l[1] == "all" {
		topic = l[1]
	} else {
		var i int64
		i, err = strconv.ParseInt(l[1], 10, 64)
		if err != nil {
			return
		}
		topic = i
	}

	if measurement == "runid" {
		value = l[2]
	} else {
		value, err = strconv.ParseFloat(l[2], 64)
		if err != nil {
			return
		}
	}
	return
}

// RunsFromReader creates a run file (the runs and result summary) from a reader.
func RunsFromReader(reader io.Reader) (rf RunFile, err error) {
	var (
		run       *Run
		runs            = make(map[int64]Run)
		result          = NewResult()
		prevTopic int64 = -1
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var (
			measurement string
			topicV      interface{}
			valueV      interface{}
		)
		measurement, topicV, valueV, err = readLine(scanner.Text())
		if err != nil {
			return
		}

		if prevTopic < 0 {
			prevTopic = topicV.(int64)
		}

		switch v := valueV.(type) {
		case string:
			result.RunId = v
			continue
		case float64:
			if run == nil {
				t := topicV.(int64)
				run = NewRun(t)
			}
			run.Add(measurement, v)
		}

		switch t := topicV.(type) {
		case string:
			result.Add(measurement, valueV.(float64))
			continue
		case int64:
			if prevTopic != t {
				runs[prevTopic] = *run
				run = NewRun(t)
				prevTopic = t
			}
		}
	}

	return RunFile{
		runs,
		*result,
	}, nil
}
