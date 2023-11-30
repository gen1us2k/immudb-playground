# immudb-playground

A simple demo inspired by immudb blog post to store k8s events immutable. 

## Description
The operator reads all events in a Kubernetes cluster and stores them in immudb vault.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).


## Running it outside of the K8s cluster

It uses the default kubeconfig stored in `~/.kube/config` for authentication and requires reading events permission across all namespaces 

```
export IMMUDB_API_KEY=your_key
go run main.go
```

Visit vault.immudb.io to get the data 



## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

