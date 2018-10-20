package phz

import (
	"github.com/russross/blackfriday"
)

func ParseMarkdown(input []byte) []byte {
	return blackfriday.Run(input)
}
