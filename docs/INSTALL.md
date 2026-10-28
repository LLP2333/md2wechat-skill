# 安装指南

md2wechat 社区版只提供一种安装方式：npm。

当前安装说明对应 v3.9.1。

## 目录

- [系统要求](#系统要求)
- [可选的浏览器要求](#可选的浏览器要求)
- [安装](#安装)
- [验证安装](#验证安装)
- [升级](#升级)
- [卸载](#卸载)

---

## 系统要求

- **操作系统**：macOS（`amd64` / `arm64`）、Linux（`amd64` / `arm64`）、Windows（`amd64`）
- **Node 版本**：18 或更高
- **网络**：安装时需要访问 npm 和 GitHub Release；默认本地排版不需要联网

## 可选的浏览器要求

`sync prepare` 是纯本地命令，不需要浏览器。知乎、CSDN、头条、腾讯云开发者社区的草稿操作由宿主 Agent 复用用户已登录的浏览器完成；宿主需提供页面读取、正常编辑、HTML 粘贴与本地文件上传能力。

CLI 不启动或安装浏览器，不管理登录资料，也不要求固定扩展。宿主能力不足时，按内置流程检查可选工具；仍缺少必要能力时交接给用户。完整流程和备选工具验证范围见 [多平台草稿](SYNC.md)。

---

## 安装

```bash
npm install -g @qvqw2233/wechat-md@3.9.1
```

安装后的命令名是 `md2wechat`。

这个包不会在安装时编译源码。它会从 [LLP2333/md2wechat-skill](https://github.com/LLP2333/md2wechat-skill/releases) 下载与包版本同号的 GitHub Release 二进制，并用同一个 Release 里的 `checksums.txt` 校验。

如果你的 npm 默认指向 `https://registry.npmmirror.com`，而新版本刚发布，镜像 tarball 可能会短暂 `404`。先用官方源安装：

```bash
npm install -g @qvqw2233/wechat-md@3.9.1 --registry=https://registry.npmjs.org/
```

维护者在 npm 发布新版本后，应手动执行一次：

```bash
npx cnpm sync @qvqw2233/wechat-md
```

这样可以尽快把新版本同步到 `npmmirror`。

---

## 验证安装

```bash
md2wechat version --json
md2wechat skills read md2wechat --json
```

如果提示找不到命令，说明 npm 全局 bin 目录不在 `PATH` 中。用下面的命令查看目录，加入 `PATH` 后重新打开终端：

```bash
npm prefix -g
```

macOS / Linux 的命令在该目录下的 `bin/` 里，Windows 的命令就在该目录里。

---

## 升级

```bash
npm install -g @qvqw2233/wechat-md@latest
```

---

## 卸载

```bash
npm uninstall -g @qvqw2233/wechat-md
```

---

## 下一步

安装完成后，请继续阅读 [配置指南](CONFIG.md) 和 [使用教程](USAGE.md)。
