# Contributing

This repository includes a [GitHub Action] that builds the
colorschemes once a week. This keeps the colorschemes up-to-date
automatically.

## Building

### Dependencies

- `>=0.2.6` [base16-builder-go]

### Usage for template editing

1. Clone the template repository onto your system
1. Download the [latest base16-builder-go binary]
1. Execute the binary with the `-template-dir` arg and provide
   `/path/to/base16-xfce4-terminal` 

Or the above steps represented in shell commands:

```shell 
/path/to/base16-builder-go-binary \
  -template-dir "/path/to/base16-xfce4-terminal
```

### Usage for adding or editing a colorscheme

1. Clone the base16-xfce4-terminal
1. Clone [base16-schemes]
1. Download the [latest base16-builder-go binary]
1. Execute the binary with 
  - `-schemes-dir` arg - provide `/path/to/base16-scehemes`
  - `-template-dir` arg - provide `/path/to/base16-xfce4-terminal`
    `base16-xfce4-terminal` repository)

```shell
/path/to/base16-builder-go-binary \
  -schemes-dir /path/to/base16-schemes \
  -template-dir "/path/to/base16-xfce4-terminal
```

If you have more questions about [base16-builder-go], have a look at
the information on the GitHub page.

[base16-builder-go]: https://github.com/tinted-theming/base16-builder-go
[base16-schemes]: https://github.com/tinted-theming/base16-schemes
[GitHub Action]: .github/workflows/update.yml
[latest base16-builder-go binary]: https://github.com/tinted-theming/base16-builder-go/releases
