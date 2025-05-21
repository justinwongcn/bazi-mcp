# 客户端配置参考

```json
{
   "mcpServers": {
      "bazi": {
         "command": "your path",
         "env": {
            "API_KEY": "your key"
         }
      }
   }
}
```

## API 地址

[缘分居](https://doc.yuanfenju.com)

## 构建方式

### 使用构建脚本

项目提供两种构建脚本：

1. `build.sh` - Linux/macOS系统使用，支持构建多个平台和架构的版本：
   - Mac Intel (amd64)
   - Mac Apple Silicon (arm64)
   - Windows 64位 (amd64)
   - Windows ARM64
   - Windows 32位 (386)
   - Linux amd64
   - Linux ARM64
   - Linux ARMv7
   - Linux 32位 (386)

   使用方法：
   ```bash
   ./build.sh
   ```

2. `build.bat` - Windows系统使用，功能与build.sh类似，支持相同的平台和架构组合。

   使用方法：
   ```bat
   build.bat
   ```

执行任一构建脚本后，所有构建产物都会输出到`dist/`目录下。

## 注意事项

在 `/Users/john/GolandProjects/bazi_mcp/cmd/bazi-mcp/main.go` 文件的第 83 到 117 行，注册资源的部分存在问题，无法根据省份查询城市。如有相关 example 可提供一个从 JSON 文件注册资源的 example 以供参考。
