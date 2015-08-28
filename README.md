# Snail Server

## Quick Start
***

###1 [Download](https://github.com/osfx/snail/releases/tag/v0.1.0-beta) Snail for your platform.

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
