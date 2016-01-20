# README #

Telegraph plugin to push metrics on Warp10

### Telegraph output for Warp10 ###

* Execute a post http on Warp10 at every flush time configured in telegraph in order to push the metrics collected

### Install ###

* Download telegraph source files (https://github.com/influxdb/telegraf)

* Copy directory warp in the output directory (github.com/influxdb/telegraf/outputs)

* do the 'make' command

* Add following instruction in the config file (Output part) 

```
[[outputs.warp10]]
warpUrl = "https://warp1.cityzendata.net/api/v0/update" 
token = "token"
prefix = "telegraf."
debug = false

```

### Contact ###

* contact@cityzendata.com