package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var gcc = "/s/gcc-3.4.4/bin/gcc"

var granularity = flag.String("gran", "word", "granularity of deltas (line, word, or char)")

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [input-file] [err-file]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

type GccTester struct {
	expectedErrs [][]byte
}

func NewGccTester(expectedErr io.Reader) (*GccTester, error) {
	data, err := ioutil.ReadAll(expectedErr)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(data, []byte("\n"))
	return &GccTester{expectedErrs: lines}, nil
}

func (t *GccTester) Test(input []byte) Outcome {
	var stderr bytes.Buffer
	cmd := exec.Command(gcc, "-c", "-O3", "-xc", "-")
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = &stderr
	_ = cmd.Run()

	errput := stderr.Bytes()

	lines := bytes.Split(errput, []byte("\n"))
	if len(lines) == 0 {
		return Passed
	} else if len(lines) != len(t.expectedErrs) {
		return Undetermined
	}

	for _, line := range t.expectedErrs {
		if !bytes.Contains(errput, line) {
			return Undetermined
		}
	}

	return Failed
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		return
	}

	infile, errfile := flag.Arg(0), flag.Arg(1)

	// create builder/deltas for test input
	f, err := os.Open(infile)
	if err != nil {
		log.Fatal("oops: ", err)
	}
	defer f.Close()

	var builder Builder
	switch *granularity {
	case "word":
		builder, err = ByWord(f)
	case "line":
		builder, err = ByLine(f)
	case "char":
		builder, err = ByChar(f)
	default:
		flag.Usage()
		return
	}

	if err != nil {
		log.Fatal("oops: ", err)
	}

	// load expected failure/error output
	ef, err := os.Open(errfile)
	if err != nil {
		log.Fatal("oops: ", err)
	}
	defer ef.Close()

	gcctest, err := NewGccTester(ef)
	if err != nil {
		log.Fatal("oops: ", err)
	}

	//fmt.Println(string(builder.BuildInput(IntRange(builder.Len()))))

	// run minimization test case
	tcase := &TestCase{B: builder, T: gcctest}

	run, err := MinFail(tcase)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(builder.BuildInput(run.Minimal)))
}
