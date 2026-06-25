import { useState, useMemo } from 'react';
import { useApp } from '../store/AppContext';
import { useApi } from '../api/api';
import { Icon } from '../components/icons';
import { Empty, PriorityBadge } from '../components/components';

const PRIORITY_COLOR = { critical: '#EF4444', high: '#F97316', medium: '#F59E0B', low: '#10B981' };
const MONTH_SHORT = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];

function addMonths(date, n) {
  const d = new Date(date);
  d.setMonth(d.getMonth() + n);
  return d;
}

function diffDays(a, b) {
  return Math.round((b - a) / 86400000);
}

function parseDate(s) {
  if (!s) return null;
  const d = new Date(s);
  return isNaN(d) ? null : d;
}

// Build month columns between viewStart and viewEnd
function buildMonths(viewStart, viewEnd) {
  const months = [];
  let cur = new Date(viewStart.getFullYear(), viewStart.getMonth(), 1);
  while (cur <= viewEnd) {
    months.push(new Date(cur));
    cur = addMonths(cur, 1);
  }
  return months;
}

const ROW_H = 40;
const LABEL_W = 280;
const COL_W = 120;

export function RoadmapView({ nav }) {
  const { activeProjectId, projects } = useApp();
  const proj = projects.find((p) => p.id === activeProjectId);
  const { data, loading } = useApi(
    activeProjectId ? `/projects/${activeProjectId}/roadmap` : null,
    [activeProjectId]
  );

  const items = data?.items || data || [];

  // View window: 3 months back → 9 months forward from today
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const [viewOffset, setViewOffset] = useState(0); // months offset
  const viewStart = useMemo(() => {
    const d = new Date(today.getFullYear(), today.getMonth() - 3 + viewOffset, 1);
    return d;
  }, [viewOffset]);
  const viewEnd = useMemo(() => addMonths(viewStart, 12), [viewStart]);

  const months = useMemo(() => buildMonths(viewStart, viewEnd), [viewStart, viewEnd]);
  const totalDays = diffDays(viewStart, viewEnd);

  function xForDate(date) {
    const d = diffDays(viewStart, date);
    return Math.round((d / totalDays) * (months.length * COL_W));
  }

  const todayX = xForDate(today);
  const svgW = months.length * COL_W;

  // Flatten items (epics + children)
  const rows = useMemo(() => {
    const result = [];
    (items || []).forEach((item) => {
      result.push({ ...item, depth: 0 });
      (item.children || []).forEach((child) => {
        result.push({ ...child, depth: 1 });
      });
    });
    return result;
  }, [items]);

  if (!activeProjectId) {
    return (
      <div style={{ display: 'grid', placeItems: 'center', height: '60vh' }}>
        <Empty icon="map" title="No project selected" hint="Select a project from the sidebar." />
      </div>
    );
  }

  if (loading) {
    return (
      <div style={{ padding: 40 }}>
        <div className="skel" style={{ height: 32, width: 220, marginBottom: 24 }} />
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="skel" style={{ height: ROW_H - 6, marginBottom: 8, borderRadius: 6 }} />
        ))}
      </div>
    );
  }

  if (!rows.length) {
    return (
      <div>
        <RoadmapHeader proj={proj} viewOffset={viewOffset} setViewOffset={setViewOffset} />
        <div style={{ display: 'grid', placeItems: 'center', height: '50vh' }}>
          <Empty icon="calendar" title="No roadmap items" hint="Create epics with start and due dates to see them here." />
        </div>
      </div>
    );
  }

  const svgH = rows.length * ROW_H + 48;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', overflow: 'hidden' }}>
      <RoadmapHeader proj={proj} viewOffset={viewOffset} setViewOffset={setViewOffset} />

      <div style={{ flex: 1, overflow: 'auto', display: 'flex' }}>
        {/* Left: issue labels */}
        <div style={{ flexShrink: 0, width: LABEL_W, borderRight: '1px solid var(--border)', background: 'var(--bg)' }}>
          {/* header spacer */}
          <div style={{ height: 48, borderBottom: '1px solid var(--border)' }} />
          {rows.map((row, i) => (
            <div
              key={row.id}
              style={{
                height: ROW_H,
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                padding: `0 12px 0 ${12 + row.depth * 20}px`,
                borderBottom: '1px solid var(--border)',
                fontSize: 13,
                cursor: 'pointer',
                background: i % 2 === 0 ? 'transparent' : 'var(--bg-subtle)',
              }}
              onClick={() => nav('issue', row.id)}
            >
              <span style={{
                width: 8, height: 8, borderRadius: 2, flexShrink: 0,
                background: PRIORITY_COLOR[row.priority] || '#94A3B8',
              }} />
              <span style={{
                overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                color: 'var(--text)', fontWeight: row.depth === 0 ? 600 : 400,
              }}>
                {row.title}
              </span>
              {row.progress > 0 && (
                <span style={{ marginLeft: 'auto', fontSize: 11, color: 'var(--text-muted)', flexShrink: 0 }}>
                  {Math.round(row.progress)}%
                </span>
              )}
            </div>
          ))}
        </div>

        {/* Right: SVG timeline */}
        <div style={{ flex: 1, overflow: 'auto', position: 'relative' }}>
          <svg width={svgW} height={svgH} style={{ display: 'block' }}>
            {/* Month headers */}
            {months.map((m, i) => (
              <g key={i}>
                <rect x={i * COL_W} y={0} width={COL_W} height={48}
                  fill={m.getMonth() % 2 === 0 ? 'var(--bg)' : 'var(--bg-subtle)'}
                  stroke="var(--border)" strokeWidth={0.5} />
                <text x={i * COL_W + COL_W / 2} y={28} textAnchor="middle"
                  fontSize={12} fontWeight={600} fill="var(--text-muted)">
                  {MONTH_SHORT[m.getMonth()]}
                </text>
                <text x={i * COL_W + COL_W / 2} y={42} textAnchor="middle"
                  fontSize={10} fill="var(--text-muted)">
                  {m.getFullYear()}
                </text>
              </g>
            ))}

            {/* Row backgrounds + grid lines */}
            {rows.map((_, i) => (
              <g key={i}>
                <rect
                  x={0} y={48 + i * ROW_H} width={svgW} height={ROW_H}
                  fill={i % 2 === 0 ? 'transparent' : 'var(--bg-subtle)'}
                />
                {months.map((_, mi) => (
                  <line key={mi} x1={mi * COL_W} y1={48 + i * ROW_H} x2={mi * COL_W} y2={48 + i * ROW_H + ROW_H}
                    stroke="var(--border)" strokeWidth={0.5} />
                ))}
              </g>
            ))}

            {/* Today line */}
            {todayX >= 0 && todayX <= svgW && (
              <>
                <line x1={todayX} y1={0} x2={todayX} y2={svgH}
                  stroke="var(--indigo-500, #6366F1)" strokeWidth={1.5} strokeDasharray="4 3" />
                <text x={todayX + 4} y={14} fontSize={10} fill="var(--indigo-600, #4F46E5)" fontWeight={600}>
                  Today
                </text>
              </>
            )}

            {/* Bars */}
            {rows.map((row, i) => {
              const start = parseDate(row.start_date);
              const end = parseDate(row.due_date);
              if (!start && !end) return null;

              const effectiveStart = start || end;
              const effectiveEnd = end || start;

              const x1 = Math.max(0, xForDate(effectiveStart));
              const x2 = Math.min(svgW, xForDate(effectiveEnd) + COL_W / 30);
              const barW = Math.max(6, x2 - x1);
              const y = 48 + i * ROW_H + (ROW_H - 22) / 2;
              const barColor = PRIORITY_COLOR[row.priority] || '#6366F1';
              const barH = row.depth === 0 ? 22 : 16;
              const radius = row.depth === 0 ? 5 : 3;

              return (
                <g key={row.id}>
                  {/* Background track */}
                  <rect x={x1} y={48 + i * ROW_H + (ROW_H - barH) / 2}
                    width={barW} height={barH} rx={radius}
                    fill={barColor} opacity={0.18} />
                  {/* Progress fill */}
                  {row.progress > 0 && (
                    <rect
                      x={x1} y={48 + i * ROW_H + (ROW_H - barH) / 2}
                      width={Math.max(radius * 2, barW * (row.progress / 100))}
                      height={barH} rx={radius}
                      fill={barColor} opacity={0.75}
                    />
                  )}
                  {/* Label inside bar */}
                  {barW > 60 && (
                    <text
                      x={x1 + 8} y={48 + i * ROW_H + ROW_H / 2 + 1}
                      fontSize={11} fontWeight={600}
                      fill={barColor} dominantBaseline="middle"
                      style={{ pointerEvents: 'none' }}
                    >
                      {row.title.length > 20 ? row.title.slice(0, 20) + '…' : row.title}
                    </text>
                  )}
                </g>
              );
            })}
          </svg>
        </div>
      </div>
    </div>
  );
}

function RoadmapHeader({ proj, viewOffset, setViewOffset }) {
  return (
    <div className="page-head" style={{ paddingBottom: 12 }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}>
        <div>
          <div className="crumbs">
            {proj?.name} <Icon name="chevronRight" size={11} /> <span>Roadmap</span>
          </div>
          <h1 style={{ margin: 0 }}>Roadmap</h1>
          <p style={{ margin: '4px 0 0', color: 'var(--text-muted)', fontSize: 13 }}>
            Epic timeline — items with start or due dates appear here.
          </p>
        </div>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <button className="btn btn-ghost" data-size="sm" onClick={() => setViewOffset(0)}>
            Today
          </button>
          <button className="btn btn-ghost" data-size="sm" onClick={() => setViewOffset((o) => o - 3)}>
            <Icon name="chevronLeft" size={13} />
          </button>
          <button className="btn btn-ghost" data-size="sm" onClick={() => setViewOffset((o) => o + 3)}>
            <Icon name="chevronRight" size={13} />
          </button>
        </div>
      </div>
    </div>
  );
}
