# goppt

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
