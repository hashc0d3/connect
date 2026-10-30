import { redirect } from "next/navigation";
import type { ReactNode } from "react";
import { Shell } from "@/components/shell";
import { currentUser } from "@/lib/guard";

export default async function CabinetLayout({ children }: { children: ReactNode }) {
  const user = await currentUser();
  if (!user) redirect("/login");
  return (
    <Shell user={{ email: user.email, name: user.name, balance: user.balance.toString(), role: user.role }}>
      {children}
    </Shell>
  );
}
