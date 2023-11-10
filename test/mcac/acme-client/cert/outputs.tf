output "acme_client_keys" {
  value = {
    "privkey.pem"   = acme_certificate.certificate.private_key_pem
    "cert.pem"      = acme_certificate.certificate.certificate_pem
    "chain.pem"     = acme_certificate.certificate.issuer_pem
    "fullchain.pem" = <<-EOF
    ${acme_certificate.certificate.certificate_pem}

    ${acme_certificate.certificate.issuer_pem}
    EOF
    "full.pem"      = <<-EOF
    ${acme_certificate.certificate.certificate_pem}

    ${acme_certificate.certificate.issuer_pem}

    ${acme_certificate.certificate.private_key_pem}
    EOF
  }
}
