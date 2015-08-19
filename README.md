# caladbolg
caladbolg is a tool to send a go test coverage report to slack.

## Usage

```
$ go test -cover ./... | caladbolg --channel "#random" --name gopher --icon "https://blog.golang.org/gopher/gopher.png"
```

## Sample

![sample](https://cloud.githubusercontent.com/assets/1456047/9276533/829791fa-42df-11e5-8f43-1de16f541beb.png)

## Auther

[kyokomi](https://github.com/kyokomi)

## License

MIT
