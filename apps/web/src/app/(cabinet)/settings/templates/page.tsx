"use client";

import { FormEvent, useEffect, useState } from "react";

type Template = { id: number; name: string; body: string };

export default function TemplatesPage() {
  const [rows, setRows] = useState<Template[]>([]);
  const [name, setName] = useState("");
  const [body, setBody] = useState("");
  function load() { fetch("/api/templates").then((response) => response.json()).then(setRows); }
  useEffect(load, []);
  async function save(event: FormEvent) {
    event.preventDefault();
    await fetch("/api/templates", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name, body }) });
    setName("");
    setBody("");
    load();
  }
  return (
    <section>
      <h2 className="h">Шаблоны</h2>
      <form onSubmit={save}>
        <input className="field" placeholder="Название" value={name} onChange={(event) => setName(event.target.value)} required />
        <textarea className="area" value={body} onChange={(event) => setBody(event.target.value)} required />
        <button className="btn green">Сохранить</button>
      </form>
      {rows.map((row) => (
        <article key={row.id}>
          <b>{row.name}</b>
          <p>{row.body}</p>
          <button className="btn" onClick={async () => { await fetch("/api/templates", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ id: row.id, remove: true }) }); load(); }}>Удалить</button>
        </article>
      ))}
    </section>
  );
}
