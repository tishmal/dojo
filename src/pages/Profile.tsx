import { Dumbbell, Brain, Flame, Eye, Trophy, Zap, Coins } from "lucide-react";
import PageHeader from "@/components/PageHeader";
import { useUser } from "@/store/user";

export default function Profile() {
  const { user } = useUser();
  if (!user) return null;
  const xpPct = Math.min(100, Math.round((user.xp / user.xp_to_next_level) * 100));
  return (
    <div>
      <PageHeader title="Профиль" subtitle="Твоя статистика охотника" />

      <div className="glass-card relative overflow-hidden p-5">
        <div className="absolute inset-0 bg-gradient-aurora opacity-40" />
        <div className="relative">
          <div className="flex items-center gap-4">
            <div className="relative h-20 w-20 overflow-hidden rounded-2xl border border-border/80 bg-muted">
              {user.photo_url ? (
                <img src={user.photo_url} alt="" className="h-full w-full object-cover" />
              ) : (
                <div className="flex h-full w-full items-center justify-center font-display text-3xl font-bold">
                  {user.first_name?.[0]}
                </div>
              )}
            </div>
            <div className="flex-1">
              <div className="font-display text-xl font-bold">{user.first_name}</div>
              <div className="text-sm text-muted-foreground">@{user.username || "охотник"}</div>
              <div className="mt-2 flex items-center gap-2">
                <span className="rounded-md bg-gradient-primary px-2 py-0.5 text-xs font-bold text-primary-foreground">
                  Ур. {user.level}
                </span>
                <span className="stat-chip"><Trophy className="h-3 w-3 text-warning" /> Ранг {rank(user.level)}</span>
              </div>
            </div>
          </div>

          <div className="mt-5">
            <div className="mb-1 flex items-center justify-between text-xs text-muted-foreground">
              <span>Опыт</span>
              <span>{user.xp} / {user.xp_to_next_level}</span>
            </div>
            <div className="h-2.5 overflow-hidden rounded-full bg-muted">
              <div
                className="h-full bg-gradient-primary shadow-glow transition-all"
                style={{ width: `${xpPct}%` }}
              />
            </div>
          </div>
        </div>
      </div>

      <div className="mt-3 grid grid-cols-2 gap-3">
        <ResourceCard icon={Coins} label="Золото" value={user.gold} color="text-warning" />
        <ResourceCard icon={Zap} label="Энергия" value={`${user.energy}/${user.max_energy}`} color="text-primary-glow" />
      </div>

      <h2 className="mb-3 mt-6 font-display text-lg font-semibold">Характеристики</h2>
      <div className="grid grid-cols-2 gap-3">
        <StatCard icon={Dumbbell} label="Сила" value={user.strength} color="text-danger" />
        <StatCard icon={Flame} label="Ловкость" value={user.agility} color="text-warning" />
        <StatCard icon={Brain} label="Интеллект" value={user.intelligence} color="text-primary-glow" />
        <StatCard icon={Eye} label="Интуиция" value={user.insight} color="text-accent" />
      </div>
    </div>
  );
}

function StatCard({ icon: Icon, label, value, color }: any) {
  return (
    <div className="glass-card flex items-center gap-3 p-4">
      <div className={`rounded-xl bg-muted p-2.5 ${color}`}><Icon className="h-5 w-5" /></div>
      <div>
        <div className="text-xs text-muted-foreground">{label}</div>
        <div className="font-display text-xl font-bold">{value}</div>
      </div>
    </div>
  );
}
function ResourceCard(p: any) { return <StatCard {...p} />; }
function rank(level: number) {
  if (level >= 50) return "S"; if (level >= 40) return "A"; if (level >= 30) return "B";
  if (level >= 20) return "C"; if (level >= 10) return "D"; return "E";
}