# ----------------------------------------------------------------------------------
# BECAUSE GOOGLE KSM KEYRING DESTROY BY SCHEDULER, SO WE DON'T USE ENCRYPTION FIRST.
# ----------------------------------------------------------------------------------

# data "google_storage_project_service_account" "gcs_account" {
# }

# data "google_iam_policy" "cloudkms" {
#   binding {
#     role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
#     members = [
#       "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
#     ]
#   }
# }

##################################################

# resource "google_kms_key_ring" "keyring" {
#   name     = "${var.brand}-${var.environment}-acme-client-${lower(var.region)}"
#   location = var.keyring_location
# }
#
# resource "google_kms_crypto_key" "kms_key" {
#   name            = "${var.brand}-${var.environment}-acme-client-${lower(var.region)}"
#   key_ring        = google_kms_key_ring.keyring.id
#   rotation_period = "100000s"
# }

# https://stackoverflow.com/questions/56320241/permission-denied-on-cloud-kms-key-when-using-cloud-storage
# resource "google_kms_key_ring_iam_policy" "key_ring" {
#   key_ring_id = google_kms_key_ring.keyring.id
#   policy_data = data.google_iam_policy.cloudkms.policy_data
# }

resource "google_storage_bucket" "storage_bucket" {
  location                    = var.region
  name                        = "${var.resource_name}-${var.ca}-${lower(var.region)}"
  uniform_bucket_level_access = true
  force_destroy               = true

  versioning {
    enabled = false
  }

  # encryption {
  #   default_kms_key_name = google_kms_crypto_key.kms_key.id
  # }

  labels = {
    app = var.app_tag
    ca  = var.ca
  }
}
