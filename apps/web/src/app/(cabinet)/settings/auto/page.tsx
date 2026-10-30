"use client";

import { useEffect, useState } from "react";

export default function AutoPage() {
  const [enabled, setEnabled] = useState(false);
  const [signName, setSignName] = useState("");
  const [templateId, setTemplateId] = useState("");
  const [templates, setTemplates] = useState<{ id: number; name: string }[]>([]);
  const [signs, setSigns] = useState<{ name: string }[]>([]);
  useEffect(() => {
    fetch("/api/automation").then((response) => response.json()).then((rule) => {
      if (!rule) return;
      setEnabled(Boolean(rule.enabled));
      setSignName(rule.signName || "");
      setTemplateId(rule.templateId ? String(rule.templateId) : "");
    });
    fetch("/api/templates").then((response) => response.json()).then(setTemplates);
    fetch("/api/signs").then((response) => response.json()).then(setSigns);
  }, []);
  return (
    <section>
      <h2 className="h">Авторассылка на день рождения</h2>
      <label><input type="checkbox" checked={enabled} onChange={(event) => setEnabled(event.target.checked)} /> Включена</label>
      <select className="select" value={templateId} onChange={(event) => setTemplateId(event.target.value)}>
        <option value="">Шаблон</option>
        {templates.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
      </select>
      <select className="select" value={signName} onChange={(event) => setSignName(event.target.value)}>
        <option value="">Подпись</option>
        {signs.map((item) => <option key={item.name}>{item.name}</option>)}
      </select>
      <button className="btn green" onClick={() => fetch("/api/automation", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ enabled, templateId, signName }) })}>Сохранить</button>
    </section>
  );
}
