package output

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

// Influxdb output struct
type Influxdb struct {
	bucket  string
	token   string
	org     string
	address string
	client  influxdb2.Client
	options map[string]string
}

const (
	tagKey = "host"
)

// StackErr add a new error
func (o *Influxdb) StackErr(key string, pre string, content string) {
	o.push(key, pre, content)
}

// StackOk add a new success
func (o *Influxdb) StackOk(key string, pre string, content string) {
	o.push(key, pre, content)
}

// Flush closes this output
func (o *Influxdb) Flush(string) {
	o.client.Close()
}

// Push pushes output
// https://www.influxdata.com/blog/getting-started-with-the-influxdb-go-client/
// https://docs.influxdata.com/influxdb/v1.8/write_protocols/line_protocol_tutorial/
func (o *Influxdb) push(key string, pre string, content string) {
	if len(key) < 1 {
		log.Errorf("empty key/value for %s", key)
		return
	}

	// add tag for host
	tags := make(map[string]string)
	tags[tagKey] = key

	log.Debugf("influxdb push %s/%s ===", pre, key)
	log.Debugf("influxdb push %s/%s - measurement: \"%s\"", pre, key, content)

	// TODO
	//api := o.client.WriteAPIBlocking(o.org, o.bucket)

	//// write new point
	//p := influxdb2.NewPoint(checkName, // measurement
	//	map[string]string{
	//		tagKey: hostName,
	//	}, // tag(s)
	//	keyval, // key-value fields
	//	time.Now())
	//api.WritePoint(context.Background(), p)

	//return nil
}

// NewInfluxdb creates a new output for influxdb
func NewInfluxdb(options map[string]string) (*Influxdb, error) {
	bucket, ok := options["bucket"]
	if !ok {
		return nil, fmt.Errorf("\"bucket\" option required")
	}
	token, ok := options["token"]
	if !ok {
		return nil, fmt.Errorf("\"token\" option required")
	}
	org, ok := options["org"]
	if !ok {
		return nil, fmt.Errorf("\"org\" option required")
	}
	address, ok := options["address"]
	if !ok {
		return nil, fmt.Errorf("\"address\" option required")
	}

	o := &Influxdb{
		bucket:  bucket,
		token:   token,
		org:     org,
		address: address,
		options: options,
		client:  influxdb2.NewClient(address, token),
	}
	return o, nil
}
