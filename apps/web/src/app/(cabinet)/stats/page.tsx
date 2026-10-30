"use client";

import { useEffect, useState } from "react";

type Summary = { month: string; delivered: number; sent: number; undelivered: number; errors: number; total: number };
type Row = { id: number; dateCreate: string; fromSign: string; text: string; cost: string; extendStatus: string };

export default function StatsPage() {
  const [tab, setTab] = useState<"summary" | "detail">("summary");
  const [summary, setSummary] = useState<Summary[]>([]);
  const [rows, setRows] = useState<Row[]>([]);
  const [text, setText] = useState("");
  function load() {
    fetch(`/api/stats?text=${encodeURIComponent(text)}`).then((response) => response.json()).then((payload) => {
      setSummary(payload.summary || []);
      setRows(payload.rows || []);
    });
  }
  useEffect(load, [text]);
  return (
    <section>
      <div className="stats-switch">
        <button className={`btn ${tab === "summary" ? "pink" : ""}`} onClick={() => setTab("summary")}>Общая</button>
        <button className={`btn ${tab === "detail" ? "pink" : ""}`} onClick={() => setTab("detail")}>Детальная</button>
      </div>
      {tab === "summary" ? (
        <table>
          <thead><tr><th>Месяц</th><th>Доставлено</th><th>Передано</th><th>Не доставлено</th><th>Ошибок</th><th>Всего</th></tr></thead>
          <tbody>{summary.length === 0 ? <tr><td className="empty" colSpan={6}>Вы пока не делали рассылок.</td></tr> : summary.map((row) => <tr key={row.month}><td>{row.month}</td><td>{row.delivered}</td><td>{row.sent}</td><td>{row.undelivered}</td><td>{row.errors}</td><td>{row.total}</td></tr>)}</tbody>
        </table>
      ) : (
        <>
          <input className="field" placeholder="Поиск по тексту" value={text} onChange={(event) => setText(event.target.value)} />
          <table>
            <thead><tr><th>Дата</th><th>Отправитель</th><th>Текст</th><th>Стоимость</th><th>Статус</th></tr></thead>
            <tbody>{rows.length === 0 ? <tr><td className="empty" colSpan={5}>Вы пока не делали рассылок.</td></tr> : rows.map((row) => <tr key={row.id}><td>{new Date(row.dateCreate).toLocaleString("ru-RU")}</td><td>{row.fromSign}</td><td>{row.text}</td><td>{row.cost}</td><td>{row.extendStatus}</td></tr>)}</tbody>
          </table>
        </>
      )}
    </section>
  );
}
