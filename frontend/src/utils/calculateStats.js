// Calculate the % for clicks and submit
export function calculatePercentage(total, num) {
    try {
        if (total === 0) {
            return 0;
        }
        return Math.round((num / total) * 100);
    } catch (error) {
        console.error('Error calculating percentage:', error);
        return 0;
    }
}
