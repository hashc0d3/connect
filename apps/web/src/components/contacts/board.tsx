"use client";

import { useEffect, useState } from "react";
import { ImportModal } from "./import-modal";

type Group = { id: number; name: string; _count: { members: number } };
type Contact = { id: number; number: string; lname: string; fname: string; sname: string; operatorName: string; hlrStatus: number };

export function ContactsBoard() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [counts, setCounts] = useState({ all: 0, own: 0 });
  const [groupId, setGroupId] = useState(0);
  const [personal, setPersonal] = useState(false);
  const [q, setQ] = useState("");
  const [pageSize, setPageSize] = useState(20);
  const [name, setName] = useState("");
  const [adding, setAdding] = useState(false);
  const [importing, setImporting] = useState(false);
  const [picked, setPicked] = useState<number[]>([]);
  const [draft, setDraft] = useState({ number: "", fname: "", lname: "" });

  function load() {
    const params = new URLSearchParams({ q, pageSize: String(pageSize) });
    if (groupId) params.set("groupId", String(groupId));
    if (personal) params.set("personal", "1");
    fetch(`/api/contacts?${params}`).then((response) => response.json()).then((payload) => {
      setGroups(payload.groups || []);
      setContacts(payload.contacts || []);
      setCounts({ all: payload.all || 0, own: payload.own || 0 });
    });
  }

  useEffect(() => { load(); }, [groupId, personal, q, pageSize]);

  return (
    <div className="contacts">
      <aside className="side-col">
        <div className="row">
          <input className="field" value={name} placeholder="Новая группа" onChange={(event) => setName(event.target.value)} />
          <button className="btn green" onClick={async () => { await fetch("/api/groups", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name }) }); setName(""); load(); }}>✓</button>
        </div>
        <button className={`pill ${!groupId && !personal ? "on" : ""}`} onClick={() => { setGroupId(0); setPersonal(false); }}>Все контакты {counts.all}</button>
        <button className={`pill ${personal ? "on" : ""}`} onClick={() => { setPersonal(true); setGroupId(0); }}>Личные контакты {counts.own}</button>
        {groups.map((group) => (
          <button key={group.id} className={`pill ${groupId === group.id ? "on" : ""}`} onClick={() => { setGroupId(group.id); setPersonal(false); }}>{group.name} {group._count.members}</button>
        ))}
      </aside>
      <section style={{ flex: 1 }}>
        <div className="row">
          <button className="btn" onClick={() => setAdding(true)}>Добавить</button>
          <button className="btn" onClick={() => setImporting(true)}>Импорт</button>
          <button className="btn" onClick={async () => { await fetch("/api/contacts/hlr", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ ids: picked.length ? picked : contacts.map((item) => item.id) }) }); load(); }}>HLR</button>
          <button className="btn" onClick={async () => { await fetch(`/api/contacts?ids=${picked.join(",")}`, { method: "DELETE" }); load(); }}>Удалить</button>
          <input className="field search" placeholder="Поиск по контактам" value={q} onChange={(event) => setQ(event.target.value)} />
        </div>
        <table>
          <thead><tr><th></th><th>Номер</th><th>Фамилия</th><th>Имя</th><th>Отчество</th><th>Оператор</th><th>HLR</th></tr></thead>
          <tbody>
            {contacts.length === 0 ? <tr><td className="empty" colSpan={7}>В данной группе пока отсутствуют контакты.</td></tr> : contacts.map((contact) => (
              <tr key={contact.id}>
                <td><input type="checkbox" checked={picked.includes(contact.id)} onChange={(event) => setPicked((current) => event.target.checked ? [...current, contact.id] : current.filter((id) => id !== contact.id))} /></td>
                <td>{contact.number}</td><td>{contact.lname}</td><td>{contact.fname}</td><td>{contact.sname}</td><td>{contact.operatorName}</td><td>{contact.hlrStatus === 1 ? "доступен" : "не проверялся"}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <div className="pager">Отображать контакты по: {[20, 50, 100, 200].map((size) => <button key={size} className={size === pageSize ? "on" : ""} onClick={() => setPageSize(size)}>{size}</button>)}</div>
      </section>
      {importing ? <ImportModal onClose={() => setImporting(false)} onDone={() => { setImporting(false); load(); }} /> : null}
      {adding ? (
        <div className="modal-back">
          <form className="modal" onSubmit={async (event) => { event.preventDefault(); await fetch("/api/contacts", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ ...draft, groupId: groupId || undefined }) }); setAdding(false); load(); }}>
            <h3>Новый контакт</h3>
            <input className="field" placeholder="Номер" value={draft.number} onChange={(event) => setDraft({ ...draft, number: event.target.value })} required />
            <input className="field" placeholder="Фамилия" value={draft.lname} onChange={(event) => setDraft({ ...draft, lname: event.target.value })} />
            <input className="field" placeholder="Имя" value={draft.fname} onChange={(event) => setDraft({ ...draft, fname: event.target.value })} />
            <div className="row"><button className="btn green">Сохранить</button><button className="btn" type="button" onClick={() => setAdding(false)}>Закрыть</button></div>
          </form>
        </div>
      ) : null}
    </div>
  );
}
