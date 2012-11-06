
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
  "github.com/rwcarlsen/godd/byteinp"
)

var gcc = "/s/gcc-3.4.4/bin/gcc"

func init() {
  log.SetPrefix("mylog:")
  log.SetFlags(log.Lshortfile)
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
