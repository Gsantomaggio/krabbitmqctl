## krabbitmqctl

krabbitmqctl is a [`rabbitmqctl`](https://www.rabbitmq.com/rabbitmqctl.8.html)  command line wrapper for [kubernetes](https://kubernetes.io/) 

`rabbitmqctl` is not accesible outside the POD, with `krabbitmqctl` is possible to execute all the `rabbitmqctl` calls

## Install

Download the binary:

### Linux

```bash
curl -L https://github.com/Gsantomaggio/krabbitmqctl/releases/download/$(curl -sL https://raw.githubusercontent.com/Gsantomaggio/krabbitmqctl/master/version.txt)/krabbitmqctl_linux_amd64 -o krabbitmqctl
chmod +x krabbitmqctl
```
move it some `bin` location:

```bash
 sudo mv krabbitmqctl /usr/local/bin
```

### Mac

```bash
curl -L https://github.com/Gsantomaggio/krabbitmqctl/releases/download/$(curl -sL https://raw.githubusercontent.com/Gsantomaggio/krabbitmqctl/master/version.txt)/krabbitmqctl_darwin_amd64 -o krabbitmqctl
chmod +x krabbitmqctl
```
move it some `bin` location


### Windows

Download the windows binary from the releases: https://github.com/Gsantomaggio/krabbitmqctl/releases 




## How it works

```
rabbitmqctl kubernetes interface command

Usage:
  rabbitmqctl command [flags]

Flags:
      --context string      Kubernetes context to use. Default to current context configured in kubeconfig.
  -h, --help                help for rabbitmqctl
      --kubeconfig string   Path to kubeconfig file to use
  -n, --namespace string    Kubernetes namespace to use. Default to namespace configured in Kubernetes context (default "default")
  -p, --podname string      Pod where execute the command. Default is "" pick one random
  -s, --service string      RabbitMQ Service (default "rabbitmq")
  -v, --version             Print the version and exit
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

