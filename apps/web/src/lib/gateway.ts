export async function gateway(path: string, userId: string, body?: unknown) {
  const response = await fetch(`${process.env.GATEWAY_HTTP_URL || "http://localhost:8080"}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Internal-Token": process.env.INTERNAL_TOKEN || "dev-internal-token",
      "X-User-Id": userId,
    },
    body: body === undefined ? undefined : JSON.stringify(body),
    cache: "no-store",
  });
  const payload = (await response.json()) as { success?: boolean; message?: string; data?: unknown };
  if (!response.ok || payload.success === false) {
    throw new Error(payload.message || "Шлюз отклонил запрос");
  }
  return payload;
}
