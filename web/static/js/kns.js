/** KNS 4.0 browser client. Talks to this app's indexer proxy. */
export async function resolve(q, base = "") {
  const r = await fetch(base + "/api/v1/resolve?q=" + encodeURIComponent(q));
  const j = await r.json();
  if (!j.ok) throw new Error(j.error || "resolve failed");
  return j.data;
}

export async function names(owner, base = "") {
  const r = await fetch(base + "/api/v1/names?owner=" + encodeURIComponent(owner));
  const j = await r.json();
  if (!j.ok) throw new Error(j.error || "names failed");
  return j.data;
}

export async function batch(list, base = "") {
  const r = await fetch(base + "/api/v1/batch", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ names: list }),
  });
  return (await r.json()).data;
}

export function payUri(owner, name) {
  if (!owner) return "";
  return name ? owner + "?label=" + encodeURIComponent(name) : owner;
}

export async function agentCard(name, base = "") {
  const r = await fetch(base + "/agent/" + encodeURIComponent(name) + ".json");
  return r.json();
}

export async function mcpCall(tool, args, base = "") {
  const r = await fetch(base + "/mcp", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ jsonrpc: "2.0", id: 1, method: "tools/call", params: { name: tool, arguments: args } }),
  });
  return r.json();
}
