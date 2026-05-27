import { Crown, Check, Zap, Sparkles, Shield, Flame } from "lucide-react";
import PageHeader from "@/components/PageHeader";
import { haptic } from "@/lib/telegram";

const PERKS = [
  { icon: Zap, text: "+50 энергии в день" },
  { icon: Sparkles, text: "Безлимит запросов к Сенсею" },
  { icon: Shield, text: "Защита от рейдов" },
  { icon: Flame, text: "x2 опыт за задания" },
];

export default function DojoPlus() {
  return (
    <div>
      <PageHeader title="Dojo+" subtitle="Премиум-лицензия охотника" />

      <div className="glass-card relative overflow-hidden p-6 text-center">
        <div className="absolute inset-0 bg-gradient-aurora opacity-50" />
        <div className="relative">
          <div className="mx-auto mb-3 grid h-14 w-14 place-items-center rounded-2xl bg-gradient-primary shadow-glow">
            <Crown className="h-7 w-7 text-primary-foreground" />
          </div>
          <h2 className="font-display text-2xl font-bold">Стань S-ранговым</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Открой все возможности Додзё
          </p>

          <div className="mt-5 space-y-2 text-left">
            {PERKS.map(({ icon: Icon, text }) => (
              <div key={text} className="flex items-center gap-3 rounded-xl bg-muted/60 p-3">
                <div className="rounded-lg bg-gradient-primary p-1.5 text-primary-foreground">
                  <Icon className="h-4 w-4" />
                </div>
                <span className="flex-1 text-sm">{text}</span>
                <Check className="h-4 w-4 text-success" />
              </div>
            ))}
          </div>

          <div className="mt-5 grid grid-cols-2 gap-3">
            <PlanCard title="Месяц" price="299₽" />
            <PlanCard title="Год" price="2 490₽" badge="-30%" highlight />
          </div>

          <button className="btn-primary mt-4 w-full" onClick={() => haptic("success")}>
            Активировать Dojo+
          </button>
        </div>
      </div>
    </div>
  );
}

function PlanCard({ title, price, badge, highlight }: { title: string; price: string; badge?: string; highlight?: boolean }) {
  return (
    <div className={
      "relative rounded-2xl border p-4 text-center " +
      (highlight ? "border-primary bg-primary/10 shadow-glow" : "border-border bg-muted/40")
    }>
      {badge && (
        <span className="absolute -top-2 right-2 rounded-full bg-gradient-primary px-2 py-0.5 text-[10px] font-bold text-primary-foreground">
          {badge}
        </span>
      )}
      <div className="text-xs text-muted-foreground">{title}</div>
      <div className="font-display text-xl font-bold">{price}</div>
    </div>
  );
}