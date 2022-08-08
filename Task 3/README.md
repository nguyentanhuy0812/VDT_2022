# I. Note từ Task 2.
## i. Làm rõ `staging path` và `target path`
- Phương thức `NodeStageVolume` thực hiện việc mount tạm thời volume vào `staging path`. Thường thường staging path này là một `global directory` trong một node. `Global directory` là một thiết bị khối được mount vào hệ thống Linux. Trong K8s, một PV có thể được gắn vào nhiều phiên bản pod trong một node. Một thiết bị khối đã được định dạng được gắn vào `global directory` trong một node. 

Ví dụ: global directory : ``/var/lib/kubelet/pods/[pod uid]/volumes/kubernetes.io~iscsi/[PV
name]``
- Phương thức `NodePublishVolume` thực hiện việc mount volume từ `staging path` tới `target path`. Trong K8s, `target path` là `pod directory`  
# II. Task 3
## 1. Custom Resources
- `Resource` là điểm cuối trong `Kubernetes API` lưu trữ một tập các đối tượng API của một loại nhất định, ví dụ như resource tích hợp pods chứa một tập các đối tượng Pod.
- Trong K8s có 2 loại resource, một loại là resource có sẵn khi cài đặt K8s gọi là native còn 1 loại resource được tạo mới bằng kỹ thuật CRD(CustomResourceDefinition).
- `Custom resource` là một extension của Kubernetes API mà không nhất thiết phải có sẵn trong bộ cài đặt Kubernetes. 
- `Custom resource` có thể xuất hiện và biến mất trong quá trình chạy cluster thông qua `dynamic registration`.
- Khi một `custom resource` được cài đặt, người dùng có thể tạo và truy cập các đối tượng của nó bằng cách sử dụng `kubectl`, giống như cách họ làm đối với các tài nguyên tích hợp sẵn như Pod.
- Bên cạnh những resource mặc định như Pod, Deployment, ReplicaSet, ... thì K8s cho phép ta tạo thêm những custom resource để đáp ứng nhu cầu của ta trong dự án.
- Từng custom resource sẽ phục vụ cho 1 mục đích cụ thể. 
## 1.1. CRD là gì?
- Viết tắt của: Custom Resource Definition
- `Definition`: các lệnh khai báo với máy chủ API ở dạng cấu trúc YAML.
- `CRD` là tính năng trong K8s cho phép người dùng có thế thêm các đối tượng (objects) vào K8s cluster và sử dụng nó như các đối tượng K8s khác.
- CRD API resource cho phép định nghĩa custom resource. Việc xác định CRD object tạo custom resource mới với tên chỉ định. K8s API phục vụ và xử lý việc lưu trữ custom resource.
- Các hoạt động CRD được xử lý bên trong kube-apiserver bởi module apiextensions-apiserver (module được tích hợp vào kube-apiserver)
## 1.2. Custom Controllers
- Để mang lại nhiều chắc năng nâng cao hơn cho custom resource thì cần triển khai bộ điều khiển: `Custom controller`.
- `Custom resource` cho phép bạn lưu trữ và truy xuất dữ liệu được cấu trúc. Khi kết hợp một `custom resource` với một `custom controller`, custom resource cung cấp một "declarative API".
- "Declarative API" K8s thực thi việc phân quyền. Bạn khai báo trạng thái mong muốn của resource. K8s controller giữ trạng thái hiện tại của K8s objects đồng bộ với trạng thái mong muốn đó.
- Trong K8s, Controller Manager có nhiệm vụ theo dõi API Server và tạo ra các resource liên quan tới nó. Thì bên cạnh những Controller Manager có sẵn bên trong K8s thì ta cần tạo thêm custom controller để phục vụ cho một mục đích khác nào đó.
- Để tạo một custom controller, đầu tiên ta cần viết code theo dõi API Server với resource ta muốn. Sau đó build code thành image. Sau đó tạo Deployment sử dụng image vừa tạo rồi deploy nó lên K8s. 
- Thực chất, customer controller chỉ là một Deployment bình thường, khác ở chỗ ta sẽ tự viết code để tương tác với API server.

### Cách viết Custom Controller.

## 1.3. Thêm custom resources
- K8s cung cấp 2 cách để thêm custom resource vào cluster:
    - CRDs : không yêu cầu lập trình. Dễ sử dụng hơn. 
    - API Aggregation (tổng hợp API): yêu cầu lập trình, nhưng cho phép kiểm soát nhiều hơn hành vi API như: cách dữ liệu được lưu trữ và chuyển đổi giữa các phiên bản API. Linh hoạt hơn.
- API tổng hợp là các máy chủ API cấp dưới nằm phía sau máy chủ API chính, hoạt động giống như proxy. Sự sắp xếp này được gọi là API Aggregation.
- CRD cho phép người dùng tạo các tài nguyên mới mà không cần thêm máy chủ API khác.
- Các resource mới đều đều được gọi là custom resource để phân biệt với các resources tích hợp sẵn của K8s.
- Custom resources sau khi tạo sẽ được lưu trữ ở etcd trong master node.
## Demo mở rộng K8s API bằng CRDs.
- Khi ta tạo một API CustomResourceDefinition mới, máy chủ Kubernetes API sẽ tạo một đường dẫn tài nguyên RESTful mới cho từng phiên bản được chỉ định. CRD có thể là không gian tên hoặc phạm vi cụm, như được chỉ định trong trường phạm vi của CRD.
- Create a CustomResourceDefinition.

    - Tạo file `resourcedefinition.yaml`

    - Dùng câu lệnh để tạo : 

        ```
        kubectl apply -f resourcedefinition.yaml
        ```
- Create custom objects

    - Tạo file `my-crontab.yaml`

    - Dùng câu lệnh để tạo:
        ```
        kubectl apply -f my-crontab.yaml
        ```

# 2. Kubebuilder

- Kubebuilder là một framework để xây dựng K8s API sử dụng CRDs.

# 3. Cluster API: Kiến trúc, mô hình mạng, các loại controller: core, boostrap, control-plane, infrastructure và chức năng của từng cái, sơ đồ tương tác giữa các thành phần, future work
- Cluster API là một dự án con của Kubernetes tập trung vào việc cung cấp các declarative APIs và công cụ để đơn giản hóa việc cung cấp, nâng cấp và vận hành nhiều cụm Kubernetes.
- 
## 3.1 Kiến trúc
## 3.2 Mô hình mạng
## 3.3 Các loại Controller
## 3.4 Sơ đồ tương tác
## 3.5 Future work.

# 4. Gardener architecture





- Ref:
    - [1] [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#:~:text=The%20CustomResourceDefinition%20API%20resource%20allows,a%20valid%20DNS%20subdomain%20name.) 
    - [2] [Kubebuilder] (https://github.com/kubernetes-sigs/kubebuilder)
    - [3] [Custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
    - [4] [Custom Controller](https://github.com/kubernetes/sample-controller)
    - [5] [Extend K8s API with CRD](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)

