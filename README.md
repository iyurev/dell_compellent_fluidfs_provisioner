# dell_compellent_fluidfs_provisioner
Dell Compellent Kubernetes dynamic storage provisioner 

Show provisioner logs:

```oc exec <provisioner_pod_name> -- sh -c 'tail /tmp/*'```

Install provisioner :

1. Inside this repository move to ```k8s_deployment``` directory :
2. Run ```ansible-playbook  install_plb.yaml  -e dsm_username=$DSM_USERNAME  -e dsm_password=$DSM_PASSWORD -e image_pull_secret=$PULL_SECRET```

```PULL_SECRET - imagePullSecret in base64 encoding```

 

