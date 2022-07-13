# Các Thuật Ngữ
- Volume: Một đơn vị lưu trữ sẽ được cung cấp bên trong thùng chứa được quản lý bằng CO, thông qua CSI.
- CO (Container Orchestration): Hệ thống điều phối vùng chứa, giao tiếp với các Plugin sử dụng các RPC dịch vụ CSI.
- SP(Storage Provider): Nhà cung cấp bộ nhớ, nhà cung cấp triển khai plugin CSI.
- RPC: Gọi thủ tục từ xa. Là cơ chế giao tiếp giữa 2 tiến trình. Thực hiện lời gọi thủ tục trên tiến trình khác giống như lời gọi thủ tục trong một tiến trình cục bộ. 
- Node: Máy chủ lưu trữ nơi khối lượng công việc của người dùng sẽ chạy, có thể nhận dạng duy nhất từ ​​quan điểm của một Plugin bằng ID nút.
- Plugin: một điểm cuối gRPC triển khai Dịch vụ CSI.
- gRPC : một RPC framework giúp bạn kết nối giữa các service trong hệ thống, nó hỗ trợ load balancing, tracing, health checking và authentication hỗ trợ từ mobile, trình duyệt cho tới back-end service.
- Workload: Đơn vị của "công việc" được lên lịch bởi một CO. Đây CÓ THỂ là một thùng chứa hoặc một tập hợp các thùng chứa.
# 1. Container Storage Interface (CSI) là gì?
- Giao diện lưu trữ vùng chứa (CSI) được phát triển như một tiêu chuẩn để hiển thị các hệ thống lưu trữ khối và lưu trữ tệp tùy theo khối lượng công việc được chứa trên các hệ thống điều phối vùng chứa (CO) như Kubernetes, Mesos, Nomad, Cloud Foundry.
- CSI về cơ bản là một giao diện (interface) giữa khối lượng công việc vùng chứa (container workloads) và bộ lưu trữ của bên thứ ba hỗ trợ việc `tạo` và `cấu hình` bộ lưu trữ liên tục bên ngoài bộ điều phối, đầu vào/đầu ra (I / O) của nó và các chức năng nâng cao như snapshots và cloning.
- CSI thay thế các plugin được phát triển trước đó trong quá trình phát triển Kubernetes, chẳng hạn như plugin khối lượng trong cây và plugin FlexVolume. CSI là Kubernetes API framework. CSI cho phép Kubernetes hỗ trợ linh hoạt các ứng dụng yêu cầu bộ nhớ liên tục. Nhà cung cấp dịch vụ lưu trữ sẽ cần viết CSI drivers áp dụng cho các thông số kỹ thuật của CSI framework và có thể lựa chọn các chức năng mà họ có thể và sẽ hỗ trợ. CSI là nền tảng của lưu trữ Kubernetes linh hoạt và có thể quản lý. 
- Tổng quan: 
    - CSI là về các hoạt động lưu trữ trong môi trường gốc của vùng chứa.
    - CSI là một đặc điểm kỹ thuật API để quản lý lưu trữ container.
    - CSI cũng là một giao diện vùng chứa tiêu chuẩn công nghiệp
    - Vì vậy, CSI là một API và phương pháp luận chung để giải quyết vấn đề lưu trữ như một tiêu chuẩn công nghiệp
    Nó là một lớp trừu tượng để quản lý khối lượng lưu trữ mà không phải lo lắng về các nội dung hoặc chi tiết cụ thể của bộ nhớ cơ bản phức tạp.

    <img src= "./imgs/csi.png">

- CSI cho phép các nhà cung cấp dịch vụ lưu trữ:
    - Tự động tạo lưu trữ khi có yêu cầu.
    - Cung cấp bộ nhớ cho các vùng chứa ở bất cứ nơi nào chúng được lên lịch.
    - Tự động xóa bộ nhớ khi không còn cần thiết.  

# 2. Kiến trúc của CSI?
- CSI được chia thành một số kiến trúc thiết kế xác định cách triển khai các plugin như thế nào. Các kiến ​​trúc phù hợp với các triển khai CO điển hình có máy chủ chính (master host) và máy chủ nút (node host). Có 3 tình huống xảy ra như sau:
    - Master/Node với các plugin riêng rẽ cho cả bộ điều khiển và chức năng node.
     <img src="./imgs/scenario1.png">
    - "Headless" trong đó các plugin vẫn chạy riêng rẽ cho bộ điều khiển và node nhưng chỉ chạy trên máy chủ node.
    <img src="./imgs/scenario2.png">
    - Headless kết hợp, trong đó một plugin cung cấp bộ điều khiển và các khả năng node cùng nhau.
# 3. Container Storage Interface
*Phần này chủ yếu mô tả giao diện giữa các CO và các Plugin*
- RPC Interface:
    - Một CO tương tác với một Plugin thông qua RPCs. Mỗi SP phải cung cấp:
        - Node Plugin: Điểm cuối gRPC phục vụ CSI RPCs phải được chạy trên Node, sau đó một volume do SP cung cấp sẽ được xuất bản. 
        - Controller Plugin: Một điểm cuối gRPC phục vụ CSI RPCs có thể được chạy ở bất cứ đâu.
        - Trong một số trường hợp, một điểm cuối gRPC có thể phục vụ tất cả CSI RPCs.
- Bên trong mỗi plugin có 3 dịch vụ xử lý các cuộc gọi RPC: 
    - Dịch vụ nhận dạng (Identity service) - được thực hiện bởi cả Plugin điều khiển và Plugin Node.
    - Dịch vụ điều khiển (Controller service) - chạy trên Plugin điều khiển.
    - Dịch vụ Node (Node service) - được thực hiện bởi Plugin Node.

```
service Identity {
  rpc GetPluginInfo(GetPluginInfoRequest)
    returns (GetPluginInfoResponse) {}

  rpc GetPluginCapabilities(GetPluginCapabilitiesRequest)
    returns (GetPluginCapabilitiesResponse) {}

  rpc Probe (ProbeRequest)
    returns (ProbeResponse) {}
}

service Controller {
  rpc CreateVolume (CreateVolumeRequest)
    returns (CreateVolumeResponse) {}

  rpc DeleteVolume (DeleteVolumeRequest)
    returns (DeleteVolumeResponse) {}

  rpc ControllerPublishVolume (ControllerPublishVolumeRequest)
    returns (ControllerPublishVolumeResponse) {}

  rpc ControllerUnpublishVolume (ControllerUnpublishVolumeRequest)
    returns (ControllerUnpublishVolumeResponse) {}

  rpc ValidateVolumeCapabilities (ValidateVolumeCapabilitiesRequest)
    returns (ValidateVolumeCapabilitiesResponse) {}

  rpc ListVolumes (ListVolumesRequest)
    returns (ListVolumesResponse) {}

  rpc GetCapacity (GetCapacityRequest)
    returns (GetCapacityResponse) {}

  rpc ControllerGetCapabilities (ControllerGetCapabilitiesRequest)
    returns (ControllerGetCapabilitiesResponse) {}

  rpc CreateSnapshot (CreateSnapshotRequest)
    returns (CreateSnapshotResponse) {}

  rpc DeleteSnapshot (DeleteSnapshotRequest)
    returns (DeleteSnapshotResponse) {}

  rpc ListSnapshots (ListSnapshotsRequest)
    returns (ListSnapshotsResponse) {}

  rpc ControllerExpandVolume (ControllerExpandVolumeRequest)
    returns (ControllerExpandVolumeResponse) {}

  rpc ControllerGetVolume (ControllerGetVolumeRequest)
    returns (ControllerGetVolumeResponse) {
        option (alpha_method) = true;
    }
}

service Node {
  rpc NodeStageVolume (NodeStageVolumeRequest)
    returns (NodeStageVolumeResponse) {}

  rpc NodeUnstageVolume (NodeUnstageVolumeRequest)
    returns (NodeUnstageVolumeResponse) {}

  rpc NodePublishVolume (NodePublishVolumeRequest)
    returns (NodePublishVolumeResponse) {}

  rpc NodeUnpublishVolume (NodeUnpublishVolumeRequest)
    returns (NodeUnpublishVolumeResponse) {}

  rpc NodeGetVolumeStats (NodeGetVolumeStatsRequest)
    returns (NodeGetVolumeStatsResponse) {}


  rpc NodeExpandVolume(NodeExpandVolumeRequest)
    returns (NodeExpandVolumeResponse) {}


  rpc NodeGetCapabilities (NodeGetCapabilitiesRequest)
    returns (NodeGetCapabilitiesResponse) {}

  rpc NodeGetInfo (NodeGetInfoRequest)
    returns (NodeGetInfoResponse) {}
}
```
## 3.1 Controller Service RPC
### 1. CreateVolume
- Một Plugin điều khiển phải thực hiện lệnh gọi RPC nếu nó có khả năng điều khiển CREATE_DELETE_VOLUME. RPC này sẽ được CO gọi để cung cấp một volume mới thay mặt cho người dùng. 
- Plugin có thể tạo 3 loại volume:
    - Volume trống. Khi plugin hỗ trợ khả năng tùy chọn CREATE_DELETE_VOLUME.
    - Từ một snapshot hiện có. Khi plugin hỗ trợ khả năng tùy chọn CREATE_DELETE_VOLUME và CREATE_DELETE_SNAPSHOT.
    - Từ một volume hiện có. Khi plugin hỗ trợ clone và báo cáo khả năng tùy chọn CREATE_DELETE_VOLUME và CLONE_VOLUME.
- Nếu như CO yêu cầu một volume được tạo ra từ snapshot hay volume hiện có và kích thước của volume yêu cầu lớn hơn bản gốc thì Plugin có thể từ chối cuộc gọi với lỗi OUT_OF_RANGE hoặc phải cung cấp một volume, khi được trình bày với workload bởi cuộc gọi NodePublish, có cả kích thước yêu cầu và chứa dữ liệu từ snapshot.
- Nếu như plugin không thể hoàn thành CreateVolume, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. Các mã lỗi gRPC sẽ được gửi tương ứng với từng trường hợp và CO phải thực hiện các hành vi khôi phục lỗi. 
### 2. ControllerPublishVolume
- Một plugin điều khiển phải thực hiện cuộ gọi RPC nếu nó có khả năng điều khiển PUBLISH_UNPUBLISH_VOLUME. RPC sẽ được gọi bởi CO khi nó muốn đặt 1 workload sử dụng volume trong một node. Plugin NÊN thực hiện công việc cần thiết để làm cho volume có sẵn trên nút nhất định. Plugin KHÔNG PHẢI cho rằng RPC này sẽ được thực thi trên nút nơi ổ đĩa sẽ được sử dụng. 
- Nếu hoạt động lỗi hoặc CO không biết liệu hoạt động có lỗi hay không, nó có thể chọn gọi ControllerPublishVolume lần nữa hoặc chọn gọi ControllerUnpublishVolume.
- Nếu như plugin không thể hoàn thành ControllerPublishVolume, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. Các mã lỗi gRPC sẽ được gửi tương ứng với từng trường hợp và CO phải thực hiện các hành vi khôi phục lỗi. 
### 3. ControllerUnpublishVolume
- Plugin điều khiển phải thực hiện cuộc gọi RPC nếu nó có khả năng điều khiển PUBLISH_UNPUBLISH_VOLUME. RPC này là hoạt động ngược lại của ControllerPublishVolume. Nó phải được gọi sau NodeUnstageVolume và NodeUnpublishVolume trong volume được gọi và thành công. Plugin NÊN thực hiện công việc cần thiết để làm cho volume sẵn sàng được sử dụng bởi một nút khác. Plugin KHÔNG PHẢI giả định rằng RPC này sẽ được thực thi trên nút mà volume đã được sử dụng trước đó.
- RPC này thường được gọi bởi CO khi khối lượng công việc sử dụng volume đang được chuyển đến một nút khác hoặc tất cả khối lượng công việc sử dụng volume trên một nút đã kết thúc.
- Nếu như plugin không thể hoàn thành ControllerUnpublishVolume, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. Các mã lỗi gRPC sẽ được gửi tương ứng với từng trường hợp và CO phải thực hiện các hành vi khôi phục lỗi. 
### 4. ValidateVolumeCapabilities.
- Một Plugin điều khiển phải thực hiện cuộc gọi RPC. RPC này được gọi bởi CO để kiểm tra xem một volume được cấp phép trước có tất cả các khả năng mà CO muốn hay không. Cuộc gọi RPC nên trả về `confirmed` chỉ nếu tất cả các khả năng được chỉ định được hỗ trợ.  Hành động này phải là `idempotent`.
- Nếu như plugin không thể hoàn thành ValidateVolumeCapabilities, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. Các mã lỗi gRPC sẽ được gửi tương ứng với từng trường hợp và CO phải thực hiện các hành vi khôi phục lỗi. 
### 5. ListVolumes
- Plugin điều khiển phải thực hiện cuộc gọi RPC nếu nó có khả năng điều khiển LIST_VOLUMES. Plugin nên trả về thông tin về thất cả các volume mà nó biết. Nếu các volume được tạo hoặc bị xóa trong khi CO đang đồng thời phân trang thông qua các kết quả ListVolumes thì CO CÓ THỂ thấy ​​các tập trùng lặp trong danh sách, không chứng kiến ​​các tập hiện có hoặc cả hai.
- Nếu như plugin không thể hoàn thành ListVolumes, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. Các mã lỗi gRPC sẽ được gửi tương ứng với từng trường hợp và CO phải thực hiện các hành vi khôi phục lỗi. 
### 6. GetCapacity
- Plugin điều khiển phải thực hiện cuộc gọi RPC nếu nó có khả năng điều khiển GET_CAPACITY. RPC cho phép CO truy vẫn khả năng lưu trữ mà bộ điều khiển cấp volumes.
- Nếu như plugin không thể hoàn thành GetCapacity, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC
### 7. ControllerGetCapabilities
- Plugin trình điều khiển PHẢI thực hiện lệnh gọi RPC này. RPC này cho phép CO kiểm tra các khả năng được hỗ trợ của dịch vụ bộ điều khiển do Plugin cung cấp.
### 8. CreateSnapshot
- Plugin điều khiển phải thực hiện cuộc gọi RPC nếu nó có khả năng điều khiển CREATE_DELETE_SNAPSHOT. RPC này sẽ được gọi bởi CO để tạo một snapshot mới từ một volume nguồn thay mặt cho người dùng.
- Hành động này phải được `idempotent`. 
- Nếu lỗi xảy ra trước khi snapshot, CreateSnapshot NÊN trả về mã lỗi gRPC tương ứng phản ánh tình trạng lỗi.
- Snapshot CÓ THỂ được sử dụng làm nguồn để cung cấp volume mới. Thông báo CreateVolumeRequest CÓ THỂ chỉ định một tham số snapshot nguồn TÙY CHỌN. Hoàn nguyên snapshot, trong đó dữ liệu trong volume gốc bị xóa và được thay thế bằng dữ liệu trong snapshot, là một chức năng nâng cao không phải hệ thống lưu trữ nào cũng có thể hỗ trợ và do đó hiện đang nằm ngoài phạm vi.
## 3.2 Node Service RPC
### 1. NodeStageVolume.
- Một Plugin Node phải thực hiện cuộc gọi RPC nếu nó có khả năng Node STAGE_UNSTAGE_VOLUME.
- RPC này được gọi bởi CO trước khi volume được sử dụng bởi bất kỳ khối lượng công việc nào trên node bởi NodePublishVolume. RPC NÊN được gọi bởi CO khi một khối lượng công việc muốn sử dụng volume được chỉ định được đặt (đã lên lịch) trên node được chỉ định lần đầu tiên hoặc lần đầu tiên kể từ khi lệnh gọi NodeUnstageVolume cho volume được chỉ định được gọi và trả về thành công trên node đó.
- nếu RPC này thất bại, hoặc CO không biết liệu thất bại hay không, nó có thể chọn gọi lại NodeStageVolume hoặc chọn gọi NodeUnstageVolume.
- Nếu như plugin không thể hoàn thành NodeStageVolume, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC.
### 2. NodeUnstageVolume
- Một Plugin Node phải thực hiện cuộc gọi RPC nếu nó có khả năng Node STAGE_UNSTAGE_VOLUME.

### 3. NodePublishVolume
### 4. NodeUnpublishVolume
### 5. NodeGetVolumeStats
### 6. NodeGetCapabilities
### 7. NodeGetInfo
### 8. NodeExpandVolume
## 3.3 Identity Service RPC
`Identity Service RPC` cho phép CO truy vấn một plugin về các khả năng, tình trạng và siêu dữ liệu khác. Quy trình chung của trường hợp thành công CÓ THỂ như sau:
- 1. CO truy vấn siêu dữ liệu qua Identity RPC.
- 2. CO truy vấn khả năng có sẵn của plugin.
- 3. CO truy vấn mức độ sẵn sàng của plugin
### 1. GetPluginInfo
- Nếu như plugin không thể hoàn thành GetPluginInfo, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. 
### 2. GetPluginCapabilities
- RPC BẮT BUỘC này cho phép CO truy vấn các khả năng được hỗ trợ của Plugin "nói chung": nó là tổng thể của tất cả các khả năng của tất cả các phiên bản của phần mềm Plugin, vì nó được dự định triển khai.
-  Tất cả các bản sao của cùng một phiên bản của Plugin nên trả về cùng một tập hợp các khả năng, bất kể cả hai: (a) trong đó các bản sao cũng được triển khai trên cluster; (b) RPC mà một cá thể đang phân phát.
- Nếu như plugin không thể hoàn thành GetPluginCâpbilities, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC. 
### 3. Probe(ProbeRequest)
- Một Plugin phải thực hiện cuộc gọi RPC này. Tiện ích chính của Probe RPC là xác minh rằng plugin đang ở trạng thái khỏe mạnh và sẵn sàng.  Nếu một trạng thái không "khỏe mạnh" được báo cáo, thông qua một phản hồi không thành công, CO có thể thực hiện hành động với mục đích đưa plugin về trạng thái khỏe mạnh, có thể là:
    - Khởi động lại plugin container.
    - Thông báo cho người giám sát plugin.
- Plugin CÓ THỂ xác minh rằng nó có cấu hình, thiết bị, phụ thuộc và trình điều khiển phù hợp để chạy và trả về thành công nếu xác thực thành công. CO có thể gọi RPC này bắt cứ khi nào. 
- Nếu như plugin không thể hoàn thành Probe, nó phải trả về mã non-ok gRPCtrong trạng thái gRPC.
# 4. Mục đích của CSI trong MVP (Minimum Viable Product).
Để xác định một tiêu chuẩn công nghiệp “Giao diện lưu trữ vùng chứa” (CSI) sẽ cho phép các nhà cung cấp dịch vụ lưu trữ (SP) phát triển một plugin một lần và để nó hoạt động trên một số hệ thống điều phối vùng chứa (CO).
- Cho phép các tác giả SP viết một Plugin tuân thủ CSI “chỉ hoạt động” trên tất cả các CO triển khai CSI.
- Xác định API (RPC) cho phép: 
    - Cấp phép động và hủy cấp phép một volume.
    - Đính kèm hoặc tách một volume từ một node.
    - Gắn hoặc ngắt kết nối một volume từ một node.
    - Nhà cung cấp bộ nhớ cục bộ. 
    - Tạo và xóa một snapshot ( nguồn của snapshot là một volume).
    - Cung cấp một volume mới từ một snapshot.
- Xác định giao thức plugin RECOMMENDATIONS:
    - Mô tả một quy trình người giám sát cấu hình một Plugin.
    - Cân nhắc triển khai vùng chứa. 
# 5. CSI trong Kubernetes
