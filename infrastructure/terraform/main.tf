terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }

  backend "s3" {
    endpoints = {
      s3 = "https://storage.yandexcloud.net"
    }
    bucket     = "s3-terraform-state-rusakov"
    region     = "ru-central1"
    key        = "prod/k8s.tfstate"
    
    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true 
    skip_s3_checksum            = true 
  }
}

provider "yandex" {
  zone = var.zone
}

resource "yandex_vpc_network" "k8s_network" {
  name = var.network
}

resource "yandex_vpc_subnet" "k8s_subnet" {
  network_id     = yandex_vpc_network.k8s_network.id
  name           = var.k8s_subnet_a_name
  v4_cidr_blocks = var.k8s_subnet_a_v4_cidr_blocks
  zone           = var.zone
}



module "k8s" {
  source = "./modules/terraform-yc-kubernetes-master"

  network_id = yandex_vpc_network.k8s_network.id

  master_locations = [
    {
      zone      = var.zone
      subnet_id = yandex_vpc_subnet.k8s_subnet.id
    }
  ]
  master_maintenance_windows = [] 
  public_access              = true            

  node_groups = {
    "worker-group" = {
      description = "Test node group"
      fixed_scale = {
        size = 1 
      }
      
      platform_id   = "standard-v3"
      node_cores    = 2
      node_memory   = 4
      disk_type     = "network-hdd"
      disk_size     = 30
      
      nat = true
    }
  }

}

resource "yandex_vpc_security_group_rule" "allow_http" {
  count = length(module.k8s.nodes_security_group_ids)
  security_group_binding = module.k8s.nodes_security_group_ids[count.index]
  direction              = "ingress"
  description            = "Allow Load Balancer"
  v4_cidr_blocks         = ["0.0.0.0/0"]
  port                   = 80 
  protocol               = "TCP"
}

resource "null_resource" "ci_cd_test" {
  triggers = {
    run_id = timestamp()
  }
}
