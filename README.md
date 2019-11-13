# adif2cloudlog
Utility to tail WSJT-X ADIF logs and post the results to Cloudlog's API. By default it only posts new entries added to the log. When started with the `-b` argument it posts the entire log file before tailing.

## Usage
```
./adif2cloudlog [-b] <ADIF log> <cloudlog url>

-b If true, load entire log file from the beginning, otherwise tail the file, only posting new entries to Cloudlog.
```

## Example 
```
adif2logcloud ~/.local/share/WSJT-X/wsjtx_log.adi http://localhost/Cloudlog
```