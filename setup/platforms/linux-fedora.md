# Sandpiper Setup on Fedora

# Environment

**Host OS:** Windows 10

**Target OS:** Fedora31 (VirtualBox/Vagrant)

# Steps

(1) **Install VirtualBox 6.1.8 (&quot;Windows hosts&quot;)**

[https://www.virtualbox.org/wiki/Downloads](https://www.virtualbox.org/wiki/Downloads)

(2) **Install Vagrant 2.2.9 (&quot;Windows 64-bit&quot;)**

[https://www.vagrantup.com/downloads.html](https://www.vagrantup.com/downloads.html)

(3) **Create a new vagrant file for Fedora31**

Open Windows PowerShell Terminal (PS prompts not shown)
```
cd $HOME
mkdir Vagrant/Fedora31
cd Vagrant/Fedora31
vagrant init generic/fedora31
vagrant up
```
(4) **Enable SSH** (I already had it setup, so not sure what is required here)
```
Vagrant ssh-config
Vagrant ssh
```
We're now in the virtual machine and have a command prompt
```
$ [vagrant@fedora31 ~]$
```
(5) **Install Git (Fedora)**
```
$ git version  # if a version shown, skip install
https://git-scm.com/downloads
$ sudo dnf install git
```
(6) **Install Go (Fedora 32 is required for go1.14)**
```
$ sudo dnf install -y golang
$ go version
go version go1.13.10 linux/amd64
$ export PATH=$PATH:$HOME/go/bin # add this to your $HOME/.bashrc
```
(7) **Install PostgreSQL 12**
```
$ sudo dnf upgrade -y
$ sudo reboot

$ sudo dnf install https://download.postgresql.org/pub/repos/yum/reporpms/F-31-x86\_64/pgdg-fedora-repo-latest.noarch.rpm
$ sudo dnf install postgresql12 postgresql12-server
$ sudo /usr/pgsql-12/bin/postgresql-12-setup initdb

$ sudo systemctl start postgresql-12
$ sudo systemctl enable postgresql-12
$ sudo systemctl status postgresql-12

$ sudo -u postgres psql -c "alter user postgres with password 'strongpass'"
$ sudo passwd postgres   # enter the same strong password as above

$ sudo -u postgres psql -l # should list the databases
$ sudo egrep '^(local|host).*[^5]$' /var/lib/pgsql/12/data/pg_hba.conf

if this displays any rows, you will need to edit the "pg_hba.conf" file...

$ sudo vi /var/lib/pgsql/12/data/pg_hba.conf # you know "vi", right?
Change "ident" or "peer" to "md5" in these lines:

local       all           all                          md5
host        all           all       127.0.0.1/32       md5

$ sudo systemctl restart postgresql-12
```
(8) **Install Taskfile.dev**
```
$ cd $HOME
$ curl -sL https://taskfile.dev/install.sh | sh
```
(9) **Get Sandpiper from GitHub**
```
$ cd $HOME
$ git clone https://github.com/sandpiper-framework/sandpiper.git
```
(10) **Compile Sandpiper**
```
$ cd $HOME/sandpiper
$ go mod download
$ task build
```
(11) **Create and Initialize database**
```
$ task init
$ mv cmd/cli/api-primary.yaml cmd/api/config.yaml
```
(12) **Test Server**
```
$ task server
you should see `http server started on ...`
ctrl-c  # to stop server
```
(13) **Follow instructions in Testing Workbook**

(14) **Finally, be sure to close the vagrant VM**
```
$ exit # to log out of linux
C:\> vagrant halt
```
