import { StarIcon } from "./icons";

const GH = "https://github.com/Vinyl-Davyl/hermes";

export function Header({ page }: { page: "home" | "docs" }) {
  return (
    <header className="nav">
      <div className="wrap nav-in">
        <a className="brand" href="/" aria-label="Hermes home">
          <img src="/favicon.svg" alt="" width={28} height={28} />
          <span>hermes</span>
        </a>
        <nav aria-label="Main">
          <a href="/docs" aria-current={page === "docs" ? "page" : undefined}>Docs</a>
          <a href={page === "home" ? "#install" : "/#install"}>Install</a>
          <a className="gh" href={GH} target="_blank" rel="noopener noreferrer">
            <StarIcon />
            <span>GitHub</span>
          </a>
        </nav>
      </div>
    </header>
  );
}

export function Footer({ page }: { page: "home" | "docs" }) {
  return (
    <footer className="foot dark">
      <div className="wrap">
        <div className="foot-top">
          <div className="foot-brand">
            <img src="/favicon.svg" alt="" width={48} height={48} />
            <div>
              <p className="nm">hermes</p>
              <p>Named for the messenger of the gods. The work travels. The agents do not merge.</p>
            </div>
          </div>
          <div className="foot-actions">
            <a className="btn cream" href={GH} target="_blank" rel="noopener noreferrer">
              <StarIcon />
              Star on GitHub
            </a>
            {page === "home" ? (
              <a className="btn ghost" href="/docs">Read the docs</a>
            ) : (
              <a className="btn ghost" href="/">Back to home</a>
            )}
          </div>
        </div>
        <div className="foot-meta">
          <span>MIT · local-first messenger for coding agents</span>
          <nav aria-label="Footer">
            {page === "home" ? (
              <>
                <a href="/docs">Docs</a>
                <a href="#install">Install</a>
              </>
            ) : (
              <>
                <a href="/">Home</a>
                <a href="/#install">Install</a>
              </>
            )}
            <a href={GH} target="_blank" rel="noopener noreferrer">GitHub</a>
          </nav>
        </div>
      </div>
      {page === "home" ? <div className="wordmark" aria-hidden="true">hermes</div> : null}
    </footer>
  );
}

export function Room({ className = "room" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 1440 820" preserveAspectRatio="xMidYMin slice" aria-hidden="true">
      <path d="M392 96H1048V504H392Z M254.2 10.3H1185.8V589.7H254.2Z M47.6 -118.2H1392.4V718.2H47.6Z M-264 -312H1704V912H-264Z M392 96L-1248 -924 M392 504L-1248 1524 M501.3 96L-592 -924 M501.3 504L-592 1524 M610.7 96L64 -924 M610.7 504L64 1524 M720 96L720 -924 M720 504L720 1524 M829.3 96L1376 -924 M829.3 504L1376 1524 M938.7 96L2032 -924 M938.7 504L2032 1524 M1048 96L2688 -924 M1048 504L2688 1524 M392 198L-1248 -312 M1048 198L2688 -312 M392 300L-1248 300 M1048 300L2688 300 M392 402L-1248 912 M1048 402L2688 912" />
    </svg>
  );
}

export function Lights() {
  return (
    <span className="lights" aria-hidden="true"><i /><i /><i /></span>
  );
}

const INSTALL = "curl -fsSL https://tryhermes.pages.dev/install | sh";

export function CommandChip({
  command = INSTALL,
  act,
  label = "Copy install command",
  textKey,
  copyAttr,
  className = "cmd",
}: {
  command?: string;
  act?: string;
  label?: string;
  textKey?: string;
  copyAttr?: boolean;
  className?: string;
}) {
  return (
    <div className={className}>
      <span className="p">$</span>
      <code>{command}</code>
      <button
        type="button"
        className="copy"
        data-act={act}
        data-copy={copyAttr ? "" : undefined}
        aria-label={label}
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <rect x="9" y="9" width="12" height="12" rx="2" />
          <path d="M5 15H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v1" />
        </svg>
        <span data-text={textKey}>Copy</span>
      </button>
    </div>
  );
}
