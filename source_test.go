package main

import (
	"strings"
	"testing"
)

func TestLoadGalleryAnchorHREF(t *testing.T) {
	testHTML := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<title>Gallery</title>
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
	<meta name="keywords" content="Gallery" />
	<meta name="description" content="Gallery" />
	<!-- Start WOWSlider.com HEAD section -->
	<link rel="stylesheet" type="text/css" href="engine1/style.css" />
	<script type="text/javascript" src="engine1/jquery.js"></script>
	<!-- End WOWSlider.com HEAD section -->
</head>
<body background="Backround1.gif">
<body style="background-color:#000000">
	<!-- Start WOWSlider.com BODY section -->
	<div id="wowslider-container1">
	<div class="ws_images"><ul>
<li><img src="data/images/picture_1280_001.jpg" alt="Picture_1280_001" title="Picture_1280_001" id="wows1_0"/></li>
</ul></div>
<div class="ws_bullets"><div>
<a href="data/images/picture_1280_001.jpg" title="Picture_1280_001"><img src="data/tooltips/picture_1280_001.jpg" alt="Picture_1280_001"/>1</a>
<a href="data/images/picture_1280_002.jpg" title="Picture_1280_002">2</a>
<a href="data/images/picture_1280_003.jpg" title="Picture_1280_003">3</a>
<a href="data/images/picture_1280_004.jpg" title="Picture_1280_004">4</a>
<a href="data/images/picture_1280_005.jpg" title="Picture_1280_005">5</a>
<a href="data/images/picture_1280_006.jpg" title="Picture_1280_006">6</a>
<a href="data/images/picture_1280_007.jpg" title="Picture_1280_007">7</a>
<a href="data/images/picture_1280_008.jpg" title="Picture_1280_008">8</a>
<a href="data/images/picture_1280_009.jpg" title="Picture_1280_009">9</a>
<a href="data/images/picture_1280_010.jpg" title="Picture_1280_010">10</a>
</div></div>
<a href="/home">home</a>
<a href="other.png">PNG</a>
<!-- Generated by WOWSlider.com v5.5 -->
	<a href="#" class="ws_frame"></a>
	</div>
	<script type="text/javascript" src="engine1/wowslider.js"></script>
	<script type="text/javascript" src="engine1/script.js"></script>
	<!-- End WOWSlider.com BODY section -->
</body>
</html>
`

	pictures, err := loadGalleryAnchorHREF(strings.NewReader(testHTML))
	if err != nil {
		t.Fatalf("loadGalleryAnchorHREF returned an error: %v", err)
	}

	if len(pictures) != 10 {
		t.Errorf("%d should have been detected, but found %d", 10, len(pictures))
	}
}