# Gallery Downloader

Download all images from a web page

It simply downloads all the linked pictures from these HTML tags:
```html
<a href="picture1.jpg" title="picture title">picture 1</a>
<a href="picture2.jpg" title="picture title">picture 2</a>
```

## Usage

### From a remote web page
```
gallery-downloader -source https://website.example.com -output ~/all-images/
```

### From a local web page
```
gallery-downloader -source ./example.html -referer https://website.example.com/ -output ~/all-images/
```
