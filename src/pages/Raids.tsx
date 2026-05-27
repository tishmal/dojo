import { Flame, Swords, Clock, Coins } from "lucide-react";
import PageHeader from "@/components/PageHeader";
import { haptic } from "@/lib/telegram";

const TARGETS = [
  { id: 1, name: "ленивый_дракон", rank: "C", level: 18, gold: 420, inactive: "5 дн.", cost: 50 },
  { id: 2, name: "ushedshiy_v_otpusk", rank: "D", level: 11, gold: 180, inactive: "4 дн.", cost: 25 },
  { id: 3, name: "spящий_воин", rank: "B", level: 27, gold: 980, inactive: "7 дн.", cost: 120 },
];

export default function Raids() {
  return (
    <div>
      <PageHeader title="Рейды" subtitle="Налетай на неактивных охотников" />

      <div className="glass-card mb-4 flex items-center gap-3 p-4">
        <div className="rounded-xl bg-danger/10 p-2.5 text-danger"><Flame className="h-5 w-5" /></div>
        <div className="flex-1 text-sm">
          <div className="font-semibold">Активный сезон рейдов</div>
          <div className="text-muted-foreground">Лут x2 первые 24 часа</div>
        </div>
      </div>

      <div className="space-y-3">
        {TARGETS.map((t) => (
          <div key={t.id} className="glass-card p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="grid h-12 w-12 place-items-center rounded-xl bg-gradient-primary font-display text-lg font-bold text-primary-foreground shadow-glow">
                  {t.rank}
                </div>
                <div>
                  <div className="font-semibold">{t.name}</div>
                  <div className="text-xs text-muted-foreground">Ур. {t.level} · <Clock className="inline h-3 w-3" /> {t.inactive}</div>
                </div>
              </div>
              <span className="stat-chip text-warning"><Coins className="h-3 w-3" /> {t.gold}</span>
            </div>
            <button
              onClick={() => haptic("heavy")}
              className="btn-primary mt-3 w-full text-sm"
            >
              <Swords className="h-4 w-4" /> Рейд за {t.cost}g
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}