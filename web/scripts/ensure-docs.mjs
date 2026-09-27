import { cpSync, existsSync, mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const out = join(dirname(fileURLToPath(import.meta.url)), "..", "out");
const flat = join(out, "docs.html");
const nested = join(out, "docs", "index.html");

// Static export may emit either docs.html or docs/index.html. Keep both so
// /docs works in Next, on file servers, and via the Pages rewrite.
if (!existsSync(flat) && existsSync(nested)) {
  cpSync(nested, flat);
}
if (!existsSync(nested) && existsSync(flat)) {
  mkdirSync(dirname(nested), { recursive: true });
  cpSync(flat, nested);
}
