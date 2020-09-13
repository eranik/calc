# calc
Simple calculator doing add and subtract with rest endpoints
## Notes

* Added calc binary to repo as well.
* By default binds to "127.0.0.1:8000"
* No command line argument parsing or help is provided at this time. (TODO: use gflags)


## Dependencies
Uses gorilla/mux

## api endpoints
* Health check:
    * /health
    * /add/health
    * /subtract/health
    * all endpoints above with trailing / as well

* add and subtract:
    * /add?x=NUM&y=NUM
    * /add/?x=NUM&y=NUM
        * x and y are int
        * returns x + y as int
        * ignores overflow/underflow
    * /subtract?x=NUM&y=NUM
    * /subtract/?x=NUM&y=NUM
        * x and y are int
        * returns x - y as int
        * ignores overflow/underflow
    

## Usage and test
`go test`

`go build & ./calc &`

`curl '127.0.0.1:8000/add?x=1&y=2'`

`curl '127.0.0.1:8000/subtract?x=11&y=2'`

`curl '127.0.0.1:8000/health'`

`curl '127.0.0.1:8000/add/health'`

`curl '127.0.0.1:8000/subtract/health'`

`fg`

`Ctr + c`
