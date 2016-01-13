freeSP - SGE - Signal Graph Editor
==================================

SGE is the new frontend tool for freeSP. It lets you create and edit all
the artifacts signal graphs, signal processing libraries and platforms,
which are the input files for the freeSP toolchain. Details about the
concepts, tagets and methodologies of freeSP can be read
[here](https://github.com/axel.freesp/freesp/overview.html)

SGE has been written from scratch using Go and its
[GTK+3 bindings](https://github.com/gotk3/gotk3/). Please report any
bugs to [...]

## Getting Started

### Installation and Compilation

gotk3 currently requires GTK 3.6-3.16, GLib 2.36-2.40, and
Cairo 1.10 or 1.12.  A recent Go (1.3 or newer) is also required. See
also the documentation of [GTK+3 bindings](https://github.com/gotk3/gotk3/).

To install the latest SGE version:

```bash
$ go get github.com/gotk3/gotk3/gtk
$ go get github.com/axel-freesp/sge
```

Compilation runs best with

```bash
$ go install github.com/axel-freesp/sge/sge
```

### Environment Variables

SGE requires some environment variables

- *SGE_ICON_PATH* must point to the icon folder (it is planned to integrate
  the icons with the binary)
- *FREESP_PATH* is used for the file dialogs when opening or saving
  artifacts.
- *FREESP_SEARCH_PATH* lists all paths that may contain freeSP-libraries
  (see [freeSP overview](https://github.com/axel.freesp/freesp/overview.html)
  for details)

```bash
SGE_PATH=$GOPATH/src/github.com/axel-freesp/sge
export FREESP_PATH=$GOPATH/src/github.com/axel-freesp/part-2.0
export SGE_ICON_PATH=$SGE_PATH/icons
export FREESP_SEARCH_PATH="$FREESP_PATH"
```

### Example Session



## License

Package sge is licensed under the BSD 2-Clause License.
