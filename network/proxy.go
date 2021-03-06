/*
Copyright 2016, Cossack Labs Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package network

import (
	"context"
	"github.com/cossacklabs/acra/logging"
	"go.opencensus.io/trace"
	"net"
)

// ProxyWithTracing reads data from connFrom, writes data to connTo and trace with OpenCensus
func ProxyWithTracing(context context.Context, connFrom, connTo net.Conn, errCh chan<- error) {
	_, span := trace.StartSpan(context, "ProxyWithTracing")
	defer span.End()
	logger := logging.NewLoggerWithTrace(context)

	buf := make([]byte, 8192)
	for {
		n, err := connFrom.Read(buf)
		if err != nil {
			errCh <- err
			return
		}
		if n == 0 {
			logger.Warningln("Read 0 bytes")
			continue
		}
		for nTo := 0; nTo < n; {
			nn, err := connTo.Write(buf[nTo:n])
			nTo += nn
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}
