# Concurrent CSV Read

We have generated a number of (csv) code files that should be guaranteed to have unique codes (e.g. X59J) not only within themselves, but also among each other.  Write code that can take the names of these files and verify the code uniqueness across all csv files in parallel. Bonus points if you can make it immediately stop itself once it has found a duplicate code.

## CLI tool

Initiate the program with a CLI command that takes a directory containing the CSV files as an arg.

## Data & Approach

### A Map container to hold the codes

I think we can rely on maps in Go having unique keys to help solve the problem. A package level type will be a map `map[string]bool`. This map will use each code as a key. The values aren't important, we're just depenedent on the property of maps that will allow us to determine whether a key (ie: code) exists in the map. If we check the map for each key before adding it, when we find one existing already, we'll know it's a duplicate.

Update: Using a sync.Map for this is a little cleaner and gets rid of the mutex locking/unlocking.

~~### Channels to hold state between goroutines [^1]~~

~~`chan bool` `quit`. Written to when neither `done` or `errc` get input, signifies there's nothing left to do because all rows of all files have been processed.~~
~~`chan error` `done`. Each time this receives increment a counter. When the counter value is the same as the number of files we're all done.~~
~~`chan error` `errc`. This channel gets an err when the map is found to already contain a key. When its listener reads a value~~

Channels were a little unwieldy and ultimately not the best tool.

The `golang.org/x/sync/errgroup` package provides a simmpler way to propogate errors up through goroutines and halt execution when one is found.[^2]

### A Reader function to read codes from files

A function `checkFilesForDuplicateCodes(ctx, filename string)` will run for all of the files

- Open an errgroup with context
- For each file, call the groups `Go()` to kick off a goroutine for that file, each will:
    - Open a CSV file for reading
    - Read each line to get its code
    - Check for code as a key in the codes map, if found write a `DuplicateKey` error to `errc`
    - If no duplicate is found, write the code to the map
    - When no more rows, write to the <-ctx.Done()

[^1]: https://stackoverflow.com/questions/40809504/idiomatic-goroutine-termination-and-error-handling

[^2]: https://www.oreilly.com/learning/run-strikingly-fast-parallel-file-searches-in-go-with-sync-errgroup