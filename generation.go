package genetic

import (
	"bytes"
)

func generateChromosome(nextChromosome, nextGene chan string, geneSet string, numberOfGenesPerChromosome int, quit chan bool) {
	defer func() { close(nextChromosome) }()

	for {
		c := bytes.NewBuffer(make([]byte, 0, numberOfGenesPerChromosome))
		for i := 0; i < numberOfGenesPerChromosome; i++ {
			select {
			case <-quit:
				quit <- true
				return
			default:
				gene := <-nextGene
				if len(gene) == 0 {
					return
				}
				c.WriteString(gene)
			}
		}
		select {
		case <-quit:
			quit <- true
			return
		default:
			nextChromosome <- c.String()
		}
	}
}

func generateGene(nextGene chan string, geneSet string, quit chan bool) {
	localRand := createRandomNumberGenerator()
	defer func() { close(nextGene) }()
	for {
		index := localRand.Intn(len(geneSet))
		select {
		case <-quit:
			quit <- true
			return
		default:
			nextGene <- geneSet[index : index+1]
		}
	}
}

func generateParent(nextChromosome chan string, geneSet string, numberOfChromosomes, numberOfGenesPerChromosome int) string {
	s := bytes.NewBuffer(make([]byte, 0, numberOfChromosomes*numberOfGenesPerChromosome))
	for i := 0; i < numberOfChromosomes; i++ {
		s.WriteString(<-nextChromosome)
	}
	return s.String()
}
