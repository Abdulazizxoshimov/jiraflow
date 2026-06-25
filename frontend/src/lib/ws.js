const WS_URL = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/v1/ws`

const RECONNECT_DELAYS = [1000, 2000, 5000, 10000, 30000]

class WSClient {
  constructor() {
    this._socket = null
    this._token = null
    this._handlers = {}       // type → Set<fn>
    this._rooms = new Set()   // subscribed rooms
    this._reconnectAttempt = 0
    this._reconnectTimer = null
    this._stopped = false
  }

  connect(token) {
    this._token = token
    this._stopped = false
    this._openSocket()
  }

  disconnect() {
    this._stopped = true
    clearTimeout(this._reconnectTimer)
    if (this._socket) {
      this._socket.close()
      this._socket = null
    }
    this._rooms.clear()
    this._handlers = {}
  }

  subscribe(room) {
    this._rooms.add(room)
    if (this._socket?.readyState === WebSocket.OPEN) {
      this._send({ type: 'subscribe', room })
    }
  }

  unsubscribe(room) {
    this._rooms.delete(room)
    if (this._socket?.readyState === WebSocket.OPEN) {
      this._send({ type: 'unsubscribe', room })
    }
  }

  on(type, handler) {
    if (!this._handlers[type]) this._handlers[type] = new Set()
    this._handlers[type].add(handler)
    return () => this._handlers[type]?.delete(handler)
  }

  off(type, handler) {
    this._handlers[type]?.delete(handler)
  }

  _openSocket() {
    if (!this._token) return
    const url = `${WS_URL}?token=${encodeURIComponent(this._token)}`
    const socket = new WebSocket(url)

    socket.onopen = () => {
      this._reconnectAttempt = 0
      // Re-subscribe to all rooms after reconnect
      for (const room of this._rooms) {
        this._send({ type: 'subscribe', room })
      }
    }

    socket.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        this._dispatch(msg.type, msg)
      } catch (_) {}
    }

    socket.onclose = () => {
      this._socket = null
      if (!this._stopped) {
        this._scheduleReconnect()
      }
    }

    socket.onerror = () => {
      socket.close()
    }

    this._socket = socket
  }

  _send(data) {
    if (this._socket?.readyState === WebSocket.OPEN) {
      this._socket.send(JSON.stringify(data))
    }
  }

  _dispatch(type, msg) {
    const fns = this._handlers[type]
    if (fns) {
      for (const fn of fns) {
        try { fn(msg) } catch (_) {}
      }
    }
    // wildcard handlers
    const all = this._handlers['*']
    if (all) {
      for (const fn of all) {
        try { fn(msg) } catch (_) {}
      }
    }
  }

  _scheduleReconnect() {
    const delay = RECONNECT_DELAYS[Math.min(this._reconnectAttempt, RECONNECT_DELAYS.length - 1)]
    this._reconnectAttempt++
    this._reconnectTimer = setTimeout(() => this._openSocket(), delay)
  }
}

export const ws = new WSClient()
