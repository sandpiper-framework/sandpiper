# Sandpiper Development Setup (Windows)

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

Scoop provides an easy way to install `task`. You could also just download the correct Task binary from its [release page](https://github.com/go-task/task/releases) and add to your PATH. If you take that approach, skip this step and the next one. (Scoop can also be used to install other common utilities like 7zip, nodejs and yarn.)

To install Scoop, open a [PowerShell](https://docs.microsoft.com/en-us/powershell/) window and execute:

```
PS> set-executionpolicy remotesigned -scope currentuser
PS> iwr -useb get.scoop.sh | iex
PS> scoop bucket add extras
```

(5) **Install Taskfile.dev**

[Task](https://taskfile.dev/#/) is a cross-platform build tool designed to be easier than [Make](https://www.gnu.org/software/make/).

```
PS> scoop install task
```

(6) **Get Sandpiper from GitHub**

```
PS> cd $HOME
PS> git clone https://github.com/sandpiper-framework/sandpiper.git
```

(7) **Compile Sandpiper**

```
PS> cd $HOME/sandpiper
PS> go mod download
PS> task build
```

(8) **Create and Initialize database**

Follow the instructions found with the `sandpiper` CLI utility. The command is `sandpiper init` which is also included as a `task` command (as shown below).

```
PS> task init
```
Copy the server and command config files for default use (so you won't need a runtime parameter to select the correct file).
```
PS> copy cmd/cli/api-primary.yaml cmd/api/api-config.yaml
PS> copy cmd/cli/cli-primary.yaml cmd/api/cli-config.yaml
```

(9) **Test Server**

```
PS> task server
you should see `http server started on ...`
ctrl-c  # to stop server
```

# PowerShell Tips

## Environment Variables

To set a session-level environment variable in PowerShell

```
PS> $env:SANDPIPER_USER='admin'
```
