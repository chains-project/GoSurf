 # GoSurf 🏄

GoSurf is a tool that aims to analyze the potential attack surface of open-source Go packages and modules. It looks for occurrences of various features and constructs that could potentially introduce security risks, known as attack vectors.

## Repository Structure

- **attack_vectors**: This folder contains an analysis of 12 different attack vectors in Go, along with their respective proof-of-concept implementations.
- **experiments**: This folder contains scripts and results for attack surface analysis of different Go modules.
    - *popular10* contains experiments on the 10 most popular Go modules.
    - *top500* contains experiments on the 500 most imported Go modules.
- **libs**: This folder contains utility functions used by the GoSurf tool.
- **template**: This folder contains HTML templates used by the experiment scripts to print results.
- **gosurf.go**: The file `gosurf.go` file is the entry point for the GoSurf tool, which allows you to analyze a Go module and identify all the defined attack vectors, effectively framing the attack surface through Abstract Syntax Tree (AST) analysis.




## Simple Usage
To use the GoSurf tool, follow these steps:

```bash
# Clone the repository
git clone https://github.com/chains-project/GoSurf.git

# Navigate to the gosurf directory
cd gosurf

# Build the tool
go build

# Analyze the github.com/ethereum/go-ethereum module
./gosurf $GOPATH/pkg/mod/github.com/ethereum/go-ethereum@v1.13.14

```
The tool will analyze the specified module and its direct dependencies,
identifying occurrences of the defined attack vectors, and print results on the CLI.


## Experiments

#### Analyze Top 500 most imported modules
The `top500/run_exp.go` script in the experiments folder allows for automating large-scale analysis on 500 Go (most imported) modules using the GoSurf library. To use this script, simply run:

```bash
cd experiments/top500
go run run_exp.go 
```

The results for the analysis will be reported in the `experiments/top500/results` folder in HTML format.

#### Analyze custom list of modules
The `popular10/run_exp.go` script in the experiments folder allows for customized analysis on a set of selected packages. To use this script, insert a list of "go_module_name version" entries in a text file.

Two experiments are pre-configured to run:

- **Experiment 1**: Analyzes 10 popular Go projects. The project names and versions are contained in the `urls_exp1.txt` file. To run this experiment, execute 

```bash
    cd experiments/popular10
    go run run_exp.go exp1
```

- **Experiment 2**: Performs a differential analysis over versions for a single Go project (Kubernetes). The project name and versions to be analyzed are contained in the `urls_exp2.txt` file. To run this experiment, execute 

```bash
cd experiments/popular10
go run run_exp.go exp2
```

The results for the analysis will be reported in the `experiments/popular10/results` folder in HTML format.
 
>[!NOTE]
>These programs assume a [**Libraries.io API token**](https://libraries.io/api) stored in the environment variable `LIBRARIESIO_TOKEN`.