## krabbitmqctl

krabbitmqctl is a [`rabbitmqctl`](https://www.rabbitmq.com/rabbitmqctl.8.html)  command line wrapper for [kubernetes](https://kubernetes.io/) 

`rabbitmqctl` is not accesible outside the POD, with `krabbitmqctl` is possible to execute all the `rabbitmqctl` calls


## How it works

```
Usage:
  rabbitmqctl command [flags]

Flags:
      --context string      Kubernetes context to use. Default to current context configured in kubeconfig.
  -h, --help                help for rabbitmqctl
      --kubeconfig string   Path to kubeconfig file to use
  -n, --namespace string    Kubernetes namespace to use. Default to namespace configured in Kubernetes context (default "default")
  -s, --service string      RabbitMQ Service (default "rabbitmq")
```

for example:

```
➜ krabbitmqctl list_queues
Timeout: 60.0 seconds ...
Listing queues for vhost / ...
name    messages
test2   0
test1   0
test    0
```

or

```
➜ krabbitmqctl list_queues -n default -s rabbitmq
Timeout: 60.0 seconds ...
Listing queues for vhost / ...
name    messages
test2   0
test1   0
test    0
```

