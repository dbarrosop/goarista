// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// The occli tool is a simple client to dump in JSON or text format the
// protobufs returned by the OpenConfig gRPC interface.
package main

import (
	"flag"
	"sync"

	"github.com/aristanetworks/goarista/openconfig"
	"github.com/aristanetworks/goarista/openconfig/client"
	"github.com/golang/protobuf/proto"

	"github.com/aristanetworks/glog"
)

const defaultPort = "6042"

var jsonFlag = flag.Bool("json", false,
	"Print the output in JSON instead of protobuf")

func main() {
	username, password, subscriptions, addrs, opts := client.ParseFlags()

	publish := func(resp *openconfig.SubscribeResponse) {
		if resp.GetHeartbeat() != nil && !glog.V(1) {
			return // Log heartbeats with verbose logging only.
		}
		var respTxt string
		var err error
		if *jsonFlag {
			respTxt, err = openconfig.SubscribeResponseToJSON(resp)
			if err != nil {
				glog.Fatal(err)
			}
		} else {
			respTxt = proto.MarshalTextString(resp)
		}
		glog.Info(respTxt)
	}

	wg := new(sync.WaitGroup)
	for _, addr := range addrs {
		wg.Add(1)
		go client.Run(publish, wg, username, password, addr, subscriptions, opts)
	}
	wg.Wait()
}
