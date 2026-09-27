import { createEngine, type Engine } from "./engine";

function qs(sel: string, root: ParentNode = document): HTMLElement[] {
  return Array.from(root.querySelectorAll(sel));
}

function copyText(text: string) {
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).catch(() => {});
    return;
  }
  const ta = document.createElement("textarea");
  ta.value = text;
  ta.setAttribute("readonly", "");
  ta.style.position = "fixed";
  ta.style.opacity = "0";
  document.body.appendChild(ta);
  ta.select();
  try {
    document.execCommand("copy");
  } catch {
    /* ignore */
  }
  document.body.removeChild(ta);
}

function splitTyped(root: ParentNode) {
  qs("[data-type]", root).forEach((el) => {
    if (el.querySelector(".c")) return;
    const t = el.textContent || "";
    const parts: string[] = [];
    for (let i = 0; i < t.length; i++) {
      if (t[i] === " " && parts.length) parts[parts.length - 1] += " ";
      else parts.push(t[i]);
    }
    el.textContent = "";
    parts.forEach((p, i) => {
      const s = document.createElement("span");
      s.className = "c";
      s.style.setProperty("--i", String(i));
      s.textContent = p;
      el.appendChild(s);
    });
  });
}

type Binding =
  | { el: HTMLElement; t: "c"; keys: string[]; base: string; prev: string | null }
  | { el: HTMLElement; t: "t"; k: string }
  | { el: HTMLElement; t: "s"; pairs: string[][] }
  | { el: HTMLInputElement | HTMLSelectElement; t: "v"; k: string };

export function bindSite(root: ParentNode = document): () => void {
  const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  splitTyped(root);

  const B: Binding[] = [];
  qs("[data-cls]", root).forEach((el) => {
    B.push({
      el,
      t: "c",
      keys: (el.getAttribute("data-cls") || "").split(/\s+/),
      base: el.className,
      prev: null,
    });
  });
  qs("[data-text]", root).forEach((el) => {
    B.push({ el, t: "t", k: el.getAttribute("data-text") || "" });
  });
  qs("[data-sty]", root).forEach((el) => {
    B.push({
      el,
      t: "s",
      pairs: (el.getAttribute("data-sty") || "").split(";").map((p) => p.split(":")),
    });
  });
  qs("[data-val]", root).forEach((el) => {
    B.push({ el: el as HTMLSelectElement, t: "v", k: el.getAttribute("data-val") || "" });
  });

  const raf = window.requestAnimationFrame.bind(window);
  let queued = false;
  let E: Engine;

  function render() {
    queued = false;
    const S = E.S;
    for (const b of B) {
      if (b.t === "c") {
        const v = b.keys.map((k) => S[k] || "").join(" ").replace(/\s+/g, " ").trim();
        if (v !== b.prev) {
          b.el.className = (b.base + " " + v).trim();
          b.prev = v;
        }
      } else if (b.t === "t") {
        const tv = S[b.k];
        if (tv != null && b.el.textContent !== String(tv)) b.el.textContent = tv;
      } else if (b.t === "s") {
        b.pairs.forEach((p) => {
          if (S[p[1]] != null) b.el.style.setProperty(p[0], S[p[1]]);
        });
      } else if (S[b.k] != null && b.el.value !== S[b.k]) {
        b.el.value = S[b.k];
      }
    }
  }

  E = createEngine({
    reduced,
    gate: true,
    copy: copyText,
    onChange() {
      if (!queued) {
        queued = true;
        raf(render);
      }
    },
  });
  render();

  function onClick(e: Event) {
    const target = e.target as HTMLElement | null;
    const c = target?.closest?.("[data-copy]");
    if (c) {
      const code = c.parentElement?.querySelector("code");
      const lbl = c.querySelector("span");
      if (code) copyText(code.textContent?.trim() || "");
      if (lbl) {
        lbl.textContent = "Copied";
        setTimeout(() => { lbl.textContent = "Copy"; }, 1400);
      }
      return;
    }
    const a = target?.closest?.("[data-act]");
    if (!a) return;
    const p = (a.getAttribute("data-act") || "").split(":");
    E.act(p[0], p[1]);
  }

  document.addEventListener("click", onClick);

  const changeCleanups: Array<() => void> = [];
  qs("[data-change]", root).forEach((el) => {
    const handler = () => E.act(el.getAttribute("data-change") || "", (el as HTMLSelectElement).value);
    el.addEventListener("change", handler);
    changeCleanups.push(() => el.removeEventListener("change", handler));
  });

  const stage = document.getElementById("stage");
  function fit() {
    if (!stage) return;
    const mobile = window.matchMedia("(max-width: 760px)").matches;
    const base = mobile ? 358 : 1120;
    const s = Math.min(mobile ? 1.2 : 1, stage.clientWidth / base);
    stage.style.setProperty("--s", s.toFixed(4));
  }
  fit();
  window.addEventListener("resize", fit);

  const faqCleanups: Array<() => void> = [];
  root.querySelectorAll(".faq").forEach((faq) => {
    const handler = (e: Event) => {
      const t = e.target as HTMLDetailsElement | null;
      if (!t || t.tagName !== "DETAILS" || !t.open) return;
      faq.querySelectorAll("details").forEach((d) => {
        if (d !== t) d.open = false;
      });
    };
    faq.addEventListener("toggle", handler, true);
    faqCleanups.push(() => faq.removeEventListener("toggle", handler, true));
  });

  E.start();

  const observers: IntersectionObserver[] = [];
  function watch(id: string, name: "rule" | "tour") {
    const el = document.getElementById(id);
    if (!el) return;
    if (!("IntersectionObserver" in window)) {
      E.visible(name);
      return;
    }
    const io = new IntersectionObserver((es) => {
      es.forEach((en) => {
        if (en.isIntersecting) {
          E.visible(name);
          io.disconnect();
        }
      });
    }, { threshold: 0.3 });
    io.observe(el);
    observers.push(io);
  }
  watch("rule", "rule");
  watch("cli", "tour");

  const toc = qs(".toc a[href^='#']", root);
  if (toc.length && "IntersectionObserver" in window) {
    const map: Record<string, HTMLElement> = {};
    toc.forEach((a) => {
      map[a.getAttribute("href")?.slice(1) || ""] = a;
    });
    const tio = new IntersectionObserver((es) => {
      es.forEach((en) => {
        if (!en.isIntersecting) return;
        toc.forEach((a) => a.classList.remove("cur"));
        map[en.target.id]?.classList.add("cur");
      });
    }, { rootMargin: "-20% 0px -70% 0px" });
    Object.keys(map).forEach((id) => {
      const s = document.getElementById(id);
      if (s) tio.observe(s);
    });
    observers.push(tio);
  }

  return () => {
    E.stop();
    document.removeEventListener("click", onClick);
    changeCleanups.forEach((fn) => fn());
    faqCleanups.forEach((fn) => fn());
    window.removeEventListener("resize", fit);
    observers.forEach((o) => o.disconnect());
  };
}
