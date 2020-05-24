# Sandpiper Setup on Windows

# Environment

**Host OS:** Windows 10

Not virtualized

# Steps

(1) **Install GitHub Desktop**

One easy way to install Git under Windows is with GitHub Desktop. Their installer includes the command line tools as well as a Windows GUI.

https://desktop.github.com/

(2) **Install Go**

```
https://golang.org/dl/
download go1.14.3.windows-amd64.msi (or similar)
double-click on the installer
```

(3) PostgreSQL

https://www.postgresql.org/download/windows/

todo: add more setup instructions here...

(4) **Install Scoop (package manager)**

Prerequisites are [PowerShell 5](https://aka.ms/wmf5download) (or later, include PowerShell Core) and [.NET Framework 4.5](https://www.microsoft.com/net/download) (or later).
```
open a powershell windows and execute:
PS> set-executionpolicy remotesigned -scope currentuser
PS> iwr -useb get.scoop.sh | iex
PS> scoop bucket add extras
```

(5) **Install Taskfile.dev**

```
PS> scoop install task
```

(6) **Get Sandpiper from GitHub**

```
$ cd $HOME
$ git clone https://github.com/sandpiper-framework/sandpiper.git
```

(7) **Compile Sandpiper**

```
$ cd $HOME/sandpiper
$ go mod download
$ task build
```

(8) **Create and Initialize database**

```
$ task init
$ mv cmd/cli/api-primary.yaml cmd/api/config.yaml
```

(9) **Test Server**

```
$ task server
you should see `http server started on ...`
ctrl-c  # to stop server
```

# PowerShell Tips

## Environment Variables

To set a session-level environment variable in PowerShell

```
$env:SANDPIPER_USER='admin'
```
