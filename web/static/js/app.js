(function () {
  const $ = (s, r = document) => r.querySelector(s);

  document.querySelectorAll("[data-copy]").forEach((el) => {
    el.addEventListener("click", () => {
      navigator.clipboard.writeText(el.getAttribute("data-copy") || "");
      const prev = el.textContent;
      el.textContent = "Copied";
      setTimeout(() => { el.textContent = prev; }, 1200);
    });
  });

  const ins = $("#inscribeBtn");
  if (ins) {
    ins.addEventListener("click", async () => {
      const w = window.kasware;
      const payload = ins.getAttribute("data-payload");
      if (!w || !w.buildScript) {
        await navigator.clipboard.writeText(payload || "");
        alert("Payload copied. Broadcast from app.knsdomains.org — this app does not submit the reveal.");
        return;
      }
      try {
        const { p2shAddress } = await w.buildScript({ type: "KNS", data: payload });
        alert("Commit P2SH " + p2shAddress + "\nReveal output 0 must pay the KNS fee address. This UI does not broadcast.");
      } catch (e) {
        alert(e.message || String(e));
      }
    });
  }

  const gw = $("#gwForm");
  if (gw) {
    gw.addEventListener("submit", (e) => {
      e.preventDefault();
      const n = ($("#gwName")?.value || "kns").replace(/\.kas$/i, "");
      location.href = "/site/" + n + ".kas";
    });
  }

  const q = $("#q");
  const box = $("#suggest");
  let t = 0;
  if (q && box) {
    q.addEventListener("input", () => {
      clearTimeout(t);
      const v = q.value.trim();
      if (v.length < 1) {
        box.hidden = true;
        return;
      }
      t = setTimeout(async () => {
        try {
          const r = await fetch("/api/v1/suggest?q=" + encodeURIComponent(v));
          const j = await r.json();
          if (!j.ok || !j.data) {
            box.hidden = true;
            return;
          }
          const d = j.data;
          const line = d.available
            ? d.name + " is free on the indexer · " + d.priceKas + " KAS"
            : d.name + " taken · " + (d.owner || "").slice(0, 18) + "…";
          box.innerHTML = "<a href=\"/app?q=" + encodeURIComponent(d.name) + "\">" + line + "</a>";
          box.hidden = false;
        } catch {
          box.hidden = true;
        }
      }, 280);
    });
  }

  async function simPost(body) {
    const r = await fetch("/api/v4/sim", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(body),
    });
    const j = await r.json();
    const err = document.getElementById("simErr");
    if (!j.ok) {
      if (err) err.textContent = j.error || "failed";
      return;
    }
    location.reload();
  }
  function bindSim(id, build) {
    const f = document.getElementById(id);
    if (!f) return;
    f.addEventListener("submit", (e) => {
      e.preventDefault();
      const fd = new FormData(f);
      const o = {};
      fd.forEach((v, k) => { o[k] = String(v); });
      simPost(build(o));
    });
  }
  bindSim("simReg", (o) => ({ op: "register", name: o.name, owner: o.owner }));
  bindSim("simSpawn", (o) => ({ op: "spawn", name: o.name, actor: o.actor, child: o.child, childOwner: o.childOwner }));
  bindSim("simVault", (o) => ({ op: "bind_vault", name: o.name, actor: o.actor, commit: o.commit }));

  const rankBox = document.getElementById("rankLive");
  if (rankBox && rankBox.dataset.name) {
    fetch("/api/v1/rank?q=" + encodeURIComponent(rankBox.dataset.name))
      .then((r) => r.json())
      .then((j) => {
        if (!j.ok || !j.rank) return;
        rankBox.textContent = j.rank.glyph + " " + j.rank.title + " · " + Number(j.rank.kas).toFixed(2) + " KAS";
      })
      .catch(() => {});
  }

  const pf = document.getElementById("postForm");
  if (pf) {
    pf.addEventListener("submit", async (e) => {
      e.preventDefault();
      const fd = new FormData(pf);
      const body = { name: fd.get("name"), owner: fd.get("owner"), text: fd.get("text") };
      const r = await fetch("/api/v1/kaposts", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(body) });
      const j = await r.json();
      if (!j.ok) { alert(j.error || "fail"); return; }
      location.reload();
    });
  }
})();
