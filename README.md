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
- `-t` or `--time` waiting time for executing in miliseconds (default 5000)<br>

###Example
`owl --run 'echo \"some file was changed\"' --ignore 'vendor' --ignore ".git" --ignore .glide -t 10000`
