---
title: Server provisioning
permalink: /docs/provisioning/
---

This guide installs the Icaro server components on a fresh Rocky Linux 9 host.
The deployment uses the `generic/rocky9` Vagrant box locally and the
`rockylinux-9-x64` image on DigitalOcean. CentOS hosts, existing-database
migrations, and in-place operating system upgrades are not supported.

## Getting the code

Clone the repository and enter the deployment directory:

```console
git clone https://github.com/nethesis/icaro.git
cd icaro/deploy
```

## Requirements

Install a recent Vagrant release and a currently supported Ansible Core release
(2.14 or newer) on the machine where you run the deployment. Install the
deployment's Ansible collection requirements with:

```console
ansible-galaxy collection install -r ansible/requirements.yml
```

For local provisioning, install VirtualBox. For DigitalOcean provisioning,
install the Vagrant provider plugin instead:

```console
vagrant plugin install vagrant-digitalocean
```

## Configure Icaro and TLS

Edit `ansible/group_vars/all.yml` and replace the example values. Set at least
`icaro.hostname`, every database password, and one of the following `icaro.tls`
modes:

- `auto` enables Caddy v2 automatic HTTPS. The hostname's public DNS record
  must point to the server and ports 80 and 443 must be reachable so Caddy can
  obtain and renew its certificate.
- `self_signed` uses Caddy's internal certificate authority. Clients must trust
  that authority or accept a certificate warning.
- `manual` reads `/opt/icaro/<hostname>_cert.pem` and
  `/opt/icaro/<hostname>_private_key.pem` without managing their renewal.

In `manual` mode, install both certificate files with group `caddy`. The
certificate can be world-readable, but the private key must only be readable by
root and the Caddy group:

```console
hostname=hotspot.example.com
sudo install -o root -g caddy -m 0644 certificate.pem \
  "/opt/icaro/${hostname}_cert.pem"
sudo install -o root -g caddy -m 0640 private_key.pem \
  "/opt/icaro/${hostname}_private_key.pem"
```

Caddy v2 keeps the dashboard, API proxies, and `/survey` on HTTPS. The
`/wings` captive portal remains available only over HTTP; all other HTTP paths
redirect to HTTPS.

The playbook preserves the deployment's existing SELinux policy: it disables
SELinux in `/etc/selinux/config` and immediately leaves enforcing mode when
necessary. The final Vagrant reboot starts Caddy and the other services in the
normal multi-user boot target.

## Provision locally

Start a Rocky Linux 9 VirtualBox guest:

```console
vagrant up --provider=virtualbox
```

Guest HTTPS port 443 is forwarded to host port 8080, and guest HTTP port 80 is
forwarded to host port 8081. The default `self_signed` mode is intended for
local testing.

## Provision on DigitalOcean

In `Vagrantfile`, replace `YOUR TOKEN` with a DigitalOcean API token. Upload an
SSH key to the DigitalOcean account, then replace `YOUR KEY NAME` with its name.
Create the Rocky Linux 9 droplet with:

```console
vagrant up --provider=digital_ocean
```

Destroy the managed machine with:

```console
vagrant destroy
```

Rebuild it from scratch with the command below. This deletes all data stored on
the machine:

```console
vagrant rebuild
```
