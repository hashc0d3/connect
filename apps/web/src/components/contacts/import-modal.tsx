"use client";

import { useState } from "react";

export type ImportRow = { number: string; lname?: string; fname?: string; sname?: string; param1?: string; param2?: string; param3?: string };

export function ImportModal({ onClose, onDone }: { onClose: () => void; onDone: () => void }) {
  const [mode, setMode] = useState<"text" | "file">("file");
  const [text, setText] = useState("");
  const [rows, setRows] = useState<ImportRow[]>([]);
  const [name, setName] = useState("");

  async function readFile(file: File) {
    setName(file.name);
    if (file.name.endsWith(".csv") || file.name.endsWith(".txt")) {
      const raw = await file.text();
      setRows(parseLines(raw));
      return;
    }
    const XLSX = await import("xlsx");
    const book = XLSX.read(await file.arrayBuffer());
    const sheet = book.Sheets[book.SheetNames[0]];
    const table = XLSX.utils.sheet_to_json<string[]>(sheet, { header: 1 });
    setRows(table.slice(1).map((line) => ({
      number: String(line[0] || ""), lname: String(line[1] || ""), fname: String(line[2] || ""), sname: String(line[3] || ""),
      param1: String(line[5] || ""), param2: String(line[6] || ""), param3: String(line[7] || ""),
    })).filter((row) => row.number));
  }

  async function submit() {
    const contacts = mode === "text" ? parseLines(text) : rows;
    await fetch("/api/contacts/import", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ contacts }) });
    onDone();
  }

  return (
    <div className="modal-back">
      <div className="modal">
        <h3>Импорт контактов</h3>
        <label><input type="radio" checked={mode === "text"} onChange={() => setMode("text")} /> Текстом</label>
        <label><input type="radio" checked={mode === "file"} onChange={() => setMode("file")} /> Из файла</label>
        {mode === "text" ? <textarea className="area" value={text} onChange={(event) => setText(event.target.value)} placeholder="79990000000" /> : (
          <label className="btn">{name || "Выбрать"}<input hidden type="file" accept=".csv,.txt,.xlsx" onChange={(event) => event.target.files?.[0] && readFile(event.target.files[0])} /></label>
        )}
        <div className="row">
          <button className="btn green" onClick={submit}>Продолжить</button>
          <button className="btn" onClick={onClose}>Закрыть</button>
        </div>
      </div>
    </div>
  );
}

function parseLines(raw: string): ImportRow[] {
  return raw.split(/\r?\n/).map((line) => line.trim()).filter(Boolean).map((line) => {
    const [number, lname, fname, sname, , param1, param2, param3] = line.split(/[;,\t]/);
    return { number, lname, fname, sname, param1, param2, param3 };
  });
}
