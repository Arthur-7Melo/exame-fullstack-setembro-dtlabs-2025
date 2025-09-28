import { useEffect, useRef, useState, useCallback } from "react";
import { useAuthContext } from "@/contexts/authContext";

type WebSocketMessage = {
  id: string;
  user_id: string;
  name: string;
  description: string;
  device_sn: string;
  triggered_value: number;
  timestamp: string;
  heartbeat_data: {
    cpu: number;
    ram: number;
    disk_free: number;
    temperature: number;
    latency: number;
    connectivity: number;
  };
};

type Listener = (msg: WebSocketMessage) => void;

type ManagerEntry = {
  ws: WebSocket | null;
  listeners: Set<Listener>;
  refCount: number;
  isConnected: boolean;
  reconnectTimer?: number | null;
  intentionallyClosed: boolean;
};

const socketManager = new Map<string, ManagerEntry>();

function buildWsUrl(userId: string) {
  const envUrl = (import.meta.env.REACT_APP_WS_URL || "").trim();
  if (envUrl) {
    return envUrl.endsWith("?") || envUrl.includes("?")
      ? `${envUrl}&user_id=${userId}`
      : `${envUrl}?user_id=${userId}`;
  }

  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const host = window.location.hostname || "localhost";
  const port = import.meta.env.REACT_APP_WS_PORT || "8080";
  return `${protocol}://${host}:${port}/ws/notifications?user_id=${userId}`;
}

function ensureEntry(userId: string): ManagerEntry {
  let entry = socketManager.get(userId);
  if (!entry) {
    entry = {
      ws: null,
      listeners: new Set(),
      refCount: 0,
      isConnected: false,
      reconnectTimer: null,
      intentionallyClosed: false,
    };
    socketManager.set(userId, entry);
  }
  return entry;
}

function createSocket(userId: string) {
  const entry = ensureEntry(userId);
  if (entry.ws && (entry.ws.readyState === WebSocket.OPEN || entry.ws.readyState === WebSocket.CONNECTING)) {
    return;
  }

  const url = buildWsUrl(userId);
  entry.intentionallyClosed = false;
  const ws = new WebSocket(url);
  entry.ws = ws;

  ws.onopen = () => {
    entry.isConnected = true;
    console.debug("[WS Manager] connected", userId);
  };

  ws.onmessage = (ev) => {
    try {
      const data = JSON.parse(ev.data);
      entry.listeners.forEach((l) => {
        try { l(data); } catch (e) { console.error("listener error", e); }
      });
    } catch (err) {
      console.error("[WS Manager] failed to parse message", err, ev.data);
    }
  };

  ws.onclose = (ev) => {
    entry.isConnected = false;
    entry.ws = null;
    console.debug("[WS Manager] closed", userId, "code:", ev.code, "reason:", ev.reason);

    if (!entry.intentionallyClosed && entry.refCount > 0) {
      if (entry.reconnectTimer) {
        window.clearTimeout(entry.reconnectTimer);
      }
      entry.reconnectTimer = window.setTimeout(() => {
        entry.reconnectTimer = null;
        console.debug("[WS Manager] reconnecting...", userId);
        createSocket(userId);
      }, 5000) as unknown as number;
    }
  };

  ws.onerror = (ev) => {
    console.error("[WS Manager] socket error", userId, ev);
  };
}

function closeSocketIfUnneeded(userId: string) {
  const entry = socketManager.get(userId);
  if (!entry) return;
  if (entry.refCount <= 0) {
    entry.intentionallyClosed = true;
    if (entry.reconnectTimer) {
      window.clearTimeout(entry.reconnectTimer);
      entry.reconnectTimer = null;
    }
    if (entry.ws) {
      try {
        entry.ws.close(1000, "no listeners");
      } catch (e) {
      }
    }
    socketManager.delete(userId);
    console.debug("[WS Manager] socket closed and entry removed", userId);
  }
}

function addManagerListener(userId: string, listener: Listener) {
  const entry = ensureEntry(userId);
  entry.listeners.add(listener);
  entry.refCount = entry.refCount + 1;

  if (!entry.ws) createSocket(userId);

  return () => {
    const e = socketManager.get(userId);
    if (!e) return;
    e.listeners.delete(listener);
    e.refCount = Math.max(0, e.refCount - 1);
    if (e.refCount === 0) {
      closeSocketIfUnneeded(userId);
    }
  };
}

export const useWebSocket = () => {
  const { user } = useAuthContext();
  const [messages, setMessages] = useState<WebSocketMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const cleanupRef = useRef<(() => void) | null>(null);

  const onMessage = useCallback((msg: WebSocketMessage) => {
    setMessages(prev => [msg, ...prev.slice(0, 9)]);
  }, []);

  useEffect(() => {
    if (!user?.id) {
      if (cleanupRef.current) {
        cleanupRef.current();
        cleanupRef.current = null;
      }
      setIsConnected(false);
      return;
    }

    const userId = user.id;
    const entry = ensureEntry(userId);

    setIsConnected(entry.isConnected);

    const listener: Listener = (m) => {
      onMessage(m);
    };

    const cleanup = addManagerListener(userId, listener);
    cleanupRef.current = cleanup;

    const interval = window.setInterval(() => {
      const e = socketManager.get(userId);
      setIsConnected(Boolean(e?.isConnected));
    }, 500);

    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
        cleanupRef.current = null;
      }
      window.clearInterval(interval);
    };
  }, [user?.id, onMessage]);

  useEffect(() => {
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
        cleanupRef.current = null;
      }
    };
  }, []);

  return { messages, isConnected };
};
