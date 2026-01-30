variable "zone" {
  description = "Зона доступности"
  type    = string
  default = "ru-central1-b"
}

variable "network" {
  description = "Имя VPC сети"
  type    = string
  default = "k8s_network"
}

variable "k8s_subnet_a_name" {
  description = "Имя подсети"
  type    = string
  default = "k8s-subnet-a"
}

variable "k8s_subnet_a_v4_cidr_blocks" {
  description = "CIDR блок для подсети"
  type    = list(string)
  default = ["10.0.0.0/16"]
}

