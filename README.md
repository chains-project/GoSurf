# GoSurface

GoSurface is a tool that aims to analyze the potential attack surface of open-source Go packages and modules. It looks for occurrences of various features and constructs that could potentially introduce security risks, known as attack vectors.

## Repository Structure

- **attack_vectors**: This folder contains an analysis of 12 different attack vectors in Go, along with their respective proof-of-concept implementations.
- **gosurface.go**: The file `gosurface.go` file is the entry point for the GoSurface tool, which allows you to analyze a Go module and identify all the defined attack vectors, effectively framing the attack surface through Abstract Syntax Tree (AST) analysis.
- **libs**: This folder contains utility functions used by the GoSurface tool.
- **experiments**: This folder contains scripts to perform an analysis of the attack surface for different popular Go modules.

## Simple Usage
To use the GoSurface tool, follow these steps:

```bash
# Clone the repository
git clone https://github.com/chains-project/GoSurface.git

# Navigate to the gosurface directory
cd gosurface

# Build the tool
go build

# Analyze the github.com/ethereum/go-ethereum module
./gosurface $GOPATH/pkg/mod/github.com/ethereum/go-ethereum@v1.13.14

```
The tool will analyze the specified module and its direct dependencies,
identifying occurrences of the defined attack vectors, and print results on the CLI.


## Experiments
The `run_exp.go` script in the experiments folder allows for automating large-scale analysis on Go modules using the GoSurface library. To use this script, simply insert a list of "go_module_name version" entries in a text file.

Two experiments are pre-configured to run:

- **Experiment 1**: Analyzes 10 popular Go projects. The project names and versions are contained in the `urls_exp1.txt` file. To run this experiment, execute 

```bash
    go run run_exp.go exp1
```

- **Experiment 2**: Performs a differential analysis over versions for a single Go project (Kubernetes). The project name and versions to be analyzed are contained in the `urls_exp2.txt` file. To run this experiment, execute 

```bash
go run run_exp.go exp2
```

The results for the analysis will be reported in the `experiments/results` folder in HTML format.
 
