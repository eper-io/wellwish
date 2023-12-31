package burst

import (
	"bytes"
	"fmt"
	"gitlab.com/eper.io/engine/billing"
	"gitlab.com/eper.io/engine/drawing"
	"gitlab.com/eper.io/engine/englang"
	"gitlab.com/eper.io/engine/management"
	"gitlab.com/eper.io/engine/mesh"
	"gitlab.com/eper.io/engine/metadata"
	"gitlab.com/eper.io/engine/stateful"
	"io"
	"net/http"
	"time"
)

// This document is Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
// to this document to the public domain worldwide.
// This document is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along with this document.
// If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.

// The main design behind burst runners is that they are designed to be scalable.
// Data locality means that data is co-located with burst containers.
// Data locality is important in some cases, especially UI driven code like ours.
// However, bursts are designed to handle the longer running processes with chaining.

// UI should be low latency using just bags and direct code.
// Bursts should scale out. They are okay to be located elsewhere than the data bags.
// The reason is that large computation will require streaming, and
// streaming is driven by pipelined steps without replies and feedbacks.
// Streaming bandwidth is not affected by co-location of data and code.
// Example: 1million 100ms reads followed by 100ms compute will last 200000 seconds.
// Example: 1million 100ms reads streamed into 100ms compute will last 100000 seconds,
// even if there is an extra network latency of 100ms per burst.

var startTime = time.Now()
var code = make(chan chan string)
var firstRun = true

func Setup() {
	stateful.RegisterModuleForBackup(&BurstSession)

	http.HandleFunc("/run", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		_, call := BurstSession[apiKey]
		if !call {
			management.QuantumGradeAuthorization()
			writer.WriteHeader(http.StatusPaymentRequired)
			drawing.NoErrorWrite(writer.Write([]byte("Payment required with a PUT to /run.coin")))
			return
		}

		input := drawing.NoErrorString(io.ReadAll(request.Body))
		callChannel := make(chan string)

		select {
		case <-time.After(MaxBurstRuntime):
			break
		case code <- callChannel:
			break
		}

		select {
		case <-time.After(MaxBurstRuntime):
			break
		case callChannel <- input:
			break
		}

		select {
		case <-time.After(MaxBurstRuntime + MaxBurstRuntime):
			break
		case output := <-callChannel:
			drawing.NoErrorWrite64(io.Copy(writer, bytes.NewBuffer([]byte(output))))
			break
		}
	})
	http.HandleFunc("/idle", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		if request.Method == "GET" {
			if apiKey == metadata.ActivationKey {
				// We may live without activation key
				// but this allows restricting the office cluster endpoint
				// to internal 127.0.0.1 addresses that was easier with udp.
				lock.Lock()
				idle := drawing.GenerateUniqueKey()
				ContainerRunning[apiKey] = fmt.Sprintf("Burst box %s registered at %s second.", idle, englang.DecimalString(int64(time.Now().Sub(startTime).Seconds())))
				ret := bytes.NewBufferString(idle)
				drawing.NoErrorWrite64(io.Copy(writer, ret))
				lock.Unlock()
				go func(key string) {
					if !firstRun {
						time.Sleep(MaxBurstRuntime * 2)
					}
					lock.Lock()
					ContainerRunning[key] = fmt.Sprintf("Burst box %s registered at %s second is ready.", idle, englang.DecimalString(int64(time.Now().Sub(startTime).Seconds())))
					lock.Unlock()
				}(idle)
				return
			}
			lock.Lock()
			v, ok := ContainerRunning[apiKey]
			if ok {
				var key, started string
				if nil == englang.Scanf1(v, "Burst box %s registered at %s second is ready.", &key, &started) {
					firstRun = false
					select {
					case <-time.After(MaxBurstRuntime):
						break
					case callChannel := <-code:
						request := <-callChannel
						delete(ContainerRunning, apiKey)
						ContainerResults[apiKey] = callChannel
						go func(key string) {
							time.Sleep(MaxBurstRuntime * 2)
							lock.Lock()
							delete(ContainerResults, key)
							lock.Unlock()
						}(apiKey)
						ret := bytes.NewBufferString(request)
						drawing.NoErrorWrite64(io.Copy(writer, ret))
						break
					}
				}
			} else {
				// Not ready
			}
			lock.Unlock()
			return
		}
		if request.Method == "PUT" {
			// TODO Get container result
			result := drawing.NoErrorString(io.ReadAll(request.Body))
			lock.Lock()
			replyCh, ok := ContainerResults[apiKey]
			if ok {
				select {
				case <-time.After(10 * time.Millisecond):
					break
				case replyCh <- result:
					break
				}
				delete(ContainerResults, apiKey)
			}
			lock.Unlock()
			return
		}
	})
	http.HandleFunc("/run.coin", func(w http.ResponseWriter, r *http.Request) {
		// Setup burst sessions, a range of time, when a coin can be used for bursts.
		if r.Method == "PUT" {
			coinToUse := billing.ValidatedCoinContent(w, r)
			if coinToUse != "" {
				func() {
					lock.Lock()
					defer lock.Unlock()
					// TODO generate new?
					burst := coinToUse
					// TODO cleanup
					BurstSession[burst] = englang.Printf(fmt.Sprintf("Burst chain api created from %s is %s/run.coin?apikey=%s. Chain is valid until %s.", coinToUse, metadata.Http11Port, burst, time.Now().Add(24*time.Hour).String()))
					mesh.SetExpiry(burst, ValidPeriod)
					mesh.RegisterIndex(burst)
					// TODO cleanup
					// mesh.SetIndex(burst, mesh.WhoAmI)
					management.QuantumGradeAuthorization()
					_, _ = w.Write([]byte(burst))
				}()
				return
			}
			management.QuantumGradeAuthorization()
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		if r.Method == "GET" {
			apiKey := r.URL.Query().Get("apikey")
			session, sessionValid := BurstSession[apiKey]
			if !sessionValid {
				management.QuantumGradeAuthorization()
				w.WriteHeader(http.StatusPaymentRequired)
				return
			}
			management.QuantumGradeAuthorization()
			_, _ = w.Write([]byte(session))
			return
		}
	})

	for i := 0; i < BurstRunners; i++ {
		// Normally this will be done by external docker containers
		// This is good for local in container testing
		go func() {
			time.Sleep(10 * time.Millisecond)
			// TODO docker
			_ = RunBox()
		}()
	}
}
