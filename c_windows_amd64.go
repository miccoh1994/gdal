//go:build windows && amd64
// +build windows,amd64

package gdal

/*
#cgo windows LDFLAGS: -Lc:/gdal/3.9.1/lib -lgdal_i
#cgo windows CFLAGS: -Ic:/gdal/3.9.1/include
*/
import "C"
