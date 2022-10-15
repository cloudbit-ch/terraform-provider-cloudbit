---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudbit_compute_security_group_rule Data Source - terraform-provider-cloudbit"
subcategory: ""
description: |-
  
---

# cloudbit_compute_security_group_rule (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) unique identifier of the security group rule
- `security_group_id` (Number) unique identifier of the security group

### Read-Only

- `direction` (String) direction of the security group rule (ingress or egress)
- `icmp` (Attributes) ICMP message of the security group rule (see [below for nested schema](#nestedatt--icmp))
- `ip_range` (String) ip range of the security group rule
- `port_range` (Attributes) port range of the security group rule (see [below for nested schema](#nestedatt--port_range))
- `protocol` (Attributes) protocol of the security group rule (see [below for nested schema](#nestedatt--protocol))
- `remote_security_group_id` (Number) unique identifier of the remote security group

<a id="nestedatt--icmp"></a>
### Nested Schema for `icmp`

Read-Only:

- `code` (Number) code of the ICMP message
- `type` (Number) type of the ICMP message


<a id="nestedatt--port_range"></a>
### Nested Schema for `port_range`

Read-Only:

- `from` (Number) starting port of the security group rule
- `to` (Number) ending port of the security group rule


<a id="nestedatt--protocol"></a>
### Nested Schema for `protocol`

Read-Only:

- `name` (String) protocol name of the security group rule
- `number` (Number) iana protocol number of the security group rule

