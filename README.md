# Stargazer

Tool to extract STAR files from the PSX used by its package manager "PackmanJr".

## Usage

```bash
stargazer <file> [output dir (optional)]
```

If no output directory is given, the file is extracted to the file name minus the extension plus "`_extracted`" (e.g. `xPackmanJr_0.105.star` -> `xPackmanJr_0.105_extracted`).

## Credits
Thanks to @martravi for helping with reverse-engineering!
