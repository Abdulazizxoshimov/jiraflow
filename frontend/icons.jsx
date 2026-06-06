// icons.jsx — Tabler-style line icons (stroke only)
// Usage: <Icon name="home" size={16} />

const ICONS = {
  // brand
  forge: <><path d="M12 3l8 4.5v9L12 21l-8-4.5v-9L12 3z"/><path d="M12 12l8-4.5"/><path d="M12 12v9"/><path d="M12 12L4 7.5"/></>,
  // nav
  home: <><path d="M5 12L12 5l7 7"/><path d="M5 10v10h14V10"/></>,
  checkbox: <><rect x="4" y="4" width="16" height="16" rx="3"/><path d="M9 12l2 2 4-4"/></>,
  folder: <><path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z"/></>,
  kanban: <><rect x="4" y="4" width="6" height="16" rx="2"/><rect x="14" y="4" width="6" height="10" rx="2"/></>,
  list: <><path d="M9 6h11M9 12h11M9 18h11"/><circle cx="5" cy="6" r="1"/><circle cx="5" cy="12" r="1"/><circle cx="5" cy="18" r="1"/></>,
  rocket: <><path d="M4 13c2-5 6-9 12-9 0 6-4 10-9 12l-3-3z"/><path d="M9 15l-4 4 1-5"/><circle cx="14" cy="10" r="1.5"/></>,
  notes: <><rect x="5" y="3" width="14" height="18" rx="2"/><path d="M9 7h6M9 11h6M9 15h4"/></>,
  users: <><circle cx="9" cy="9" r="3"/><path d="M3 19c.5-3 3-5 6-5s5.5 2 6 5"/><circle cx="17" cy="8" r="2"/><path d="M16 14c2 0 4 1.5 5 4"/></>,
  chart: <><path d="M4 4v16h16"/><path d="M8 16l3-4 3 2 4-6"/></>,
  settings: <><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.8-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1.1-1.5 1.7 1.7 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.8 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1A1.7 1.7 0 0 0 4.6 9a1.7 1.7 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.8.3H9a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.8V9a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1z"/></>,
  briefcase: <><rect x="3" y="7" width="18" height="13" rx="2"/><path d="M9 7V5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/><path d="M3 13h18"/></>,
  bell: <><path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"/><path d="M10 21a2 2 0 0 0 4 0"/></>,
  shield: <><path d="M12 3l8 3v6c0 5-4 8-8 9-4-1-8-4-8-9V6l8-3z"/></>,
  // top
  search: <><circle cx="11" cy="11" r="7"/><path d="M21 21l-4-4"/></>,
  plus: <><path d="M12 5v14M5 12h14"/></>,
  filter: <><path d="M4 5h16l-6 8v6l-4-2v-4L4 5z"/></>,
  more: <><circle cx="12" cy="6" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="18" r="1"/></>,
  moreH: <><circle cx="6" cy="12" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="18" cy="12" r="1"/></>,
  x: <><path d="M6 6l12 12M18 6L6 18"/></>,
  check: <><path d="M5 12l5 5L20 7"/></>,
  chevronDown: <><path d="M6 9l6 6 6-6"/></>,
  chevronRight: <><path d="M9 6l6 6-6 6"/></>,
  chevronLeft: <><path d="M15 6l-6 6 6 6"/></>,
  arrowUp: <><path d="M12 19V5M5 12l7-7 7 7"/></>,
  arrowDown: <><path d="M12 5v14M5 12l7 7 7-7"/></>,
  arrowRight: <><path d="M5 12h14M13 6l6 6-6 6"/></>,
  // brand-telegram (paper plane)
  telegram: <><path d="M21 4L3 11l6 2 2 6 3-4 5 4 2-15z"/><path d="M9 13l8-6"/></>,
  // issue types
  bug: <><rect x="6" y="8" width="12" height="11" rx="6"/><path d="M9 8V6a3 3 0 0 1 6 0v2M3 13h3M3 8h2M3 18h3M21 13h-3M21 8h-2M21 18h-3M12 12v6"/></>,
  task: <><rect x="4" y="4" width="16" height="16" rx="3"/><path d="M9 12l2 2 4-4"/></>,
  story: <><path d="M5 4h11l3 3v13H5z"/><path d="M16 4v4h3"/><path d="M8 11h7M8 15h7M8 7h4"/></>,
  epic: <><path d="M13 2l-1 9h6l-8 11 1-9H5l8-11z"/></>,
  // priority arrows
  prHigh: <><path d="M12 19V5M5 12l7-7 7 7"/></>,
  prLow: <><path d="M12 5v14M5 12l7 7 7-7"/></>,
  prMed: <><path d="M5 12h14"/></>,
  // misc
  link: <><path d="M10 14l-2 2a4 4 0 0 1-5.6-5.6L5 8M14 10l2-2a4 4 0 0 1 5.6 5.6L19 16M9 15l6-6"/></>,
  attach: <><path d="M21 12l-8 8a5 5 0 0 1-7-7l9-9a4 4 0 0 1 5 5l-9 9a3 3 0 0 1-4-4l8-8"/></>,
  calendar: <><rect x="4" y="5" width="16" height="16" rx="2"/><path d="M16 3v4M8 3v4M4 11h16"/></>,
  clock: <><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></>,
  comment: <><path d="M21 12a8 8 0 1 1-3.2-6.4L21 4l-1 4a8 8 0 0 1 1 4z"/></>,
  history: <><path d="M3 12a9 9 0 1 0 3-6.7"/><path d="M3 4v5h5"/><path d="M12 8v5l3 2"/></>,
  pencil: <><path d="M4 20h4l11-11-4-4L4 16v4z"/></>,
  trash: <><path d="M4 7h16M9 7V5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2M6 7l1 13a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-13"/></>,
  copy: <><rect x="8" y="8" width="12" height="12" rx="2"/><path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2"/></>,
  eye: <><path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12z"/><circle cx="12" cy="12" r="3"/></>,
  eyeOff: <><path d="M3 3l18 18"/><path d="M10 6c.6-.1 1.3-.1 2-.1 6 0 10 6.1 10 6.1s-1 1.7-3 3.5"/><path d="M6.5 7.5C3.5 9.5 2 12 2 12s4 6 10 6c1.8 0 3.4-.5 4.7-1.3"/><circle cx="12" cy="12" r="3"/></>,
  bolt: <><path d="M13 3L4 14h7l-1 7 9-11h-7l1-7z"/></>,
  flag: <><path d="M5 21V4l8 2-2 5 2 5-8-2"/></>,
  tag: <><path d="M3 12V4h8l10 10-8 8-10-10z"/><circle cx="7.5" cy="7.5" r="1.5"/></>,
  upload: <><path d="M21 15v3a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-3M7 9l5-5 5 5M12 4v12"/></>,
  download: <><path d="M21 15v3a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-3M7 11l5 5 5-5M12 16V4"/></>,
  sun: <><circle cx="12" cy="12" r="4"/><path d="M12 3v2M12 19v2M3 12h2M19 12h2M5.5 5.5l1.4 1.4M17.1 17.1l1.4 1.4M5.5 18.5l1.4-1.4M17.1 6.9l1.4-1.4"/></>,
  moon: <><path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z"/></>,
  grid: <><rect x="4" y="4" width="7" height="7" rx="1.5"/><rect x="13" y="4" width="7" height="7" rx="1.5"/><rect x="4" y="13" width="7" height="7" rx="1.5"/><rect x="13" y="13" width="7" height="7" rx="1.5"/></>,
  send: <><path d="M22 2L11 13"/><path d="M22 2l-7 20-4-9-9-4 20-7z"/></>,
  user: <><circle cx="12" cy="8" r="4"/><path d="M4 20c1-4 4.5-6 8-6s7 2 8 6"/></>,
  exit: <><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><path d="M16 17l5-5-5-5M21 12H9"/></>,
  database: <><ellipse cx="12" cy="5" rx="8" ry="3"/><path d="M4 5v6c0 1.7 3.6 3 8 3s8-1.3 8-3V5M4 11v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6"/></>,
  server: <><rect x="3" y="4" width="18" height="6" rx="2"/><rect x="3" y="14" width="18" height="6" rx="2"/><circle cx="7" cy="7" r="1"/><circle cx="7" cy="17" r="1"/></>,
  cloud: <><path d="M6 17a4 4 0 1 1 .8-7.9 6 6 0 0 1 11.6 1.4A4 4 0 0 1 18 17H6z"/></>,
  branch: <><circle cx="6" cy="5" r="2"/><circle cx="18" cy="5" r="2"/><circle cx="12" cy="19" r="2"/><path d="M6 7v3a3 3 0 0 0 3 3h6a3 3 0 0 0 3-3V7M12 13v4"/></>,
  star: <><path d="M12 2l3 6.5 7 .8-5.2 4.7 1.5 7L12 17.5 5.7 21l1.5-7L2 9.3l7-.8L12 2z"/></>,
  refresh: <><path d="M21 12a9 9 0 1 1-3-6.7"/><path d="M21 4v5h-5"/></>,
  sparkle: <><path d="M12 3l1.5 4.5L18 9l-4.5 1.5L12 15l-1.5-4.5L6 9l4.5-1.5L12 3zM19 3l.7 2L22 6l-2.3.5L19 9l-.7-2.5L16 6l2.3-1L19 3z"/></>,
  layout: <><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/></>,
  inbox: <><path d="M3 13l3-8h12l3 8M3 13v6a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-6M3 13h5l2 3h4l2-3h5"/></>,
  archive: <><rect x="3" y="4" width="18" height="4" rx="1"/><path d="M5 8v11a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8M10 12h4"/></>,
  externalLink: <><path d="M14 5h5v5"/><path d="M19 5L10 14"/><path d="M19 13v4a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2h4"/></>,
  globe: <><circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3a14 14 0 0 1 0 18M12 3a14 14 0 0 0 0 18"/></>,
  code: <><path d="M9 8l-5 4 5 4M15 8l5 4-5 4M14 4l-4 16"/></>,
  picture: <><rect x="3" y="5" width="18" height="14" rx="2"/><circle cx="9" cy="11" r="2"/><path d="M21 17l-5-5-9 9"/></>,
  table: <><rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 10h18M9 5v14"/></>,
  bold: <><path d="M7 5h6a3.5 3.5 0 0 1 0 7H7zM7 12h7a3.5 3.5 0 0 1 0 7H7z"/></>,
  italic: <><path d="M19 4h-9M14 20H5M15 4L9 20"/></>,
  heading: <><path d="M6 5v14M18 5v14M6 12h12"/></>,
  paperclip: <><path d="M21 12l-8 8a5 5 0 0 1-7-7l9-9a4 4 0 0 1 5 5l-9 9a3 3 0 0 1-4-4l8-8"/></>,
  at: <><circle cx="12" cy="12" r="4"/><path d="M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-4 8"/></>,
  lock: <><rect x="5" y="11" width="14" height="10" rx="2"/><path d="M8 11V8a4 4 0 0 1 8 0v3"/></>,
  mail: <><rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 7l9 6 9-6"/></>,
  collapse: <><path d="M9 6l-6 6 6 6M3 12h11M21 4v16"/></>,
};

function Icon({ name, size = 16, color = "currentColor", strokeWidth = 1.75, style, className }) {
  const path = ICONS[name];
  if (!path) return <span style={{display:"inline-block",width:size,height:size}}/>;
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size} height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke={color}
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      style={style}
      className={className}
      aria-hidden="true"
    >{path}</svg>
  );
}

window.Icon = Icon;
