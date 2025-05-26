output "ecr_repository_urls" {
  value = {
    for name, repo in aws_ecr_repository.services :
    name => repo.repository_url
  }
}
