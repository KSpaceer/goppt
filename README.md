# goppt [![Go Report Card](https://goreportcard.com/badge/github.com/KSpaceer/goppt)](https://goreportcard.com/report/github.com/KSpaceer/goppt) [![Go Reference](https://pkg.go.dev/badge/github.com/KSpaceer/goppt.svg)](https://pkg.go.dev/github.com/KSpaceer/goppt)

Native Go text extractor from the legacy MS PPT (Microsoft PowerPoint) binary files.

## Example
```go
f, err := os.Open("testdata/simplepres.ppt")
if err != nil {
  handleErr(err)
}
text, err := goppt.ExtractText(f)
if err != nil {
  handleErr(err)
}
fmt.Println(text)
```
