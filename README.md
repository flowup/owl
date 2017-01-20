# Owl
## Overview
Owl watches all files in current folder and all files in its subfolders. Every time file changes it runs the command.

##Installation
```bash
$ go get github.com/flowup/owl/cmd/owl
```

## Usage

###Flags 
- `-i` or `--ignore` to ignore folder <br>  
- `-r` or `--run` for specific command <br>
- `-t` or `--time` waiting time for executing in miliseconds (default 500)<br>

###Config File
If no flags are present Owl tries to read config file `owl.yaml` like in example:
```
run: "echo \"Hello Owl!\""
time: 5000
verbose: true
ignore:
 - ".git"
 - "bin"
```

###Example
`owl --run 'echo \"some file was changed\"' --ignore 'vendor' --ignore ".git" --ignore .glide -t 10000`
