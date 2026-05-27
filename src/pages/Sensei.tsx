import { Sparkles, Send } from "lucide-react";
import { useState } from "react";
import PageHeader from "@/components/PageHeader";
import { useUser } from "@/store/user";
import { haptic } from "@/lib/telegram";

interface Msg { role: "user" | "ai"; text: string; }

const SUGGESTIONS = [
  "Помоги составить план тренировок",
  "Дай задание на интеллект",
  "Как восстановить энергию?",
];

export default function Sensei() {
  const { user } = useUser();
  const [msgs, setMsgs] = useState<Msg[]>([
    { role: "ai", text: "Я твой Сенсей. Спроси что угодно — помогу составить путь." },
  ]);
  const [input, setInput] = useState("");

  const send = (text?: string) => {
    const t = (text ?? input).trim();
    if (!t) return;
    haptic("light");
    setMsgs((m) => [...m, { role: "user", text: t }]);
    setInput("");
    setTimeout(() => {
      setMsgs((m) => [...m, { role: "ai", text: "Хм. Я обдумаю это и предложу задание. (демо-ответ)" }]);
    }, 600);
  };

  return (
    <div className="flex h-[calc(100vh-12rem)] flex-col">
      <PageHeader
        title="Сенсей"
        subtitle="ИИ-наставник"
        action={
          <span className="stat-chip text-accent">
            <Sparkles className="h-3 w-3" /> {user?.sensei_requests ?? 0} запр.
          </span>
        }
      />

      <div className="glass-card flex-1 space-y-3 overflow-y-auto p-4">
        {msgs.map((m, i) => (
          <div key={i} className={"flex " + (m.role === "user" ? "justify-end" : "justify-start")}>
            <div
              className={
                "max-w-[85%] rounded-2xl px-4 py-2.5 text-sm " +
                (m.role === "user"
                  ? "bg-gradient-primary text-primary-foreground shadow-glow"
                  : "bg-muted text-foreground")
              }
            >
              {m.text}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-3 flex flex-wrap gap-2">
        {SUGGESTIONS.map((s) => (
          <button key={s} className="stat-chip" onClick={() => send(s)}>{s}</button>
        ))}
      </div>

      <div className="mt-3 flex gap-2">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && send()}
          placeholder="Спроси Сенсея..."
          className="flex-1 rounded-xl border border-border bg-muted px-4 py-3 text-sm outline-none focus:border-primary"
        />
        <button className="btn-primary px-4" onClick={() => send()}>
          <Send className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
}