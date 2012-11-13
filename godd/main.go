package main

import (
  "regexp"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var gcc = "/s/gcc-3.4.4/bin/gcc"

var granularity = flag.String("gran", "word", "granularity of deltas (line, word, or char)")

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [input-file] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

type GccTester struct {
	expectedErrs [][]byte
  errput []byte
}

func NewGccTester(input []byte) *GccTester {
	t := &GccTester{}
  t.Test(input)
	t.expectedErrs = bytes.Split(t.errput, []byte("\n"))

  reg := regexp.MustCompile("^<stdin>:[0-9]*:*")
  for i, line := range t.expectedErrs {
    if match := reg.Find(line); match != nil {
      t.expectedErrs[i] = line[len(match):]
    }
  }

  return t
}

func (t *GccTester) Test(input []byte) Outcome {
	var stderr bytes.Buffer
	cmd := exec.Command(gcc, "-c", "-O3", "-xc", "-")
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = &stderr
	_ = cmd.Run()

	t.errput = stderr.Bytes()

	lines := bytes.Split(t.errput, []byte("\n"))
	if len(lines) == 0 {
		return Passed
	} else if len(lines) != len(t.expectedErrs) {
		return Undetermined
	}

	for _, line := range t.expectedErrs {
		if !bytes.Contains(t.errput, line) {
			return Undetermined
		}
	}

	return Failed
}

func main() {
	if flag.Parse(); len(flag.Args()) != 1 {
		flag.Usage()
		return
	}


	// create builder/deltas for test input
	infile := flag.Arg(0)
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

  gcctest := NewGccTester(builder.BuildInput(IntRange(builder.Len())))

	// run minimization test case
	tcase := &TestCase{B: builder, T: gcctest}

	run, err := MinFail(tcase)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(builder.BuildInput(run.Minimal)))
}
