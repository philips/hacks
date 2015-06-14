# Go RPC Comparison

## Goal

The goal of this repository is to compare RPC options for Golang and their ability to be useful.
The metrics used for usefulness include: 

- ease of use in Go
- compatibility with frontend web (XHR/websocket/etc)
- bandwidth
- latency
- memory usage
- language interoperability
- client generation
 
The methods of analysis will be benchmarks, research, and educated guesses.

## Options Considered

Any engineer can do a little Sprintf and Scanf, throw in some sockets and BOOM RPC method; and in fact many people have done just that (myself included).
But, we need to narrow the field and so we will focus on just a handful of options.

name             | transport | encoding 
-----------------|-----------|-----------
RESTish          | http      | json
RESTish+protobuf | http      | protobuf
net/rpc          | http      | gob
net/rpc/jsonrpc  | http      | json
capnproto        | tcp       | capnproto
grpc             | http2     | protobuf3

## RESTish



## net/rpc

## net/rpc/jsonrpc

## capnproto

## grpc

Package: google.golang.org/grpc
Version: 0acf7


