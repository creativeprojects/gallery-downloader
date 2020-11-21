package scan

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedAnchorHREF = []string{
		"data/images/picture_1280_001.jpg",
		"data/images/picture_1280_002.jpg",
		"data/images/picture_1280_003.jpg",
		"data/images/picture_1280_004.jpg",
		"data/images/picture_1280_005.jpg",
		"data/images/picture_1280_006.jpg",
		"data/images/picture_1280_007.jpg",
		"data/images/picture_1280_008.jpg",
		"data/images/picture_1280_009.jpg",
		"data/images/picture_1280_010.jpg",
	}

	expectedListItem = []string{
		"data1/images/picture001.jpg",
		"data1/images/picture002.jpg",
		"data1/images/picture003.jpg",
		"data1/images/picture004.jpg",
		"data1/images/picture005.jpg",
		"data1/images/picture006.jpg",
		"data1/images/picture007.jpg",
		"data1/images/picture008.jpg",
		"data1/images/picture009.jpg",
		"data1/images/picture010.jpg",
	}
)

func getTestData(t *testing.T, name string) []byte {
	testfile := fmt.Sprintf("test_data/%s.html", name)
	file, err := os.Open(testfile)
	if err != nil {
		t.Fatalf("cannot open test file %q: %s", testfile, err)
	}
	defer file.Close()

	output, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("cannot read test file %q: %s", testfile, err)
	}
	return output
}

func TestLoadGalleryAnchorHREF(t *testing.T) {
	galleryAnchorHREF := getTestData(t, "anchor_href")
	gallery, err := NewLegacyAnchorGallery(galleryAnchorHREF)
	require.NoError(t, err)

	pictures := gallery.Find()
	if len(pictures) != 10 {
		t.Errorf("%d pictures should have been detected, but found %d", 10, len(pictures))
	}
	assert.ElementsMatch(t, expectedAnchorHREF, pictures)
}

func TestEmptyGalleryAnchorHREF(t *testing.T) {
	galleryListItem := getTestData(t, "list_item")
	gallery, err := NewLegacyAnchorGallery(galleryListItem)
	require.NoError(t, err)

	pictures := gallery.Find()
	if len(pictures) != 0 {
		t.Fatal("no picture should have been detected")
	}
}

func TestLoadGalleryListItem(t *testing.T) {
	galleryListItem := getTestData(t, "list_item")
	gallery, err := NewLegacyListItemGallery(galleryListItem)
	require.NoError(t, err)

	pictures := gallery.Find()
	if len(pictures) != 10 {
		t.Errorf("%d pictures should have been detected, but found %d", 10, len(pictures))
	}
	assert.ElementsMatch(t, expectedListItem, pictures)
}
