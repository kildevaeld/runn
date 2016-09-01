# Runn
Bundle commands

usage:

```bash

mkdir example
echo "console.log('Hello,', process.env.WHO, "!")" > example/javascript.js
echo "echo Hello shell" > example/shell.sh

runn gen example
runn add example

runn run example:javascript -e WHO=World
// Hello, World


```