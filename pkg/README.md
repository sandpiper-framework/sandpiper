# `/pkg`

This folder contains all source code for the Sandpiper project (except `main` entry-points). Each `cmd` has a matching directory under `pkg`.
Besides these matching directories, there are also `internal` and `shared` directories.

There's a fine distinction between these two directories. Both contain general (non-cmd specific) code, but `shared` includes functions called
from `/cmd` (which cannot make calls to functions in `/pkg/internal`. We could have put these routines under a `/cmd/internal`, but wanted
to keep only executable directories under `cmd`.) All the rest of the shared code is under `pkg/internal`, indicating that these are not intended
for use outside Sandpiper.
