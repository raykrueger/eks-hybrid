package api

import "fmt"

type validator struct {
	validatorFuncs []func(*NodeConfig) error
}

func NewValidator(cfg *NodeConfig) *validator {
	v := newWithCommonValidations()
	if cfg.IsOutpostNode() {
		v.withOutpostValidations()
	} else if cfg.IsHybridNode() {
		v.withHybridValidations()
	} else {
		v.withEc2UserdataValidations()
	}
	return v
}

func (v *validator) Validate(cfg *NodeConfig) error {
	for _, f := range v.validatorFuncs {
		if err := f(cfg); err != nil {
			return err
		}
	}
	return nil
}

func newWithCommonValidations() *validator {
	v := new(validator)
	v.validatorFuncs = append(v.validatorFuncs, func(cfg *NodeConfig) error {
		if cfg.Spec.Cluster.Name == "" {
			return fmt.Errorf("Name is missing in cluster configuration")
		}
		return nil
	})
	return v
}

func (v *validator) withEc2UserdataValidations() {
	v.validatorFuncs = append(v.validatorFuncs, func(cfg *NodeConfig) error {
		if cfg.Spec.Cluster.APIServerEndpoint == "" {
			return fmt.Errorf("Apiserver endpoint is missing in cluster configuration")
		}
		if cfg.Spec.Cluster.CertificateAuthority == nil {
			return fmt.Errorf("Certificate authority is missing in cluster configuration")
		}
		if cfg.Spec.Cluster.CIDR == "" {
			return fmt.Errorf("CIDR is missing in cluster configuration")
		}
		return nil
	})
}

func (v *validator) withOutpostValidations() {
	v.validatorFuncs = append(v.validatorFuncs, func(cfg *NodeConfig) error {
		if cfg.Spec.Cluster.ID == "" {
			return fmt.Errorf("CIDR is missing in cluster configuration")
		}
		return nil
	})
}

func (v *validator) withHybridValidations() {
	v.validatorFuncs = append(v.validatorFuncs, func(cfg *NodeConfig) error {
		if cfg.Spec.Cluster.Region == "" {
			return fmt.Errorf("Region is missing in cluster configuration")
		}
		if !cfg.IsIAMRolesAnywhere() && !cfg.IsSSM() {
			return fmt.Errorf("Either IAMRolesAnywhere or SSM must be provided for hybrid node configuration")
		}
		if cfg.IsIAMRolesAnywhere() && cfg.IsSSM() {
			return fmt.Errorf("Only one of IAMRolesAnywhere or SSM must be provided for hybrid node configuration")
		}
		if cfg.IsIAMRolesAnywhere() {
			if cfg.Spec.Hybrid.IAMRolesAnywhere.AssumeRoleARN == "" {
				return fmt.Errorf("AssumeRoleARN is missing in hybrid iam roles anywhere configuration")
			}
			if cfg.Spec.Hybrid.IAMRolesAnywhere.RoleARN == "" {
				return fmt.Errorf("RoleARN is missing in hybrid iam roles anywhere configuration")
			}
			if cfg.Spec.Hybrid.IAMRolesAnywhere.ProfileARN == "" {
				return fmt.Errorf("ProfileARN is missing in hybrid iam roles anywhere configuration")
			}
			if cfg.Spec.Hybrid.IAMRolesAnywhere.TrustAnchorARN == "" {
				return fmt.Errorf("TrustAnchorARN is missing in hybrid iam roles anywhere configuration")
			}
		}
		if cfg.IsSSM() {
			if cfg.Spec.Hybrid.SSM.ActivationCode == "" {
				return fmt.Errorf("ActivationCode is missing in hybrid ssm configuration")
			}
			if cfg.Spec.Hybrid.SSM.ActivationID == "" {
				return fmt.Errorf("ActivationID is missing in hybrid ssm configuration")
			}
		}
		return nil
	})
}
