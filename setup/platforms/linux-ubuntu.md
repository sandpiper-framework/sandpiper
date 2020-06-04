# Sandpiper Development Setup (Linux Ubuntu)

# Environment

**Host OS:** Ubuntu 18.04

**Digital Ocean Droplet**

# Steps

(1) **Install Git**
```
$ git version  # if a version shown, skip install (should already be installed on DigitalOcean)
$ sudo apt-get update
$ sudo apt-get install git

$ git config --global user.name "Your Name"
$ git config --global user.email "youremail@domain.com"
```
(2) **Install Go**
```
$ cd ~
$ curl -O https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
$ sha256sum go1.14.4.linux-amd64.tar.gz
(check hash against ones listed on download page)

$ tar xvf go1.14.4.linux-amd64.tar.gz

$ sudo chown -R root:root ./go
$ sudo mv go /usr/local
$ export PATH=$PATH:/usr/local/go/bin  # add this to your $HOME/.profile
$ source $HOME/.profile

$ go version
go version go1.14.4 linux/amd64

```
(3) **Install PostgreSQL 12**

This process is more complicated than it should be because we want v12 instead of v10 (as explained [here](https://itsfoss.com/install-postgresql-ubuntu/))
```
$ sudo apt-get install wget ca-certificates
$ wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
$ sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
```
Now we are good to go with a normal `apt` install.
```
$ sudo apt-get update
$ sudo apt-get install postgresql postgresql-contrib

$ sudo -u postgres psql -c "alter user postgres with password 'strongpass'"
$ sudo passwd postgres   # enter the same strong password as above

$ sudo -u postgres psql -l # should list the databases
$ sudo egrep '^(local|host).*[^5]$' /etc/postgresql/12/main/pg_hba.conf

if this displays any rows, you will need to edit the "pg_hba.conf" file...

$ sudo vi /etc/postgresql/12/main/pg_hba.conf   # or "nano" instead of "vi"
Change "ident" or "peer" to "md5" in these lines:

local       all           all                          md5
host        all           all       127.0.0.1/32       md5

$ sudo systemctl restart postgresql.service
```
(4) **Install Taskfile.dev**
```
$ cd $HOME
$ curl -sL https://taskfile.dev/install.sh | sh
$ mv bin/task /usr/local/bin
```
(5) **Get Sandpiper from GitHub**
```
$ cd $HOME
$ git clone https://github.com/sandpiper-framework/sandpiper.git
```
(6) **Compile Sandpiper**
```
$ cd $HOME/sandpiper
$ go mod download
$ task build
```
(7) **Create and Initialize Database**
```
$ task init
(follow separate instructions)

$ mv cmd/cli/api-primary.yaml cmd/api/config.yaml
```
(8) **Test Server**
```
$ task server
you should see `http server started on ...`
ctrl-c  # to stop server
```
(9) **Follow instructions in Testing Workbook**
