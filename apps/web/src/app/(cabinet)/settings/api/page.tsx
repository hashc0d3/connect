"use client";

import { useEffect, useState } from "react";

type Key = { id: number; name: string; prefix: string; createdAt: string };

export default function ApiPage() {
  const [keys, setKeys] = useState<Key[]>([]);
  const [fresh, setFresh] = useState("");
  function load() { fetch("/api/keys").then((response) => response.json()).then(setKeys); }
  useEffect(load, []);
  return (
    <section className="doc">
      <h2 className="h">API и gRPC</h2>
      <p>REST повторяет сценарии SMS Aero и доступен на <code>http://localhost:8080/v2</code>. gRPC слушает <code>localhost:50051</code>, сервис <code>connect.v1.SmsGateway</code>.</p>
      <p>Авторизация: <code>Authorization: Bearer ключ</code> или Basic <code>email:ключ</code>. Кабинет ходит во внутренний шлюз отдельно, ключ наружу не нужен.</p>
      <button className="btn green" onClick={async () => {
        const response = await fetch("/api/keys", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name: "Ключ" }) });
        const payload = await response.json();
        setFresh(payload.apiKey);
        load();
      }}>Выпустить ключ</button>
      {fresh ? <p>Сохраните ключ, он больше не покажется: <b>{fresh}</b></p> : null}
      <table>
        <thead><tr><th>Имя</th><th>Префикс</th><th></th></tr></thead>
        <tbody>{keys.map((key) => <tr key={key.id}><td>{key.name}</td><td>{key.prefix}</td><td><button className="btn" onClick={async () => { await fetch(`/api/keys?id=${key.id}`, { method: "DELETE" }); load(); }}>Удалить</button></td></tr>)}</tbody>
      </table>
    </section>
  );
}
