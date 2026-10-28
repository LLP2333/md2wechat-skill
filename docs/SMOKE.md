# 真实烟雾测试记录

本文档记录社区版每次发布前的最小闭环验证。只记录结论和关键观察，不记录凭证、草稿 ID 或素材 ID。

## 发布前必须执行

### 1. 本地排版渲染

```bash
make build
./md2wechat convert examples/layout-e2e-test.md --output /tmp/layout-smoke.html
./md2wechat layout validate --file examples/layout-e2e-test.md --json
```

通过标准：

- `convert` 成功，模式为 `local`，不需要任何 API Key
- 输出 HTML 中没有 `:::hero`、`:::toc`、`:::verdict` 等原始模块语法残留
- `layout validate` 返回 `LAYOUT_VALIDATED`，errors 和 warnings 都为 0

### 2. npm 安装

用本次打好的发布包模拟一次全局安装，确认安装脚本能下载、校验并运行二进制：

```bash
npm pack
MD2WECHAT_RELEASE_BASE_URL=/path/to/release-bundle \
  npm install -g --prefix /tmp/md2wechat-npm-smoke ./qvqw2233-wechat-md-*.tgz
/tmp/md2wechat-npm-smoke/bin/md2wechat version --json
/tmp/md2wechat-npm-smoke/bin/md2wechat skills read md2wechat --json
```

发布到 npm 后，再从官方源真实安装一次：

```bash
npm install -g @qvqw2233/wechat-md --registry=https://registry.npmjs.org/
md2wechat version --json
md2wechat skills read md2wechat --json
```

### 3. 微信草稿（有凭证时）

只在具备真实 `WECHAT_APPID` / `WECHAT_SECRET`、IP 白名单已放行时执行：

```bash
md2wechat inspect article.md --draft --cover cover.jpg --json
md2wechat convert article.md --draft --cover cover.jpg --json
```

## 记录

### v3.9.0（2026-10-28）

- 本地模式渲染 `examples/layout-e2e-test.md`：发布清单要求的 19 个核心模块全部渲染，无原始语法残留；`layout validate` 返回 `LAYOUT_VALIDATED`。
- npm：本地打包模拟安装通过；发布后从 npm 官方源真实安装，从 `LLP2333/md2wechat-skill` 的 v3.9.0 Release 下载二进制并通过 `checksums.txt` 校验，`version` 返回 3.9.0，未配置 API Key 时 `convert` 成功。
- 未验证：没有真实微信凭证，未执行草稿创建；未在 Windows 和 Linux 实机运行。
