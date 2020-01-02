# `/cmd`

This folder contains all applications for the Sandpiper project. Each sub-directory should include a `main` function and be named the same as the
executable it creates (e.g. `/cmd/api`).

Keep code in this directory as small as possible. It's common to have a small `main` function that imports and invokes code from functions from
the `/internal` and `/pkg` directories and nothing else. If you think code can be used in other projects, it should live in the `/pkg` directory.
If the code is not reusable, or if you don't want others to reuse it, put that code in the `/internal` directory.

