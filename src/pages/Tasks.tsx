import { useEffect, useState } from "react";
import { Plus, CheckCircle2, Play, Flame, Brain, Dumbbell, Eye } from "lucide-react";
import { motion } from "framer-motion";
import PageHeader from "@/components/PageHeader";
import { api } from "@/lib/api";
import { haptic } from "@/lib/telegram";
import { useUser } from "@/store/user";

const DEMO_TASKS = [
  { id: 1, title: "Утренняя пробежка", description: "5 км в умеренном темпе", task_type: "strength", xp_reward: 80, gold_reward: 25, energy_cost: 15, status: "active", difficulty: 3 },
  { id: 2, title: "Прочитать главу книги", description: "Атомные привычки, гл. 4", task_type: "intelligence", xp_reward: 50, gold_reward: 15, energy_cost: 5, status: "active", difficulty: 2 },
  { id: 3, title: "Медитация 15 мин", description: "Сфокусироваться на дыхании", task_type: "insight", xp_reward: 40, gold_reward: 10, energy_cost: 5, status: "in_progress", difficulty: 2 },
  { id: 4, title: "Тренировка ловкости", description: "Скакалка 10 минут", task_type: "agility", xp_reward: 60, gold_reward: 20, energy_cost: 10, status: "active", difficulty: 3 },
];

const typeMeta: Record<string, { icon: any; label: string; color: string }> = {
  strength: { icon: Dumbbell, label: "Сила", color: "text-danger" },
  agility: { icon: Flame, label: "Ловкость", color: "text-warning" },
  intelligence: { icon: Brain, label: "Интеллект", color: "text-primary-glow" },
  insight: { icon: Eye, label: "Интуиция", color: "text-accent" },
};

export default function Tasks() {
  const [tasks, setTasks] = useState<any[]>(DEMO_TASKS);
  const [openNew, setOpenNew] = useState(false);
  const { refresh } = useUser();

  useEffect(() => {
    api.getTasks().then((d) => d?.tasks?.length && setTasks(d.tasks)).catch(() => {});
  }, []);

  const onStart = async (id: number) => {
    haptic("medium");
    setTasks((t) => t.map((x) => (x.id === id ? { ...x, status: "in_progress" } : x)));
    api.startTask(id).catch(() => {});
  };
  const onComplete = async (id: number) => {
    haptic("success");
    setTasks((t) => t.filter((x) => x.id !== id));
    try { await api.completeTask(id); await refresh(); } catch {}
  };

  return (
    <div>
      <PageHeader
        title="Задания"
        subtitle="Прокачивай характеристики реальным действием"
        action={
          <button className="btn-primary px-3 py-2 text-sm" onClick={() => { haptic("light"); setOpenNew(true); }}>
            <Plus className="h-4 w-4" /> Новое
          </button>
        }
      />

      <div className="space-y-3">
        {tasks.map((t, i) => {
          const meta = typeMeta[t.task_type] ?? typeMeta.intelligence;
          const Icon = meta.icon;
          return (
            <motion.div
              key={t.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.04 }}
              className="glass-card p-4"
            >
              <div className="flex items-start gap-3">
                <div className={`rounded-xl bg-muted p-2.5 ${meta.color}`}>
                  <Icon className="h-5 w-5" />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-start justify-between gap-2">
                    <h3 className="truncate font-display text-base font-semibold">{t.title}</h3>
                    <span className="text-[10px] uppercase tracking-wider text-muted-foreground">
                      {"★".repeat(t.difficulty || 1)}
                    </span>
                  </div>
                  {t.description && (
                    <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">{t.description}</p>
                  )}
                  <div className="mt-2 flex flex-wrap items-center gap-1.5">
                    <span className="stat-chip text-primary-glow">+{t.xp_reward} XP</span>
                    <span className="stat-chip text-warning">+{t.gold_reward}g</span>
                    <span className="stat-chip">⚡ {t.energy_cost}</span>
                  </div>
                  <div className="mt-3 flex gap-2">
                    {t.status === "in_progress" ? (
                      <button className="btn-primary flex-1 py-2 text-sm" onClick={() => onComplete(t.id)}>
                        <CheckCircle2 className="h-4 w-4" /> Завершить
                      </button>
                    ) : (
                      <button className="btn-ghost flex-1 py-2 text-sm" onClick={() => onStart(t.id)}>
                        <Play className="h-4 w-4" /> Начать
                      </button>
                    )}
                  </div>
                </div>
              </div>
            </motion.div>
          );
        })}
        {tasks.length === 0 && (
          <div className="glass-card p-8 text-center text-muted-foreground">
            Нет активных заданий. Создай первое!
          </div>
        )}
      </div>

      {openNew && <NewTaskSheet onClose={() => setOpenNew(false)} onCreated={(t) => setTasks((p) => [t, ...p])} />}
    </div>
  );
}

function NewTaskSheet({ onClose, onCreated }: { onClose: () => void; onCreated: (t: any) => void }) {
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [type, setType] = useState("intelligence");
  const [saving, setSaving] = useState(false);

  const save = async () => {
    if (!title.trim()) return;
    setSaving(true);
    try {
      const created = await api.createTask({ title, description: desc, task_type: type });
      onCreated(created ?? { id: Date.now(), title, description: desc, task_type: type, xp_reward: 30, gold_reward: 10, energy_cost: 5, status: "active", difficulty: 2 });
      haptic("success");
      onClose();
    } catch {
      onCreated({ id: Date.now(), title, description: desc, task_type: type, xp_reward: 30, gold_reward: 10, energy_cost: 5, status: "active", difficulty: 2 });
      onClose();
    } finally { setSaving(false); }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-background/70 backdrop-blur" onClick={onClose}>
      <motion.div
        initial={{ y: "100%" }} animate={{ y: 0 }} exit={{ y: "100%" }}
        transition={{ type: "spring", damping: 26, stiffness: 220 }}
        className="w-full max-w-md rounded-t-3xl border-t border-border bg-card p-5 pb-8"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mx-auto mb-4 h-1.5 w-12 rounded-full bg-muted" />
        <h2 className="font-display text-xl font-bold">Новое задание</h2>
        <p className="mb-4 text-sm text-muted-foreground">ИИ-Сенсей рассчитает награду</p>
        <input
          className="mb-3 w-full rounded-xl border border-border bg-muted px-4 py-3 text-sm outline-none focus:border-primary"
          placeholder="Что сделаешь?"
          value={title} onChange={(e) => setTitle(e.target.value)}
        />
        <textarea
          className="mb-3 w-full resize-none rounded-xl border border-border bg-muted px-4 py-3 text-sm outline-none focus:border-primary"
          placeholder="Описание (необязательно)"
          rows={3} value={desc} onChange={(e) => setDesc(e.target.value)}
        />
        <div className="mb-4 grid grid-cols-4 gap-2">
          {(["strength", "agility", "intelligence", "insight"] as const).map((t) => {
            const m = typeMeta[t];
            const Icon = m.icon;
            const active = type === t;
            return (
              <button
                key={t}
                onClick={() => setType(t)}
                className={
                  "flex flex-col items-center gap-1 rounded-xl border p-2 text-[11px] transition " +
                  (active ? "border-primary bg-primary/10 text-foreground" : "border-border bg-muted/40 text-muted-foreground")
                }
              >
                <Icon className={`h-4 w-4 ${active ? m.color : ""}`} />
                {m.label}
              </button>
            );
          })}
        </div>
        <button className="btn-primary w-full" disabled={saving || !title.trim()} onClick={save}>
          {saving ? "Сохраняем..." : "Создать задание"}
        </button>
      </motion.div>
    </div>
  );
}