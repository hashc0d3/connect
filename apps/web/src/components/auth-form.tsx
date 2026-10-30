"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

export function AuthForm({ mode }: { mode: "login" | "register" }) {
  const router = useRouter();
  const [error, setError] = useState("");
  const [key, setKey] = useState("");
  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const response = await fetch(mode === "login" ? "/api/auth/login" : "/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: form.get("email"), password: form.get("password"), name: form.get("name") }),
    });
    const payload = await response.json();
    if (!response.ok) {
      setError(payload.message || "Ошибка");
      return;
    }
    if (payload.apiKey) setKey(payload.apiKey);
    else router.push("/campaigns");
  }
  return (
    <main className="auth-page">
      <form className="login" onSubmit={onSubmit}>
        <a className="brand" href="/">Connect<i>→</i></a>
        <h1>{mode === "login" ? "Вход" : "Регистрация"}</h1>
        {mode === "register" ? <input className="field" name="name" placeholder="Имя" required /> : null}
        <div className="label">Email</div>
        <input className="field" name="email" type="email" required />
        <div className="label">Пароль</div>
        <input className="field" name="password" type="password" minLength={8} required />
        {error ? <p className="error">{error}</p> : null}
        {key ? <p>Ключ API, сохраните его: <b>{key}</b></p> : null}
        <button className="btn green" type="submit">{mode === "login" ? "Войти" : "Создать аккаунт"}</button>
        {key ? <button className="btn" type="button" onClick={() => router.push("/campaigns")}>В кабинет</button> : null}
      </form>
    </main>
  );
}
