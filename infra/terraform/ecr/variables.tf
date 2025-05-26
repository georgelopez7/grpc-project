variable "service_names" {
  description = "Names of the services to create ECR repositories for"
  type        = list(string)
  default     = ["gateway", "fraud", "validation"]
}
