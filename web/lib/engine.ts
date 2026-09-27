export type AgentId =
  | "claude"
  | "codex"
  | "cursor"
  | "cline"
  | "kimi"
  | "antigravity"
  | "opencode"
  | "pi"
  | "copilot"
  | "zcode"
  | "deepseek";

export type Agent = {
  id: string;
  gui: boolean;
  name: string;
};

export const AGENTS: Record<AgentId, Agent> = {
  claude: { id: "7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44", gui: false, name: "Claude Code" },
  codex: { id: "91aa4c07-5d2e-4f61-8b3a-7e2d9c0f1a22", gui: false, name: "Codex" },
  cursor: { id: "e1884976-2f0c-4b7e-9a51-3d0c1b2a7f10", gui: true, name: "Cursor" },
  cline: { id: "4b7d02e1-8a3c-4d9f-b2e6-5c1f0a9d3e77", gui: false, name: "Cline" },
  kimi: { id: "c3e9a118-6f2d-4b0a-9c75-2e8d1f4a6b03", gui: false, name: "Kimi" },
  antigravity: { id: "5f0b7d3a-1e4c-4a8b-8d2f-9b6c0e3a7d51", gui: true, name: "Antigravity" },
  opencode: { id: "a82c4e6f-0b9d-4c1e-a7f3-6d2b8e0c4f19", gui: false, name: "OpenCode" },
  pi: { id: "0d6e2b9c-7a1f-4e3d-b5c8-3f9a1d7e2c64", gui: false, name: "Pi Agent" },
  copilot: { id: "9e4a1c7d-3b8f-4d2a-a6e0-1c5b9f3d7a28", gui: false, name: "Copilot CLI" },
  zcode: { id: "6b3f8d0e-2c7a-4f9b-9e1d-4a8c2f6b0e35", gui: false, name: "ZCode" },
  deepseek: { id: "2a9d5f1b-8e3c-4b7a-8f0d-7c1e5a9b3d42", gui: false, name: "DeepSeek Harness" },
};

export const AGENT_IDS = Object.keys(AGENTS) as AgentId[];

export const COPY: Record<string, string> = {
  install: "curl -fsSL https://tryhermes.pages.dev/install | sh",
  install2: "curl -fsSL https://tryhermes.pages.dev/install | sh",
  doctor: "hermes doctor",
  path: 'export PATH="$HOME/.local/bin:$PATH"',
};

export type State = Record<string, string>;

type EventValue = State | (() => State);
type TimelineEvent = [number, EventValue];

export type Engine = {
  S: State;
  act: (name: string, arg?: string) => void;
  visible: (which: "rule" | "tour") => void;
  start: () => void;
  stop: () => void;
};

export type EngineOptions = {
  reduced?: boolean;
  gate?: boolean;
  copy?: (text: string) => void;
  onChange?: () => void;
};

function assign(a: State, b: State): State {
  for (const k of Object.keys(b)) a[k] = b[k];
  return a;
}

export function createEngine(opt: EngineOptions = {}): Engine {
  const S: State = {};
  const T: Record<string, ReturnType<typeof setTimeout>[]> = {};
  const reduced = !!opt.reduced;
  let queued = false;
  let nonce = 0;
  let tnonce = 0;
  const started = { rule: false, tour: false };

  function changed() {
    if (queued) return;
    queued = true;
    setTimeout(() => {
      queued = false;
      opt.onChange?.();
    }, 0);
  }

  function later(g: string, ms: number, fn: () => void) {
    (T[g] = T[g] || []).push(setTimeout(fn, ms));
  }

  function clear(g: string) {
    (T[g] || []).forEach(clearTimeout);
    T[g] = [];
  }

  function val(p: EventValue): State {
    return typeof p === "function" ? p() : p;
  }

  function set(patch: State, ff?: boolean) {
    for (const k of Object.keys(patch)) {
      const v = patch[k];
      S[k] = ff && typeof v === "string" ? v.replace(/\bon\b/g, "done") : v;
    }
  }

  const STEP_AT = [0, 7800, 12200, 16300];
  const H_LEN = 26000;
  const H_BASE: State = {
    hFade: "",
    hTerm: "front",
    hEd: "back",
    hPaneTop: "focus",
    hPaneBot: "",
    hStallLbl: "",
    hT0: "",
    hL1: "",
    hL2: "",
    hL3: "",
    hL4: "",
    hL5: "",
    hRun: "",
    hL6: "",
    hL7: "",
    hStall: "",
    hZc: "on",
    hT1: "",
    hO1: "",
    hO2: "",
    hO3: "",
    hO4: "",
    hO5: "",
    hToast: "",
    hTree: "",
    hM: "",
    hOld: "done",
    hOldT: "done",
    hNewT: "",
    hPlus: "",
    hPtr: "",
    hKey: "",
    hIn: "",
    hPh: "done",
    hPaste: "",
    hSend: "",
    hU: "",
    hA0: "",
    hA1: "",
    hA2: "",
    hA3: "",
    hA4: "",
    hA5: "",
    hA6: "",
    hA7: "",
    hA8: "",
    hS0: "",
    hS1: "",
    hS2: "",
    hS3: "",
  };

  function stepPatch(i: number): State {
    const p: State = {};
    for (let j = 0; j < 4; j++) p["hS" + j] = j < i ? "past" : j === i ? "cur " + (nonce % 2 ? "a" : "b") : "";
    return p;
  }

  const H_EV: TimelineEvent[] = [
    [0, () => stepPatch(0)],
    [400, { hT0: "on" }],
    [1650, { hT0: "done" }],
    [1900, { hL1: "on" }],
    [2500, { hL2: "on" }],
    [3100, { hL3: "on" }],
    [3700, { hL4: "on" }],
    [4300, { hL5: "on", hRun: "on" }],
    [5300, { hRun: "", hL6: "on" }],
    [6000, { hL7: "on" }],
    [6900, { hStall: "on", hStallLbl: "err" }],
    [7800, () => assign(stepPatch(1), { hPaneTop: "", hPaneBot: "focus", hZc: "", hT1: "on" })],
    [9500, { hT1: "done", hO1: "on" }],
    [9750, { hO2: "on" }],
    [10000, { hO3: "on", hTree: "on" }],
    [10250, { hO4: "on", hToast: "on" }],
    [10600, { hO5: "on" }],
    [12200, () => assign(stepPatch(2), { hTerm: "back", hEd: "front", hToast: "", hPtr: "p0" })],
    [12450, { hPtr: "p1" }],
    [13250, { hPtr: "p1 click", hPlus: "press" }],
    [13450, { hPtr: "p1", hPlus: "", hOld: "", hOldT: "", hNewT: "on", hIn: "focus" }],
    [13650, { hPtr: "p2" }],
    [14400, { hPtr: "p2 click", hKey: "on" }],
    [14650, { hPtr: "p2", hPh: "", hPaste: "on", hSend: "ready" }],
    [15100, { hKey: "" }],
    [15200, { hPtr: "p3" }],
    [15950, { hPtr: "p3 click", hSend: "ready press" }],
    [16150, { hPtr: "p3", hU: "on", hPaste: "", hPh: "done", hSend: "", hIn: "" }],
    [16300, () => assign(stepPatch(3), { hPtr: "", hA0: "on" })],
    [17000, { hA0: "", hA1: "on" }],
    [17350, { hA2: "on" }],
    [17700, { hA3: "on" }],
    [18300, { hA4: "on" }],
    [19300, { hA5: "on", hM: "on" }],
    [20300, { hA6: "on" }],
    [21200, { hA6: "done ok", hA7: "on" }],
    [21900, { hA8: "on" }],
    [25400, { hFade: "fade" }],
  ];

  function heroPlay(from: number, still?: boolean) {
    clear("hero");
    nonce++;
    set(H_BASE);
    H_EV.forEach((e) => {
      if (e[0] <= from) set(val(e[1]), true);
    });
    changed();
    if (still) return;
    H_EV.forEach((e) => {
      if (e[0] > from) later("hero", e[0] - from, () => { set(val(e[1])); changed(); });
    });
    later("hero", H_LEN - from, () => { heroPlay(0); });
  }

  const EX: [string, AgentId][] = [
    ["", "claude"],
    ["", "cursor"],
    ["cursor", "claude"],
    ["claude", "cursor"],
  ];
  let ruleIdx = 0;

  function ruleTexts(from: string, to: AgentId): State {
    const src = (from || (to === "claude" ? "cursor" : "claude")) as AgentId;
    const gui = AGENTS[to].gui;
    const name = (k: AgentId) => AGENTS[k].name;
    return {
      rFrom: from,
      rTo: to,
      rCmd: "hermes handoff " + (from ? from + " " : "") + to,
      rFromTxt: src + " id=" + AGENTS[src].id,
      rToTxt: to,
      rDelK: gui ? "copied" : "opened",
      rDelTxt: gui ? "PROMPT.md → clipboard" : to + " (seeded from the pack)",
      rNext: "hermes resume " + to,
      rMean: from
        ? "From " + name(from as AgentId) + " to " + name(to) + "."
        : "To " + name(to) + ". Source: the newest chat in this folder.",
      rHowTxt: gui
        ? "Editor: PROMPT.md goes to your clipboard. Open a new chat and paste."
        : "CLI: Hermes starts it seeded from the pack when it is on your PATH.",
    };
  }

  function ruleRun(from: string, to: AgentId, auto?: boolean, still?: boolean) {
    clear("rule");
    const t = ruleTexts(from, to);
    set(t);
    S.rN = String(t.rCmd.length);
    for (let i = 0; i < 4; i++) S["rE" + i] = EX[i][0] === from && EX[i][1] === to ? "cur" : "";
    ["rO1", "rO2", "rO3", "rO4", "rO5", "rHow"].forEach((k) => { S[k] = still ? "done" : ""; });
    S.rTw = still ? "done" : "";
    changed();
    if (still) return;
    const go = t.rCmd.length * 34 + 60 + 260;
    later("rule", 60, () => { S.rTw = "on"; changed(); });
    later("rule", go, () => { S.rTw = "done"; S.rO1 = "on"; changed(); });
    later("rule", go + 190, () => { S.rO2 = "on"; changed(); });
    later("rule", go + 380, () => { S.rO3 = "on"; changed(); });
    later("rule", go + 570, () => { S.rO4 = "on"; changed(); });
    later("rule", go + 820, () => { S.rO5 = "on"; S.rHow = "on"; changed(); });
    if (auto) {
      later("rule", go + 820 + 3000, () => {
        ruleIdx = (ruleIdx + 1) % EX.length;
        ruleRun(EX[ruleIdx][0], EX[ruleIdx][1], true);
      });
    }
  }

  const T_DUR = [9200, 7000, 7000, 7400];
  const T_TITLE = ["app · zsh", "blog · zsh", "app · zsh", "app · zsh"];
  const SESS = [
    { ag: "claude" as AgentId, to: "cursor" as AgentId, id: AGENTS.claude.id },
    { ag: "cursor" as AgentId, to: "claude" as AgentId, id: AGENTS.cursor.id },
    { ag: "cursor" as AgentId, to: "claude" as AgentId, id: "b03d11e2-9c4a-4e0b-8f7d-61a2c3e4f5a6" },
  ];
  const T_KEYS = [
    "tA_c", "tA_h1", "tA_r0", "tA_h2", "tA_r1", "tA_r2", "tA_s0", "tA_s1", "tA_s2", "tA_p", "tA_tw",
    "tA_q1", "tA_q2", "tA_q3", "tA_q4",
    "tB_c", "tB_o1", "tB_o2", "tB_o3", "tB_o4", "tB_o5", "tB_o6", "tC_c", "tD_c",
  ];
  for (let q = 1; q <= 10; q++) T_KEYS.push("tC_o" + q);
  for (let q = 0; q <= 17; q++) T_KEYS.push("tD_o" + q);

  function pickTexts(i: number): State {
    const s = SESS[i];
    const gui = AGENTS[s.to].gui;
    const cmd = "hermes handoff " + s.ag + " " + s.to + " --id " + s.id.slice(0, 4);
    return {
      tA_s0: i === 0 ? "sel" : "",
      tA_s1: i === 1 ? "sel" : "",
      tA_s2: i === 2 ? "sel" : "",
      tA_cmd: cmd,
      tA_n: String(cmd.length),
      tA_from: s.ag + " id=" + s.id,
      tA_to: s.to,
      tA_dk: gui ? "copied" : "opened",
      tA_dt: gui ? "PROMPT.md → clipboard" : s.to + " (seeded from the pack)",
    };
  }

  function tourPick(i: number, still?: boolean) {
    clear("pick");
    set(pickTexts(i));
    S.tA_p = still ? "done" : "on";
    S.tA_tw = still ? "done" : "";
    ["tA_q1", "tA_q2", "tA_q3", "tA_q4"].forEach((k) => { S[k] = still ? "done" : ""; });
    changed();
    if (still) return;
    const go = S.tA_cmd.length * 34 + 60 + 260;
    later("pick", 60, () => { S.tA_tw = "on"; changed(); });
    later("pick", go, () => { S.tA_tw = "done"; S.tA_q1 = "on"; changed(); });
    later("pick", go + 190, () => { S.tA_q2 = "on"; changed(); });
    later("pick", go + 380, () => { S.tA_q3 = "on"; changed(); });
    later("pick", go + 570, () => { S.tA_q4 = "on"; changed(); });
  }

  function lines(prefix: string, from: number, to: number, start: number, step: number): TimelineEvent[] {
    const ev: TimelineEvent[] = [];
    for (let n = from; n <= to; n++) {
      const p: State = {};
      p[prefix + n] = "on";
      ev.push([start + (n - from) * step, p]);
    }
    return ev;
  }

  const SCN: TimelineEvent[][] = [
    [
      [0, { tA_c: "on" }],
      [760, { tA_c: "done", tA_h1: "on" }],
      [900, { tA_r0: "on" }],
      [1080, { tA_h2: "on" }],
      [1220, { tA_r1: "on" }],
      [1360, { tA_r2: "on" }],
      [2400, () => { tourPick(1); return {}; }],
    ],
    [
      [0, { tB_c: "on" }],
      [820, { tB_c: "done", tB_o1: "on" }],
      [1010, { tB_o2: "on" }],
      [1200, { tB_o3: "on" }],
      [1390, { tB_o4: "on" }],
      [1800, { tB_o5: "on" }],
      [2300, { tB_o6: "on" }],
    ],
    [[0, { tC_c: "on" }], [860, { tC_c: "done" }], ...lines("tC_o", 1, 8, 860, 120)],
    [[0, { tD_c: "on" }], [620, { tD_c: "done", tD_o0: "on" }], ...lines("tD_o", 1, 15, 780, 90), [2250, { tD_o16: "on" }], [2350, { tD_o17: "on" }]],
  ];

  function tourPlay(i: number, auto?: boolean, still?: boolean) {
    clear("tour");
    clear("pick");
    tnonce++;
    T_KEYS.forEach((k) => { S[k] = ""; });
    for (let j = 0; j < 4; j++) {
      S["tT" + j] = j === i ? "cur " + (tnonce % 2 ? "a" : "b") : "";
      S["tV" + j] = j === i ? "done" : "";
    }
    S.tTitle = T_TITLE[i];
    set(pickTexts(1));
    S.tA_s1 = "";
    if (still) {
      SCN[i].forEach((e) => {
        if (typeof e[1] !== "function") set(e[1], true);
      });
      if (i === 0) tourPick(1, true);
      changed();
      return;
    }
    changed();
    SCN[i].forEach((e) => {
      later("tour", e[0], () => { set(val(e[1])); changed(); });
    });
    if (auto) later("tour", T_DUR[i], () => { tourPlay((i + 1) % 4, true); });
  }

  ["install", "install2", "doctor", "path", "rule"].forEach((k) => { S["cp_" + k] = "Copy"; });

  function copy(key: string) {
    const text = COPY[key] || (key === "rule" ? S.rCmd : "");
    if (!text) return;
    opt.copy?.(text);
    S["cp_" + key] = "Copied";
    changed();
    later("copy", 1400, () => { S["cp_" + key] = "Copy"; changed(); });
  }

  const eng: Engine = {
    S,
    act(name, arg) {
      if (name === "seek") {
        const st = +(arg || 0);
        if (reduced) heroPlay(st < 3 ? STEP_AT[st + 1] - 1 : H_LEN - 1200, true);
        else heroPlay(STEP_AT[st]);
      } else if (name === "ex") {
        ruleIdx = +(arg || 0);
        ruleRun(EX[ruleIdx][0], EX[ruleIdx][1], false, reduced);
      } else if (name === "from") {
        const f = arg || "";
        let t = S.rTo as AgentId;
        if (f && f === t) t = f === "claude" ? "cursor" : "claude";
        ruleRun(f, t, false, reduced);
      } else if (name === "to") {
        ruleRun(S.rFrom === arg ? "" : S.rFrom, (arg || "cursor") as AgentId, false, reduced);
      } else if (name === "tab") {
        tourPlay(+(arg || 0), false, reduced);
      } else if (name === "pick") {
        clear("tour");
        ["tA_c", "tA_h1", "tA_r0", "tA_h2", "tA_r1", "tA_r2"].forEach((k) => { S[k] = "done"; });
        tourPick(+(arg || 0), reduced);
      } else if (name === "copy") {
        copy(arg || "");
      }
    },
    visible(which) {
      if (started[which]) return;
      started[which] = true;
      if (which === "rule") ruleRun(EX[0][0], EX[0][1], !reduced, reduced);
      if (which === "tour") tourPlay(0, !reduced, reduced);
    },
    start() {
      heroPlay(reduced ? H_LEN - 1200 : 0, reduced);
      if (!opt.gate) {
        eng.visible("rule");
        eng.visible("tour");
      }
    },
    stop() {
      Object.keys(T).forEach(clear);
    },
  };

  heroPlay(10700, true);
  ruleRun(EX[1][0], EX[1][1], false, true);
  tourPlay(0, false, true);
  return eng;
}
