"use client";

import { useEffect, useState } from "react";

type Queue = { messages: { id: number; text: string; number: string }[]; signs: { id: number; name: string }[]; profiles: { id: number; title: string }[] };

export default function AdminPage() {
  const [data, setData] = useState<Queue>({ messages: [], signs: [], profiles: [] });
  function load() { fetch("/api/admin").then((response) => response.json()).then(setData); }
  useEffect(load, []);
  async function act(body: object) {
    await fetch("/api/admin", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
    load();
  }
  return (
    <section>
      <h2 className="h">Модерация</h2>
      <h3>Сообщения</h3>
      {data.messages?.map((row) => <div className="row" key={row.id}><span>{row.number}: {row.text}</span><button className="btn green" onClick={() => act({ kind: "message", id: row.id, approve: true })}>Выпустить</button><button className="btn" onClick={() => act({ kind: "message", id: row.id, approve: false })}>Отклонить</button></div>)}
      <h3>Подписи</h3>
      {data.signs?.map((row) => <div className="row" key={row.id}><span>{row.name}</span><button className="btn green" onClick={() => act({ kind: "sign", id: row.id, approve: true })}>Одобрить</button></div>)}
      <h3>Профили</h3>
      {data.profiles?.map((row) => <div className="row" key={row.id}><span>{row.title}</span><button className="btn green" onClick={() => act({ kind: "profile", id: row.id, approve: true })}>Одобрить</button></div>)}
    </section>
  );
}
