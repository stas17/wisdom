#Kozma Prutkov says
This is a service that sends one aphorism of Kuzma Prutkov.

####How it works
The service implements a few concepts for safety work.
- *The challenge-response protocol*. The first client must get a signature for the request.
- *Pow.* Then the client needs to compute the hash from signature, TimeStamp, and nonce by sha256 algorithm. That hash must contain zero value in first 24 bytes. Wisdom is never free.

####How to launch and test
For easy work, there is a makefile.
```shell

➜  wisdom git:(main) ✗ make help

Usage:
  make <target>

Targets:
  Go:
    dep                 Tidy up mod file
  Build server:
    bin-server          Build the project
  Build client:
    bin-client          Build the project
    clean               Remove build related file
  Test:
    test                Run the tests of the project
    generate            Generate mocks for testing
    coverage            Run the tests of the project with coverage
    lint                Lint with revive
  Local:
    local               Run a local instance of the application
  Help:
    help                Show this help.

```