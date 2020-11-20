# WASM
A go package for WebAssembly. 
This package contains some functions to be simple WebAssembly applications developed with go language.
Allows to access and manipulate the LocalStorage, IndexeDB and DOM.

For now the code documentation in only available in Spanish language.

Warning: this package only has been tried with go version 1.15.2 or high.

## Tests

If you want run tests, only go to test directory and execute  `./run_tests.sh`.
Then, open your browser in direction http://localhost:8080 .

For wasm compilations you can't use the default tests package. For this, 
I was built a simple bash script to execute custom tests functions. Essentially, 
this script copy the wasm_exec.js file (from go instalation files), build the 
file test.go and run a server in :8080 port.