provider "uniteddomains" {

}

data "uniteddomains_dns_domains" "domains" {

}

output "foo" {
  value = data.uniteddomains_dns_domains.domains
}
