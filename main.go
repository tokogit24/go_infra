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
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDkyzYvvsYz/VAzWSDiHgL1vO0DHHkvpHrsbBb2teBrJ51vL5w7tIuMfJLbAsACwNYc65VKCBIo0wfRu2x+D+Uwk2xIXsn6Ex7N43LfNfuJjcBkh4aLILYgGWt9v0bnPMXkxiJwOawTIRGUc2xJvddRl95cy7D0B4P8pnR3iyKesCtxMWvce2xL3rTiuRqFOcAkWK+0ui9wTtxDBsmnealRZsG5OHXkAsqPeUzkd2Clt3Z/cTgcdgvqJoRrubgm5MpGLK6/m0D4h89v6B9MUR5hbZbik2vj1Z/81hnUHQjdHVXbxaINL7d/gA2deyljpIZn+OCmuGn9OTYd/hYVUtjDlAS8xX5O/bqHP3dRk9yTYbWn4BEAn3WmvJTC55PoUShKBwYWl+cmEGt4Yr2woHOq2l9wRDEtY2bVSnGsGHwEgh6fRJzLyndsftQ4jw6SBH0GM9HVEa6TOXFay/Ahpfg4rS8sHGX0+mEvK8xf3GuJbS6Kk8kkw782PRFACbyHmJM= jasmine@DESKTOP-SJH3J0E"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-02bf8ce06a8ed6092"),
			KeyName:             kp.KeyName,
		})

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicHostName", jenkinsServer.PublicDns)

		return nil
	})
}
