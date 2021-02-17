# IMG 
## Introduction
img is image batch processing tool. Currently only feature is to crop pictures to a given aspect ratio. However, it uses face detection to ensure the cropped picture keeps the faces in the picture
 
## How Does it work
I use my other project https://github.com/abhinababasu/facethumbnail to generate face detected cropped images.

## Build
Get the sources and then run the following command the first time
```
go get
```

Subsequently build using 
```
go build .
```

## Run
The tool needs the facefinder binary that ships with https://github.com/abhinababasu/facethumbnail in the local folder. `go get` run above 
already fetches that repo into your GOPATH. So use the following command to copy facefinder over

```
copy %GOPATH%\src\github.com\abhinababasu\facethumbnail\test\facefinder .
```
Now you can run img. Check usage using
```
img -h
```

Example below runs the generator with facedetection and verbose mode
```
img -src c:\Users\abhin\OneDrive\Frame -dst c:\temp\img -ratio 9:16 -face -v
```

You can skip the `-face` flag when running this tool for non-portrait albums like Landscapes

## Sample

| Image                               | Cropping 9:16                      | Cropping with face-detection       |
| :---------------------------------: | :--------------------------------: | :--------------------------------: |
| ![Portrait](./sample/portrait.JPG) | ![Crop](./sample/crop2.JPG) | ![Thumnail](./sample/crop1.JPG) |


