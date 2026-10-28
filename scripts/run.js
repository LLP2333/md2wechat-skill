#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");

const ext = process.platform === "win32" ? ".exe" : "";
const binaryPath = path.join(__dirname, "..", "bin", `md2wechat${ext}`);

if (!fs.existsSync(binaryPath)) {
  console.error(
    "md2wechat binary is missing. Reinstall with `npm install -g @qvqw2233/wechat-md`."
  );
  process.exit(1);
}

try {
  execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
} catch (error) {
  if (typeof error.status === "number") {
    process.exit(error.status);
  }

  console.error(`Failed to launch md2wechat: ${error.message}`);
  process.exit(1);
}
