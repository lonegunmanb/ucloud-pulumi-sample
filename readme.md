# UCloud Pulumi Sample Go

To run this sample, you need a UCloud account and a Pulumi account first, and Pulumi client [should be installed](https://www.pulumi.com/docs/get-started/install/), Hashicorp Packer [should be installed](https://packer.io/downloads.html).

Then, compile main.go first:

```bash
go build -o ucloud-pulumi-sample
```

which "ucloud-pulumi-sample" corresponding to name in Pulumi.yaml

You can config **required** UCloud Provider as [environment variable](https://github.com/lonegunmanb/pulumi-ucloud/blob/master/README.md#configuration).

Then login to pulumi:

```bash
pulumi login
```

Then let's create a new-empty stack:

```bash
pulumi stack init <STACK_NAME_HERE>
```

Before we create a new stack, we **must** create a new UCloud nginx image:

```bash
packer build nginx.json
```

Now we can create a new stack by executing pulumi up:

```bash
pulumi up
Previewing update (dev):

     Type                           Name                       Plan       Info
 +   pulumi:pulumi:Stack            ucloud-pulumi-sample-dev   create     
 +   ├─ ucloud:unet:Eip             ulb_eip                    create     
 +   ├─ ucloud:vpc:Vpc              pulumi_test_vpc            create     
 +   ├─ ucloud:vpc:Subnet           pului_test_subnet          create     
 +   ├─ ucloud:ulb:Lb               publiculb                  create     1 warning
 +   ├─ ucloud:ulb:LbListener       nginx                      create     
 +   ├─ ucloud:uhost:Instance       pulumi_host0               create     
 +   ├─ ucloud:uhost:Instance       pulumi_host1               create     
 +   ├─ ucloud:uhost:Instance       pulumi_host2               create     
 +   ├─ ucloud:unet:EipAssociation  eip_association            create     
 +   ├─ ucloud:ulb:LbAttachment     nginx-1                    create     
 +   ├─ ucloud:ulb:LbAttachment     nginx-0                    create     
 +   └─ ucloud:ulb:LbAttachment     nginx-2                    create     
 
Diagnostics:
  ucloud:ulb:Lb (publiculb):
    warning: urn:pulumi:dev::ucloud-pulumi-sample::ucloud:ulb/lb:Lb::publiculb verification warning: "charge_type": [DEPRECATED] attribute `charge_type` is deprecated for optimizing parameters
 

Do you want to perform this update? yes
Updating (dev):

     Type                           Name                       Status      Info
 +   pulumi:pulumi:Stack            ucloud-pulumi-sample-dev   created     
 +   ├─ ucloud:unet:Eip             ulb_eip                    created     
 +   ├─ ucloud:vpc:Vpc              pulumi_test_vpc            created     
 +   ├─ ucloud:vpc:Subnet           pului_test_subnet          created     
 +   ├─ ucloud:ulb:Lb               publiculb                  created     1 warning
 +   ├─ ucloud:uhost:Instance       pulumi_host1               created     
 +   ├─ ucloud:uhost:Instance       pulumi_host0               created     
 +   ├─ ucloud:uhost:Instance       pulumi_host2               created     
 +   ├─ ucloud:ulb:LbListener       nginx                      created     
 +   ├─ ucloud:unet:EipAssociation  eip_association            created     
 +   ├─ ucloud:ulb:LbAttachment     nginx-0                    created     
 +   ├─ ucloud:ulb:LbAttachment     nginx-2                    created     
 +   └─ ucloud:ulb:LbAttachment     nginx-1                    created     
 
Diagnostics:
  ucloud:ulb:Lb (publiculb):
    warning: urn:pulumi:dev::ucloud-pulumi-sample::ucloud:ulb/lb:Lb::publiculb verification warning: "charge_type": [DEPRECATED] attribute `charge_type` is deprecated for optimizing parameters
 
Outputs:
    public_ip: "1.1.1.1"

Resources:
    + 13 created

Duration: 26s

Permalink: https://app.pulumi.com/lonegunmanb/ucloud-pulumi-sample/dev2/updates/1

```

You can cleanup stack by destroy them:

```bash
pulumi destroy
```

## Incorrect Pulumi Destroy Order

For now, pulumi cannot figure out that all UCloud Instances are dependent on subnet, pulumi engine always try to delete subnet first, which causes destruction failure. I've submit an [issue](https://github.com/pulumi/pulumi/issues/3856), hope we can solve this problem soon.

If you are experiencing the very same problem as me, you can delete pulumi stack by:

```bash
pulumi stack rm <STACK_NAME_HERE> --force
```

Then delete all resources manually on UCloudd console website.

Enjoy.
