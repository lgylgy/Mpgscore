# Sheet extractor

```
$ ./sheet-extractor -h
sheet-extractor extracts player grades from mpg googlesheet.
The result is serialized in players.json.

  -credentials string
        google authorization file
  -jobs uint
        number of concurrent jobs (default 1)
  -output string
        output directory
  -spreadsheet string
        google spreadsheet identifier
```