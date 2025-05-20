# go-slogio

Same as [zapio.Writer](https://github.com/uber-go/zap/blob/fcf8ee58669e358bbd6460bef5c2ee7a53c0803a/zapio/writer.go), but for [slog](https://pkg.go.dev/log/slog). That's it, that's the library.

```shell
go get -u github.com/utkuozdemir/go-slogio
```

## Explanation

Sometimes you need to give an [io.Writer](https://pkg.go.dev/io#Writer) to an external API/library, so it can write its logs etc. into it.

But you have an `slog.Logger` at hand, not an `io.Writer`.

You can use this to fill the gap.

### Example

```golang
logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

writer := &slogio.Writer{Log: logger, Level: slog.LevelWarn}
defer writer.Close()

cmd := exec.Command("rm", "-rf", "-v", "/")
cmd.Stdout = writer
cmd.Stderr = writer

if err := cmd.Run(); err != nil {
	panic(err)
}

logger.Info("done")
```

### Disclaimer

This project is based on [uber-go/zap](https://github.com/uber-go/zap), licensed under the [MIT License](https://github.com/uber-go/zap/blob/fcf8ee58669e358bbd6460bef5c2ee7a53c0803a/LICENSE).

Original copyright:
Copyright (c) 2016-2017 Uber Technologies, Inc.

### License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
