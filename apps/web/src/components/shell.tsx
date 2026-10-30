"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, type ReactNode } from "react";

function money(value: string) {
  const amount = Number(value);
  if (!Number.isFinite(amount)) return value;
  return amount.toLocaleString("ru-RU", { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

const links = [
  ["/campaigns", "Рассылки"],
  ["/contacts", "Контакты"],
  ["/signs", "Подписи"],
  ["/stats", "Статистика"],
  ["/billing", "Счета"],
  ["/promos", "Акции"],
  ["/settings", "Настройки"],
];

export function Shell({
  user,
  children,
}: {
  user: { email: string; name: string; balance: string; role: string };
  children: ReactNode;
}) {
  const path = usePathname();
  const [open, setOpen] = useState(false);
  const [text, setText] = useState("");
  return (
    <div className="shell">
      <header className="top">
        <Link className="brand" href="/campaigns">Connect<i>→</i></Link>
        <div className="promo">Дарим 10% от суммы платежа при пополнении от 3 000 ₽ или по коду CONNECT10</div>
        <div className="account">
          <span>{user.email}</span>
          <div className="balance"><span>Баланс</span><b>{money(user.balance)} ₽</b></div>
          <Link className="btn green" href="/billing">Пополнить баланс</Link>
          {user.role === "admin" ? <Link href="/admin">Админ</Link> : null}
        </div>
      </header>
      <nav className="nav">
        {links.map(([href, label]) => (
          <Link key={href} href={href} className={path.startsWith(href) ? "active" : ""}>{label}</Link>
        ))}
      </nav>
      <main className="page">{children}</main>
      <footer className="footer">
        <div>
          <div>© {new Date().getFullYear()} Connect</div>
          <Link href="/docs">Документация API</Link>
          <button className="btn" onClick={() => setOpen(true)}>Пожаловаться руководителю</button>
        </div>
        <div>support@connect.local</div>
      </footer>
      {open ? (
        <div className="modal-back">
          <form className="modal" onSubmit={async (event) => {
            event.preventDefault();
            await fetch("/api/feedback", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ text }) });
            setOpen(false);
            setText("");
          }}>
            <h3>Сообщение руководителю</h3>
            <textarea className="area" value={text} onChange={(event) => setText(event.target.value)} required />
            <div className="row"><button className="btn green">Отправить</button><button type="button" className="btn" onClick={() => setOpen(false)}>Закрыть</button></div>
          </form>
        </div>
      ) : null}
    </div>
  );
}
