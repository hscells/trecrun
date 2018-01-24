# trecrun

[![GoDoc](https://godoc.org/github.com/hscells/trecrun?status.svg)](https://godoc.org/github.com/hscells/trecrun)
[![Go Report Card](https://goreportcard.com/badge/github.com/hscells/trecrun)](https://goreportcard.com/report/github.com/hscells/trecrun)

```
go get github.com/hscells/trecrun
```

trecrun deals with the deserialization of output from trec_eval-style run files. This package is inspired by
the companion go library https://github.com/TimothyJones/trecresults.

## Usage

```go
rf, err := trecrun.RunsFromReader(f)
if err != nil {
    log.Fatal(err)
}

// p@5 for topic 1.
fmt.Println(rf.Runs[1].Measurement["P_5"])

// p@5 for all.
fmt.Println(rf.Measurement["P_5"])
```