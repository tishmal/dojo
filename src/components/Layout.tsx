import { NavLink, Outlet, useLocation } from "react-router-dom";
import { Swords, User, Flame, Sparkles, Crown } from "lucide-react";
import { motion } from "framer-motion";
import { haptic } from "@/lib/telegram";
import StatusBar from "./StatusBar";

const tabs = [
  { to: "/tasks", label: "Задания", icon: Swords },
  { to: "/raids", label: "Рейды", icon: Flame },
  { to: "/profile", label: "Профиль", icon: User },
  { to: "/sensei", label: "Сенсей", icon: Sparkles },
  { to: "/plus", label: "Dojo+", icon: Crown },
];

export default function Layout() {
  const { pathname } = useLocation();
  return (
    <div className="mx-auto flex h-full max-w-md flex-col">
      <StatusBar />
      <main className="flex-1 overflow-y-auto px-4 pb-28 pt-2">
        <motion.div
          key={pathname}
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.25 }}
        >
          <Outlet />
        </motion.div>
      </main>
      <nav className="fixed bottom-0 left-1/2 z-50 w-full max-w-md -translate-x-1/2 px-3 pb-[max(0.5rem,env(safe-area-inset-bottom))] pt-2">
        <div className="glass-card flex items-stretch gap-1 px-2 py-1.5">
          {tabs.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              onClick={() => haptic("light")}
              className="tab-btn group"
              data-active={pathname === to}
            >
              {({ isActive }) => (
                <>
                  <div
                    className={
                      "relative rounded-xl p-2 transition " +
                      (isActive
                        ? "bg-gradient-primary text-primary-foreground shadow-glow"
                        : "text-muted-foreground")
                    }
                  >
                    <Icon className="h-4 w-4" strokeWidth={2.4} />
                  </div>
                  <span>{label}</span>
                </>
              )}
            </NavLink>
          ))}
        </div>
      </nav>
    </div>
  );
}