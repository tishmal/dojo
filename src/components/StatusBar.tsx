import { Coins, Zap, Shield, ShieldAlert } from "lucide-react";
import { useUser } from "@/store/user";

export default function StatusBar() {
  const { user } = useUser();
  if (!user) return <div className="h-14" />;
  const energyPct = Math.round((user.energy / user.max_energy) * 100);
  return (
    <header className="sticky top-0 z-30 px-4 pb-3 pt-3 backdrop-blur-xl">
      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <div className="relative h-10 w-10 overflow-hidden rounded-full border border-border/80 bg-muted">
            {user.photo_url ? (
              <img src={user.photo_url} alt="" className="h-full w-full object-cover" />
            ) : (
              <div className="flex h-full w-full items-center justify-center font-display font-bold text-foreground/80">
                {user.first_name?.[0] ?? "?"}
              </div>
            )}
            <span className="absolute -bottom-0.5 -right-0.5 grid h-5 w-5 place-items-center rounded-full bg-gradient-primary text-[10px] font-bold text-primary-foreground shadow-glow">
              {user.level}
            </span>
          </div>
          <div className="leading-tight">
            <div className="text-sm font-semibold">{user.first_name || user.username}</div>
            <div className="text-[11px] text-muted-foreground">Ранг {rank(user.level)}</div>
          </div>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="stat-chip text-warning"><Coins className="h-3.5 w-3.5" /> {user.gold}</span>
          <span className="stat-chip"><Zap className="h-3.5 w-3.5 text-primary-glow" /> {energyPct}%</span>
          <span className="stat-chip">
            {user.license_active ? <Shield className="h-3.5 w-3.5 text-success" /> : <ShieldAlert className="h-3.5 w-3.5 text-danger" />}
          </span>
        </div>
      </div>
    </header>
  );
}

function rank(level: number) {
  if (level >= 50) return "S";
  if (level >= 40) return "A";
  if (level >= 30) return "B";
  if (level >= 20) return "C";
  if (level >= 10) return "D";
  return "E";
}