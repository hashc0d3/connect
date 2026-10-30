"use client";

import { FormEvent, useEffect, useState } from "react";

type Entry = { id: number; number: string; comment: string };

export default function BlacklistPage() {
  const [rows, setRows] = useState<Entry[]>([]);
  const [number, setNumber] = useState("");
  function load() { fetch("/api/blacklist").then((response) => response.json()).then(setRows); }
  useEffect(load, []);
  async function add(event: FormEvent) {
    event.preventDefault();
    await fetch("/api/blacklist", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ number }) });
    setNumber("");
    load();
  }
  return (
    <section>
      <h2 className="h">Чёрный список</h2>
      <form className="row" onSubmit={add}><input className="field" value={number} onChange={(event) => setNumber(event.target.value)} placeholder="79990000000" /><button className="btn pink">Добавить</button></form>
      <table>
        <tbody>{rows.map((row) => <tr key={row.id}><td>{row.number}</td><td><button className="btn" onClick={async () => { await fetch(`/api/blacklist?id=${row.id}`, { method: "DELETE" }); load(); }}>Удалить</button></td></tr>)}</tbody>
      </table>
    </section>
  );
}
