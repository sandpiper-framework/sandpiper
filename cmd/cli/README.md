## Introduction

Level 1 implementation of Sandpiper means syncing "file-based" data-objects. The `sandpiper` CLI tool supports the adding and pulling of Level 1 files.
This functionality will be duplicated in the Admin screens, but is easier to implement than a full admin interface (since there is no UI and it just calls the exposed API)
and allows for hands-free automation.

## Add File-Based Objects
Implement the "add" command to add "file" data-objects to a slice in the pool. This command could be called by an internal PIM, for example, to "publish" completed delivery files.

### For example:

```
sandpiper add \
-u user                  \ # database user
-p password              \ # database password, prompt if not supplied
-slice "aap-brake-pads"  \ # slice-name
-type "aces-file"        \ # grain-type
-key  "brakes"           \ # grain-key
-noprompt                \ # don't prompt before over-writing
-f "acme_brakes_full_2019-12-12.xml" # file to add
```

This command adds the ACES xml file as a grain as defined by the supplied request body (see below).

The sandpiper server url is pulled from the "config" file. If the password is not provided, it would prompt for it.

Here is an example of a config.yml file.

```
server:
  port: :4040
  debug: true
  read_timeout_seconds: 10
  write_timeout_seconds: 5
```
  
### Sandpiper API

The following api is used to add an object with the sandpiper add command:

```
POST /slice/{slice-name}
```

The "response" would look like this:

```
{
	"id": "0d5e171e-d3c2-4ddb-bd37-92fda5eca8a1",
	"slice_id": "2bea8308-1840-4802-ad38-72b53e31594c",
	"grain_type": "aces-file",
	"grain_key": "disc-brakes",
	"encoding": "gzipb64"
}
```
 
## Pull File-Based Objects

Implement the "pull" command to retrieve "file" data-objects from an optional slice in the pool. If the slice is not supplied it will create a sub-directory for each one it finds.

#### For example:

```
sandpiper pull \
-s "aap-slice" \ # optional slice
-t "aces-file" \  # optional file-type
-d "publish"     # required output directory
```

The generated file structure might be something like:

```
publish
|-- aces-file
    |-- aap-slice
```
    
Leaving off an option will pull all available.