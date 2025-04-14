export const formatDate = (date) => {
  const now = new Date(date);
  const utcDateOnly = now.toISOString().split('T')[0];
  const year = utcDateOnly.split('-')[0];
  const month = utcDateOnly.split('-')[1];
  const day = utcDateOnly.split('-')[2];
  return `${month}/${day}/${year}`;
} 