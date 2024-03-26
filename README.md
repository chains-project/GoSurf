# run_capslock.py

## Overview
`run_capslock.py` analyzes capabilities within an entire Go module installed on your system. It utilizes the Capslock tool, natively used for analyzing one package or a list of packages. However, `run_capslock.py` extends this functionality by automatically identifying all packages within a Go module and their respective locations. Subsequently, it executes the Capslock tool on these packages.

Additionally, the repository includes analysis results for various popular Go modules.

## Table of Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Example](#example)

## Requirements
- Python 3.x
- Capslock tool installed (Refer to [Capslock GitHub repository](https://github.com/google/capslock/tree/main/docs) for installation instructions)

## Usage
To view the usage and all available parameters, run the script with the `--help` option:

```sh
python3 run_capslock.py --help
```

Here are the parameters you need to specify when using the script:

- `-output`: Specify the output format for the capabilities list (*'v'* for verbose, *'json'* for JSON) 
- `-module`: Specify the module path (relative or absolute)
- `-packages`: Specify the path and filename for the packages output file 
- `-packages`: Specify the path and filename for the capabilities ouput file 

`run_capslock.py` automatically infers the packages contained within the provided Go module and execute the Capslock tool on each of these packages.

## Example
Install the ethereum client go module:
```sh
go mod install github.com/ethereum/go-ethereum/ethclient@1.13.14
```

Generally, go packages and modules are installed in `$GOPATH/pkg/mod`. To analyze the installed go module:
```sh
python3 run_capslock.py -output=json -module=$GOPATH/pkg/mod/github.com/ethereum/go-ethereum\@v1.13.14/ -packages=results/pkgs.list -capabilities=results/caps.json
```


