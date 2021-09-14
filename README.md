# tszero
Filter for tar and zip archives that sets all timestamps to zero. For zip files,
it also removes extra data from file headers. This allows filtered versions of
archives with content that is the same except for timestamps to compare as 
equal and to have the same hash.

To use tszero on a tar file, the command syntax is
```shell
tszero -format tar fileName
```

For a zip file, the command syntax is
```shell
tszero -format zip fileName
```

The filtered archive file is written to stdout.

No binary distribution of tszero is currently available.

You can install tszero from source. To do this you must first have go 
installed. If it is not already installed you can download it from 
https://golang.org/dl/

Once go is installed, check out this repository and issue the command 
```shell
go install
```
