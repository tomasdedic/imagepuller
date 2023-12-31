# imagepuller

1. define images and tags in YAML format to read from
2. define repositories TO which we will push images
3. get SHA256 checksum for each image tag defined in YAML file
4. if SHA256 chacksum is different pull the new image under same tag
5. create apppropriate tags TO repository
6. populate YAML file with new SHA1 checksum, tags
7. create history file with checksum and appropriate tags

Digest  "Digest": "sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f", for each container
will be saved 

out of app scope (pipeline scope):
- scan each image with Aquascan after pull
- clean ACR registry for untagged images after some time

## input YAML file format
```yaml
---
images:
- name: docker.io/library/alpine 
  tags:
  - 3.9.6
  - 3.9.7
  - 3.18.3
  - 3
- name: docker.io/library/ubuntu
  tags:
  - 22.04
  - focal-20231003
  #- name: c1pltdevopsbase  
  #  custombuild: true
  #  dockerfile:|
  #    FROM 
  #  tags:
  #  - 0.1
  #  - latest
custombuild:
- name: c1pltdevopsbase.acurecr.io/mcr.microsoft.com/azurefunction/nodes
  tags:
  - 3
  - 3.1
  dockerfile:
    from: mcr.microsoft.com/azurefunction/nodes
    tag: 3.1 
    inline: |
    #dockerfile path except FROM

```

## OUTPUT YAML file
We will write only new SHA256 and tags
---
date: 2023-10-20
images:
- name: c1pltdevopsbase.azurecr.io/docker.io/library/alpine
  tag: 3.9.6
  SHA256: sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f
- name: c1pltdevopsbase.acurecr.io/mcr.microsoft.com/azurefunction/nodes
  SHA256: sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f
  tag: 3
---
