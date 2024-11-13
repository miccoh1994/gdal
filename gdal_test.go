// Copyright 2011 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdal

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiffDriver(t *testing.T) {
	_, err := GetDriverByName("GTiff")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestMissingMetadata(t *testing.T) {
	ds, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}

	metadata := ds.Metadata("something-that-wont-exist")
	if len(metadata) != 0 {
		t.Errorf("got %d items, want 0", len(metadata))
	}
}

func TestRegenerateOverviews(t *testing.T) {
	memDrv, err := GetDriverByName("MEM")
	if err != nil {
		t.Errorf(err.Error())
	}
	// 4x4 pixel size
	ds, err := Open("testdata/smallgeo.tif", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	redband := ds.RasterBand(1)

	// create 1x1 pixelled tiff which would store overview of  redband
	dstile := memDrv.Create("", 1, 1, 3, Byte, nil)
	redbandDsTile := dstile.RasterBand(1)
	greenbandDsTile := dstile.RasterBand(2)
	blueBandDsTile := dstile.RasterBand(3)
	bands := []RasterBand{
		redbandDsTile,
		greenbandDsTile,
		blueBandDsTile,
	}
	err = redband.RegenerateOverviews(3, &bands[0], "average", DummyProgress, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	// generate PNG
	outDrv, err := GetDriverByName("PNG")
	if err != nil {
		panic(err)
	}
	outDrv.CreateCopy(makeTempFilePath("temp.png"), dstile, 0, nil, nil, nil)

	imgFile, err := os.Open(makeTempFilePath("temp.png"))
	if err != nil {
		t.Errorf(err.Error())
	}
	imgData, err := png.Decode(imgFile)
	if err != nil {
		t.Errorf(err.Error())
	}

	// assert that png is of 1x1 pixel
	assert.Equal(t, imgData.Bounds().Min.X, 0)
	assert.Equal(t, imgData.Bounds().Min.Y, 0)
	assert.Equal(t, imgData.Bounds().Max.X, 1)
	assert.Equal(t, imgData.Bounds().Max.Y, 1)

	// assert that image contains values RGB(32896, 32896, 32896)
	// 32896 is the computed average of r, g, b values from smallgeo.tif
	r, g, b, _ := imgData.At(0, 0).RGBA()
	assert.Equal(t, int(r), 32896)
	assert.Equal(t, int(g), 32896)
	assert.Equal(t, int(b), 32896)
}

func TestReadBlock(t *testing.T) {
	ds, err := Open("testdata/smallgeo.tif", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	redband := ds.RasterBand(1)
	x, y:= redband.BlockSize()
	block := make([]byte, x*y)
	err = redband.Read(0, 0, block)
	if err != nil {
		t.Errorf(err.Error())
	}
	println(fmt.Sprintf("%v", block))
	assert.Equal(t, len(block), 16)
	assert.Equal(t, int(block[0]), 255)
}

func TestGetSpatialRef(t *testing.T) {
	ds, err := Open("testdata/smallgeo.tif", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	srs := ds.SpatialRef()
	assert.Equal(t, srs.IsGeographic(), false)
}
