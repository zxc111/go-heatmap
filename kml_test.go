package heatmap

import (
	"archive/zip"
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/zxc111/go-heatmap/schemes"
)

const testKmlImgURL = "http://www.example.com/thing.png"

const expKml = `<?xml version="1.0" encoding="UTF-8"?>
		<kml xmlns="http://www.opengis.net/kml/2.2">
		<Folder>
		    <GroundOverlay>
		      <Icon>
		        <href>http://www.example.com/thing.png</href>
		      </Icon>
		      <LatLonBox>
		        <north>1.0844535751966644</north>
		        <south>-0.0815174737483185</south>
		        <east>1.0688449455069537</east>
		        <west>-0.0780794272341640</west>
		        <rotation>0</rotation>
		      </LatLonBox>
		    </GroundOverlay></Folder></kml>`

const expKmz = `<?xml version="1.0" encoding="UTF-8"?>
		<kml xmlns="http://www.opengis.net/kml/2.2">
		<Folder>
		    <GroundOverlay>
		      <Icon>
		        <href>heatmap.png</href>
		      </Icon>
		      <LatLonBox>
		        <north>1.0844535751966644</north>
		        <south>-0.0815174737483185</south>
		        <east>1.0688449455069537</east>
		        <west>-0.0780794272341640</west>
		        <rotation>0</rotation>
		      </LatLonBox>
		    </GroundOverlay></Folder></kml>`

func xsimilar(a, b string) bool {
	return strings.Join(strings.Fields(a), " ") ==
		strings.Join(strings.Fields(b), " ")
}

func TestKMLOutOfRange(t *testing.T) {
	kmlBuf := &bytes.Buffer{}

	_, err := KML(image.Rect(0, 0, 1024, 1024),
		append(testPoints, P(-200, 0)), 150, 128, schemes.AlphaFire,
		testKmlImgURL, kmlBuf)
	if err == nil {
		t.Fatalf("Expected error with bad input")
	}
}

func TestKMZBadInput(t *testing.T) {
	err := KMZ(image.Rect(0, 0, 1024, 1024),
		append(testPoints, P(-200, 0)), 150, 128, schemes.AlphaFire,
		ioutil.Discard)
	if err == nil {
		t.Fatalf("Expected error with bad input")
	}
}

type writeFailer struct {
	n int
}

func (w *writeFailer) Write(p []byte) (int, error) {
	towrite := len(p)
	var err error
	if towrite > w.n {
		towrite = w.n
		err = io.EOF
	}
	w.n -= towrite
	return towrite, err
}

func TestKMZBadWriter(t *testing.T) {
	err := KMZ(image.Rect(0, 0, 1024, 1024),
		testPoints, 150, 128, schemes.AlphaFire,
		&writeFailer{514})
	if err == nil {
		t.Fatalf("Expected error with bad input")
	}
}

func rzd(t *testing.T, zf *zip.File) []byte {
	r, err := zf.Open()
	if err != nil {
		t.Fatalf("Error opening %v: %v", zf.Name, err)
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Error reading %v: %v", zf.Name, err)
	}
	return data
}