package utils

import (
	"math/rand"
)

// Perceptron struct
type Perceptron struct {
	Weights      []float64
	bias         float64
	learningRate float64
}

type PerceptronItem struct {
	Flags      []float64
	IsSelected bool
}

type PerceptronItemList []*PerceptronItem

// NewPerceptron initializes a new Perceptron with a given number of inputs
func NewPerceptron(inputSize int, learningRate float64) *Perceptron {
	weights := make([]float64, inputSize)
	for i := range weights {
		weights[i] = rand.Float64()*2 - 1 // Initialize weights with random values between -1 and 1
	}
	return &Perceptron{
		Weights:      weights,
		bias:         rand.Float64()*2 - 1,
		learningRate: learningRate,
	}
}

// Train function to train the Perceptron with a dataset
func (p *Perceptron) Train(inputs PerceptronItemList, epochs int) {
	for epoch := 0; epoch < epochs; epoch++ {
		for i := range inputs {
			prediction := p.Predict(inputs[i].Flags)
			e := float64(IfThen(inputs[i].IsSelected, 1, 0)) - prediction
			// Update weights and bias
			for j := range p.Weights {
				p.Weights[j] += p.learningRate * e * inputs[i].Flags[j]
			}
			p.bias += p.learningRate * e
		}
	}
}

// Predict function to classify a given input vector
func (p *Perceptron) Predict(input []float64) float64 {
	sum := p.bias
	for i := range input {
		sum += p.Weights[i] * input[i]
	}
	if sum >= 0 {
		return 1
	}
	return 0
}

func (p *Perceptron) Save(file string) error {
	return SaveToJSON(file, p.Weights, true)
}

func (p *Perceptron) Load(file string) error {
	return LoadFromJSON(file, &p.Weights)
}
