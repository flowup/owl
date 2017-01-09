# Owl
## Overview
Owl watches all files in current folder and all files in its subfolders. Every time file changes it runs the command.
## Usage
###Flags 
- `-i` or `--ignore` to ignore folder <br>
- `-r` or `--run` for specific command <br>

###Example
`owl -r 'echo \"some file was changed\"' -r 'ls -a' --ignore \"vendor\" --ignore \".git\" --ignore \".glide\"`
