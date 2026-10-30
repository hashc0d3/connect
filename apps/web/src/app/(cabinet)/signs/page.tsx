"use client";

import { FormEvent, useEffect, useState } from "react";

export default function SignsPage() {
  const [rows, setRows] = useState<{ id: number; name: string; status: number; extendStatus: string; channel: string }[]>([]);
  const [name, setName] = useState("");
  function load() { fetch("/api/signs").then((response) => response.json()).then(setRows); }
  useEffect(load, []);
  async function add(event: FormEvent) {
    event.preventDefault();
    await fetch("/api/signs", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name }) });
    setName("");
    load();
  }
  return (
    <section>
      <h2 className="h">Подписи отправителя</h2>
      <form className="row" onSubmit={add}><input className="field" value={name} onChange={(event) => setName(event.target.value)} placeholder="CONNECT" /><button className="btn green">Запросить</button></form>
      <table>
        <thead><tr><th>Имя</th><th>Канал</th><th>Статус</th></tr></thead>
        <tbody>{rows.map((row) => <tr key={row.id}><td>{row.name}</td><td>{row.channel}</td><td>{row.extendStatus}</td></tr>)}</tbody>
      </table>
    </section>
  );
}
