// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

/*
 * This example file shows how to configure the oci provider to target a single region.
 */

// These variables would commonly be defined as environment variables or sourced in a .env file
variable "tenancy_ocid" {}

variable "user_ocid" {}
variable "fingerprint" {}
variable "private_key_path" {}
variable "compartment_ocid" {}

variable "region" {
  default = "us-ashburn-1"
}

variable "notification_topic_defined_tags_value" {
  default = "value"
}

variable "notification_topic_description" {
  default = "description"
}

variable "notification_topic_freeform_tags" {
  default = {
    "Department" = "Finance"
  }
}

variable "notification_topic_name" {
  default = "name"
}

variable "notification_topic_state" {
  default = "ACTIVE"
}

variable "tag_namespace_description" {
  default = "Just a test"
}

variable "tag_namespace_name" {
  default = "exampletagns"
}

provider "oci" {
  region           = "${var.region}"
  tenancy_ocid     = "${var.tenancy_ocid}"
  user_ocid        = "${var.user_ocid}"
  fingerprint      = "${var.fingerprint}"
  private_key_path = "${var.private_key_path}"
}
