output "cluster_id" {
  description = "ID созданного кластера Kubernetes"
  value       = module.k8s.cluster_id
}

output "external_v4_endpoint" {
  description = "Публичный endpoint кластера"
  value       = module.k8s.external_v4_endpoint
}

output "k8s_connect_command" {
  description = "Команда для подключения к кластеру"
  value       = "yc managed-kubernetes cluster get-credentials --id ${module.k8s.cluster_id} --external"
}

output "ci_cd_test_run" {
  value = null_resource.ci_cd_test.triggers.run_id
}