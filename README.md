# Go Radiko

[![godoc](https://godoc.org/github.com/yyoshiki41/go-radiko?status.svg)](https://godoc.org/github.com/yyoshiki41/go-radiko)
[![build](https://travis-ci.org/yyoshiki41/go-radiko.svg?branch=master)](https://travis-ci.org/yyoshiki41/go-radiko)
[![codecov](https://codecov.io/gh/yyoshiki41/go-radiko/branch/master/graph/badge.svg)](https://codecov.io/gh/yyoshiki41/go-radiko)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/go-radiko)](https://goreportcard.com/report/github.com/yyoshiki41/go-radiko)

The __unofficial__ [radiko.jp](https://radiko.jp/) APIs Client Library for Go

## About this fork

This repository is a fork of
[yyoshiki41/go-radiko](https://github.com/yyoshiki41/go-radiko).
It fixes timefree (timeshift) downloads, which broke when radiko changed its
timefree API.
radiko no longer serves `POST /v2/api/ts/playlist.m3u8`.

Differences from upstream:

- `TimeshiftPlaylistM3U8` builds the playlist URL from the new timefree
  endpoint instead of calling the removed v2 API.
- `GetTimeshiftPlaylistInfo` returns a `TimeshiftPlaylistInfo`.
  `TimeshiftPlaylistInfo.ChunkURL` builds a playlist URL for one time range.
  Callers download a program in chunks because the new endpoint limits the
  number of segments per request.
- `GetStreamPcHtml5URL` reads `/v3/station/stream/pc_html5/{station}.xml`.
- `GetTimefreePlaylistBaseURL` picks the timefree playlist base URL from that
  response, for either the area-free or the regular stream.
- `GenerateLsid` generates the random session ID (`lsid`) that the timefree
  endpoint requires.

The module path stays `github.com/yyoshiki41/go-radiko`.
To use this fork, add a `replace` directive to your `go.mod`:

```
replace github.com/yyoshiki41/go-radiko => github.com/garaemon/go-radiko <version>
```

[garaemon/radigo](https://github.com/garaemon/radigo), a fork of
[yyoshiki41/radigo](https://github.com/yyoshiki41/radigo), uses this fork.
garaemon/radigo records timefree programs in chunks with ffmpeg through
`GetTimeshiftPlaylistInfo`.

## Installation


- Go 1.7 or newer

```bash
$ go get github.com/yyoshiki41/go-radiko
```

## Usage

### ■ Default

```go
// authentication token is not necessary.
client, err := radiko.New("")
if err != nil {
	panic(err)
}
// Get programs data
stations, err := client.GetNowPrograms(context.Background())
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%v", stations)
```

### ■ Get & Set authentication token

```go
// 1. Create a new Client.
client, err := radiko.New("")
if err != nil {
	panic(err)
}

// 2. Enables and sets the auth_token.
// After client.AuthorizeToken() has succeeded,
// the client has the enabled auth_token internally.
authToken, err := client.AuthorizeToken(context.Background(), authKeyPath)
if err != nil {
	log.Fatal(err)
}
```

#### Premium member (Enable to use the [area free](http://radiko.jp/rg/premium/).)

```go
// 1. Create a new Client.
client, err := radiko.New("")
if err != nil {
	panic(err)
}

// 2. Login as the premium member
// After client.Login() has succeeded,
// the client has the valid cookie internally.
ctx := context.Background()
login, err := client.Login(ctx, "example@mail.com", "example_password")
if err != nil {
	log.Fatal(err)
}
if login.StatusCode() != 200 {
	log.Fatalf("Failed to login premium member.\nInvalid status code: %d",
		login.StatusCode())
}

// 3. Enables and sets the auth_token.
// After client.AuthorizeToken() has succeeded,
// the client has the enabled auth_token internally.
authToken, err := client.AuthorizeToken(context.Background(), authKeyPath)
if err != nil {
	log.Fatal(err)
}
```


### ■ Use your authentication token

```go
// If the auth_token is cached, set your token in HTTP Header like below.
client, err = radiko.New("auth_token")
if err != nil {
	panic(err)
}
```

## Examples

It is possible to try [examples](https://github.com/yyoshiki41/go-radiko/tree/master/examples).

```bash
# Get programs data
$ go run ./examples/main.go
# Get & Set auth_token
$ go run ./examples/auth/main.go
```

## Projects using go-radiko

- [yyoshiki41/radigo](https://github.com/yyoshiki41/radigo) - Record a radiko program.
- [garaemon/radigo](https://github.com/garaemon/radigo) - A fork of radigo that uses this fork for timefree recording.
- [juntaki/toyotomimi](https://github.com/juntaki/toyotomimi)

## License 

Licensed under the GPLv3 license for all open source applications.

Please do not use this project for commercial use.

## Author

Yoshiki Nakagawa
