"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { previewText, smsParts } from "@/lib/text";

type Sign = { id: number; name: string; status: number };
type Group = { id: number; name: string; _count?: { members: number } };
type Template = { id: number; name: string; body: string };

const tokens = ["{Фамилия}", "{Имя}", "{Отчество}", "{Параметр1}", "{Параметр2}", "{Параметр3}"];

export function CampaignForm() {
  const [text, setText] = useState("");
  const [sign, setSign] = useState("");
  const [signs, setSigns] = useState<Sign[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [picked, setPicked] = useState<number[]>([]);
  const [templates, setTemplates] = useState<Template[]>([]);
  const [numbers, setNumbers] = useState("");
  const [translit, setTranslit] = useState(false);
  const [shortLink, setShortLink] = useState(false);
  const [when, setWhen] = useState<"now" | "later" | "smooth">("now");
  const [start, setStart] = useState("2026-10-30T12:00");
  const [end, setEnd] = useState("2026-10-30T18:00");
  const [message, setMessage] = useState("");
  const router = useRouter();

  useEffect(() => {
    fetch("/api/signs").then((response) => response.json()).then((rows: Sign[]) => {
      setSigns(rows);
      const active = rows.find((row) => row.status === 1);
      if (active) setSign(active.name);
    });
    fetch("/api/templates").then((response) => response.json()).then(setTemplates);
    fetch("/api/contacts?pageSize=20").then((response) => response.json()).then((payload) => setGroups(payload.groups || []));
  }, []);

  const shown = previewText(text, { translit, shortLink });
  const parts = useMemo(() => smsParts(shown), [shown]);

  async function send() {
    const dateSend = when === "now" ? 0 : Math.floor(new Date(`${start}:00+03:00`).getTime() / 1000);
    const response = await fetch("/api/campaigns", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        text, sign, translit, shortLink, groupIds: picked,
        numbers: numbers.split(/[\s,;]+/).filter(Boolean),
        dateSend,
        smoothUntil: when === "smooth" ? Math.floor(new Date(`${end}:00+03:00`).getTime() / 1000) : 0,
      }),
    });
    const payload = await response.json();
    setMessage(payload.message || payload.data?.message || "Готово");
    if (response.ok) router.refresh();
  }

  return (
    <div className="grid-2">
      <div>
        <div className="note">Перед рассылкой проверьте текст, подпись и баланс. Сообщение проходит правила Connect и только потом уходит в EasySMS.</div>
        <h2 className="h">Текст рассылки</h2>
        <div className="row">
          <select className="select" onChange={(event) => setText(templates.find((item) => String(item.id) === event.target.value)?.body || text)}>
            <option value="">Шаблон не выбран</option>
            {templates.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
          </select>
        </div>
        <div className="checks">
          <label><input type="checkbox" checked={translit} onChange={(event) => setTranslit(event.target.checked)} /> Преобразовать кириллицу в латиницу</label>
          <label><input type="checkbox" checked={shortLink} onChange={(event) => setShortLink(event.target.checked)} /> Сокращать ссылки</label>
        </div>
        <textarea className="area" value={text} onChange={(event) => setText(event.target.value)} />
        <div className="tokens">{tokens.map((token) => <button key={token} type="button" onClick={() => setText((value) => `${value}${token}`)}>{token}</button>)}</div>
        <p className="meta">Символов: {parts.length} / {parts.limit}. SMS: {parts.parts}</p>
        <h2 className="h">Получатели</h2>
        {groups.map((group) => (
          <label key={group.id} className="row">
            <input type="checkbox" checked={picked.includes(group.id)} onChange={(event) => setPicked((current) => event.target.checked ? [...current, group.id] : current.filter((id) => id !== group.id))} />
            {group.name} ({group._count?.members || 0})
          </label>
        ))}
        <div className="label">Быстрая загрузка номеров</div>
        <textarea className="area" value={numbers} onChange={(event) => setNumbers(event.target.value)} placeholder="79990000000" />
        <h2 className="h">Время отправки</h2>
        <label><input type="radio" checked={when === "now"} onChange={() => setWhen("now")} /> Немедленно</label>
        <label><input type="radio" checked={when === "later"} onChange={() => setWhen("later")} /> В запланированное время</label>
        <label><input type="checkbox" checked={when === "smooth"} onChange={(event) => setWhen(event.target.checked ? "smooth" : "later")} /> Плавная отправка</label>
        {when !== "now" ? (
          <div className="row">
            <input className="field" type="datetime-local" value={start} onChange={(event) => setStart(event.target.value)} />
            {when === "smooth" ? <input className="field" type="datetime-local" value={end} onChange={(event) => setEnd(event.target.value)} /> : null}
          </div>
        ) : null}
        <p className="meta">Время указывается по Москве.</p>
        <div className="label">Подпись отправителя</div>
        <select className="select" value={sign} onChange={(event) => setSign(event.target.value)}>
          {signs.map((item) => <option key={item.id} value={item.name}>{item.name} · {item.status === 1 ? "активна" : "на модерации"}</option>)}
        </select>
        <div className="row"><button className="btn green" onClick={send}>Отправить</button>{message ? <span>{message}</span> : null}</div>
      </div>
      <aside className="phone">
        <header><span>9:41</span><b>{sign || "Connect"}</b></header>
        <div className="bubble">{shown || "Текст сообщения появится здесь"}</div>
      </aside>
    </div>
  );
}
