
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

var expectedFail []byte

func init() {
  log.SetPrefix("mylog:")
  log.SetFlags(0)

  expectedFail, err := ioutil.ReadFile("expected-fail.txt")
  if err != nil {
    log.Fatal(err)
  }
}

type WordInput struct {
  words [][]byte
}

func NewWordInput(r io.Reader) (*WordInput, error) {
  data, err := ioutil.ReadAll(r)
  if err != nil {
    return nil, err
  }


  return &WordInput{words: bytes.Fields(data)}, nil
}

func (wi *WordInput) Test(set godd.Set) bool {
  sort.Ints([]int(set))

  inputWords := make([][]byte, len(set))
  for i, index := range set {
    inputWords[i] = wi.words[index]
  }

  input := bytes.Join(inputWords, []byte(" "))
  return gccPasses(input)
}

func (wi *WordInput) Len() int {
  return len(wi.words)
}

func gccPasses(input []byte) bool {
  var out, stderr bytes.Buffer
  cmd := exec.Command("gcc", "-c", "-O3", "-xc", "-")
  cmd.Stdin = bytes.NewReader(input)
  cmd.Stdout = &out
  cmd.Stderr = &stderr

  if err := cmd.Run(); err != nil {
    log.Println("execution err: ", err)
  }

  //log.Println("input file: \n", string(input), "\n")
  //output := out.Bytes()
  //errput := stderr.Bytes()
  //log.Println("output:\n", string(output))
  //log.Println("errput:\n", string(errput))
  return !bytes.Equal(output, expectedFail)
}

func main() {
  f, err := os.Open("./gcc-tests/nested-1.c")
  if err != nil {
    log.Fatal("oops: ", err)
  }

  wi, err := NewWordInput(f)
  if err != nil {
    log.Fatal("oops: ", err)
  }

  run, err := godd.MinFail(wi)
  if err != nil {
    log.Fatal(err)
  }

  return 
}
