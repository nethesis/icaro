# Provisioning Rocky Linux 9 with Vagrant and Ansible

This deployment provisions a fresh Rocky Linux 9 host. It supports the
`generic/rocky9` Vagrant box for local virtualization and the
`rockylinux-9-x64` image on DigitalOcean. CentOS installations and in-place OS
or database upgrades are not supported.

Run the commands below from the `deploy` directory.

## Requirements

- A recent Vagrant installation and either VirtualBox or a DigitalOcean
  account.
- A currently supported Ansible Core release (2.14 or newer) on the host that
  runs Vagrant.
- The Ansible collections declared by this deployment:

  ```console
  ansible-galaxy collection install -r ansible/requirements.yml
  ```

For DigitalOcean, install the provider plugin:

```console
vagrant plugin install vagrant-digitalocean
```

## Configuration

Edit `ansible/group_vars/all.yml` and replace the example values. At minimum,
set `icaro.hostname`, all database passwords, and the `icaro.tls` mode.

Caddy v2 supports these TLS modes:

- `auto`: Caddy obtains and renews a publicly trusted certificate. Public DNS
  for `icaro.hostname` must point to the server, and ports 80 and 443 must be
  reachable.
- `self_signed`: Caddy uses its internal certificate authority. Clients must
  trust that authority or accept a certificate warning.
- `manual`: Caddy loads an existing certificate and private key from
  `/opt/icaro/<hostname>_cert.pem` and
  `/opt/icaro/<hostname>_private_key.pem`.

For `manual` mode, both files must belong to the `caddy` group. Keep the
private key unreadable by other users:

```console
hostname=hotspot.example.com
sudo install -o root -g caddy -m 0644 certificate.pem \
  "/opt/icaro/${hostname}_cert.pem"
sudo install -o root -g caddy -m 0640 private_key.pem \
  "/opt/icaro/${hostname}_private_key.pem"
```

Caddy serves the APIs, dashboard, and survey over HTTPS. The `/wings` captive
portal remains HTTP-only, and other HTTP requests are redirected to HTTPS.

Caddy access logs are written to standard output in Apache Common Log Format
by default. Set `icaro.caddy_log_format: "json"` to use Caddy's native JSON
format instead. The default `common` setting installs the
[`transform-encoder`](https://github.com/caddyserver/transform-encoder) plugin
during provisioning.

The playbook intentionally disables SELinux persistently. If the fresh host is
enforcing during provisioning, it switches to permissive mode immediately and
becomes disabled after reboot.

## Local virtual machine

Start the Rocky Linux 9 VirtualBox guest with:

```console
vagrant up --provider=virtualbox
```

Guest ports 443 and 80 remain forwarded to host ports 8080 and 8081,
respectively. The default `self_signed` TLS mode is suitable for local testing.

## DigitalOcean

In `Vagrantfile`, replace `YOUR TOKEN` with a DigitalOcean API token. Upload an
SSH key to the account and replace `YOUR KEY NAME` with that key's name. Then
create the Rocky Linux 9 droplet:

```console
vagrant up --provider=digital_ocean
```

Provisioning ends with the existing reboot step so all services, including
Caddy, start in their normal multi-user boot target.

Destroy the managed machine with:

```console
vagrant destroy
```

Rebuild it from scratch with the following command. All data on the machine is
lost:

```console
vagrant rebuild
```
