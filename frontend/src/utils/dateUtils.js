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

export const toUserDate = (date) => {
  try {
    const localDate = new Date(date)
    if (Number.isNaN(localDate.getTime())) {
      return ''
    }

    const month = String(localDate.getMonth() + 1).padStart(2, '0')
    const day = String(localDate.getDate()).padStart(2, '0')
    const year = localDate.getFullYear()
    return `${month}/${day}/${year}`
  } catch (error) {
    console.error('Error converting to user date:', error)
    return ''
  }
}