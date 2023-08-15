# adif2cloudlog
Utility to tail WSJT-X ADIF logs and post the results to [Cloudlog](https://github.com/magicbug/Cloudlog)'s API. By default it only posts new entries added to the log. When started with the `-b` argument it posts the entire log file before tailing.

This version of the tool works with Cloudlog version 2 and later.

## Usage
```
./adif2cloudlog [-b] <ADIF log> <station location ID number> <cloudlog url>

-b If true, load entire log file from the beginning, otherwise tail the file, only posting new entries to Cloudlog.
```
`adif2cloudlog` requires an API key to call Cloudlog's REST API. The API key may be obtained through the Admin|API menu options in the Cloudlog dashboard. Choose the _Generate Key with Read & Write Access_ button. The key is passed to `adif2cloudlog` in the `CLOUDLOG_API_KEY` environment variable so that it will not be visible using the `ps` command.

## Example 
```
adif2logcloud ~/.local/share/WSJT-X/wsjtx_log.adi 1 http://cloudlog.example.com
```