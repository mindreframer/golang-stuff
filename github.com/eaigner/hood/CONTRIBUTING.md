If you want to contribute, please follow these rules.

## General

1. Code must be run through `go fmt`
2. Pull requests containing commit messages with unappropriate content (e.g. smilies) will be rejected
3. Pull requests with pending `TODO(xyz):`'s will be rejected

## Dialects

1. The dialect file name must reflect the actual dialect name, e.g. `sqlite3.go`
2. Unit and live test cases must be updated to run properly on the new dialect (see top of `dialects_test.go`)
3. Live tests should not run alongside real unit tests and should be commented out before pushing (again, see top of `dialects_test.go`).

