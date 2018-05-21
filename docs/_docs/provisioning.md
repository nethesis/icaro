---
title: Server provisioning
permalink: /docs/provisioning/
---

This document will guide you on installing all server-side components: Sun, Wings, Wax.

Provisioning scripts will use latest [Icaro release](https://github.com/nethesis/icaro/releases).

## Getting the code

Clone the git repository:
```
git clone https://github.com/nethesis/icaro.git
```

Then enter the provisioning directory:
```
cd deploy
```

## Provisioning with Vagrant and Ansible

This procedure uses Vagrant to provision a Digital Ocean (DO) droplet.
If you prefere to use another cloud provider, edit ``Vagrantfile`` accordingly.


1. Install Vagrant along with Digital Ocean plugin.

   On Fedora 27:
   ```
   dnf install vagrant-digitalocean
   ```
   On Ubuntu 17.10
   ```
   apt-get install vagrant-digitalocean
   ```
2. Install ansible (>= 2.4)

   On Fedora 27:
   ```
   sudo dnf install ansible
   ```
   On Ubuntu 17.10:
   ```
   sudo apt-get install python-pip
   pip install ansible
   ```
   add ``$HOME/.local/bin/`` to your ``$PATH``.

3. Replace ``YOUR TOKEN`` with your DO token inside the ``Vagrantfile``

4. Make sure to upload an SSH key to yuor DO account, then replace ``YOUR KEY NAME``
   with you SSH key name inside the ``Vagrantfile``

5. Modify ``deploy/ansible/group_vars/all.yml`` file by replacing all variables.
   Make sure to customize at least ``icaro.hostname``.
   Set ``icaro.tls: "auto"`` if you want https certicate provisioned by Let's Encrypt.

6. Create the ``Icaro`` droplet:
   ```
   vagrant up --provider=digital_ocean
   ```

7. At the end, reboot the droplet:
   ```
   vagrant reload
   ```

You can destroy the droplet using:
```
vagrant destroy
```

The machine can be rebuilt using (all data will be lost):
```
vagrant rebuild
```
