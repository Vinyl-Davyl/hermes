import { cpSync, existsSync } from "node:fs";
import { join } from "node:path";

const out = join(import.meta.dirname, "..", "out");
const flat = join(out, "docs.html");
const nested = join(out, "docs", "index.html");
if (!existsSync(flat) && existsSync(nested)) {
  cpSync(nested, flat);
}
