import { existsSync, rmSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const out = join(dirname(fileURLToPath(import.meta.url)), "..", "out");
const flat = join(out, "docs.html");
const nestedDir = join(out, "docs");
const nested = join(nestedDir, "index.html");

// Next emits docs.html when trailingSlash is false. A sibling docs/index.html
// makes Cloudflare Pages bounce /docs ↔ /docs/ forever.
if (existsSync(flat) && existsSync(nested)) {
  rmSync(nestedDir, { recursive: true, force: true });
}
