"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";

const items = [
  ["/settings", "Общие настройки"],
  ["/settings/templates", "Шаблоны"],
  ["/settings/profile", "Реквизиты и договоры"],
  ["/settings/api", "API и gRPC"],
  ["/settings/blacklist", "Чёрный список"],
  ["/settings/auto", "Авторассылка"],
  ["/settings/integrations", "Интеграции"],
];

export default function SettingsLayout({ children }: { children: ReactNode }) {
  const path = usePathname();
  return (
    <div className="settings">
      <aside className="side-col">
        {items.map(([href, label]) => <Link key={href} className={`pill ${path === href ? "on" : ""}`} href={href}>{label}</Link>)}
      </aside>
      <section>{children}</section>
    </div>
  );
}
