package gdal

import (
	"strings"
	"testing"
)



func TestVectorTranslate(t *testing.T) {
	srcDS, err := OpenEx("testdata/test.shp", OFReadOnly, nil, nil, nil)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-t_srs", "epsg:4326", "-f", "GeoJSON"}

	dstDS, err := VectorTranslate(makeTempFilePath("test4326.geojson"), []Dataset{srcDS}, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()
	dstDS, err = OpenEx(makeTempFilePath("test4326.geojson"), OFReadOnly|OFVector, []string{"geojson"}, nil, nil)
	if err != nil {
		t.Errorf("Open after translate: %v", err)
	}
	dstDS.Close()

}
func TestRasterize(t *testing.T) {
	srcDS, err := OpenEx("testdata/test.shp", OFReadOnly, nil, nil, nil)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-a", "code", "-tr", "10", "10"}

	dstDS, err := Rasterize(makeTempFilePath("rasterized.tif"), srcDS, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()
	dstDS, err = Open(makeTempFilePath("rasterized.tif"), ReadOnly)
	if err != nil {
		t.Errorf("Open after vector translate: %v", err)
	}
	dstDS.Close()

}

func TestWarp(t *testing.T) {
	srcDS, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-t_srs", "epsg:3857", "-of", "GPKG"}

	dstDS, err := Warp(makeTempFilePath("tiles-3857.gpkg"), nil, []Dataset{srcDS}, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}

	pngdriver, err := GetDriverByName("PNG")
	pngdriver.CreateCopy(makeTempFilePath("foo.png"), dstDS, 0, nil, nil, nil)
	dstDS.Close()
}

func TestTranslate(t *testing.T) {
	srcDS, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-of", "GTiff"}

	dstDS, err := Translate(makeTempFilePath("tiles.tif"), srcDS, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()

	dstDS, err = Open(makeTempFilePath("tiles.tif"), ReadOnly)
	if err != nil {
		t.Errorf("Open after raster translate: %v", err)
	}
	dstDS.Close()
}

func TestDEMProcessingColorRelief(t *testing.T) {
	srcDS, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-of", "GTiff"}

	dstDS, err := DEMProcessing(makeTempFilePath("demproc_output.tif"), srcDS, "color-relief", "testdata/demproc_colors.txt", opts)
	if err != nil {
		t.Errorf("DEMProcessing: %v", err)
	}
	dstDS.Close()

	dstDS, err = Open(makeTempFilePath("demproc_output.tif"), ReadOnly)
	if err != nil {
		t.Errorf("Open after raster DEM Processing: %v", err)
	}
	dstDS.Close()
}

func TestDEMProcessing(t *testing.T) {
	srcDS, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	opts := []string{"-of", "GTiff"}

	dstDS, err := DEMProcessing(makeTempFilePath("demproc_output_hillshade.tif"), srcDS, "hillshade", "", opts)
	if err != nil {
		t.Errorf("DEMProcessing: %v", err)
	}
	dstDS.Close()

	dstDS, err = Open(makeTempFilePath("demproc_output_hillshade.tif"), ReadOnly)
	if err != nil {
		t.Errorf("Open after raster DEM Processing: %v", err)
	}
	dstDS.Close()
}

func TestVectorInfo(t *testing.T) {
	ds, err := OpenEx("testdata/test.shp", OFReadOnly, nil, nil, nil)
	if err != nil {
		t.Errorf("Open: %v", err)
	}

	sql := "SELECT SUM(ST_Area(geometry)) AS TOTAL_AREA FROM test"
	opts := []string{"-sql", sql, "-dialect", "SQLite"}
	info := VectorInfo(ds, opts)
	
	expected := "INFO: Open of `testdata/test.shp' using driver `ESRI Shapefile' successful. Layer name: SELECT Geometry: None Feature Count: 1 Layer SRS WKT: (unknown) TOTAL_AREA: Real (0.0) OGRFeature(SELECT):0 TOTAL_AREA (Real) = 63754.5521128553"
	expected = strings.ReplaceAll(expected, " ", "")
	expected = strings.ReplaceAll(expected, "\n", "")
	info = strings.ReplaceAll(info, " ", "")
	info = strings.ReplaceAll(info, "\n", "")
	if info != expected {
		println("Expected: ", expected)
		println("Got: ", info)
		println("----")
		t.Errorf("Expected %s, got %s", expected, info)
	}
}
