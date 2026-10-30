import { gateway } from "@/lib/gateway";
import { requireUser } from "@/lib/guard";

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  try {
    const body = await request.json();
    const result = await gateway("/v2/sms/send", user.id, body);
    return Response.json(result);
  } catch (error) {
    return Response.json({ message: error instanceof Error ? error.message : "Ошибка отправки" }, { status: 422 });
  }
}
