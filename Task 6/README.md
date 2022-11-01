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

# 2. Building CronJob
## 2.1 Chuẩn bị
- Cài môi trường Golang
- Cài Docker
- Cài kubectl
- Cài kubebuilder
- Tạo kết nối tới 1 cluster K8s
## 2.2 Thực hiện
### Step 1: Tạo folder Project
```
    mkdir project
    cd project
    kubebuilder init --domain tutorial.kubebuilder.io --repo tutorial.kubebuilder.io/project
```
### Step 2: Tạo CRD mới (Adding a New API)
```
    kubebuilder create api --group batch --version v1 --kind CronJob
``` 
- Tiến hành hoàn thiện file cronjob_types.go (nơi định nghĩa model của object)
- Cần quan tâm đến CronJobSpec và CronJobStatus
### Step 3: Tạo logic Controller
- Hoàn thiện file   `cronjob_controller.go`

- Tập trung vào 2 func() chính:
    - Reconcile(): là hàm sẽ chạy khi có message trong queue mà controller đã đăng ký. Nó giống như một hàm callback khi có sự kiện xảy ra với object mà ta đã đăng ký ở hàm SetupWithManager() ở dưới
    - SetupWithManager(): là hàm đăng ký manager, có chức năng định nghĩa ra các object mà controller sẽ theo dõi sự thay đổi. Mỗi controller có thể watch một object chính và các object liên quan. Nhưng khi trigger sự kiện thì input của hệ thống sẽ chỉ trả về tên object chính và namespace của nó vì thế phải dùng rất cẩn thận, tránh việc theo dõi các loại object không cần thiết, hoặc care quá nhiều object trong một controller.

### Step 4: Triển khai.
- Project được gen ra đã chứa một file Makefile chứa rất nhiều lựa chọn để từng bước build, gen manifest đến triển khai trên K8S.
- Trong đó:

    - make manifests: sẽ gen tất cả các manifest của CRD mà mình đã định nghĩa trong các model
    - make install:  apply các file manifest do quá trình make manifests tạo ra
    - make docker-build: build controller thành docker image
    - make deploy: triển khai controller lên K8S
    ...

- Ví dụ:

    - Khi debug chương trình:  make manifests  ==> make install  ==> make run
    - Khi release lên k8s: make manifests  ==> make install ==> make docker-build ==> make deploy

# Ref
[1] [CronJob](https://book.kubebuilder.io/cronjob-tutorial/new-api.html)