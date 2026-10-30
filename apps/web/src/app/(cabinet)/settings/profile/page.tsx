"use client";

import { FormEvent, useEffect, useState } from "react";

const types = [
  ["legal", "Юридическое лицо"],
  ["entrepreneur", "Индивидуальный предприниматель"],
  ["person", "Физическое лицо"],
];

export default function ProfilePage() {
  const [rows, setRows] = useState<{ id: number; title: string; type: string; status: string }[]>([]);
  const [open, setOpen] = useState(false);
  const [type, setType] = useState("legal");
  function load() { fetch("/api/profiles").then((response) => response.json()).then(setRows); }
  useEffect(load, []);
  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    await fetch("/api/profiles", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ type, title: form.get("title"), inn: form.get("inn"), address: form.get("address") }) });
    setOpen(false);
    load();
  }
  return (
    <section>
      <button className="btn pink" onClick={() => setOpen(true)}>Добавить профиль</button>
      <table>
        <thead><tr><th>Платёжный профиль</th><th>Тип</th><th>Статус</th></tr></thead>
        <tbody>{rows.length === 0 ? <tr><td className="empty" colSpan={3}>Добавьте платёжный профиль для заключения договора.</td></tr> : rows.map((row) => <tr key={row.id}><td>{row.title}</td><td>{row.type}</td><td>{row.status}</td></tr>)}</tbody>
      </table>
      {open ? (
        <div className="modal-back"><form className="modal" onSubmit={save}>
          <h3>Тип профиля</h3>
          {types.map(([value, label]) => <label key={value}><input type="radio" checked={type === value} onChange={() => setType(value)} /> {label}</label>)}
          <input className="field" name="title" placeholder="Название" required />
          <input className="field" name="inn" placeholder="ИНН" />
          <input className="field" name="address" placeholder="Адрес" />
          <button className="btn green">Продолжить</button>
        </form></div>
      ) : null}
    </section>
  );
}
