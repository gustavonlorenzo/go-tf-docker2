terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Define Docker container resource
resource "docker_container" "postgres_pets" {
  name  = "postgres_pets"
  image = "postgres-pets"  # Specify your Docker image here
  ports {
    internal = 5432
    external = 5432
  }
  #env = [  # Define environment variables if needed
  #  "POSTGRES_DB=pets",
  #  "POSTGRES_USER=pets_user",
  #  "POSTGRES_PASSWORD=12345"
  #]
}

# Output container ID
output "container_id" {
  value = docker_container.postgres_pets.id
}
