# SpaceExtract

Extract and analyze space data from PDF files, including calculating areas and perimeters of spaces.

## Install

To install SpaceExtract, use the following command:

```
go install github.com/Bamorph/SpaceExtract@latest
```

## Basic Usage

Run SpaceExtract with the path to your PDF file as an argument:

```
SpaceExtract.exe "C:\\PATHTOPDF"
```


### Verbose Mode

To enable verbose output, add the `-v` flag:


## Output

SpaceExtract will extract spaces from the PDF and export them to a CSV file. The output will include the title, area (in square meters), and perimeter (in meters) of each space.

## Example

Given a PDF file with space data, running SpaceExtract will produce the following output:

The data will also be saved to a CSV file named according to the PDF file's base name.
