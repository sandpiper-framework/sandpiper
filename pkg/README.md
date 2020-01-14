# `/pkg`

This folder contains all source code for the Sandpiper project (except `main` entry-points). There are `pkg` directories for each `cmd` as well as `internal`
and `shared` directories for helper functions.

There's a fine distinction between these two helper directories. Both contain general (non-cmd specific) code, but `shared` includes functions called
from `/cmd` (which cannot make calls to functions in `/pkg/internal`). We could have put these routines under a `/cmd/internal` folder, but wanted
to keep only executable directories under `cmd`. All the rest of the shared code is under `/pkg/internal`, indicating that these are not intended
for use outside Sandpiper.
