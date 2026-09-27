import { CommandChip, Footer, Header, Lights, Room } from "./chrome";
import { SiteBindings } from "./site-bindings";

function Out({ title, html }: { title: string; html: string }) {
  return (
    <div className="out">
      <div className="tbar"><Lights /><span className="ttl">{title}</span></div>
      <pre dangerouslySetInnerHTML={{ __html: html }} />
    </div>
  );
}

function Chip({ command }: { command: string }) {
  return <CommandChip command={command} copyAttr label="Copy" />;
}

export function DocsPage() {
  return (
    <>
      <SiteBindings />
      <Header page="docs" />
      <section className="docs-hero dark">
        <Room />
        <div className="wrap">
          <h1>Docs</h1>
          <p className="lede">Everything Hermes does, on one page. It all runs on your machine.</p>
          <div className="cta">
            <Chip command="curl -fsSL https://tryhermes.pages.dev/install | sh" />
          </div>
        </div>
      </section>

      <div className="docs light">
        <div className="wrap docs-in">
          <nav className="toc" aria-label="On this page">
            <p>Start</p>
            <a href="#quickstart" className="cur">Quickstart</a>
            <a href="#rule">The rule</a>
            <div className="gp">
              <p>Commands</p>
              <a href="#handoff">handoff</a>
              <a href="#destinations">Where it lands</a>
              <a href="#sessions">list and --id</a>
              <a href="#scope">Folder scope</a>
              <a href="#resume">resume</a>
              <a href="#doctor">doctor</a>
              <a href="#more">More commands</a>
            </div>
            <div className="gp">
              <p>Reference</p>
              <a href="#flags">Flags</a>
              <a href="#agents">Agents</a>
              <a href="#limits">What it won’t do</a>
            </div>
          </nav>

          <article className="doc">
            <section id="quickstart">
              <h2>Quickstart</h2>
              <p>macOS and Linux. One binary. No npm, no plugin, no account.</p>
              <div className="cmds"><Chip command="curl -fsSL https://tryhermes.pages.dev/install | sh" /></div>
              <p>Open a new terminal, then check what Hermes can see:</p>
              <div className="cmds"><Chip command="hermes doctor" /></div>
              <p>If <code>hermes</code> is not found, add the install folder to your PATH:</p>
              <div className="cmds"><Chip command={'export PATH="$HOME/.local/bin:$PATH"'} /></div>
              <p>Then, in the project you were editing:</p>
              <div className="cmds"><Chip command={'hermes handoff claude cursor -m "what you were doing"'} /></div>
            </section>

            <section id="rule">
              <h2>The rule</h2>
              <p>One name is where you are going. Two names are from, then to.</p>
              <div className="tbl"><table>
                <thead><tr><th>Command</th><th>Means</th></tr></thead>
                <tbody>
                  <tr><td><code>hermes handoff claude</code></td><td>Go to Claude</td></tr>
                  <tr><td><code>hermes handoff cursor</code></td><td>Go to Cursor. Not “leave Cursor”.</td></tr>
                  <tr><td><code>hermes handoff cursor claude</code></td><td>From Cursor to Claude</td></tr>
                  <tr><td><code>hermes handoff claude cursor</code></td><td>From Claude to Cursor</td></tr>
                </tbody>
              </table></div>
              <p>With one name, the source is the newest session that belongs to this folder. <code>hermes to &lt;agent&gt;</code> is the same as <code>hermes handoff &lt;agent&gt;</code>.</p>
            </section>

            <section id="handoff">
              <h2>handoff</h2>
              <p>Packs this repo and prepares the next agent.</p>
              <Out title="app · zsh" html={`<span class="c1">~/src/app</span> <span class="c5">main*</span> <span class="c2">❯</span> hermes handoff claude cursor -m "finish upload rate limit"
<span class="c3">packed</span> handoff-20260927T141205
<span class="c3">from</span> claude id=7c21a9f0-3b1e-4c2a-9e51-0d3c1b2a7f44
<span class="c3">to</span> cursor
<span class="c3">copied</span> PROMPT.md → clipboard

<span class="c3">next</span> hermes resume cursor`} />
              <p>The pack is a <code>handoff-&lt;timestamp&gt;/</code> folder next to your project:</p>
              <Out title="pack" html={`handoff-20260927T141205/
├── PROMPT.md            <span class="c3">the brief the next agent reads</span>
├── SUMMARY.md
├── manifest.json
├── git/
│   ├── branch.txt
│   ├── changed-files.txt
│   ├── diff.patch
│   └── log.txt
├── sessions/claude.md   <span class="c3">when Hermes can read the chat</span>
├── DECISIONS.md         <span class="c3">with --decisions</span>
└── files/               <span class="c3">with --include-files</span>`} />
              <div className="note warn"><b>Never type <code>./handoff-*</code>.</b> zsh expands the glob. Use <code>hermes resume</code>.</div>
            </section>

            <section id="destinations">
              <h2>Where it lands</h2>
              <div className="tbl"><table>
                <thead><tr><th>Destination</th><th>What happens</th></tr></thead>
                <tbody>
                  <tr><td>CLI agents on your PATH<br /><span className="muted">claude, codex, opencode…</span></td><td>Starts seeded from the pack: <code>opened claude (seeded from the pack)</code></td></tr>
                  <tr><td>Editors and GUIs<br /><span className="muted">Cursor, Antigravity, Claude desktop or web, VS Code chat</span></td><td><code>PROMPT.md</code> goes to your clipboard. Open a <b>new</b> chat and paste, or type <code>@PROMPT.md</code>.</td></tr>
                </tbody>
              </table></div>
              <p>Force either path with <code>--open</code>, <code>--no-open</code> or <code>--copy</code>. There is no <code>vscode</code> agent: for VS Code or the Claude app, use <code>--no-open --copy</code> and paste.</p>
              <p>In Cursor this is a new chat, not the old composer id. Continuity is the work, not the thread id.</p>
            </section>

            <section id="sessions">
              <h2>list and --id</h2>
              <p>Pick a specific chat instead of the newest one.</p>
              <div className="cmds">
                <Chip command="hermes list --here" />
                <Chip command="hermes list --agent cursor --here -q keyword" />
                <Chip command="hermes handoff cursor claude --id e188" />
              </div>
              <Out title="app · zsh" html={`<span class="c1">~/src/app</span> <span class="c5">main*</span> <span class="c2">❯</span> hermes list --agent cursor --here

# cursor

 2026-09-27 14:10 <span class="c4">e1884976-2f0c-4b7e-9a51-3d0c1b2a7f10</span> upload rate limit
 <span class="c3">~/Library/Application Support/Cursor/User/globalStorage/state.vscdb</span>
 2026-09-27 11:02 <span class="c4">b03d11e2-9c4a-4e0b-8f7d-61a2c3e4f5a6</span> fix flaky upload test
 <span class="c3">~/Library/Application Support/Cursor/User/globalStorage/state.vscdb</span>`} />
              <p>A prefix of the id is enough. Without <code>--id</code>, Hermes takes the newest session for this folder. <code>hermes list</code> without <code>--here</code> shows every session Hermes can see on this machine.</p>
            </section>

            <section id="scope">
              <h2>Folder scope</h2>
              <p>Hermes only packs chats that belong to the folder you are in. Another project’s chat on the same machine is ignored. If nothing matches, you still get a git-only pack and a warning:</p>
              <Out title="blog · zsh" html={`<span class="c1">~/src/blog</span> <span class="c5">main</span> <span class="c2">❯</span> hermes handoff codex
<span class="c3">packed</span> handoff-20260927T142030
<span class="c3">from</span> git only
<span class="c3">to</span> codex
<span class="c3">opened</span> codex (seeded from the pack)
<span class="c4">warn</span> no agent session for this folder; pack is git-only (this is fine). Use --id to pick a chat from another project

<span class="c3">next</span> hermes resume codex`} />
              <p>Pass <code>--id</code> to pick any chat that <code>hermes list</code> shows. For GUI chats Hermes cannot read, use <code>--from none</code> for a git-only pack.</p>
            </section>

            <section id="resume">
              <h2>resume</h2>
              <p>Reuses the latest pack. No glob.</p>
              <div className="cmds">
                <Chip command="hermes resume" />
                <Chip command="hermes resume cursor" />
                <Chip command="hermes resume --print" />
              </div>
              <p><code>--print</code> prints <code>PROMPT.md</code>, the brief the next agent reads.</p>
            </section>

            <section id="doctor">
              <h2>doctor</h2>
              <p>Shows what git sees and which agents Hermes can read on this machine. You do not need every agent installed to try Hermes.</p>
              <Out title="app · zsh" html={`<span class="c1">~/src/app</span> <span class="c5">main*</span> <span class="c2">❯</span> hermes doctor
hermes doctor

 <span class="c2">✓</span> git binary             ok
 <span class="c2">✓</span> git repo               ~/src/app (branch main, 3 changed files)
 <span class="c2">✓</span> sqlite3 (Cursor DBs)   ok
 <span class="c2">✓</span> claude                 3 session(s)
 <span class="c2">✓</span> cursor                 2 session(s)
 <span class="c4">✗</span> opencode               0 session(s)
 <span class="c2">✓</span> home                   /Users/you

Hermes never writes into Cursor/Claude databases.
You do not need every agent installed to test Hermes.`} />
            </section>

            <section id="more">
              <h2>More commands</h2>
              <div className="tbl"><table>
                <thead><tr><th>Command</th><th>Does</th></tr></thead>
                <tbody>
                  <tr><td><code>hermes to &lt;agent&gt;</code></td><td>Same as <code>hermes handoff &lt;agent&gt;</code></td></tr>
                  <tr><td><code>hermes from &lt;agent&gt;</code></td><td>Packs a session from one agent. The destination stays Cursor unless you name one.</td></tr>
                  <tr><td><code>hermes validate [pack]</code></td><td>Checks a handoff pack</td></tr>
                  <tr><td><code>hermes site</code></td><td>Serves this landing page locally</td></tr>
                  <tr><td><code>hermes version</code></td><td>Prints the version</td></tr>
                </tbody>
              </table></div>
            </section>

            <section id="flags">
              <h2>Flags</h2>
              <div className="tbl"><table>
                <thead><tr><th>Flag</th><th>What it does</th></tr></thead>
                <tbody>
                  <tr><td><code>-m, --message</code></td><td>What you were doing</td></tr>
                  <tr><td><code>--from</code></td><td>Source agent: <code>auto</code>, <code>none</code> or an agent id</td></tr>
                  <tr><td><code>--to</code></td><td>Destination agent (default cursor)</td></tr>
                  <tr><td><code>--id</code></td><td>Session id from <code>hermes list</code> (a prefix is enough)</td></tr>
                  <tr><td><code>--session</code></td><td>Explicit session file</td></tr>
                  <tr><td><code>--include-files</code></td><td>Copy changed files into the pack</td></tr>
                  <tr><td><code>--decisions</code></td><td>Notes the next agent must not redo</td></tr>
                  <tr><td><code>-o, --output</code></td><td>Output directory</td></tr>
                  <tr><td><code>--zip</code></td><td>Also write a .zip</td></tr>
                  <tr><td><code>--copy</code></td><td>Copy PROMPT.md to the clipboard</td></tr>
                  <tr><td><code>--open</code></td><td>Start the destination CLI (default when it is on PATH)</td></tr>
                  <tr><td><code>--no-open</code></td><td>Do not start the next agent; only write the pack</td></tr>
                  <tr><td><code>--print</code></td><td>Print PROMPT.md</td></tr>
                  <tr><td><code>--agent</code> <span className="muted">list</span></td><td>Filter by agent id</td></tr>
                  <tr><td><code>-q, --query</code> <span className="muted">list</span></td><td>Search title, id, or path</td></tr>
                  <tr><td><code>--here</code> <span className="muted">list</span></td><td>Prefer sessions that look like this directory</td></tr>
                </tbody>
              </table></div>
              <p><code>--from auto</code> prefers a Claude session, then the newest. <code>--from none</code> skips every reader and packs git only.</p>
            </section>

            <section id="agents">
              <h2>Agents</h2>
              <p>Each one can be a source, a destination, or both. Hermes reads sessions and never writes into editor databases.</p>
              <div className="tbl"><table>
                <thead><tr><th>Id</th><th>Product</th><th>Reads from</th></tr></thead>
                <tbody>
                  <tr><td><code>claude</code></td><td>Claude Code</td><td><code>~/.claude/projects</code></td></tr>
                  <tr><td><code>cursor</code></td><td>Cursor</td><td><code>state.vscdb</code>, read-only</td></tr>
                  <tr><td><code>codex</code></td><td>Codex CLI</td><td><code>~/.codex/sessions</code></td></tr>
                  <tr><td><code>opencode</code></td><td>OpenCode</td><td><code>~/.local/share/opencode</code></td></tr>
                  <tr><td><code>cline</code></td><td>Cline / Roo</td><td><code>~/.cline</code>, editor globalStorage</td></tr>
                  <tr><td><code>kimi</code></td><td>Kimi</td><td><code>~/.kimi-code</code>, <code>~/.kimi</code></td></tr>
                  <tr><td><code>antigravity</code> · <code>agy</code></td><td>Antigravity</td><td><code>~/.gemini/antigravity-cli</code></td></tr>
                  <tr><td><code>pi</code></td><td>Pi Agent</td><td><code>~/.pi/agent</code></td></tr>
                  <tr><td><code>copilot</code></td><td>Copilot CLI</td><td><code>~/.copilot</code></td></tr>
                  <tr><td><code>zcode</code></td><td>ZCode</td><td><code>~/.zcode/cli/db</code>, read-only</td></tr>
                  <tr><td><code>deepseek</code> · <code>dsh</code></td><td>DeepSeek Harness</td><td><code>~/.dsh</code></td></tr>
                </tbody>
              </table></div>
              <p>Aliases also work: <code>claude-code</code>, <code>codex-cli</code>, <code>open-code</code>.</p>
            </section>

            <section id="limits">
              <h2>What it won’t do</h2>
              <ul>
                <li>Merge two agents into one native history</li>
                <li>Write into Cursor or Claude databases</li>
                <li>Upload anything, run an MCP server, or ask for an account</li>
                <li>Restore a Cursor composer id</li>
              </ul>
              <p>The work travels. The agents do not merge.</p>
            </section>
          </article>
        </div>
      </div>
      <Footer page="docs" />
    </>
  );
}
