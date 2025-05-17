export const formatDate = (date) => {
  try {
    const now = new Date(date);
    const utcDateOnly = now.toISOString().split('T')[0];
    const year = utcDateOnly.split('-')[0];
    const month = utcDateOnly.split('-')[1];
    const day = utcDateOnly.split('-')[2];
    return `${month}/${day}/${year}`;
  } catch (error) {
    console.error('Error formatting date:', error);
    return '';
  }
} 