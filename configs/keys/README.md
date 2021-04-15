> Note! For the SSH module to work you will need to create a pair of keys and place them here. To do this simply type on your terminal, in the riotpot folder:

```bash
# Create a pair of private/public keys utilizing RSA 4096
$ ssh-keygen -t rsa -b 4096 -C "riotpot@riotpot.com" -f ./configs/keys/riopot_rsa
```

> Do not use a passphrase!