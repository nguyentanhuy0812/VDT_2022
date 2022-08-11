# 1. Cài đặt Cinder csi.
- Tạo các file cài đặt:

    ```
    cinder-csi-controllerplugin-rbac.yaml
    cinder-csi-controllerplugin.yaml
    cinder-csi-nodeplugin-rbac.yaml
    cinder-csi-nodeplugin.yaml
    csi-cinder-driver.yaml
    csi-secret-cinderplugin.yaml
    ```
- Lưu ý sửa phần `cloud.conf` trong file `csi-secret-cinderplugin.yaml`.

- Dùng lệnh `kubectl apply -f` để áp dụng các file đã tạo.
# 2. Test 