# README #

Telegraph plugin to push metrics on Warp10

### Telegraph output for Warp10 ###

* Execute a post http on Warp10 at every flush time configured in telegraph in order to push the metrics collected

### Install ###

* Git clone / go get telegraph source files (https://github.com/influxdata/telegraf)

* In the telegraf main dir, add this plugin as git submodule
```
git submodule add -b master https://github.com/clevercloud/telegraf-output-warp10.git plugins/outputs/warp10
```

* Add the plugin in the plugin list, you need to add this line to plugins/output/all/all.go
```
_ "github.com/influxdata/telegraf/plugins/outputs/warp10"
```

* do the 'make' command

* Add following instruction in the config file (Output part)

```
[[outputs.warp10]]
warpUrl = "http://127.0.0.1:8080/api/v0/update"
token = "token"
prefix = "telegraf."
debug = false

```

### Contact ###

* contact@clever-cloud.com
