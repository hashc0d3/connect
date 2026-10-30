import { many } from "@/lib/db";
import { requireUser } from "@/lib/guard";

type StatRow = { id: number; dateCreate: Date; fromSign: string; text: string; cost: string; extendStatus: string; status: number };

export async function GET(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const url = new URL(request.url);
  const from = url.searchParams.get("from");
  const to = url.searchParams.get("to");
  const text = url.searchParams.get("text") || "";
  const rows = await many<StatRow>(
    `SELECT id, date_create, from_sign, text, cost::text, extend_status, status
     FROM messages
     WHERE user_id = $1
       AND ($2::timestamptz IS NULL OR date_create >= $2)
       AND ($3::timestamptz IS NULL OR date_create <= $3)
       AND ($4 = '' OR position(lower($4) in lower(text)) > 0)
     ORDER BY id DESC LIMIT 300`,
    [user.id, from ? new Date(from) : null, to ? new Date(`${to}T23:59:59`) : null, text],
  );
  const months = new Map<string, { delivered: number; sent: number; undelivered: number; errors: number; total: number }>();
  for (const row of rows) {
    const key = new Date(row.dateCreate).toISOString().slice(0, 7);
    const bucket = months.get(key) || { delivered: 0, sent: 0, undelivered: 0, errors: 0, total: 0 };
    bucket.total += 1;
    if (row.status === 1) bucket.delivered += 1;
    else if (row.status === 3 || row.status === 0) bucket.sent += 1;
    else if (row.status === 2) bucket.undelivered += 1;
    else bucket.errors += 1;
    months.set(key, bucket);
  }
  return Response.json({
    summary: [...months.entries()].map(([month, value]) => ({ month, ...value })),
    rows,
  });
}
