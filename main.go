package main

import (
	"fmt"
	host "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/uhost"
	lb "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/ulb"
	net "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/unet"
	vpc "github.com/pulumi/pulumi-ucloud/sdk/go/ucloud/vpc"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		newVpc, err := vpc.NewVpc(ctx, "pulumi_test_vpc", &vpc.VpcArgs{
			CidrBlocks: pulumi.StringArray{pulumi.String("10.0.0.0/16")},
			Name:       pulumi.String("pulumi_test"),
			Remark:     pulumi.String("For Pulumi Test Only"),
			Tag:        pulumi.String("pulumitest"),
		})
		if err != nil {
			return err
		}

		subnet, err := vpc.NewSubnet(ctx, "pului_test_subnet", &vpc.SubnetArgs{
			CidrBlock: pulumi.String("10.0.0.0/24"),
			Name:      pulumi.String("pulumi_test_subnet"),
			Remark:    pulumi.String("For Pulumi Test Only"),
			Tag:       pulumi.String("pulumitest"),
			VpcId:     newVpc.ID(),
		})
		if err != nil {
			return err
		}

		eip, err := net.NewEip(ctx, "ulb_eip", &net.EipArgs{
			Bandwidth:    pulumi.Int(200),
			ChargeMode:   pulumi.String("traffic"),
			InternetType: pulumi.String("bgp"),
			Name:         pulumi.String("pulumi_test_eip"),
			Remark:       pulumi.String("For Pulumi Test Only"),
			Tag:          pulumi.String("pulumitest"),
		})
		if err != nil {
			return err
		}

		newLb, err := lb.NewLb(ctx, "publiculb", &lb.LbArgs{
			ChargeType: pulumi.String("dynamic"),
			Internal:   pulumi.Bool(false),
			Name:       pulumi.String("pulumitestlb"),
			Remark:     pulumi.String("ForPulumiTestOnly"),
			SubnetId:   subnet.ID(),
			Tag:        pulumi.String("pulumitest"),
			VpcId:      newVpc.ID(),
		})
		if err != nil {
			return err
		}

		listener, err := lb.NewLbListener(ctx, "nginx", &lb.LbListenerArgs{
			ListenType:     nil,
			LoadBalancerId: newLb.ID(),
			Port:           pulumi.Int(8080),
			Protocol:       pulumi.String("tcp"),
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
				AvailabilityZone: pulumi.String(az[i]),
				ChargeType:       pulumi.String("dynamic"),
				ImageId:          pulumi.String(nginxImgId),
				InstanceType:     pulumi.String("n-highcpu-1"),
				Name:             pulumi.String(fmt.Sprintf("nginx-%d", i)),
				Remark:           pulumi.String("For Pulumi Test Only"),
				RootPassword:     pulumi.String("AreallyComplictedpassw0rd"),
				SubnetId:         subnet.ID(),
				Tag:              pulumi.String("pulumitest"),
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
				Port:           pulumi.Int(80),
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

		ctx.Export("public_ip", eip.PublicIp)
		return nil
	})
}

func lookupImageId(ctx *pulumi.Context, imgNameRegex string) (string, error) {
	images, err := host.LookupImages(ctx, &host.LookupImagesArgs{
		ImageType:  String("custom"),
		MostRecent: Bool(true),
		NameRegex:  String(imgNameRegex),
	})
	if err != nil {
		return "", err
	}

	imageResults := images.Images
	nginxImgId := imageResults[0].Id
	return nginxImgId, nil
}

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}
