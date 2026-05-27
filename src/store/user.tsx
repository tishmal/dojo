import { createContext, useContext, useEffect, useState, ReactNode, useCallback } from "react";
import { api } from "@/lib/api";
import { getTgUser, initTelegram } from "@/lib/telegram";

export interface UserState {
  id: number;
  telegram_id: number;
  username: string;
  first_name: string;
  photo_url?: string;
  level: number;
  xp: number;
  xp_to_next_level: number;
  gold: number;
  energy: number;
  max_energy: number;
  strength: number;
  agility: number;
  intelligence: number;
  insight: number;
  license_active: boolean;
  sensei_requests: number;
}

interface Ctx {
  user: UserState | null;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
}

const UserCtx = createContext<Ctx>({ user: null, loading: true, error: null, refresh: async () => {} });

const DEMO_USER: UserState = {
  id: 1, telegram_id: 0, username: "охотник", first_name: "Охотник",
  level: 7, xp: 340, xp_to_next_level: 600, gold: 1250, energy: 78, max_energy: 100,
  strength: 12, agility: 9, intelligence: 14, insight: 8,
  license_active: true, sensei_requests: 4,
};

export function UserProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserState | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      // Пытаемся получить профиль с backend
      const data = await api.getProfile();
      setUser(data);
    } catch (e: any) {
      // Fallback: показываем демо-данные с данными из Telegram, если есть
      const tg = getTgUser();
      setUser({
        ...DEMO_USER,
        telegram_id: tg?.id ?? 0,
        username: tg?.username ?? DEMO_USER.username,
        first_name: tg?.first_name ?? DEMO_USER.first_name,
        photo_url: tg?.photo_url,
      });
      setError(e?.message ?? "offline");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    initTelegram();
    load();
  }, [load]);

  return (
    <UserCtx.Provider value={{ user, loading, error, refresh: load }}>
      {children}
    </UserCtx.Provider>
  );
}

export const useUser = () => useContext(UserCtx);