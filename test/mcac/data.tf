data "st-utilities_module_tmpl" "template" {
  module_info = var.module_info
  module_tmpl = var.module_tmpl
}

data "google_project" "project" {}

data "external" "docker_build" {
  program = [
    "docker", "run",
    "--rm",
    # "--pull", "always",
    "harbor.sige.la/devops-cr/20twenty:latest",
    "get-customer-domains",
    "--brand", "${var.module_info.brand}",
    "--env", "${var.module_info.env}",
    "--format", "terraform"
  ]
}
