# CRI hook plugin插入功能说明

## setRootFSWritableLayerHomeDir

在单Pod多容器下，为了能够将所有存储统一到一个Volume上，从而可以设置一个Pod Level的存储空间限制,
这里将一个Pod的所有容器的rootfs writable layer都指定到Pod使用的一个Persistent Volume在宿主机
上的挂载路径下。然后通过限制Pod的Persistent Volume Size来达到限制一个Pod总存储quota的目的。

如下Pod Yaml所示，container1,2,3都使用CSI Persistent Volume作为持久化存储，并通过annotation声明将所有容器的rootfs writable layer都放置到csi-yundisk中

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: demo-pod
  annotations:
    # Pod Level: 该pod所有容器rootfs读写层都会放在所指定的volume上
    alibabacloud.com/rootfs-writable-layer: csi-yundisk
spec:
  containers:
  # 以下3个容器共享同一个PVC binding的盘，且每个容器的rootfs writable layer
  # 也都放在该盘上，对该盘的disk quota限制或隔离将作用在整个pod上
  - name: container-1
    image: reg.docker.alibaba-inc.com/ali/os:7u2
    volumeMounts:
    - name: csi-yundisk
      mountPath: /home/admin/logs
      subPath: home/admin/logs
  - name: container-2
    image: reg.docker.alibaba-inc.com/ali/os:7u2
    volumeMounts:
    - name: csi-yundisk
      mountPath: /home/admin/logs
      subPath: home/admin/logs
  - name: container-3
    image: reg.docker.alibaba-inc.com/ali/os:7u2
    volumeMounts:
    - name: csi-yundisk
      mountPath: /home/admin/logs
      subPath: home/admin/logs
  volumes:
  # 本pod所有容器共享的盘的volume
  - name: csi-yundisk
    persistentVolumeClaim:
      claimName: csi-yundisk-pvc
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: csi-yundisk-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 60Gi
  # 通过csi plugin动态provision一块盘，CSI plugin会保证这块盘挂
  # 载到/mnt/<backend>/persistent-volume-name这个mountpoint下
  storageClassName: csi-plugin
```

在kube-apiserver的admission webhook中会将alibabacloud.com/rootfs-writable-layer的值修改为volume csi-yundisk在pod所在node上约定的挂载路径+".rootDir"(也即/mnt/<storage-backend>/volume-name/.rootDir)，volume在node上约定的挂载路径(/mnt/<storage-backend>/volume-name)有CSI plugin来保证。
