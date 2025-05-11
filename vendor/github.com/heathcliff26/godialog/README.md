[![CI](https://github.com/heathcliff26/godialog/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/godialog/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/heathcliff26/godialog/badge.svg)](https://coveralls.io/github/heathcliff26/godialog)
[![Editorconfig Check](https://github.com/heathcliff26/godialog/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/godialog/actions/workflows/editorconfig-check.yaml)
[![Generate go test cover report](https://github.com/heathcliff26/godialog/actions/workflows/go-testcover-report.yaml/badge.svg)](https://github.com/heathcliff26/godialog/actions/workflows/go-testcover-report.yaml)
[![Renovate](https://github.com/heathcliff26/godialog/actions/workflows/renovate.yaml/badge.svg)](https://github.com/heathcliff26/godialog/actions/workflows/renovate.yaml)

# GoDialog

GoDialog is a golang API for opening OS native file dialogs on linux/windows. Additionally it allows to define a fallback implementation should the native dialog not work.

## Usage

To get started run
```
go get github.com/heathcliff26/godialog@latest
```

## Examples

Please take a look at the [examples](examples).

To get started, here is a simple usage snippet:
```
fd := godialog.NewFileDialog()

res := make(chan string)

fd.Open("Test Dialog", func(s string, err error) {
	defer close(res)

	if err != nil {
		res <- fmt.Sprintf("Error: %v", err)
	} else {
		res <- fmt.Sprintf("Selected file: '%s'", s)
	}
})

println(<-res)
```
