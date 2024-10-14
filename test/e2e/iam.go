package e2e

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

const assumeRolePolicyDocument = `{
	"Version": "2012-10-17",
	"Statement": [
	  {
		"Effect": "Allow",
		"Principal": {
		  "Service": "eks.amazonaws.com"
		},
		"Action": "sts:AssumeRole"
	  }
	]
  }`

const eksClusterPolicyArn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"

func (t *TestRunner) createEKSClusterRole() error {
	svc := iam.New(t.Session)
	roleName := getRoleName(t.Spec.ClusterName)

	// Create IAM role
	role, err := svc.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
	})
	if err != nil {
		return fmt.Errorf("failed to create role: %v", err)
	}

	// Attach AmazonEKSClusterPolicy
	_, err = svc.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(eksClusterPolicyArn),
	})
	if err != nil {
		return fmt.Errorf("failed to attach policy: %v", err)
	}
	t.Status.RoleArn = *role.Role.Arn
	fmt.Printf("Successfully created IAM role: %s\n", *role.Role.Arn)
	return nil
}

// deleteIamRoles deletes the IAM roles used for the cluster.
func (t *TestRunner) deleteIamRole() error {
	roleName := getRoleName(t.Spec.ClusterName)
	svc := iam.New(t.Session)

	_, err := svc.DetachRolePolicy(&iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(eksClusterPolicyArn),
	})
	if err != nil {
		return fmt.Errorf("failed to detach AmazonEKSClusterPolicy from role %s: %v", roleName, err)
	}

	fmt.Printf("Detached AmazonEKSClusterPolicy from role %s\n", roleName)

	_, err = svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete role %s: %v", roleName, err)
	}

	fmt.Printf("Deleted IAM role: %s\n", roleName)

	return nil
}

func getRoleName(name string) string {
	return fmt.Sprintf("%s-eks-hybrid-role", name)
}