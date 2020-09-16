module "vpc" {
  source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
  version = "~> 2.0"
}

module "vpc_explicit" {
  source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
  version = "~> 2.42.0"
}

# module "vpc_error" {
#   source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
#   version = "~> 5.0"
# }
