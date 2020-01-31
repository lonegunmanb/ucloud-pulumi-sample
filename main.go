package main

import (
	"fmt"
	host "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/uhost"
	lb "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/ulb"
	net "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/unet"
	vpc "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/vpc"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"reflect"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		newVpc, err := vpc.NewVpc(ctx, "pulumi_test_vpc", &vpc.VpcArgs{
			CidrBlocks: []string{"10.0.0.0/16"},
			Name:       "pulumi_test",
			Remark:     "For Pulumi Test Only",
			Tag:        "pulumitest",
		})
		if err != nil {
			return err
		}

		subnet, err := vpc.NewSubnet(ctx, "pului_test_subnet", &vpc.SubnetArgs{
			CidrBlock: "10.0.0.0/24",
			Name:      "pulumi_test_subnet",
			Remark:    "For Pulumi Test Only",
			Tag:       "pulumitest",
			VpcId:     newVpc.ID(),
		})
		if err != nil {
			return err
		}

		eip, err := net.NewEip(ctx, "ulb_eip", &net.EipArgs{
			Bandwidth:    200,
			ChargeMode:   "traffic",
			InternetType: "bgp",
			Name:         "pulumi_test_eip",
			Remark:       "For Pulumi Test Only",
			Tag:          "pulumitest",
		})
		if err != nil {
			return err
		}

		newLb, err := lb.NewLb(ctx, "publiculb", &lb.LbArgs{
			ChargeType: "dynamic",
			Internal:   false,
			Name:       "pulumitestlb",
			Remark:     "ForPulumiTestOnly",
			SubnetId:   subnet.ID(),
			Tag:        "pulumitest",
			VpcId:      newVpc.ID(),
		})
		if err != nil {
			return err
		}

		listener, err := lb.NewLbListener(ctx, "nginx", &lb.LbListenerArgs{
			ListenType:     nil,
			LoadBalancerId: newLb.ID(),
			Port:           8080,
			Protocol:       "tcp",
		})
		if err != nil {
			return err
		}

		nginxImgId, err := lookupImageId(ctx, "nginx")
		if err != nil {
			return err
		}

		var az = []string{
			"cn-sh2-01",
			"cn-sh2-02",
			"cn-sh2-03",
		}

		var hosts []*host.Instance
		for i := 0; i < len(az); i++ {
			instance, err := host.NewInstance(ctx, fmt.Sprintf("pulumi_host%d", i), &host.InstanceArgs{
				AvailabilityZone: az[i],
				ChargeType:       "dynamic",
				ImageId:          nginxImgId,
				InstanceType:     "n-highcpu-1",
				Name:             fmt.Sprintf("nginx-%d", i),
				Remark:           "For Pulumi Test Only",
				RootPassword:     "AreallyComplictedpassw0rd",
				SubnetId:         subnet.ID(),
				Tag:              "pulumitest",
				VpcId:            newVpc.ID(),
			})
			if err != nil {
				return err
			}
			hosts = append(hosts, instance)
		}

		for i := 0; i < len(hosts); i++ {
			_, err := lb.NewLbAttachment(ctx, fmt.Sprintf("nginx-%d", i), &lb.LbAttachmentArgs{
				ListenerId:     listener.ID(),
				LoadBalancerId: newLb.ID(),
				Port:           80,
				ResourceId:     hosts[i].ID(),
			})
			if err != nil {
				return err
			}
		}

		_, err = net.NewEipAssociation(ctx, "eip_association", &net.EipAssociationArgs{
			EipId:      eip.ID(),
			ResourceId: newLb.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("public_ip", eip.PublicIp())
		return nil
	})
}

func lookupImageId(ctx *pulumi.Context, imgNameRegex string) (string, error) {
	images, err := host.LookupImages(ctx, &host.LookupImagesArgs{
		ImageType:  "custom",
		MostRecent: true,
		NameRegex:  imgNameRegex,
	})
	if err != nil {
		return "", err
	}

	imageResults := images.Images.([]interface{})
	nginxImg := imageResults[0].(map[string]interface{})
	nginxImgId := nginxImg["id"].(reflect.Value).Interface().(string)
	return nginxImgId, nil
}
