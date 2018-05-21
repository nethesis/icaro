# Provisioning with Vagrant and Ansible

This procedure uses Vagrant to provision a Digital Ocean (DO) droplet.
If you prefere to use another cloud provider, edit ``Vagrantfile`` accordingly.


1. Install Vagrant along with Digital Ocean plugin.

   On Fedora 27:
   ```
   dnf install vagrant-digitalocean
   ```

2. Replace ``YOUR TOKEN`` with your DO token inside the ``Vagrantfile``

3. Make sure to upload an SSH key to yuor DO account, then replace ``YOUR KEY NAME``
   with you SSH key name inside the ``Vagrantfile``

4. Modify ``deploy/ansible/group_vars/all.yml`` file by replacing all variables.
   Make sure to customize at least ``icaro.hostname``.
   Set ``icaro.tls: "auto"`` if you want https certicate provisioned by Let's Encrypt.

5. Create the ``Icaro`` droplet:
   ```
   vagrant up --provider=digital_ocean
   ```

You can destroy the droplet using:
```
vagrant destroy
```

The machine can be rebuilt using (all data will be lost): 
```
vagrant rebuild
```
