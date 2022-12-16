module trasys

go 1.19

require (
	github.com/golang/protobuf v1.5.2
	server v0.0.0
)

require google.golang.org/protobuf v1.28.1

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/uber/go-torch v0.0.0-20181107071353-86f327cc820e // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)

replace server => ../server

replace trasys => ./trasys
