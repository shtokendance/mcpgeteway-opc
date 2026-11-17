import React from 'react';

import iconMap from '@/icons/icon-map.json';

const icons = import.meta.glob<string>('../assets/icons/*.svg', { eager: true, query: '?raw', import: 'default' });

interface LocalIconProps {
  icon: string;
  className?: string;
  style?: React.CSSProperties;
  width?: number | string;
  height?: number | string;
  onClick?: () => void;
}

const LocalIcon: React.FC<LocalIconProps> = ({
  icon,
  className = '',
  style = {},
  width = 24,
  height = 24,
  onClick
}) => {
  const iconPath = iconMap[icon as keyof typeof iconMap];
  const svg = iconPath && icons[`../assets${iconPath}`];

  if (!svg) {
    return null;
  }

  const svgContent = svg
    .replace(/width="[^"]*"/, `width="${width}"`)
    .replace(/height="[^"]*"/, `height="${height}"`);

  return (
    <span
      className={`inline-flex items-center justify-center ${className}`}
      style={style}
      onClick={onClick}
      onKeyDown={onClick ? (e) => e.key === 'Enter' && onClick() : undefined}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      dangerouslySetInnerHTML={{ __html: svgContent }}
    />
  );
};

export default LocalIcon;