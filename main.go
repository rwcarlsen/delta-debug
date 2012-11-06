
package main

import (
  "fmt"
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
  log.SetFlags(log.Lshortfile)
}

type Builder interface {
  BuildInput(godd.Set) []byte
  Len() int
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

  input := bytes.Join(inputWords, []byte(" "))
  return append(input, byte('\n'))
}

func (wi *WordBuilder) Len() int {
  return len(wi.words)
}

type CharBuilder struct {
  data []byte
}

func NewCharBuilder(r io.Reader) (*CharBuilder, error) {
  data, err := ioutil.ReadAll(r)
  if err != nil {
    return nil, err
  }
  return &CharBuilder{data: data}, nil
}

func (wi *CharBuilder) BuildInput(set godd.Set) []byte {
  sort.Ints([]int(set))

  input := make([]byte, len(set))
  for i, index := range set {
    input[i] = wi.data[index]
  }

  return input
}

func (wi *CharBuilder) Len() int {
  return len(wi.data)
}

type Tester interface {
  Test(input []byte) bool
}

type GccTester struct {
  expectedErr []byte
}

func NewGccTester(expectedErr io.Reader) (*GccTester, error) {
  stderr, err := ioutil.ReadAll(expectedErr)
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
  _ = cmd.Run()

  errput := stderr.Bytes()
  return !bytes.Contains(errput, t.expectedErr)
}

type TestCase struct {
  T Tester
  B Builder
}

func (t *TestCase) Passes(set godd.Set) bool {
  input := t.B.BuildInput(set)
  return t.T.Test(input)
}

func (t *TestCase) Len() int {
  return t.B.Len()
}

func main() {
  testFile("./gcc-tests/nested-1.c", "./nested-1.err")
}

func testFile(name, errname string) {
  f, err := os.Open(name)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  wb, err := NewCharBuilder(f)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  f, err = os.Open(errname)
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

  fmt.Print(string(wb.BuildInput(run.Minimal)))
}
