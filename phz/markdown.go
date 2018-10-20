package phz

import (
	"github.com/russross/blackfriday"
)

func ParseMarkdown(b []byte) []byte {
	return blackfriday.Run(b)
}
