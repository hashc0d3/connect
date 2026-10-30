"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

type Card = { id: number; brand: string; last4: string };
type Invoice = { id: number; amount: string; bonus: string; status: string; createdAt: string; comment: string };

export default function BillingPage() {
  const [balance, setBalance] = useState("0");
  const [cards, setCards] = useState<Card[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [sum, setSum] = useState("1000");
  const [promo, setPromo] = useState("");
  const [cardId, setCardId] = useState(0);
  const [note, setNote] = useState("");
  const router = useRouter();
  function load() {
    fetch("/api/billing").then((response) => response.json()).then((payload) => {
      setBalance(payload.balance);
      setCards(payload.cards || []);
      setInvoices(payload.invoices || []);
    });
  }
  useEffect(load, []);
  async function topUp(event: FormEvent) {
    event.preventDefault();
    const response = await fetch("/api/billing", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ sum, promo, cardId }) });
    const payload = await response.json();
    setNote(payload.message || "Счёт пополнен");
    load();
    router.refresh();
  }
  return (
    <section>
      <h2 className="h">Счета</h2>
      <p>Баланс: <b>{balance} ₽</b></p>
      <p className="meta">Пополнение меняет внутренний счёт Connect. Платёжный провайдер подключается отдельно, номер карты целиком не сохраняется.</p>
      <form className="row" onSubmit={topUp}>
        <input className="field" value={sum} onChange={(event) => setSum(event.target.value)} />
        <input className="field" placeholder="Промокод" value={promo} onChange={(event) => setPromo(event.target.value)} />
        <select className="select" value={cardId} onChange={(event) => setCardId(Number(event.target.value))}>
          <option value={0}>Без карты</option>
          {cards.map((card) => <option key={card.id} value={card.id}>{card.brand} **** {card.last4}</option>)}
        </select>
        <button className="btn green">Пополнить</button>
      </form>
      <form className="row" onSubmit={async (event) => {
        event.preventDefault();
        const form = new FormData(event.currentTarget);
        await fetch("/api/billing", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ action: "card", brand: form.get("brand"), last4: form.get("last4") }) });
        load();
      }}>
        <input className="field" name="brand" placeholder="Visa" />
        <input className="field" name="last4" placeholder="Последние 4 цифры" maxLength={4} />
        <button className="btn">Добавить метку карты</button>
      </form>
      {note ? <p>{note}</p> : null}
      <table>
        <thead><tr><th>Дата</th><th>Сумма</th><th>Бонус</th><th>Статус</th></tr></thead>
        <tbody>{invoices.map((row) => <tr key={row.id}><td>{new Date(row.createdAt).toLocaleString("ru-RU")}</td><td>{row.amount}</td><td>{row.bonus}</td><td>{row.status}</td></tr>)}</tbody>
      </table>
    </section>
  );
}
