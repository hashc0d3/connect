export default function GeneralSettingsPage() {
  return (
    <section>
      <h2 className="h">Общие настройки</h2>
      <p>Часовой пояс кабинета — Москва. Имя отправителя и договор проверяются перед вызовом EasySMS.</p>
      <p className="meta">Пока ключ EasySMS не задан, шлюз работает в учебном режиме и не обращается к поставщику.</p>
    </section>
  );
}
