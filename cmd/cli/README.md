## Introduction

Level 1 implementation of Sandpiper means syncing "file-based" data-objects. I think it makes sense to implement two CLI tools for Level 1. This functionality would be duplicated when we create the Admin screens, but should be much easier to implement than a full admin interface (since there is no UI and it just calls the exposed API). Plus it will be useful for early testing.

## Add File-Based Objects
Implement the "add" command to add "file" data-objects to a slice in the pool. This command could be called by an internal PIM, for example, to "publish" completed delivery files.

### For example:

```
sandpiper add \
-u user         \ # database user
-p password     \ # database password, prompt if not supplied
-s "aap-slice"  \ # slice name
-t "aces-file"  \ # file-type
-f "acme_brakes_full_2019-12-12.xml" # file to add
```

This would add the ACES xml file as a data-object to the "aap-slice" slice.

The database connection information (including user and password) could be pulled from a "config" file. If the password is not provided, it would prompt for it.

Here is an example of a config.yml file.

```
database:
  psn: postgres://user:password@localhost:5432/sandpiper?sslmode=disabled
```
  
### Sandpiper API

The following api is used to add an object with the sandpiper add command:

```
POST /slice/{slice-name}
```

The "request body" would look like this:

```
{
  "token": "--jwt here--",
  "oid": "da9fd323-151a-4035-82db-7b18e4ba6c79",
  "type": "aces-file",
  "payload":"MQkxOTk1CTEJNQkJCQkJCQkJCQkJCQk3NDE4MglNR0MJNDU3M ...",
  ...
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