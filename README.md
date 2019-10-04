# zauth

## Compile

```sh
$ go get -u github.com/gobuffalo/packr/packr # Only need to run this once
$ go run build.go
```

## Run

```sh
./zauth
```

## FAQ

### I get "Forbidden - CSRF token invalid" when logging in!

Check your `config.yml` for `production`. If it is set to `true` and you're not
running via SSL/TLS, the CSRF protection will not work. Change to `false` if
you're simply developing locally, or make sure it's served behind SSL.
