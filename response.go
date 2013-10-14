package iat

import (
  "math"
  "fmt"
)

type Response struct {
  Time float64
  Dir Direction
  Correct bool
}

func (e *ReadyExperiment) CalculateScore(responses []Response) float64 {
  blocks := make(map[int][]Response)
  blockSums := map[int]float64{
    1: 0.0, 2: 0.0, 3: 0.0, 4: 0.0, 5: 0.0, 6: 0.0, 7: 0.0}

  correctBlocks := make(map[int][]Response)

  for i, r := range responses {
    blockNum := e.Frames[i].Block()
    block, ok := blocks[blockNum]
    if !ok {
      block = make([]Response, 0, 20)
    }
    blocks[blockNum] = append(block, r)
    blockSums[blockNum] += r.Time

    r.Correct = r.Dir == e.Frames[i].Correct()
    if r.Correct {
      correctBlock, ok := correctBlocks[i]
      if !ok {
        correctBlock = make([]Response, 0, 20)
      }
      correctBlocks[blockNum] = append(correctBlock, r)
    }
  }
  fmt.Println(blocks)

  correctBlockMeans := map[int]float64{
    3: calcMean(correctBlocks[3]),
    4: calcMean(correctBlocks[4]),
    6: calcMean(correctBlocks[6]),
    7: calcMean(correctBlocks[7]),
  }
  fmt.Println(correctBlockMeans)

  b3b6stddev := stdDev(append(blocks[3], blocks[6]...))
  b4b7stddev := stdDev(append(blocks[4], blocks[7]...))

  fmt.Println(b3b6stddev)
  fmt.Println(b4b7stddev)

  for i, r := range blocks[3] {
    if !r.Correct {
      blocks[3][i].Time = correctBlockMeans[3] + 600
    }
  }
  for i, r := range blocks[4] {
    if !r.Correct {
      blocks[4][i].Time = correctBlockMeans[4] + 600
    }
  }
  for i, r := range blocks[6] {
    if !r.Correct {
      blocks[6][i].Time = correctBlockMeans[6] + 600
    }
  }
  for i, r := range blocks[7] {
    if !r.Correct {
      blocks[7][i].Time = correctBlockMeans[7] + 600
    }
  }

  b3mean := calcMean(blocks[3])
  b4mean := calcMean(blocks[4])
  b6mean := calcMean(blocks[6])
  b7mean := calcMean(blocks[7])

  quotient1 := (b6mean - b3mean)/b3b6stddev
  quotient2 := (b7mean - b4mean)/b4b7stddev

  return (quotient1 + quotient2)/2
}

func stdDev(responses []Response) float64 {
  topSum := 0.0
  bottomSum := 0.0
  mean := calcMean(responses)
  for _, r := range responses {
    topSum += (r.Time - 1) * math.Pow(r.Time - mean, 2)
    bottomSum += r.Time
  }
  return math.Sqrt(topSum/(bottomSum - float64(len(responses))))
}

func calcMean(responses []Response) float64 {
  return sum(responses)/float64(len(responses))
}

func sum(responses []Response) float64 {
  sum := 0.0
  for _, r := range responses {
    sum += r.Time
  }
  return sum
}
