# Snail Server

## Quick Start
***

###1 Download Snail for your platform.

[LINUX_64](https://bitbucket.org/osfx/snail/downloads/snail_linux64.zip)

[MAC_64](https://bitbucket.org/osfx/snail/downloads/snail_darwin.zip)

[WINDOWS_64](https://bitbucket.org/osfx/snail/downloads/snail_windows.zip)

### 2 Create a config file with name Snailshell

```shell
$ echo "localhost:9999\n\nroot ./html" > Snailshell
```

### 3 Create web site root folder and add index page

```shell
$ mkdir html
$ echo "<h1>Hello Snail</h1>" > html/index.html
```
###4 Run server
```shell
$ ./snail
```
### 5 Check result
Open in your browser <http://localhost:9999>



![Placeholder](/img/snail.png)
