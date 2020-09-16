module "vpc" {
  source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
  version = "~> 2.51"
}

module "vpc_explicit" {
  source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
  version = "2.46.0"
}

# module "vpc_error" {
#   source  = "registry.development-stack.com/terraform-aws-modules/vpc/aws"
#   version = "~> 5.0"
# }
