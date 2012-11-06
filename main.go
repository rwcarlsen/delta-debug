
package main

import (
  "fmt"
  "log"
  "bytes"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "github.com/rwcarlsen/godd"
  "github.com/rwcarlsen/godd/byteinp"
)

var gcc = "/s/gcc-3.4.4/bin/gcc"

func init() {
  log.SetPrefix("mylog:")
  log.SetFlags(log.Lshortfile)
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

func (t *GccTester) Test(input []byte) bool {
  var stderr bytes.Buffer
  cmd := exec.Command(gcc, "-c", "-O3", "-xc", "-")
  cmd.Stdin = bytes.NewReader(input)
  cmd.Stderr = &stderr
  _ = cmd.Run()

  errput := stderr.Bytes()
  for _, line := range t.expectedErrs {
    if !bytes.Contains(errput, line) {
      return true
    }
  }
  if len(bytes.Split(errput, []byte("\n"))) != len(t.expectedErrs) {
    return true
  }
  return false
}

func main() {
  testFile("./gcc-tests/nested-1.c", "./nested-1.err")
}

func testFile(name, errname string) {
  f, err := os.Open(name)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  wb, err := byteinp.ByWord(f)
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

  tcase := &byteinp.TestCase{B: wb, T: gcctest}

  run, err := godd.MinFail(tcase)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Print(string(wb.BuildInput(run.Minimal)))
}
