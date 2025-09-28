export const formatNumber = (value: number): string => {
  if (value === null || value === undefined) return '0';

  const formatted = Number(value.toFixed(2));

  return formatted % 1 === 0 ? formatted.toString() : formatted.toString();
};

export const formatPercentage = (value: number): string => {
  return `${formatNumber(value)}%`;
};

export const formatTemperature = (value: number): string => {
  return `${formatNumber(value)}°C`;
};

export const formatLatency = (value: number): string => {
  return `${formatNumber(value)}ms`;
};

export const formatDiskSpace = (value: number): string => {
  return `${formatNumber(value)}GB`;
};