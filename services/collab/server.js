'use strict';

const { Server } = require('@hocuspocus/server');
const { Database } = require('@hocuspocus/extension-database');
const jwt = require('jsonwebtoken');
const { Pool } = require('pg');

const PORT     = parseInt(process.env.COLLAB_PORT || '1234', 10);
const JWT_SECRET = process.env.JWT_SECRET || '';
const DB_URL   = process.env.DATABASE_URL ||
  `postgresql://${process.env.DB_USER || 'postgres'}:${process.env.DB_PASSWORD || 'postgres'}@${process.env.DB_HOST || 'postgres'}:${process.env.DB_PORT || '5432'}/${process.env.DB_NAME || 'jiraflow'}?sslmode=disable`;

const pool = new Pool({ connectionString: DB_URL });

// Ensure the ydoc_store table exists
async function ensureTable() {
  await pool.query(`
    CREATE TABLE IF NOT EXISTS ydoc_store (
      document_name TEXT PRIMARY KEY,
      data          BYTEA NOT NULL,
      updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    )
  `);
}
ensureTable().catch((e) => console.error('DB init error:', e));

const server = Server.configure({
  port: PORT,
  timeout: 30000,

  async onAuthenticate({ token, documentName }) {
    if (!JWT_SECRET) return; // dev mode: skip auth
    try {
      const payload = jwt.verify(token, JWT_SECRET);
      if (payload.type !== 'access') {
        throw new Error('invalid token type');
      }
      // documentName format: "page:<uuid>"
      // In production you could verify the user has read access to this page.
      return { user: { id: payload.sub, role: payload.role } };
    } catch (e) {
      throw new Error('unauthorized');
    }
  },

  extensions: [
    new Database({
      fetch: async ({ documentName }) => {
        const res = await pool.query(
          'SELECT data FROM ydoc_store WHERE document_name = $1',
          [documentName]
        );
        return res.rows[0]?.data ?? null;
      },
      store: async ({ documentName, state }) => {
        await pool.query(
          `INSERT INTO ydoc_store(document_name, data, updated_at)
           VALUES($1, $2, NOW())
           ON CONFLICT (document_name)
           DO UPDATE SET data = EXCLUDED.data, updated_at = NOW()`,
          [documentName, state]
        );
      },
    }),
  ],

  async onLoadDocument(data) {
    console.log(`[collab] document loaded: ${data.documentName}`);
  },

  async onChange(data) {
    // Hook: persist on every change (handled by Database extension above).
  },

  async onDisconnect({ documentName, clientsCount }) {
    if (clientsCount === 0) {
      console.log(`[collab] all clients left: ${documentName}`);
    }
  },
});

server.listen().then(() => {
  console.log(`[collab] Hocuspocus server running on ws://0.0.0.0:${PORT}`);
}).catch((e) => {
  console.error('[collab] startup error:', e);
  process.exit(1);
});
