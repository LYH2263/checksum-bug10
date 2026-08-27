
# go-checksum

流式多算法校验聚合：CRC32 + SHA256 rolling manifest。分块由固定窗口 splitter 驱动，manifest builder 汇总。

```text
go build ./...
go test ./... -count=1
go run ./cmd/sumd -addr :8225
```
