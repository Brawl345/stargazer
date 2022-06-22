# Stargazer

Tool to extract and repack STAR files from the PSX used by its package manager "PackmanJr".

More info: <https://playstationdev.wiki/ps2devwiki/index.php/STAR_Files>

## Usage

```txt
Usage: stargazer <operation> <arguments>           
                                                   
  To extract files:                                
    stargazer x <star file> [output dir (optional)]
                                                   
  To pack a folder:                                   
    stargazer p <input dir> <star file>            
```

If no output directory is given, the file is extracted to the file name minus the extension plus "`_extracted`" (
e.g. `xPackmanJr_0.105.star` -> `xPackmanJr_0.105_extracted`). Same goes for packing (it will append `_packed.star`).

**NOTE:** Packing is experimental since I have no way to test it and I'm not sure about the limitations of the system (e.g. filenames). I also don't know if the order of the files is relevant.

## Credits

Thanks to @martravi for helping with reverse-engineering!

## Changelog

### v2.0

- Add re-packing (experimental)

### v1.0

- Initial release
