1. Đọc hiểu API cluster resource set
2. Định nghĩa các tham số của CRD Helm Chart
3. Định nghĩa các công việc của Helm Chart Controller
4. Implement helm chart CRD và Helm chart controller

# 1. Đọc hiểu API ClusterResourceSet

## 1.1 Proposal - Ý Tưởng

- Data model changes to existing API types
- ClusterResourceSet Object Definition: là CRD có set of components (resources) được apply vào cluster match the label selector.

- Sample ClusterResourceSet YAML
    ```
        ---
        apiVersion: addons.cluster.x-k8s.io/v1alpha3
        kind: ClusterResourceSet
        metadata:
        name: crs1
        namespace: default
        spec:
        mode: "ApplyOnce"
        clusterSelector:
        matchLabels:
            cni: calico
        resources:
        - name: db-secret
            kind: Secret
        - name: calico-addon
            kind: ConfigMap
    ``` 
    - Resources field: list of Secrets/ConfigMaps.
    - clusterSelector field: a K8s label selector.
    - CRS: namespace-scoped.
- CRSBinding Object Definition: lưu thông tin những resource nào được áp dụng vào cluster nào. Nó được tạo trong management cluster. Chỉ có 1 CRSB cho 1 workload cluster. CRB phải cùng tên với cluster. 
    - Khi cluster bị xóa, CRB cũng bị xóa theo. 
    - Khi CRS bị xóa, nó sẽ bị xóa trong trường `bindings`.
    ```
            apiVersion: v1
        kind: ClusterResourceBinding
        metadata:
        name: <cluster-name>
        namespace: <cluster-namespace>
        ownerReferences:
        - apiVersion: cluster.x-k8s.io/v1alpha3
            kind: Cluster
            name: <cluster-name>
            uid: e3a503a8-9be1-4264-8fa2-d536532687f9
        - apiVersion: addons.cluster.x-k8s.io/v1alpha3
            blockOwnerDeletion: true
            controller: true
            kind: ClusterResourceSet
            name: crs1
            uid: 62c77639-92d8-46d2-ba21-a880f62f7719
        spec:
        bindings:
        - clusterResourceSetName: crs1
            resources:
            - applied: true
            hash: sha256:a3473f4e92ee5a2277ff37d5c559666d61d24332a497b554e65ae18e82727245
            kind: Secret
            lastAppliedTime: "2020-07-02T05:47:38Z"
            name: db-secret
            - applied: true
            hash: sha256:c1d0dc7e51bb05945a2f99e6745dc4b1043f8a03f37ad21391fe92353a02066e
            kind: ConfigMap
            lastAppliedTime: "2020-07-02T05:47:39Z"
            name: calico-addon
    ```
- CRS sẽ sử dụng CRSB để quyết định có apply 1 resource mới hay thử lại resource cũ. Khi apply fail, CRS controller sẽ reconcile và sử dụng `controller-runtime` back-off để thử apply lại. Khi ta thêm một resource mới vào CRS, CRS đó sẽ được reconcile ngay lập tức và resource mới sẽ được apply vào tất cả các matching cluster vì resource mới kh tồn tại trong danh sách CRB. 
https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20200220-cluster-resource-set.md#proposal
