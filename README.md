# Stargazer

Library to handle STAR files from the PSX used by its package manager "PackmanJr". Comes with a CLI!

More info: <https://playstationdev.wiki/ps2devwiki/index.php/STAR_Files>

## Usage

### General usage

```txt
NAME:
   stargazer - A tool to handle PSX STAR files

USAGE:
   stargazer [global options] command [command options] [arguments...]

COMMANDS:
   unpack, u  Unpacks files from a STAR file
   info, i    Shows information about a STAR file
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --quiet, -q    Do not print any messages (default: false)
   --version, -v  print the version (default: false)    
```

### Unpack

```txt
NAME:
   stargazer unpack - Unpacks files from a STAR file

USAGE:
   stargazer unpack [command options] [arguments...]

OPTIONS:
   --input value, -i value   Path to STAR file
   --output value, -o value  Path to output directory. Defaults to '<input file without .star>_extracted'
```

If no output directory is given, the file is extracted to the file name minus the extension plus "`_extracted`" (
e.g. `xPackmanJr_0.105.star` -> `xPackmanJr_0.105_extracted`). Same goes for packing (it will append `_packed.star`).

### Info

```txt
NAME:
   stargazer info - Shows information about a STAR file

USAGE:
   stargazer info [command options] [arguments...]

OPTIONS:
   --input value, -i value  Path to STAR file
```

## Credits

Thanks to [@martravi](https://github.com/martravi) for helping with reverse-engineering!
