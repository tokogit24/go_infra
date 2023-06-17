package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDevxctOV+DNsK5Yz0MR6n3HYXh8MPvQ7Eq1yb/vkLwmyr/XNMp2IN7/wG2Tt+qY3avLIDzBoKkEBSO9o31OuCjmmNUC9frPlRaJoeSwyN0+5JBlfpkaMRG1HJG0Quaep6eWyCpFufoQCEX5+H6FE2usQB34cGaVmWsP69MrhLkY4Q/JoJYlJnFqCjX7PWgjkpUWkMuH1EeWopL4pLWTRrdutleRBZYrNCVUocqm8KBX4/jm+hm7Vfx5ULLovvNpz4d06U9epFjY0FcqYQfeWPOp5AfJdeUItegyMDL34l8KqzJrJKWtnbAbiel8NN9pYG8DdQw22y2kTMno7p1IxW864bI8mTUI5MYrMNGo6t+yZVa1k162HzBnDHa8UHx/T6kvV4tCrWLpc+hxDkdpGiEFcEHx7Dk7ZXw2Gv/sHUIoxVgShlh7TvEm/fTITxdYcCuz7g/JPNzv+m6InmWNH9tW/P0a7dbPkT478/wObuzuaR0dtXxf9SM/bJYnmxEMFc= kantan@kantans-Virtual-Machine.local"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-0ca23709ed2a0fdf9"),
			KeyName:             kp.KeyName,
		})

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicHostName", jenkinsServer.PublicDns)

		return nil
	})
}
