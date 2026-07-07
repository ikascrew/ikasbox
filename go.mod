module github.com/ikascrew/ikasbox

go 1.25.0

require (
	github.com/ikascrew/core v0.0.0-20210324041206-fb346c8e5c80
	github.com/ikascrew/plugin v0.0.0-20200715234203-87c9c5b19416
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/monochromegane/argen v0.0.0-20150711140148-c37112f9dc50
	gocv.io/x/gocv v0.38.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

require (
	github.com/fatih/color v1.19.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/monochromegane/goban v0.0.0-20141019070712-284a52313eb5 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

replace github.com/ikascrew/core => ../core

replace github.com/ikascrew/plugin => ../plugin
