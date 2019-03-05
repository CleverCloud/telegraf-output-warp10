//
//   Copyright 2019  SenX S.A.S.
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
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

// Warp10 plugin options
type Warp10 struct {
	Prefix  string
	WarpURL string
	Token   string
	Debug   bool
}

var sampleConfig = `
  # prefix for metrics class Name
  prefix = "Prefix"
  ## POST HTTP(or HTTPS) ##
  # Url name of the Warp 10 server
  warp_url = "WarpURL"
  # Token to access your app on warp 10
  token = "Token"
  # Debug true - Prints Warp communication
  debug = false
`

// MetricLine is a plain metric
type MetricLine struct {
	Metric    string
	Timestamp int64
	Value     string
	Tags      string
}

// Connect is a connection initialization to backend
// Warp10 doesn't have such connection to maintain
func (o *Warp10) Connect() error {
	return nil
}

func (o *Warp10) Write(metrics []telegraf.Metric) error {

	out := ioutil.Discard
	if o.Debug {
		out = os.Stdout
	}

	if len(metrics) == 0 {
		return nil
	}

	now := time.Now()
	collectString := make([]string, 0)
	index := 0

	for _, mm := range metrics {

		for k, v := range mm.Fields() {

			metric := &MetricLine{
				Metric:    fmt.Sprintf("%s%s", o.Prefix, mm.Name()+"."+k),
				Timestamp: now.Unix() * 1000000,
			}

			metricValue, err := buildValue(v)
			if err != nil {
				return err
			}

			metric.Value = metricValue

			tagsSlice := buildTags(mm.Tags())
			metric.Tags = strings.Join(tagsSlice, ",")

			messageLine := fmt.Sprintf("%d// %s{%s} %s\n", metric.Timestamp, metric.Metric, metric.Tags, metric.Value)

			collectString = append(collectString, messageLine)
			index++
		}
	}

	payload := fmt.Sprint(strings.Join(collectString, "\n"))
	req, err := http.NewRequest("POST", o.WarpURL, bytes.NewBufferString(payload))
	if err != nil {
		return err
	}

	req.Header.Set("X-Warp10-Token", o.Token)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Fprintf(out, "Failed to close Warp10 response body: %v", err.Error())
		}
	}()

	fmt.Fprintf(out, "response Status: %#v", resp.Status)
	fmt.Fprintf(out, "response Headers: %#v", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(out, "Failed to read Warp10 response body: %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to push data into Warp10: %v", string(body))
	}

	fmt.Fprintf(out, "response Body: %#v", string(body))

	return nil
}

func buildTags(ptTags map[string]string) []string {
	sizeTags := len(ptTags)
	sizeTags++
	tags := make([]string, sizeTags)
	index := 0
	for k, v := range ptTags {
		tags[index] = fmt.Sprintf("%s=%s", k, v)
		index++
	}
	tags[index] = fmt.Sprintf("source=telegraf")
	sort.Strings(tags)
	return tags
}

func buildValue(v interface{}) (string, error) {
	var retv string
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
		retv = fmt.Sprintf("'%s'", p)
		//		return retv, fmt.Errorf("unexpected type %T with value %v for Warp", v, v)
	}
	return retv, nil
}

// IntToString convert int64 into string
func IntToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}

// BoolToString convert bool into string
func BoolToString(inputBool bool) string {
	return strconv.FormatBool(inputBool)
}

// UIntToString convert uint64 into string
func UIntToString(inputNum uint64) string {
	return strconv.FormatUint(inputNum, 10)
}

// FloatToString convert float64 into string
func FloatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

// SampleConfig return a config example
func (o *Warp10) SampleConfig() string {
	return sampleConfig
}

// Description return plugin description
func (o *Warp10) Description() string {
	return "Configuration for Warp server to send metrics to"
}

// Close backend connection
// Warp10 doesn't have such connection to maintain
func (o *Warp10) Close() error {
	// Basically nothing to do for Warp10 here
	return nil
}

func init() {
	outputs.Add("warp10", func() telegraf.Output {
		return &Warp10{}
	})
}
