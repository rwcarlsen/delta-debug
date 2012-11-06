
package main

import (
  "log"
  "bytes"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "sort"
  "github.com/rwcarlsen/godd"
)

var gcc = "/s/gcc-3.4.4/bin/gcc"

func init() {
  log.SetPrefix("mylog:")
  log.SetFlags(0)

  expectedFail, err := ioutil.ReadFile("expected-fail.txt")
  if err != nil {
    log.Fatal(err)
  }
}

type Builder interface {
  BuildInput(godd.Set) []byte
}

type WordBuilder struct {
  words [][]byte
}

func NewWordBuilder(r io.Reader) (*WordBuilder, error) {
  data, err := ioutil.ReadAll(r)
  if err != nil {
    return nil, err
  }
  return &WordBuilder{words: bytes.Fields(data)}, nil
}

func (wi *WordBuilder) BuildInput(set godd.Set) []byte {
  sort.Ints([]int(set))

  inputWords := make([][]byte, len(set))
  for i, index := range set {
    inputWords[i] = wi.words[index]
  }

  return bytes.Join(inputWords, []byte(" "))
}

func (wi *WordBuilder) Len() int {
  return len(wi.words)
}

type Tester interface {
  Test(input []byte) bool
}

type TestCase struct {
  T Tester
  B Builder
}

func (t *TestCase) Test(set godd.Set) bool {
  input := t.B.BuildInput(set)
  return T.Test(input)
}

func (t *TestCase) Len() int {
  return t.B.Len()
}

type GccTester struct {
  expectedErr []byte
}

func NewGccTester(expectedErr io.Reader) (*GccTester, error) {
  stderr, err := ioutio.ReadAll(expectedErr)
  if err != nil {
    return nil, err
  }
  return &GccTester{expectedErr: stderr}, nil
}

func (t *GccTester) Test(input []byte) bool {
  var stderr bytes.Buffer
  cmd := exec.Command(gcc, "-c", "-O3", "-xc", "-")
  cmd.Stdin = bytes.NewReader(input)
  cmd.Stderr = &stderr

  if err := cmd.Run(); err != nil {
    log.Println("execution err: ", err)
  }

  //log.Println("input file: \n", string(input), "\n")
  //errput := stderr.Bytes()
  //log.Println("errput:\n", string(errput))
  return !bytes.Equal(output, t.expectedErr)
}

func main() {
  testFile("./gcc-tests/nested-1.c", "./nested-1.err")
}

func testFile(name, errname string) {
  f, err := os.Open(name)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  wb, err := NewWordBuilder(f)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  f, err := os.Open(errname)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  gcctest, err := NewGccTester(f)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  tcase := &TestCase{B: wb, T: gcctest}

  run, err := godd.MinFail(tcase)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("minimal:\n", wb.BuildInput(run.Minimal))
}
