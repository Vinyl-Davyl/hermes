import type { CSSProperties } from "react";
import { AGENT_IDS } from "@/lib/engine";
import { ArrowIcon, CheckIcon, ClipboardIcon, FileIcon, SendIcon } from "./icons";
import { CommandChip, Footer, Header, Lights, Room } from "./chrome";
import { SiteBindings } from "./site-bindings";

const AGENTS_LINE =
  "Claude Code · Codex · Cursor · Cline · Kimi · Antigravity · OpenCode · Pi Agent · Copilot CLI · ZCode · DeepSeek Harness";

function AgentOptions({ includeNewest = false }: { includeNewest?: boolean }) {
  return (
    <>
      {includeNewest ? <option value="">newest here</option> : null}
      {AGENT_IDS.map((id) => (
        <option key={id} value={id}>{id}</option>
      ))}
    </>
  );
}

export function HomePage() {
  return (
    <>
      <SiteBindings />
      <Header page="home" />
      <main>
        <section className="hero dark">
          <Room />
          <div className="wrap hero-copy">
            <a className="badge" href="https://github.com/Vinyl-Davyl/hermes" target="_blank" rel="noopener noreferrer">
              Open source · local-first · MIT
              <ArrowIcon />
            </a>
            <h1>Don’t start the next chat from zero.</h1>
            <p className="lede">When one coding agent stops, Hermes packs the work on your disk and hands it to the next. Same folder. Same machine.</p>
            <div className="cta">
              <CommandChip act="copy:install" textKey="cp_install" />
              <a className="btn ghost" href="/docs">Read the docs</a>
            </div>
            <p className="aside">Other tools remember forever or plug into MCP. Hermes just carries the work across: this folder, this machine, one command.</p>
          </div>

          <div className="wrap">
            <div className="stage" id="stage" aria-label="Replay of a handoff: Claude Code stops, hermes packs the work, Cursor continues">
              <div className="stage-in" data-cls="hFade">
                <div className="win term" data-cls="hTerm">
                  <div className="tbar"><Lights /><span className="ttl">app · zsh</span></div>
                  <div className="split">
                    <div className="pane top" data-cls="hPaneTop">
                      <span className="plabel" data-cls="hStallLbl">claude</span>
                      <div className="flow">
                        <div className="ln done hd"><span className="g">✻</span> Claude Code · ~/src/app</div>
                        <div className="ln done"><span className="gt">&gt;</span> <span className="tw" data-cls="hT0" data-type>add rate limiting to POST /api/upload</span></div>
                        <div className="ln bul dsk" data-cls="hL1"><span className="b">●</span><span>I’ll add a token bucket and wire it into the upload route.</span></div>
                        <div className="ln bul" data-cls="hL2"><span className="b ok">●</span><span><b>Read</b>(api/upload.ts)</span><div className="sub"><span className="el">⎿</span>Read 88 lines</div></div>
                        <div className="ln bul dsk" data-cls="hL3"><span className="b ok">●</span><span><b>Write</b>(middleware/ratelimit.ts)</span><div className="sub"><span className="el">⎿</span>Wrote 42 lines</div></div>
                        <div className="ln bul" data-cls="hL4"><span className="b ok">●</span><span><b>Update</b>(api/upload.ts)</span><div className="sub"><span className="el">⎿</span><span className="add">+6</span> <span className="del">−1</span></div></div>
                        <div className="ln bul" data-cls="hL5"><span className="b ok">●</span><span><b>Bash</b>(pnpm test ratelimit)</span>
                          <div className="sub ln" data-cls="hRun"><span className="el">⎿</span><span className="spin" aria-hidden="true" /><span>Running…</span></div>
                          <div className="sub ln" data-cls="hL6"><span className="el">⎿</span><span>5 passed, <span className="del">1 failed</span><br /><span className="del">✗ refills the bucket after the window</span></span></div>
                        </div>
                        <div className="ln bul" data-cls="hL7"><span className="b">●</span><span>The refill test fails: the bucket reads <code>Date.now()</code>. Injecting a clock…</span></div>
                        <div className="ln stall" data-cls="hStall"><span className="warn-i">!</span><span>Context limit reached. Start a new session to continue.</span></div>
                      </div>
                    </div>
                    <div className="pane bot" data-cls="hPaneBot">
                      <div className="flow">
                        <div className="ln done"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="caret" data-cls="hZc" /><span className="tw fast" data-cls="hT1" data-type>hermes handoff claude cursor -m &quot;finish upload rate limit&quot;</span></div>
                        <div className="ln" data-cls="hO1"><span className="k">packed</span> handoff-20260927T141205</div>
                        <div className="ln" data-cls="hO2"><span className="k">from</span> claude id=7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44</div>
                        <div className="ln" data-cls="hO3"><span className="k">to</span> cursor</div>
                        <div className="ln" data-cls="hO4"><span className="k">copied</span> PROMPT.md → clipboard</div>
                        <div className="ln gap" data-cls="hO5"><span className="k">next</span> hermes resume cursor</div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="win ed" data-cls="hEd">
                  <div className="tbar"><Lights /><span className="ttl">app · Cursor</span></div>
                  <div className="ed-body">
                    <div className="tree" aria-hidden="true">
                      <div className="tree-h">APP</div>
                      <div className="row dir"><span className="chev">›</span>api</div>
                      <div className="row f i1"><span className="fi ts" />upload.ts</div>
                      <div className="row dir"><span className="chev">›</span>lib</div>
                      <div className="row f i1"><span className="fi ts" />clock.ts</div>
                      <div className="row dir"><span className="chev">›</span>middleware</div>
                      <div className="row f i1 cur"><span className="fi ts" />ratelimit.ts<span className="gs m" data-cls="hM">M</span></div>
                      <div className="row dir"><span className="chev">›</span>tests</div>
                      <div className="row f i1"><span className="fi ts" />ratelimit.test.ts</div>
                      <div className="row f"><span className="fi js" />package.json</div>
                      <div className="pack" data-cls="hTree">
                        <div className="row dir u"><span className="chev">›</span>handoff-20260927T141205<span className="gs">U</span></div>
                        <div className="row f i1 u"><span className="fi md" />PROMPT.md</div>
                        <div className="row f i1 u"><span className="fi md" />SUMMARY.md</div>
                        <div className="row f i1 u"><span className="fi js" />manifest.json</div>
                        <div className="row dir i1 u"><span className="chev r">›</span>git</div>
                        <div className="row dir i1 u"><span className="chev r">›</span>sessions</div>
                      </div>
                    </div>
                    <div className="chat">
                      <div className="chat-h">
                        <span className="ct old" data-cls="hOldT">Refactor auth middleware</span>
                        <span className="ct new" data-cls="hNewT">New Chat</span>
                        <span className="ch-icons" aria-hidden="true">
                          <span className="ico plus" data-cls="hPlus"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><path d="M12 5v14M5 12h14" /></svg></span>
                          <span className="ico"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8" /><path d="M3 3v5h5M12 7v5l3 2" /></svg></span>
                          <span className="ico"><svg viewBox="0 0 24 24" fill="currentColor"><circle cx="5" cy="12" r="1.6" /><circle cx="12" cy="12" r="1.6" /><circle cx="19" cy="12" r="1.6" /></svg></span>
                        </span>
                      </div>
                      <div className="msgs">
                        <div className="flow">
                          <div className="old" data-cls="hOld">
                            <div className="um">refactor auth middleware to use the session store</div>
                            <p className="am">Done. Session checks now live in <code>middleware/auth.ts</code> and read from the store.</p>
                          </div>
                          <div className="um paste-msg" data-cls="hU">
                            <b># Hermes handoff</b>
                            <span>You are continuing work from another coding agent / window. Do not restart from scratch. Use the git state and notes below.</span>
                            <span className="goal">## Goal · finish upload rate limit</span>
                          </div>
                          <div className="thinking" data-cls="hA0">Reading the handoff</div>
                          <div className="tool" data-cls="hA1"><FileIcon /><b>Read</b><span>PROMPT.md</span></div>
                          <div className="tool" data-cls="hA2"><FileIcon /><b>Read</b><span>git/diff.patch</span></div>
                          <div className="tool" data-cls="hA3"><FileIcon /><b>Read</b><span>middleware/ratelimit.ts</span></div>
                          <p className="am" data-cls="hA4">Picking up from Claude Code on <code>main</code>. The limiter works. The refill test fails because the bucket reads <code>Date.now()</code> directly.</p>
                          <div className="diff" data-cls="hA5">
                            <div className="diff-h"><span>middleware/ratelimit.ts</span><span><span className="add">+2</span> <span className="del">−1</span></span></div>
                            <div className="dl a"><span>+</span>import {"{"} clock {"}"} from &quot;../lib/clock&quot;;</div>
                            <div className="dl d"><span>−</span>  const now = Date.now();</div>
                            <div className="dl a"><span>+</span>  const now = clock.now();</div>
                          </div>
                          <div className="tool run" data-cls="hA6"><span className="spin" aria-hidden="true" /><CheckIcon /><b>Run</b><span>pnpm test ratelimit</span></div>
                          <div className="okline" data-cls="hA7">6 passed</div>
                          <p className="am" data-cls="hA8">Done. <code>POST /api/upload</code> is limited to 20 requests a minute per key.</p>
                        </div>
                      </div>
                      <div className="input" data-cls="hIn">
                        <div className="ph" data-cls="hPh">Ask anything, @ to add files</div>
                        <div className="pasted" data-cls="hPaste"><b># Hermes handoff</b><br />You are continuing work from another coding agent / window. Do not restart from scratch…<br />## Goal<br />finish upload rate limit</div>
                        <div className="in-foot">
                          <span className="mode">Agent<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><path d="m6 9 6 6 6-6" /></svg></span>
                          <span className="send" data-cls="hSend" aria-hidden="true"><SendIcon /></span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="toast" data-cls="hToast"><ClipboardIcon />PROMPT.md copied to clipboard</div>
                <div className="kbd" data-cls="hKey" aria-hidden="true"><kbd>⌘</kbd><kbd>V</kbd></div>
                <div className="ptr" data-cls="hPtr" aria-hidden="true">
                  <svg viewBox="0 0 20 24" width="20" height="24"><path d="M2.5 2v17.6l4.5-4.2 2.9 6.6 3.1-1.4-2.9-6.4h6.2z" fill="#111" stroke="#fff" strokeWidth="1.5" strokeLinejoin="round" /></svg>
                </div>
              </div>
            </div>

            <div className="steps" role="group" aria-label="Handoff steps">
              <button type="button" className="step" data-cls="hS0" data-act="seek:0" style={{ "--dur": "7800ms" } as CSSProperties}><span className="n">01</span><span className="l">Claude stalls</span><i className="bar" /></button>
              <button type="button" className="step" data-cls="hS1" data-act="seek:1" style={{ "--dur": "4400ms" } as CSSProperties}><span className="n">02</span><span className="l"><code>hermes handoff claude cursor</code></span><i className="bar" /></button>
              <button type="button" className="step" data-cls="hS2" data-act="seek:2" style={{ "--dur": "4100ms" } as CSSProperties}><span className="n">03</span><span className="l">New Cursor chat, paste</span><i className="bar" /></button>
              <button type="button" className="step" data-cls="hS3" data-act="seek:3" style={{ "--dur": "9700ms" } as CSSProperties}><span className="n">04</span><span className="l">Cursor keeps going</span><i className="bar" /></button>
            </div>
            <p className="agents"><span>Hands off between</span> {AGENTS_LINE}</p>
          </div>
        </section>

        <section className="band light" id="rule">
          <div className="wrap rule">
            <div className="rule-copy">
              <h2>One rule.</h2>
              <p className="lede">One name is where the work goes. Two names are from, then to.</p>
              <div className="exs" role="group" aria-label="Examples">
                <button type="button" className="ex" data-cls="rE0" data-act="ex:0"><code>hermes handoff <b>claude</b></code><span>Go to Claude</span></button>
                <button type="button" className="ex" data-cls="rE1" data-act="ex:1"><code>hermes handoff <b>cursor</b></code><span>Go to Cursor. Not leave it.</span></button>
                <button type="button" className="ex" data-cls="rE2" data-act="ex:2"><code>hermes handoff <b>cursor claude</b></code><span>From Cursor to Claude</span></button>
                <button type="button" className="ex" data-cls="rE3" data-act="ex:3"><code>hermes handoff <b>claude cursor</b></code><span>From Claude to Cursor</span></button>
              </div>
              <div className="own">
                <label htmlFor="r-from">From</label>
                <div className="dd">
                  <select id="r-from" data-change="from" data-val="rFrom">
                    <AgentOptions includeNewest />
                  </select>
                </div>
                <label htmlFor="r-to">To</label>
                <div className="dd">
                  <select id="r-to" data-change="to" data-val="rTo">
                    <AgentOptions />
                  </select>
                </div>
              </div>
            </div>
            <div className="rule-demo">
              <div className="win term solo">
                <div className="tbar"><Lights /><span className="ttl">app · zsh</span></div>
                <div className="pane">
                  <div className="flow" aria-live="polite">
                    <div className="ln done"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="tw1" data-cls="rTw" data-sty="--n:rN" data-text="rCmd">hermes handoff cursor</span><span className="car" /></div>
                    <div className="ln" data-cls="rO1"><span className="k">packed</span> handoff-20260927T141205</div>
                    <div className="ln" data-cls="rO2"><span className="k">from</span> <span data-text="rFromTxt">claude id=7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44</span></div>
                    <div className="ln" data-cls="rO3"><span className="k">to</span> <span data-text="rToTxt">cursor</span></div>
                    <div className="ln" data-cls="rO4"><span className="k" data-text="rDelK">copied</span> <span data-text="rDelTxt">PROMPT.md → clipboard</span></div>
                    <div className="ln gap" data-cls="rO5"><span className="k">next</span> <span data-text="rNext">hermes resume cursor</span></div>
                  </div>
                </div>
                <div className="term-foot">
                  <span className="how" data-cls="rHow"><b data-text="rMean">To Cursor. Source: the newest chat in this folder.</b><span data-text="rHowTxt">Editor: PROMPT.md goes to your clipboard. Open a new chat and paste.</span></span>
                  <button type="button" className="copy" data-act="copy:rule" aria-label="Copy this command">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><rect x="9" y="9" width="12" height="12" rx="2" /><path d="M5 15H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v1" /></svg>
                    <span data-text="cp_rule">Copy</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="band light" id="cli">
          <div className="wrap">
            <div className="sec-head">
              <h2>The rest of the CLI.</h2>
              <p className="lede">Pick the exact chat. Stay in this folder. Resume without globs. Check what Hermes can see.</p>
            </div>
            <div className="tour">
              <div className="tabs" role="tablist" aria-label="Commands">
                <button type="button" role="tab" className="tab" data-cls="tT0" data-act="tab:0" style={{ "--dur": "9200ms" } as CSSProperties}><b>Pick the exact chat</b><span><code>hermes list --here</code> shows this folder’s sessions. Pass <code>--id</code>. A prefix is enough.</span><i className="bar" /></button>
                <button type="button" role="tab" className="tab" data-cls="tT1" data-act="tab:1" style={{ "--dur": "7000ms" } as CSSProperties}><b>Stay in this folder</b><span>Chats from other projects are ignored. No match means a git-only pack and a warning.</span><i className="bar" /></button>
                <button type="button" role="tab" className="tab" data-cls="tT2" data-act="tab:2" style={{ "--dur": "7000ms" } as CSSProperties}><b>Resume without globs</b><span><code>hermes resume</code> reuses the latest pack. Never type <code>./handoff-*</code>.</span><i className="bar" /></button>
                <button type="button" role="tab" className="tab" data-cls="tT3" data-act="tab:3" style={{ "--dur": "7400ms" } as CSSProperties}><b>Check the machine</b><span><code>hermes doctor</code> shows what git sees and which agents Hermes can read.</span><i className="bar" /></button>
              </div>
              <div className="win term solo tall">
                <div className="tbar"><Lights /><span className="ttl" data-text="tTitle">app · zsh</span></div>
                <div className="pane">
                  <div className="flow scene" data-cls="tV0">
                    <div className="ln done"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="tw" data-cls="tA_c" data-type>hermes list --here</span></div>
                    <div className="ln gap" data-cls="tA_h1"># claude</div>
                    <button type="button" className="ln srow" data-cls="tA_r0 tA_s0" data-act="pick:0"><span> 2026-09-27 13:48 <span className="id">7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44</span> refactor auth middleware</span><span className="pth"> ~/.claude/projects/-Users-you-src-app/7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44.jsonl</span></button>
                    <div className="ln gap" data-cls="tA_h2"># cursor</div>
                    <button type="button" className="ln srow" data-cls="tA_r1 tA_s1" data-act="pick:1"><span> 2026-09-27 14:10 <span className="id">e1884976-2f0c-4b7e-9a51-3d0c1b2a7f10</span> upload rate limit</span><span className="pth"> ~/Library/Application Support/Cursor/User/globalStorage/state.vscdb</span></button>
                    <button type="button" className="ln srow" data-cls="tA_r2 tA_s2" data-act="pick:2"><span> 2026-09-27 11:02 <span className="id">b03d11e2-9c4a-4e0b-8f7d-61a2c3e4f5a6</span> fix flaky upload test</span><span className="pth"> ~/Library/Application Support/Cursor/User/globalStorage/state.vscdb</span></button>
                    <div className="ln gap" data-cls="tA_p"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="tw1" data-cls="tA_tw" data-sty="--n:tA_n" data-text="tA_cmd">hermes handoff cursor claude --id e188</span><span className="car" /></div>
                    <div className="ln" data-cls="tA_q1"><span className="k">packed</span> handoff-20260927T141205</div>
                    <div className="ln" data-cls="tA_q2"><span className="k">from</span> <span data-text="tA_from">cursor id=e1884976-2f0c-4b7e-9a51-3d0c1b2a7f10</span></div>
                    <div className="ln" data-cls="tA_q3"><span className="k">to</span> <span data-text="tA_to">claude</span></div>
                    <div className="ln" data-cls="tA_q4"><span className="k" data-text="tA_dk">opened</span> <span data-text="tA_dt">claude (seeded from the pack)</span></div>
                  </div>
                  <div className="flow scene" data-cls="tV1">
                    <div className="ln done"><span className="ps">~/src/blog</span> <span className="br">main</span> <span className="pr">❯</span> <span className="tw" data-cls="tB_c" data-type>hermes handoff codex</span></div>
                    <div className="ln" data-cls="tB_o1"><span className="k">packed</span> handoff-20260927T142030</div>
                    <div className="ln" data-cls="tB_o2"><span className="k">from</span> git only</div>
                    <div className="ln" data-cls="tB_o3"><span className="k">to</span> codex</div>
                    <div className="ln" data-cls="tB_o4"><span className="k">opened</span> codex (seeded from the pack)</div>
                    <div className="ln warnl" data-cls="tB_o5"><span className="k">warn</span> no agent session for this folder; pack is git-only (this is fine). Use --id to pick a chat from another project</div>
                    <div className="ln gap" data-cls="tB_o6"><span className="k">next</span> hermes resume codex</div>
                  </div>
                  <div className="flow scene" data-cls="tV2">
                    <div className="ln done"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="tw" data-cls="tC_c" data-type>hermes resume --print</span></div>
                    <div className="ln mdh" data-cls="tC_o1"># Hermes handoff</div>
                    <div className="ln gap" data-cls="tC_o2">You are continuing work from another coding agent / window. Do **not** restart from scratch. Use the git state and notes below.</div>
                    <div className="ln gap mdh" data-cls="tC_o3">## Goal</div>
                    <div className="ln gap" data-cls="tC_o4">finish upload rate limit</div>
                    <div className="ln gap mdh" data-cls="tC_o5">## Repo</div>
                    <div className="ln gap" data-cls="tC_o6">- Path: `~/src/app`<br />- Branch: `main`<br />- Came from: `claude`</div>
                    <div className="ln gap mdh" data-cls="tC_o7">## Files in play</div>
                    <div className="ln gap" data-cls="tC_o8">Read these first. The patch is in `git/diff.patch`.</div>
                  </div>
                  <div className="flow scene" data-cls="tV3">
                    <div className="ln done"><span className="ps">~/src/app</span> <span className="br">main*</span> <span className="pr">❯</span> <span className="tw" data-cls="tD_c" data-type>hermes doctor</span></div>
                    <div className="ln" data-cls="tD_o0">hermes doctor</div>
                    <div className="ln gap pre" data-cls="tD_o1"> <span className="y">✓</span> git binary             ok</div>
                    <div className="ln pre" data-cls="tD_o2"> <span className="y">✓</span> git repo               ~/src/app (branch main, 3 changed files)</div>
                    <div className="ln pre" data-cls="tD_o3"> <span className="y">✓</span> sqlite3 (Cursor DBs)   ok</div>
                    <div className="ln pre" data-cls="tD_o4"> <span className="y">✓</span> claude                 3 session(s)</div>
                    <div className="ln pre" data-cls="tD_o5"> <span className="y">✓</span> cursor                 2 session(s)</div>
                    <div className="ln pre" data-cls="tD_o6"> <span className="y">✓</span> codex                  1 session(s)</div>
                    <div className="ln pre" data-cls="tD_o7"> <span className="x">✗</span> opencode               0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o8"> <span className="x">✗</span> cline                  0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o9"> <span className="x">✗</span> kimi                   0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o10"> <span className="y">✓</span> antigravity            1 session(s)</div>
                    <div className="ln pre" data-cls="tD_o11"> <span className="x">✗</span> pi                     0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o12"> <span className="x">✗</span> copilot                0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o13"> <span className="x">✗</span> zcode                  0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o14"> <span className="x">✗</span> deepseek               0 session(s)</div>
                    <div className="ln pre" data-cls="tD_o15"> <span className="y">✓</span> home                   /Users/you</div>
                    <div className="ln gap" data-cls="tD_o16">Hermes never writes into Cursor/Claude databases.</div>
                    <div className="ln" data-cls="tD_o17">You do not need every agent installed to test Hermes.</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="band light" id="know">
          <div className="wrap know">
            <div className="know-copy">
              <h2>Things to know.</h2>
              <p className="lede">Hermes is a one-shot messenger. It packs this folder’s work and carries it to the next chat. That is the whole job.</p>
            </div>
            <div className="faq">
              <details open>
                <summary>What goes in the pack?<span className="pm" aria-hidden="true" /></summary>
                <div className="a">
                  <p>A <code>handoff-&lt;timestamp&gt;/</code> folder next to your project. <code>PROMPT.md</code> is the brief the next agent reads. The rest is there if it wants to dig.</p>
                  <ul className="files"><li>PROMPT.md</li><li>SUMMARY.md</li><li>manifest.json</li><li>git/diff.patch</li><li>git/log.txt</li><li>sessions/claude.md</li><li>DECISIONS.md <em>--decisions</em></li><li>files/ <em>--include-files</em></li></ul>
                </div>
              </details>
              <details>
                <summary>Do I paste, or does the next agent open?<span className="pm" aria-hidden="true" /></summary>
                <div className="a"><p>CLI agents such as Claude Code or Codex start seeded from the pack when they are on your PATH. Editors such as Cursor and Antigravity get <code>PROMPT.md</code> on your clipboard: open a new chat and paste, or type <code>@PROMPT.md</code>. Force either path with <code>--open</code>, <code>--no-open</code> or <code>--copy</code>.</p></div>
              </details>
              <details>
                <summary>How do I pick one chat out of many?<span className="pm" aria-hidden="true" /></summary>
                <div className="a"><p><code>hermes list --here</code> shows the sessions for this folder. Filter with <code>--agent cursor</code> or <code>-q keyword</code>, then pass <code>--id</code>. A prefix of the id is enough. Without <code>--id</code>, Hermes takes the newest session that belongs to this folder.</p></div>
              </details>
              <details>
                <summary>Does anything leave my machine?<span className="pm" aria-hidden="true" /></summary>
                <div className="a"><p>No. Hermes reads local git and local session files, then writes a folder on disk. No upload, no MCP server, no account.</p></div>
              </details>
              <details>
                <summary>Does it merge two agents into one history?<span className="pm" aria-hidden="true" /></summary>
                <div className="a"><p>No. One agent at a time. Hermes never writes into Cursor or Claude databases and cannot restore an old Cursor composer id. The next agent starts a new chat with the brief.</p></div>
              </details>
            </div>
          </div>
        </section>

        <section className="install dark" id="install">
          <Room className="room faint" />
          <div className="wrap install-in">
            <h2>Install in one line.</h2>
            <p className="lede">macOS and Linux. One binary. No npm, no plugin, no account.</p>
            <CommandChip className="cmd big" act="copy:install2" textKey="cp_install2" />
            <div className="then">
              <div className="then-row">
                <span>Then, in a new terminal</span>
                <CommandChip className="cmd sm" command="hermes doctor" act="copy:doctor" label="Copy hermes doctor" textKey="cp_doctor" />
              </div>
            </div>
          </div>
        </section>
      </main>
      <Footer page="home" />
    </>
  );
}
