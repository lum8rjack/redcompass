export const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit' 
  })
}

export const formatProjectDate = (date) => {
  const now = new Date(date);
  const utcDateOnly = now.toISOString().split('T')[0];
  return utcDateOnly;
} 