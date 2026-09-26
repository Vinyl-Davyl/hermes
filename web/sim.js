(function () {
  var install = "curl -fsSL https://tryhermes.pages.dev/install | sh";

  document.querySelectorAll("[data-copy]").forEach(function (el) {
    var btn = el.querySelector(".copy");
    if (!btn) return;
    btn.addEventListener("click", function () {
      var text = el.getAttribute("data-copy") || install;
      navigator.clipboard.writeText(text).then(function () {
        btn.textContent = "Copied";
        setTimeout(function () { btn.textContent = "Copy"; }, 1200);
      });
    });
  });

  var out = document.getElementById("sim");
  var agent = document.getElementById("sim-agent");
  var beats = document.querySelectorAll("#beats li");
  if (!out || !agent) return;

  var scenes = [
    {
      name: "cursor",
      beat: 0,
      lines: [
        { cls: "you", text: "› fix the flaky auth test" },
        { cls: "ai", text: "● Race in the refresh path. Two goroutines hit it at once." },
        { cls: "ai", text: "● Guarding the refresh with a mutex." },
        { cls: "err", text: "✗ Limit reached. Session paused." }
      ]
    },
    {
      name: "hermes",
      beat: 1,
      lines: [
        { cls: "cmdl", text: "$ hermes handoff cursor antigravity" },
        { cls: "ok", text: "packed  ./handoff-20260926T153022" },
        { cls: "ok", text: "from    cursor → antigravity" },
        { cls: "ok", text: "copied  PROMPT.md" }
      ]
    },
    {
      name: "antigravity",
      beat: 2,
      lines: [
        { cls: "you", text: "› (pasted the Hermes pack)" },
        { cls: "ai", text: "● Mutex is in. Writing the test." },
        { cls: "ai", text: "● Auth test is green. Picking up where Cursor stopped." }
      ]
    }
  ];

  var reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  function setBeat(i) {
    beats.forEach(function (el, n) {
      el.classList.toggle("on", n === i);
    });
  }

  function sleep(ms) {
    return new Promise(function (resolve) { setTimeout(resolve, ms); });
  }

  function render(lines, caret) {
    out.innerHTML = lines.map(function (line) {
      return '<span class="' + line.cls + '">' + line.text + "</span>";
    }).join("\n") + (caret ? '<span class="caret"> </span>' : "");
  }

  function finalScene() {
    var last = scenes[scenes.length - 1];
    agent.textContent = last.name;
    setBeat(last.beat);
    render(last.lines, false);
  }

  async function typeLine(cls, text, done) {
    var shown = "";
    for (var i = 0; i < text.length; i++) {
      shown += text[i];
      render(done.concat([{ cls: cls, text: shown }]), true);
      await sleep(14);
    }
  }

  async function play() {
    if (reduce) {
      finalScene();
      return;
    }
    while (true) {
      for (var s = 0; s < scenes.length; s++) {
        var scene = scenes[s];
        agent.textContent = scene.name;
        setBeat(scene.beat);
        var done = [];
        render([], true);
        await sleep(280);
        for (var i = 0; i < scene.lines.length; i++) {
          await typeLine(scene.lines[i].cls, scene.lines[i].text, done);
          done.push(scene.lines[i]);
          render(done, true);
          await sleep(420);
        }
        await sleep(1100);
      }
    }
  }

  play();
})();
