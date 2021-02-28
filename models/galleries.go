package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImagesSplitN(n int) [][]Image {
	// Create out 2D slice
	ret := make([][]Image, n)
	// Create the inner slices - we need N of them, and we will
	// start them with a size of 0.
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	// Iterate over our images, using the index % n to determine
	// which of the slices in ret to add the image to.
	for i, img := range g.Images {
		// % is the remainder operator in Go
		// eg:
		//    0%3 = 0
		//    1%3 = 1
		//    2%3 = 2
		//    3%3 = 0
		//    4%3 = 1
		//    5%3 = 2
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}
