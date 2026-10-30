import { gateway } from "@/lib/gateway";
import { many, one } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const [cards, invoices] = await Promise.all([
    many(`SELECT id, brand, last4 FROM cards WHERE user_id = $1 ORDER BY id DESC`, [user.id]),
    many(`SELECT id, amount::text, bonus::text, status, created_at, comment FROM invoices WHERE user_id = $1 ORDER BY id DESC LIMIT 30`, [user.id]),
  ]);
  return Response.json({ balance: user.balance, cards, invoices });
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  if (body.action === "card") {
    const last4 = String(body.last4 || "").replace(/\D/g, "").slice(-4);
    if (last4.length !== 4) return Response.json({ message: "Укажите последние 4 цифры" }, { status: 422 });
    const card = await one(
      `INSERT INTO cards (user_id, brand, last4) VALUES ($1,$2,$3) RETURNING id, brand, last4`,
      [user.id, String(body.brand || "Card"), last4],
    );
    return Response.json(card);
  }
  try {
    const result = await gateway("/v2/balance/add", user.id, { sum: String(body.sum), promo: body.promo || "", cardId: Number(body.cardId || 0) });
    return Response.json(result);
  } catch (error) {
    return Response.json({ message: error instanceof Error ? error.message : "Ошибка пополнения" }, { status: 422 });
  }
}
