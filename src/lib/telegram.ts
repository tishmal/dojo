// Telegram WebApp типы и хелперы
export interface TgUser {
  id: number;
  first_name?: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
}

interface TgWebApp {
  initData: string;
  initDataUnsafe: { user?: TgUser; [k: string]: any };
  ready: () => void;
  expand: () => void;
  setHeaderColor?: (c: string) => void;
  setBackgroundColor?: (c: string) => void;
  HapticFeedback?: { impactOccurred: (s: string) => void; notificationOccurred: (s: string) => void };
  colorScheme?: "light" | "dark";
}

declare global {
  interface Window {
    Telegram?: { WebApp?: TgWebApp };
  }
}

export function getTg(): TgWebApp | null {
  return window.Telegram?.WebApp ?? null;
}

export function initTelegram() {
  const tg = getTg();
  if (!tg) return;
  try {
    tg.ready();
    tg.expand();
    tg.setHeaderColor?.("#0a0d1a");
    tg.setBackgroundColor?.("#0a0d1a");
  } catch {}
}

export function haptic(type: "light" | "medium" | "heavy" | "success" | "error" = "light") {
  const hf = getTg()?.HapticFeedback;
  if (!hf) return;
  if (type === "success" || type === "error") hf.notificationOccurred(type);
  else hf.impactOccurred(type);
}

export function getTgUser(): TgUser | null {
  return getTg()?.initDataUnsafe?.user ?? null;
}

export function getInitData(): string {
  return getTg()?.initData ?? "";
}