# Gallery Downloader

Simply download all pictures from a web page

## Usage

### From a remote web page
```
gallery-downloader -source https://website.example.com -output ~/all-images/
```

### From a local web page
```
gallery-downloader -source ./example.html -referer https://website.example.com/ -output ~/all-images/
```

## Galleries

### Type "AnchorHREF"

It downloads all the linked pictures from these HTML tags:
```html
<a href="picture1.jpg" title="picture title">picture 1</a>
<a href="picture2.jpg" title="picture title">picture 2</a>
```

### Type "ListItem"

It simply downloads all the linked pictures from these HTML tags:
```html
<li><img src="picture1.jpg" alt="picture1" title="picture1"/></li>
<li><img src="picture2.jpg" alt="picture2" title="picture2"/></li>
```

## Flags

```
  -base string
    	base URL when downloading relative images
  -config string
    	configuration file (default "config.json")
  -insecure-tls
    	Skip TLS certificate verification. Should only be enabled for testing locally
  -max-wait int
    	wait n milliseconds maximum before downloading the next image. Use 0 to deactivate (default 3000)
  -min-wait int
    	wait n milliseconds minimum before downloading the next image. Use 0 to deactivate (default 1000)
  -output string
    	output folder to store pictures
  -password string
    	password (if the http server needs basic authentication)
  -referer string
    	referer header for HTML file, or for downloading images from a local HTML file
  -source string
    	source HTML gallery
  -type string
    	type of gallery (AutoDetect, AnchorHREF, ListItem) (default "AutoDetect")
  -user string
    	user (if the http server needs basic authentication)
```