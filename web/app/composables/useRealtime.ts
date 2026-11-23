import Pusher from 'pusher-js'
import type { Channel } from 'pusher-js'

type RealtimeHandler = (data: unknown) => void

// A single websocket connection shared by the whole app.
let client: Pusher | null = null

function createClient(): Pusher | null {
  const config = useRuntimeConfig()
  const key = config.public.soketiKey as string
  if (!key) return null

  // The key is part of the WebSocket URL (/app/<key>), so characters like
  // "/" or "+" (e.g. a base64 key) break the connection.
  if (!/^[A-Za-z0-9_-]+$/.test(key)) {
    console.error(
      'SOKETI_KEY must only contain letters, numbers, "-" or "_". Generate one with: openssl rand -hex 32',
    )
    return null
  }

  if (!config.public.useSoketi) {
    return new Pusher(key, { cluster: config.public.soketiCluster as string })
  }

  // Self-hosted Soketi: nginx proxies /app/ on the same host to Soketi.
  const secure = window.location.protocol === 'https:'
  const port = window.location.port
    ? Number(window.location.port)
    : secure
      ? 443
      : 80

  return new Pusher(key, {
    cluster: '',
    wsHost: window.location.hostname,
    wsPort: port,
    wssPort: port,
    forceTLS: secure,
    disableStats: true,
    enabledTransports: ['ws', 'wss'],
  })
}

/**
 * Binds handlers to the realtime events of a user channel.
 * Must be called from a component lifecycle hook (e.g. onMounted).
 * Returns a function which removes only the handlers bound here.
 */
export function bindRealtimeEvents(
  userId: string | undefined | null,
  handlers: Record<string, RealtimeHandler>,
): () => void {
  if (import.meta.server || !userId) return () => {}

  client ??= createClient()
  if (!client) return () => {}

  const channel: Channel = client.channel(userId) ?? client.subscribe(userId)
  for (const [event, handler] of Object.entries(handlers)) {
    channel.bind(event, handler)
  }

  return () => {
    for (const [event, handler] of Object.entries(handlers)) {
      channel.unbind(event, handler)
    }
  }
}
