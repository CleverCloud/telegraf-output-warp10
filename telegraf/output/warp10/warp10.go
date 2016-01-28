//
//   Copyright 2016  Cityzen Data
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package warp10

import (
  "bytes"
  "fmt"
  "net/http"
  "sort"
  "strconv"
  "strings"
  "time"
  "io/ioutil"

  "github.com/influxdb/influxdb/client/v2"
  "github.com/influxdb/telegraf/outputs"
)

type Warp10 struct {
  Prefix string

  WarpUrl string

  Token string

  Debug bool
}

var sampleConfig = `
  # prefix for metrics class Name
  prefix = "telegraf."
  ## Telnet Mode ##
  # Url name of the Warp 10 server
  warpUrl = "localhost:4242/"
  # Token to access your app on warp 10
  token = "Token"
  # Debug true - Prints Warp communication
  debug = false
`

type MetricLine struct {
  Metric    string
  Timestamp int64
  Value     string
  Tags      string
}

func (o *Warp) Connect() error {
  // Test Connection to Warp Server
  
  return nil
}

func (o *Warp) Write(points []*client.Point) error {
  if len(points) == 0 {
    return nil
  }
  var timeNow = time.Now()
  collectString := make([]string, len(points))  
  index := 0
  for _, pt := range points {
    metric := &MetricLine{
      Metric:    fmt.Sprintf("%s%s", o.Prefix, pt.Name()),
      Timestamp: timeNow.Unix() * 1000000,
    }

    metricValue, buildError := buildValue(pt)
    if buildError != nil {
      fmt.Printf("Warp: %s\n", buildError.Error())
      continue
    }
    metric.Value = metricValue

    tagsSlice := buildTags(pt.Tags())
    metric.Tags = fmt.Sprint(strings.Join(tagsSlice, ","))

    messageLine := fmt.Sprintf("%v// %s{%s} %v \n", metric.Timestamp, metric.Metric, metric.Tags, metric.Value)
    if o.Debug {
      fmt.Print(messageLine)
    }

    collectString[index] = messageLine
    index += 1
  }

  payload := fmt.Sprint(strings.Join(collectString, "\n"))
  //defer connection.Close()
  req, err := http.NewRequest("POST", o.WarpUrl, bytes.NewBufferString(payload))
  req.Header.Set("X-Warp10-Token", o.Token)
  req.Header.Set("Content-Type", "text/plain")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))

  return nil
}

func buildTags(ptTags map[string]string) []string {
  sizeTags :=len(ptTags)
  sizeTags +=1
  tags := make([]string, sizeTags)
  index := 0
  for k, v := range ptTags {
    tags[index] = fmt.Sprintf("%s=%s", k, v)
    index += 1
  }
  tags[index] = fmt.Sprintf("source=telegraf")
  sort.Strings(tags)
  return tags
}

func buildValue(pt *client.Point) (string, error) {
  var retv string
  var v = pt.Fields()["value"]
  switch p := v.(type) {
  case int64:
    retv = IntToString(int64(p))
  case string:
    retv = fmt.Sprintf("'%s'", p)
  case bool:
    retv = BoolToString(bool(p))
  case uint64:
    retv = UIntToString(uint64(p))
  case float64:
    retv = FloatToString(float64(p))
  default:
    return retv, fmt.Errorf("unexpected type %T with value %v for Warp", v, v)
  }
  return retv, nil
}

func IntToString(input_num int64) string {
  return strconv.FormatInt(input_num, 10)
}

func BoolToString(input_bool bool) string {
  return strconv.FormatBool(input_bool)
}

func UIntToString(input_num uint64) string {
  return strconv.FormatUint(input_num, 10)
}

func FloatToString(input_num float64) string {
  return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (o *Warp) SampleConfig() string {
  return sampleConfig
}

func (o *Warp) Description() string {
  return "Configuration for Warp server to send metrics to"
}

func (o *Warp) Close() error {
  return nil
}

func init() {
  outputs.Add("warp", func() outputs.Output {
    return &Warp{}
  })
}

